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

package controller

import (
	"context"
	"testing"

	"github.com/sap/clustersecret-operator/test"
)

// test: create clustersecrets
func TestReconcile1(t *testing.T) {
	env := test.NewEnvironment()
	env.SetBasePath("testdata/1")

	env.AddObjectsFromFiles(
		"namespace.yaml",
		"clustersecret.yaml",
	)

	ctx, cancel := context.WithCancel(context.Background())
	c := NewController(ctx, env.KubernetesClient(), env.CoreClient(), env.NewSynchronizer())
	c.startInformers()
	defer cancel()

	c.reconcileClusterSecret("my-secret")
	env.MustError(t).AssertSecretFromFile("secret.yaml")
}

// test: update namespaces
func TestReconcile2(t *testing.T) {
	env := test.NewEnvironment()
	env.SetBasePath("testdata/2")

	env.AddObjectsFromFiles(
		"clustersecret.yaml",
		"namespace-1.yaml",
		"namespace-2.yaml",
		"namespace-3.yaml",
		"secret-1.yaml",
		"secret-3.yaml",
	)

	ctx, cancel := context.WithCancel(context.Background())
	c := NewController(ctx, env.KubernetesClient(), env.CoreClient(), env.NewSynchronizer())
	c.startInformers()
	defer cancel()

	c.reconcileClusterSecret("my-secret")
	env.MustError(t).AssertSecretCount("", "clustersecrets.core.cs.sap.com/name=my-secret", 2)
	env.MustError(t).AssertSecretFromFile("secret-1.yaml")
	env.MustError(t).AssertSecretFromFile("secret-2.yaml")

	env.MustFatal(t).UnlabelNamespace("my-namespace-1", "mylabel")
	env.MustFatal(t).LabelNamespace("my-namespace-2", "mylabel", "othervalue")
	env.MustFatal(t).LabelNamespace("my-namespace-3", "mylabel", "myvalue")
	c.reconcileClusterSecret("my-secret")
	env.MustError(t).AssertSecretCount("", "clustersecrets.core.cs.sap.com/name=my-secret", 1)
	env.MustError(t).AssertSecretFromFile("secret-3.yaml")
}
