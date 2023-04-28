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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/testing"
)

type listerWatcher struct {
	fake testing.FakeClient
	gvr  schema.GroupVersionResource
	gvk  schema.GroupVersionKind
}

func (env *environmentImpl) newListerWatcher(gvk schema.GroupVersionKind) *listerWatcher {
	return &listerWatcher{fake: env.client(gvk.GroupVersion()), gvr: env.groupVersionResource(gvk), gvk: gvk}
}

func (lw *listerWatcher) List(opts metav1.ListOptions) (runtime.Object, error) {
	objlist, err := lw.fake.Invokes(testing.NewRootListAction(lw.gvr, lw.gvk, opts), nil)
	if err != nil {
		return nil, err
	}
	if objlist == nil {
		panic("this cannot happen")
	}
	return objlist, err
}

func (lw *listerWatcher) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return lw.fake.InvokesWatch(testing.NewRootWatchAction(lw.gvr, opts))
}
