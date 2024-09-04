/*
SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and clustersecret-operator contributors
SPDX-License-Identifier: Apache-2.0
*/

package controller_test

import (
	"context"
	"testing"

	"github.com/sap/clustersecret-operator/internal/controller"
	"github.com/sap/clustersecret-operator/test"
)

// test: create clustersecrets
func TestController1(t *testing.T) {
	env := test.NewEnvironment()
	env.SetBasePath("testdata/1")

	env.AddObjectsFromFiles(
		"namespace.yaml",
	)

	ctx, cancel := context.WithCancel(context.Background())
	c := controller.NewController(ctx, env.KubernetesClient(), env.CoreClient(), nil)
	c.Start()
	defer c.Wait()
	defer cancel()

	clusterSecret := env.MustFatal(t).CreateClusterSecretFromFile("clustersecret.yaml")
	_ = env.MustFatal(t).WaitForClusterSecretReady(clusterSecret)
	env.MustError(t).AssertSecretFromFile("secret.yaml")
}

// test: update namespaces
func TestController2(t *testing.T) {
	env := test.NewEnvironment()
	env.SetBasePath("testdata/2")

	env.AddObjectsFromFiles(
		"clustersecret.yaml",
		"namespace-1.yaml",
		"namespace-2.yaml",
		"namespace-3.yaml",
		"namespace-4.yaml",
		"secret-1.yaml",
		"secret-3.yaml",
	)

	clusterSecret := env.LoadClusterSecretFromFile("clustersecret.yaml")

	ctx, cancel := context.WithCancel(context.Background())
	c := controller.NewController(ctx, env.KubernetesClient(), env.CoreClient(), nil)
	c.Start()
	defer c.Wait()
	defer cancel()

	clusterSecret = env.MustFatal(t).WaitForClusterSecretReady(clusterSecret)
	env.MustError(t).AssertSecretCount("", "clustersecrets.core.cs.sap.com/name=my-secret", 2)
	env.MustError(t).AssertSecretFromFile("secret-1.yaml")
	env.MustError(t).AssertSecretFromFile("secret-2.yaml")

	env.MustFatal(t).UnlabelNamespace("my-namespace-1", "mylabel")
	clusterSecret = env.MustFatal(t).WaitForClusterSecretReady(clusterSecret)
	env.MustError(t).AssertSecretCount("", "clustersecrets.core.cs.sap.com/name=my-secret", 1)
	env.MustError(t).AssertSecretFromFile("secret-2.yaml")

	env.MustFatal(t).LabelNamespace("my-namespace-2", "mylabel", "othervalue")
	clusterSecret = env.MustFatal(t).WaitForClusterSecretReady(clusterSecret)
	env.MustError(t).AssertSecretCount("", "clustersecrets.core.cs.sap.com/name=my-secret", 0)

	env.MustFatal(t).LabelNamespace("my-namespace-4", "mylabel", "myvalue")
	clusterSecret = env.MustFatal(t).WaitForClusterSecretReady(clusterSecret)
	env.MustError(t).AssertSecretCount("", "clustersecrets.core.cs.sap.com/name=my-secret", 1)
	env.MustError(t).AssertSecretFromFile("secret-4.yaml")

	env.CreateNamespaceFromFile("namespace-5.yaml")
	_ = env.MustFatal(t).WaitForClusterSecretReady(clusterSecret)
	env.MustError(t).AssertSecretCount("", "clustersecrets.core.cs.sap.com/name=my-secret", 2)
	env.MustError(t).AssertSecretFromFile("secret-4.yaml")
	env.MustError(t).AssertSecretFromFile("secret-5.yaml")
}

// test: update clustersecrets
func TestController3(t *testing.T) {
	env := test.NewEnvironment()
	env.SetBasePath("testdata/3")

	env.AddObjectsFromFiles(
		"clustersecret-a.yaml",
		"namespace-1.yaml",
		"namespace-2.yaml",
	)

	clusterSecret_a := env.LoadClusterSecretFromFile("clustersecret-a.yaml")

	ctx, cancel := context.WithCancel(context.Background())
	c := controller.NewController(ctx, env.KubernetesClient(), env.CoreClient(), nil)
	c.Start()
	defer c.Wait()
	defer cancel()

	clusterSecret_a = env.MustFatal(t).WaitForClusterSecretReady(clusterSecret_a)
	env.MustError(t).AssertSecretCount("", "clustersecrets.core.cs.sap.com/name=my-secret-a", 1)
	env.MustError(t).AssertSecretFromFile("secret-a-1.yaml")

	env.MustFatal(t).UpdateClusterSecretFromFile("clustersecret-a-updated.yaml")
	_ = env.MustFatal(t).WaitForClusterSecretReady(clusterSecret_a)
	env.MustError(t).AssertSecretCount("", "clustersecrets.core.cs.sap.com/name=my-secret-a", 2)
	env.MustError(t).AssertSecretFromFile("secret-a-updated-1.yaml")
	env.MustError(t).AssertSecretFromFile("secret-a-updated-2.yaml")

	clusterSecret_b := env.MustFatal(t).CreateClusterSecretFromFile("clustersecret-b.yaml")
	_ = env.MustFatal(t).WaitForClusterSecretReady(clusterSecret_b)
	env.MustError(t).AssertSecretCount("", "clustersecrets.core.cs.sap.com/name", 3)
	env.MustError(t).AssertSecretFromFile("secret-a-updated-1.yaml")
	env.MustError(t).AssertSecretFromFile("secret-a-updated-2.yaml")
	env.MustError(t).AssertSecretFromFile("secret-b-2.yaml")
}

// test: delete clustersecrets
func TestController4(t *testing.T) {
	env := test.NewEnvironment()
	env.SetBasePath("testdata/4")

	env.AddObjectsFromFiles(
		"clustersecret-a.yaml",
		"clustersecret-b.yaml",
		"namespace-1.yaml",
		"namespace-2.yaml",
		"secret-a-1.yaml",
		"secret-b-2.yaml",
	)

	clusterSecret_a := env.LoadClusterSecretFromFile("clustersecret-a.yaml")
	clusterSecret_b := env.LoadClusterSecretFromFile("clustersecret-b.yaml")

	secret_a_1 := env.MustFatal(t).GetSecret("my-namespace-1", "my-secret-a")
	secret_b_2 := env.MustFatal(t).GetSecret("my-namespace-2", "my-secret-b")

	ctx, cancel := context.WithCancel(context.Background())
	c := controller.NewController(ctx, env.KubernetesClient(), env.CoreClient(), nil)
	c.Start()
	defer c.Wait()
	defer cancel()

	clusterSecret_a = env.MustFatal(t).WaitForClusterSecretReady(clusterSecret_a)
	_ = env.MustFatal(t).WaitForClusterSecretReady(clusterSecret_b)
	env.MustError(t).AssertSecretCount("", "clustersecrets.core.cs.sap.com/name", 2)
	env.MustError(t).AssertSecret(secret_a_1)
	env.MustError(t).AssertSecret(secret_b_2)

	env.MustFatal(t).DeleteClusterSecret("my-secret-a")
	env.MustFatal(t).WaitForClusterSecretDeleted(clusterSecret_a)
	env.MustError(t).AssertSecretCount("", "clustersecrets.core.cs.sap.com/name", 1)
	env.MustError(t).AssertSecret(secret_b_2)
}
