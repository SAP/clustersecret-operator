/*
Copyright (c) 2023 SAP SE

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package framework

import (
	"fmt"
	"strconv"
	"sync"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

type syncItem struct {
	gvk             schema.GroupVersionKind
	namespace       string
	name            string
	resourceVersion string
}

func (item *syncItem) String() string {
	if item.namespace == "" {
		return fmt.Sprintf("%s %s %s (%s)", item.gvk.GroupVersion(), item.gvk.Kind, item.name, item.resourceVersion)
	} else {
		return fmt.Sprintf("%s %s %s/%s (%s)", item.gvk.GroupVersion(), item.gvk.Kind, item.namespace, item.name, item.resourceVersion)
	}
}

type Synchronizer struct {
	listers   map[schema.GroupVersionKind]cache.GenericLister
	items     map[types.UID]*syncItem
	observers map[chan struct{}]chan struct{}
	mutex     sync.Mutex
}

func newSynchronizer() *Synchronizer {
	return &Synchronizer{
		listers:   make(map[schema.GroupVersionKind]cache.GenericLister),
		items:     make(map[types.UID]*syncItem),
		observers: make(map[chan struct{}]chan struct{}),
	}
}

func (s *Synchronizer) Init(informers map[schema.GroupVersionKind]cache.SharedIndexInformer) {
	s.lock()
	defer s.unlock()
	for gvk, informer := range informers {
		informer.AddEventHandler(
			cache.ResourceEventHandlerFuncs{
				AddFunc:    s.handleAddEvent,
				UpdateFunc: s.handleUpdateEvent,
				DeleteFunc: s.handleDeleteEvent,
			},
		)
		// todo: discover resource for gvk (but this requires a discovery client)
		// setting it to 'unknown' is a bit hacky, but will probaably not harm since NewGenericLister will use it only when assembling errors
		s.listers[gvk] = cache.NewGenericLister(informer.GetIndexer(), schema.GroupResource{Group: gvk.Group, Resource: "unknown"})
	}
}

func (s *Synchronizer) WaitUntilSynced() {
	var ch chan struct{}
	for {
		s.lock()
		if ch == nil {
			ch = make(chan struct{}, 1)
			s.observers[ch] = ch
		}
		for uid, item := range s.items {
			lister, ok := s.listers[item.gvk]
			if !ok {
				panic("encountered unexpected object kind")
			}
			var obj runtime.Object
			var err error
			if item.namespace == "" {
				obj, err = lister.Get(item.name)
			} else {
				obj, err = lister.ByNamespace(item.namespace).Get(item.name)
			}
			if err == nil {
				objmeta := obj.(metav1.Object)
				if item.resourceVersion != "-" && objmeta.GetUID() == uid && s.compareResourceversion(objmeta.GetResourceVersion(), item.resourceVersion) >= 0 {
					klog.V(3).Infof("synchronize (wait): clearing creation/update %s (have: %s)", item, objmeta.GetResourceVersion())
					delete(s.items, uid)
				} else {
					klog.V(3).Infof("synchronize (wait): not clearing %s (have: %s)", item, objmeta.GetResourceVersion())
				}
			} else if errors.IsNotFound(err) {
				if item.resourceVersion == "-" {
					klog.V(3).Infof("synchronize (wait): clearing deletion %s", item)
					delete(s.items, uid)
				} else {
					klog.V(3).Infof("synchronize (wait): not clearing %s", item)
				}
			} else {
				klog.Warningf("synchronize (wait): error while getting %s: %s", item, err)
			}
		}
		if len(s.items) == 0 {
			delete(s.observers, ch)
			s.unlock()
			return
		}
		s.unlock()
		<-ch
	}
}

func (s *Synchronizer) lock() {
	s.mutex.Lock()
}

func (s *Synchronizer) unlock() {
	s.mutex.Unlock()
}

func (s *Synchronizer) handleCreation(new runtime.Object) {
	s.onAction("create", new)
}

func (s *Synchronizer) handleUpdate(old runtime.Object, new runtime.Object) {
	s.onAction("update", new)
}

func (s *Synchronizer) handleDeletion(old runtime.Object) {
	s.onAction("delete", old)
}

func (s *Synchronizer) handleAddEvent(new interface{}) {
	s.onEvent("add", new)
}

func (s *Synchronizer) handleUpdateEvent(old interface{}, new interface{}) {
	s.onEvent("update", new)
}

func (s *Synchronizer) handleDeleteEvent(old interface{}) {
	s.onEvent("delete", old)
}

func (s *Synchronizer) onAction(action string, obj runtime.Object) {
	s.lock()
	defer s.unlock()
	if _, ok := s.listers[obj.GetObjectKind().GroupVersionKind()]; !ok {
		return
	}
	objmeta := obj.(metav1.Object)
	uid := objmeta.GetUID()
	resourceVersion := ""
	switch action {
	case "create", "update":
		resourceVersion = objmeta.GetResourceVersion()
	case "delete":
		resourceVersion = "-"
	default:
		panic("this cannot happen")
	}
	if uid == "" {
		panic("encountered object with empty uid")
	}
	if resourceVersion == "" {
		panic("encountered object with empty resource version")
	}
	if item, ok := s.items[uid]; ok {
		if resourceVersion == "-" {
			item.resourceVersion = resourceVersion
		} else if item.resourceVersion != "-" && s.compareResourceversion(resourceVersion, item.resourceVersion) > 0 {
			item.resourceVersion = resourceVersion
		}
		klog.V(3).Infof("synchronize (action: %s): recording(updating) queue item: %s", action, item)
	} else {
		s.items[uid] = &syncItem{
			gvk:             obj.GetObjectKind().GroupVersionKind(),
			namespace:       objmeta.GetNamespace(),
			name:            objmeta.GetName(),
			resourceVersion: resourceVersion,
		}
		klog.V(3).Infof("synchronize (action: %s): recording(adding) queue item: %s", action, s.items[uid])
	}
	for ch := range s.observers {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

func (s *Synchronizer) onEvent(event string, arg interface{}) {
	s.lock()
	defer s.unlock()
	if obj, ok := arg.(runtime.Object); ok {
		objmeta := obj.(metav1.Object)
		groupVersion := obj.GetObjectKind().GroupVersionKind().GroupVersion()
		kind := obj.GetObjectKind().GroupVersionKind().Kind
		namespace := objmeta.GetNamespace()
		name := objmeta.GetName()
		resourceVersion := objmeta.GetResourceVersion()
		if namespace == "" {
			klog.V(3).Infof("synchronize (event: %s): handling event %s %s %s (%s)", event, groupVersion, kind, name, resourceVersion)
		} else {
			klog.V(3).Infof("synchronize (event: %s): handling event %s %s %s/%s (%s)", event, groupVersion, kind, namespace, name, resourceVersion)
		}
	} else {
		klog.V(3).Infof("synchronize (event: %s): handling event with unspecified object", event)
	}
	for ch := range s.observers {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

func (s *Synchronizer) compareResourceversion(x string, y string) int {
	i, err := strconv.Atoi(x)
	if err != nil {
		panic(err)
	}
	j, err := strconv.Atoi(y)
	if err != nil {
		panic(err)
	}
	if i < j {
		return -1
	}
	if i > j {
		return 1
	}
	return 0
}
