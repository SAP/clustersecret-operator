/*
SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and clustersecret-operator contributors
SPDX-License-Identifier: Apache-2.0
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
