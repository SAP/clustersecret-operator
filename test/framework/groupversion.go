/*
SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and clustersecret-operator contributors
SPDX-License-Identifier: Apache-2.0
*/

package framework

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/testing"
)

type GroupVersion struct {
	addToScheme func(*runtime.Scheme) error
	client      testing.FakeClient
}

func NewGroupVersion(addToScheme func(*runtime.Scheme) error, client testing.FakeClient) *GroupVersion {
	return &GroupVersion{addToScheme, client}
}

func (g *GroupVersion) AddToScheme(scheme *runtime.Scheme) error {
	return g.addToScheme(scheme)
}

func (g *GroupVersion) Client() testing.FakeClient {
	return g.client
}
