/*
SPDX-FileCopyrightText: 2026 SAP SE or an SAP affiliate company and clustersecret-operator contributors
SPDX-License-Identifier: Apache-2.0
*/

package framework

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/uuid"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/testing"
	watchtools "k8s.io/client-go/tools/watch"
)

type Environment interface {
	RegisterCreateCallback(CreateCallbackFunc)
	RegisterUpdateCallback(UpdateCallbackFunc)
	RegisterDeleteCallback(DeleteCallbackFunc)
	SetBasePath(string)
	Decoder() runtime.Decoder
	DynamicClient() *dynamicfake.FakeDynamicClient
	NewSynchronizer() *Synchronizer
	LoadObjectFromFile(string) runtime.Object
	AddObject(runtime.Object)
	AddObjectFromFile(string)
	AddObjectsFromFiles(...string)
	WithObject(runtime.Object) Environment
	WithObjectFromFile(string) Environment
	WithObjectsFromFiles(...string) Environment
	AssertObject(runtime.Object) error
	AssertObjectFromFile(string) error
	AssertObjectCount(schema.GroupVersionKind, string, string, int) error
	GetObject(schema.GroupVersionKind, string, string) (runtime.Object, error)
	ListObjects(schema.GroupVersionKind, string, string) ([]runtime.Object, error)
	CreateObject(runtime.Object) (runtime.Object, error)
	CreateObjectFromFile(string) (runtime.Object, error)
	UpdateObject(runtime.Object) (runtime.Object, error)
	UpdateObjectFromFile(string) (runtime.Object, error)
	PatchObject(schema.GroupVersionKind, string, string, types.PatchType, []byte) (runtime.Object, error)
	LabelObject(schema.GroupVersionKind, string, string, string, string) (runtime.Object, error)
	UnlabelObject(schema.GroupVersionKind, string, string, string) (runtime.Object, error)
	DeleteObject(schema.GroupVersionKind, string, string) error
	WaitForObject(runtime.Object, ...watchtools.ConditionFunc) (runtime.Object, error)
	WaitForObjectFromFile(string, ...watchtools.ConditionFunc) (runtime.Object, error)
}

func NewEnvironment(groupVersions []*GroupVersion) Environment {
	env := &environmentImpl{
		clients: make([]testing.FakeClient, len(groupVersions)),
		schemes: make([]*runtime.Scheme, len(groupVersions)),
		scheme:  runtime.NewScheme(),
	}
	for i, groupVersion := range groupVersions {
		env.clients[i] = groupVersion.Client()
		env.schemes[i] = runtime.NewScheme()
		if err := groupVersion.AddToScheme(env.schemes[i]); err != nil {
			panic(err)
		}
		if err := groupVersion.AddToScheme(env.scheme); err != nil {
			panic(err)
		}
		if env.clients[i] != nil {
			env.clients[i].PrependReactor("*", "*", env.createReactor(env.clients[i]))
			env.clients[i].PrependWatchReactor("*", env.createWatchReactor(env.clients[i]))
		}
	}
	env.decoder = serializer.NewCodecFactory(env.scheme).UniversalDeserializer()
	env.dynamicClient = dynamicfake.NewSimpleDynamicClient(env.scheme)
	env.dynamicClient.PrependReactor("*", "*", env.createReactor(env.dynamicClient))
	env.dynamicClient.PrependWatchReactor("*", env.createWatchReactor(env.dynamicClient))
	env.groupVersionResources = make(map[schema.GroupVersionKind]schema.GroupVersionResource)
	env.groupVersionKinds = make(map[schema.GroupVersionResource]schema.GroupVersionKind)
	for gvk := range env.scheme.AllKnownTypes() {
		// todo: guessing gvr from gvk is a bit ugly; should be solved better somehow ...
		// in addition, this method is not very correct in several cases
		// maybe we should write our own, or make the mapping partially configurable
		gvr, _ := meta.UnsafeGuessKindToResource(gvk)
		env.groupVersionResources[gvk] = gvr
		env.groupVersionKinds[gvr] = gvk
	}
	return env
}

type environmentImpl struct {
	clients               []testing.FakeClient                                    // list of clients
	schemes               []*runtime.Scheme                                       // schemes per client
	scheme                *runtime.Scheme                                         // combined scheme
	decoder               runtime.Decoder                                         // decoder
	dynamicClient         *dynamicfake.FakeDynamicClient                          // dynamic client
	groupVersionResources map[schema.GroupVersionKind]schema.GroupVersionResource // map gvk to gvr
	groupVersionKinds     map[schema.GroupVersionResource]schema.GroupVersionKind // map gvr to gvk
	resourceCounter       int64                                                   // internal counter (to mock resrouce version)
	createCallbacks       []CreateCallbackFunc                                    // create callbacks
	updateCallbacks       []UpdateCallbackFunc                                    // update callbacks
	deleteCallbacks       []DeleteCallbackFunc                                    // delete callbacks
	basePath              string                                                  // base path for file operations
}

func (env *environmentImpl) nextResourceVersion() string {
	return fmt.Sprintf("%08d", atomic.AddInt64(&env.resourceCounter, 1))
}

func (env *environmentImpl) tweakAttributes(object metav1.Object) {
	ownerRefs := object.GetOwnerReferences()
	if ownerRefs != nil {
		for i := 0; i < len(ownerRefs); i++ {
			ownerRefs[i].UID = ""
		}
		object.SetOwnerReferences(ownerRefs)
	}
}

func (env *environmentImpl) hasManagedAttributes(object metav1.Object) bool {
	if _, ok := object.GetAnnotations()["tracker/content-hash"]; ok {
		return true
	}
	return false
}

func (env *environmentImpl) initializeManagedAttributes(object metav1.Object) {
	hash := contentHash(object)
	annotations := object.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations["tracker/content-hash"] = hash
	object.SetAnnotations(annotations)
	object.SetUID(types.UID(uuid.New().String()))
	object.SetResourceVersion(env.nextResourceVersion())
	object.SetGeneration(1)
}

func (env *environmentImpl) adjustManagedAttributes(object metav1.Object, oldObject metav1.Object) {
	hash := contentHash(object)
	annotations := object.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations["tracker/content-hash"] = hash
	object.SetAnnotations(annotations)
	// todo: needed to (re-)set uid here?
	object.SetUID(oldObject.GetUID())
	object.SetResourceVersion(env.nextResourceVersion())
	if hash == oldObject.GetAnnotations()["tracker/content-hash"] {
		object.SetGeneration(oldObject.GetGeneration())
	} else {
		object.SetGeneration(oldObject.GetGeneration() + 1)
	}
}

func (env *environmentImpl) clearManagedAttributes(object metav1.Object) {
	annotations := object.GetAnnotations()
	delete(annotations, "tracker/content-hash")
	if len(annotations) == 0 {
		annotations = nil
	}
	object.SetAnnotations(annotations)
	object.SetUID("")
	object.SetResourceVersion("")
	object.SetGeneration(0)
}

func (env *environmentImpl) client(gv schema.GroupVersion) testing.FakeClient {
	for i, client := range env.clients {
		if env.schemes[i].IsVersionRegistered(gv) {
			if client == nil {
				return env.dynamicClient
			} else {
				return client
			}
		}
	}
	panic(fmt.Sprintf("group version %s not supported by test environment", gv))
}

func (env *environmentImpl) asTyped(obj runtime.Object) runtime.Object {
	if unstructured, ok := obj.(*unstructured.Unstructured); ok {
		raw, err := unstructured.MarshalJSON()
		if err != nil {
			panic(err)
		}
		obj, _, err = env.decoder.Decode(raw, nil, nil)
		if err != nil {
			panic(err)
		}
	}
	return obj
}

func (env *environmentImpl) asUnstructured(obj runtime.Object) runtime.Object {
	if _, ok := obj.(*unstructured.Unstructured); !ok {
		content, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			panic(err)
		}
		obj = &unstructured.Unstructured{Object: content}
	}
	return obj
}

func (env *environmentImpl) getObject(gvk schema.GroupVersionKind, namespace string, name string) (runtime.Object, error) {
	gvr := env.groupVersionResource(gvk)
	client := env.client(gvk.GroupVersion())
	obj, err := client.Tracker().Get(gvr, namespace, name)
	if err != nil {
		return nil, err
	}
	return env.asTyped(obj), nil
}

func (env *environmentImpl) listObjects(gvk schema.GroupVersionKind, namespace string, labelSelector string) ([]runtime.Object, error) {
	gvr := env.groupVersionResource(gvk)
	client := env.client(gvk.GroupVersion())
	objlist, err := client.Tracker().List(gvr, gvk, namespace)
	if err != nil {
		return nil, err
	}
	tmpobjs, err := meta.ExtractList(objlist)
	if err != nil {
		panic(err)
	}
	selector, err := labels.Parse(labelSelector)
	if err != nil {
		panic(err)
	}
	objs := make([]runtime.Object, 0)
	for _, obj := range tmpobjs {
		if selector.Matches(labels.Set(obj.(metav1.Object).GetLabels())) {
			objs = append(objs, env.asTyped((obj)))
		}
	}
	return objs, nil
}

func (env *environmentImpl) addObject(obj runtime.Object) error {
	gvk := obj.GetObjectKind().GroupVersionKind()
	client := env.client(gvk.GroupVersion())
	if client == env.dynamicClient {
		obj = env.asUnstructured(obj)
	}
	return client.Tracker().Add(obj)
}

func (env *environmentImpl) RegisterCreateCallback(callback CreateCallbackFunc) {
	env.createCallbacks = append(env.createCallbacks, callback)
}

func (env *environmentImpl) RegisterUpdateCallback(callback UpdateCallbackFunc) {
	env.updateCallbacks = append(env.updateCallbacks, callback)
}

func (env *environmentImpl) RegisterDeleteCallback(callback DeleteCallbackFunc) {
	env.deleteCallbacks = append(env.deleteCallbacks, callback)
}

func (env *environmentImpl) SetBasePath(basePath string) {
	env.basePath = basePath
}

func (env *environmentImpl) Decoder() runtime.Decoder {
	return env.decoder
}

func (env *environmentImpl) DynamicClient() *dynamicfake.FakeDynamicClient {
	return env.dynamicClient
}

func (env *environmentImpl) NewSynchronizer() *Synchronizer {
	s := newSynchronizer()
	env.RegisterCreateCallback(s.handleCreation)
	env.RegisterUpdateCallback(s.handleUpdate)
	env.RegisterDeleteCallback(s.handleDeletion)
	return s
}

func (env *environmentImpl) LoadObjectFromFile(path string) runtime.Object {
	// returned object will have the concrete underlying type
	if !strings.HasPrefix(path, "/") && env.basePath != "" {
		path = strings.TrimSuffix(env.basePath, "/") + "/" + path
	}
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	raw, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	obj, _, err := env.decoder.Decode(raw, nil, nil)
	if err != nil {
		panic(err)
	}
	return obj
}

func (env *environmentImpl) AddObject(obj runtime.Object) {
	obj = obj.DeepCopyObject()
	env.initializeManagedAttributes(obj.(metav1.Object))
	err := env.addObject(obj)
	if err != nil {
		panic(err)
	}
}

func (env *environmentImpl) AddObjectFromFile(path string) {
	env.AddObject(env.LoadObjectFromFile(path))
}

func (env *environmentImpl) AddObjectsFromFiles(paths ...string) {
	for _, path := range paths {
		env.AddObjectFromFile(path)
	}
}

func (env *environmentImpl) WithObject(obj runtime.Object) Environment {
	env.AddObject(obj)
	return env
}

func (env *environmentImpl) WithObjectFromFile(path string) Environment {
	env.AddObjectFromFile(path)
	return env
}

func (env *environmentImpl) WithObjectsFromFiles(paths ...string) Environment {
	env.AddObjectsFromFiles(paths...)
	return env
}

func (env *environmentImpl) AssertObject(obj runtime.Object) error {
	gvk := obj.GetObjectKind().GroupVersionKind()
	namespace := obj.(metav1.Object).GetNamespace()
	name := obj.(metav1.Object).GetName()
	existing, err := env.getObject(gvk, namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			return err
		} else {
			panic(err)
		}
	}
	if !env.hasManagedAttributes(obj.(metav1.Object)) {
		// note: it's safe here to modify the object returned by the tracker (because tracker clones it internally)
		env.clearManagedAttributes(existing.(metav1.Object))
	}
	env.tweakAttributes(obj.(metav1.Object))
	env.tweakAttributes(existing.(metav1.Object))
	if reflect.DeepEqual(obj, existing) {
		return nil
	} else {
		return fmt.Errorf("content mismatch: %s %s/%s", gvk, namespace, name)
	}
}

func (env *environmentImpl) AssertObjectFromFile(path string) error {
	return env.AssertObject(env.LoadObjectFromFile(path))
}

func (env *environmentImpl) AssertObjectCount(gvk schema.GroupVersionKind, namespace string, labelSelector string, count int) error {
	objs, err := env.listObjects(gvk, namespace, labelSelector)
	if err != nil {
		panic(err)
	}
	if len(objs) != count {
		return fmt.Errorf("expected objects: %d, found objects: %d", count, len(objs))
	}
	return nil
}

func (env *environmentImpl) GetObject(gvk schema.GroupVersionKind, namespace string, name string) (runtime.Object, error) {
	// returned objects will have the concrete underlying type
	gvr := env.groupVersionResource(gvk)
	client := env.client(gvk.GroupVersion())
	obj, err := client.Invokes(testing.NewGetAction(gvr, namespace, name), nil)
	if err != nil {
		return nil, err
	}
	if obj == nil {
		panic("this cannot happen")
	}
	return env.asTyped(obj), nil
}

func (env *environmentImpl) ListObjects(gvk schema.GroupVersionKind, namespace string, labelSelector string) ([]runtime.Object, error) {
	// returned objects will have the concrete underlying type
	gvr := env.groupVersionResource(gvk)
	client := env.client(gvk.GroupVersion())
	objlist, err := client.Invokes(testing.NewListAction(gvr, gvk, namespace, metav1.ListOptions{LabelSelector: labelSelector}), nil)
	if err != nil {
		return nil, err
	}
	if objlist == nil {
		panic("this cannot happen")
	}
	tmpobjs, err := meta.ExtractList(objlist)
	if err != nil {
		panic(err)
	}
	objs := make([]runtime.Object, len(tmpobjs))
	for i, obj := range tmpobjs {
		objs[i] = env.asTyped((obj))
	}
	return objs, nil
}

func (env *environmentImpl) CreateObject(obj runtime.Object) (runtime.Object, error) {
	// returned object will have the concrete underlying type
	gvk := obj.GetObjectKind().GroupVersionKind()
	gvr := env.groupVersionResource(gvk)
	namespace := obj.(metav1.Object).GetNamespace()
	client := env.client(gvk.GroupVersion())
	if client == env.dynamicClient {
		obj = env.asUnstructured(obj)
	}
	obj, err := client.Invokes(testing.NewCreateAction(gvr, namespace, obj), nil)
	if err != nil {
		return nil, err
	}
	if obj == nil {
		panic("this cannot happen")
	}
	return env.asTyped(obj), nil
}

func (env *environmentImpl) CreateObjectFromFile(path string) (runtime.Object, error) {
	// returned object will have the concrete underlying type
	return env.CreateObject(env.LoadObjectFromFile(path))
}

func (env *environmentImpl) UpdateObject(obj runtime.Object) (runtime.Object, error) {
	// returned object will have the concrete underlying type
	gvk := obj.GetObjectKind().GroupVersionKind()
	gvr := env.groupVersionResource(gvk)
	namespace := obj.(metav1.Object).GetNamespace()
	client := env.client(gvk.GroupVersion())
	if client == env.dynamicClient {
		obj = env.asUnstructured(obj)
	}
	obj, err := client.Invokes(testing.NewUpdateAction(gvr, namespace, obj), nil)
	if err != nil {
		return nil, err
	}
	if obj == nil {
		panic("this cannot happen")
	}
	return env.asTyped(obj), nil
}

func (env *environmentImpl) UpdateObjectFromFile(path string) (runtime.Object, error) {
	// returned object will have the concrete underlying type
	return env.UpdateObject(env.LoadObjectFromFile(path))
}

func (env *environmentImpl) PatchObject(gvk schema.GroupVersionKind, namespace string, name string, patchType types.PatchType, patch []byte) (runtime.Object, error) {
	// returned object will have the concrete underlying type
	gvr := env.groupVersionResource(gvk)
	client := env.client(gvk.GroupVersion())
	obj, err := client.Invokes(testing.NewPatchAction(gvr, namespace, name, patchType, patch), nil)
	if err != nil {
		return nil, err
	}
	if obj == nil {
		panic("this cannot happen")
	}
	return env.asTyped(obj), nil
}

func (env *environmentImpl) LabelObject(gvk schema.GroupVersionKind, namespace string, name string, key string, value string) (runtime.Object, error) {
	// returned object will have the concrete underlying type
	patch := fmt.Sprintf(`{"metadata":{"labels":{"%s": "%s"}}}`, key, value)
	return env.PatchObject(gvk, namespace, name, types.MergePatchType, []byte(patch))
}

func (env *environmentImpl) UnlabelObject(gvk schema.GroupVersionKind, namespace string, name string, key string) (runtime.Object, error) {
	// returned object will have the concrete underlying type
	patch := fmt.Sprintf(`{"metadata":{"labels":{"%s": null}}}`, key)
	return env.PatchObject(gvk, namespace, name, types.MergePatchType, []byte(patch))
}

func (env *environmentImpl) DeleteObject(gvk schema.GroupVersionKind, namespace string, name string) error {
	// todo: support prerequisite resource version (currently missing in DeleteAction)
	gvr := env.groupVersionResource(gvk)
	client := env.client(gvk.GroupVersion())
	_, err := client.Invokes(testing.NewDeleteAction(gvr, namespace, name), nil)
	if err != nil {
		return err
	}
	return nil
}

func (env *environmentImpl) WaitForObject(obj runtime.Object, conditions ...watchtools.ConditionFunc) (runtime.Object, error) {
	// returned object will have the concrete underlying type
	gvk := obj.GetObjectKind().GroupVersionKind()
	namespace := obj.(metav1.Object).GetNamespace()
	name := obj.(metav1.Object).GetName()
	client := env.client(gvk.GroupVersion())
	if client == env.dynamicClient {
		obj = env.asUnstructured(obj)
	}
	// note: it's crucial to supply only one single condition to UntilWithSync below
	// because the underlying watcher will watch all objects of the kind, not only obj
	// having multiple conditions will not work because UntilWithSync will be satisfied if each condition is fulfilled for at least one object
	// (but not necessarily requiring that all conditions are fulfilled for the same object)
	aggregatedCondition := func(event watch.Event) (bool, error) {
		if event.Object.GetObjectKind().GroupVersionKind() != gvk {
			return false, nil
		}
		if event.Object.(metav1.Object).GetNamespace() != namespace {
			return false, nil
		}
		if event.Object.(metav1.Object).GetName() != name {
			return false, nil
		}
		// todo: consider uid ?
		switch event.Type {
		case watch.Added, watch.Modified:
			if event.Object.(metav1.Object).GetResourceVersion() == obj.(metav1.Object).GetResourceVersion() {
				return false, nil
			}
		case watch.Deleted:
			// nothing to check
		default:
			panic(fmt.Sprintf("unexpected event type: %s", event.Type))
		}
		for _, condition := range conditions {
			res, err := condition(event)
			if err != nil {
				return false, err
			}
			if !res {
				return false, nil
			}
		}
		return true, nil
	}
	backoffMilliseconds := 50
	for {
		if backoffMilliseconds < 500 {
			backoffMilliseconds += 10
		}
		ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(backoffMilliseconds)*time.Millisecond)
		event, err := watchtools.UntilWithSync(ctx, env.newListerWatcher(gvk), obj, nil, aggregatedCondition)
		cancel()
		if err == nil {
			return env.asTyped(event.Object), nil
		}
		if err == wait.ErrWaitTimeout {
			if _, err := env.getObject(gvk, namespace, name); err != nil && errors.IsNotFound(err) {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
}

func (env *environmentImpl) WaitForObjectFromFile(path string, conditions ...watchtools.ConditionFunc) (runtime.Object, error) {
	// returned object will have the concrete underlying type
	return env.WaitForObject(env.LoadObjectFromFile(path), conditions...)
}
