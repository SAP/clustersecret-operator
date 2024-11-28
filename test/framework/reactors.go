/*
SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and clustersecret-operator contributors
SPDX-License-Identifier: Apache-2.0
*/

package framework

import (
	"encoding/json"
	"fmt"
	"reflect"

	jsonpatch "github.com/evanphx/json-patch/v5"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/testing"
)

/*
This reactor is inspired by the original reactor in https://github.com/kubernetes/client-go/blob/master/testing/fixture.go.
In addition it provides:
- faked metadata.uid
- faked metadata.resourceVersion
- support for server-side applypatch
*/

func (env *environmentImpl) createReactor(client testing.FakeClient) func(testing.Action) (bool, runtime.Object, error) {
	tracker := client.Tracker()
	return func(action testing.Action) (bool, runtime.Object, error) {
		gvr := action.GetResource()
		namespace := action.GetNamespace()

		switch action.GetResource() {
		case schema.GroupVersionResource{Group: "", Version: "v1", Resource: "events"}, schema.GroupVersionResource{Group: "events.k8s.io", Version: "v1", Resource: "events"}:
			if namespace == "" {
				namespace = "default"
			}
		}

		switch action := action.(type) {
		case testing.ListActionImpl:
			objlist, err := tracker.List(gvr, action.GetKind(), namespace)
			if err != nil {
				return true, nil, err
			}
			labelSelector := action.GetListRestrictions().Labels
			if labelSelector == nil {
				labelSelector = labels.Everything()
			}
			// todo: field selection is currently not possible with fake clients
			/*
				fieldSelector := action.GetListRestrictions().Fields
				if fieldSelector == nil {
					fieldSelector = fields.Everything()
				}
			*/
			objs, err := meta.ExtractList(objlist)
			if err != nil {
				return true, nil, err
			}
			var matchingobjs []runtime.Object
			for _, obj := range objs {
				if labelSelector.Matches(labels.Set(obj.(metav1.Object).GetLabels())) {
					matchingobjs = append(matchingobjs, obj)
				}
			}
			if err := meta.SetList(objlist, matchingobjs); err != nil {
				return true, nil, err
			}
			return true, objlist, nil
		case testing.GetActionImpl:
			obj, err := tracker.Get(gvr, namespace, action.GetName())
			return true, obj, err
		case testing.CreateActionImpl:
			new := action.GetObject()
			newmeta, err := meta.Accessor(new)
			if err != nil {
				return true, nil, err
			}
			if action.GetSubresource() == "" {
				env.initializeManagedAttributes(newmeta)
				if err := tracker.Create(gvr, new, namespace); err != nil {
					return true, nil, err
				}
			} else {
				old, err := tracker.Get(gvr, namespace, newmeta.GetName())
				if err != nil {
					return true, nil, err
				}
				oldmeta, err := meta.Accessor(old)
				if err != nil {
					return true, nil, err
				}
				env.adjustManagedAttributes(newmeta, oldmeta)
				if err := tracker.Update(gvr, new, namespace); err != nil {
					return true, nil, err
				}
			}
			obj, err := tracker.Get(gvr, namespace, newmeta.GetName())
			if err != nil {
				panic("this cannot happen")
			}
			for _, callback := range env.createCallbacks {
				callback(obj)
			}
			return true, obj, nil
		case testing.UpdateActionImpl:
			new := action.GetObject()
			newmeta, err := meta.Accessor(new)
			if err != nil {
				return true, nil, err
			}
			old, err := tracker.Get(gvr, namespace, newmeta.GetName())
			if err != nil {
				return true, nil, err
			}
			oldmeta, err := meta.Accessor(old)
			if err != nil {
				return true, nil, err
			}
			env.adjustManagedAttributes(newmeta, oldmeta)
			if err := tracker.Update(gvr, new, namespace); err != nil {
				return true, nil, err
			}
			obj, err := tracker.Get(gvr, namespace, newmeta.GetName())
			if err != nil {
				panic("this cannot happen")
			}
			for _, callback := range env.updateCallbacks {
				callback(old, obj)
			}
			return true, obj, nil
		case testing.DeleteActionImpl:
			old, err := tracker.Get(gvr, namespace, action.GetName())
			if err != nil {
				if !errors.IsNotFound(err) {
					return true, nil, err
				}
				old = nil
			}
			if err := tracker.Delete(gvr, namespace, action.GetName()); err != nil {
				return true, nil, err
			}
			if old == nil {
				panic("this cannot happen")
			}
			for _, callback := range env.deleteCallbacks {
				callback(old)
			}
			return true, nil, nil
		case testing.PatchActionImpl:
			if action.GetPatchType() == types.ApplyPatchType {
				exists := false
				old, err := tracker.Get(gvr, namespace, action.GetName())
				if err == nil {
					exists = true
				} else {
					if !errors.IsNotFound(err) {
						return true, nil, err
					}
				}
				new, _, err := env.decoder.Decode(action.GetPatch(), nil, nil)
				if err != nil {
					return true, nil, err
				}
				if _, ok := client.(dynamic.Interface); ok {
					content, err := runtime.DefaultUnstructuredConverter.ToUnstructured(new)
					if err != nil {
						return true, nil, err
					}
					new = &unstructured.Unstructured{Object: content}
				}
				newmeta, err := meta.Accessor(new)
				if err != nil {
					return true, nil, err
				}
				if exists {
					oldmeta, err := meta.Accessor(old)
					if err != nil {
						return true, nil, err
					}
					env.adjustManagedAttributes(newmeta, oldmeta)
					if err := tracker.Update(gvr, new, namespace); err != nil {
						return true, nil, err
					}
				} else {
					env.initializeManagedAttributes(newmeta)
					if err := tracker.Create(gvr, new, namespace); err != nil {
						return true, nil, err
					}
				}
				obj, err := tracker.Get(gvr, namespace, newmeta.GetName())
				if err != nil {
					panic("this cannot happen")
				}
				if exists {
					for _, callback := range env.updateCallbacks {
						callback(old, obj)
					}
				} else {
					for _, callback := range env.createCallbacks {
						callback(obj)
					}
				}
				return true, obj, nil
			} else {
				old, err := tracker.Get(gvr, namespace, action.GetName())
				if err != nil {
					return true, nil, err
				}
				oldmeta, err := meta.Accessor(old)
				if err != nil {
					return true, nil, err
				}
				oldJson, err := json.Marshal(old)
				if err != nil {
					return true, nil, err
				}
				new := old.DeepCopyObject()
				v := reflect.ValueOf(new)
				v.Elem().Set(reflect.New(v.Type().Elem()).Elem())
				newmeta, err := meta.Accessor(new)
				if err != nil {
					return true, nil, err
				}
				switch action.GetPatchType() {
				case types.JSONPatchType:
					patch, err := jsonpatch.DecodePatch(action.GetPatch())
					if err != nil {
						return true, nil, err
					}
					newJson, err := patch.Apply(oldJson)
					if err != nil {
						return true, nil, err
					}
					if err := json.Unmarshal(newJson, new); err != nil {
						return true, nil, err
					}
				case types.MergePatchType:
					patch := action.GetPatch()
					newJson, err := jsonpatch.MergePatch(oldJson, patch)
					if err != nil {
						return true, nil, err
					}
					if err := json.Unmarshal(newJson, new); err != nil {
						return true, nil, err
					}
				case types.StrategicMergePatchType:
					patch := action.GetPatch()
					newJson, err := strategicpatch.StrategicMergePatch(oldJson, patch, new)
					if err != nil {
						return true, nil, err
					}
					if err := json.Unmarshal(newJson, new); err != nil {
						return true, nil, err
					}
				default:
					return true, nil, fmt.Errorf("patch type %s is not supported", action.GetPatchType())
				}
				env.adjustManagedAttributes(newmeta, oldmeta)
				if err := tracker.Update(gvr, new, namespace); err != nil {
					return true, nil, err
				}
				obj, err := tracker.Get(gvr, namespace, newmeta.GetName())
				if err != nil {
					panic("this cannot happen")
				}
				for _, callback := range env.updateCallbacks {
					callback(old, obj)
				}
				return true, obj, nil
			}
		default:
			return false, nil, fmt.Errorf("no reaction implemented for %s", action)
		}
	}
}

/*
This watch reactor is inspired by the original watch reactor in https://github.com/kubernetes/client-go/blob/master/testing/fixture.go.
In addition it provides:
- repeating existing objects to new watchers
*/

func (env *environmentImpl) createWatchReactor(client testing.FakeClient) func(testing.Action) (bool, watch.Interface, error) {
	tracker := client.Tracker()
	return func(action testing.Action) (bool, watch.Interface, error) {
		gvr := action.GetResource()
		namespace := action.GetNamespace()

		switch action.(type) {
		case testing.WatchActionImpl:
			watcher, err := tracker.Watch(gvr, namespace)
			if err != nil {
				return true, nil, err
			}

			objlist, err := tracker.List(gvr, env.groupVersionKind(gvr), namespace)
			if err != nil {
				return true, nil, err
			}
			objs, err := meta.ExtractList(objlist)
			if err != nil {
				return true, nil, err
			}
			for _, obj := range objs {
				// todo: make casting safer
				watcher.(*watch.RaceFreeFakeWatcher).Add(obj)
			}
			return true, watcher, nil
		default:
			return false, nil, fmt.Errorf("no reaction implemented for %s", action)
		}
	}
}
