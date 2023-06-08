/*
SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and clustersecret-operator contributors
SPDX-License-Identifier: Apache-2.0
*/

// Code generated. DO NOT EDIT.

package test

import (
	"testing"

	_corecssapcomv1alpha1 "github.com/sap/clustersecret-operator/pkg/apis/core.cs.sap.com/v1alpha1"
	_core "github.com/sap/clustersecret-operator/pkg/client/clientset/versioned/fake"
	"github.com/sap/clustersecret-operator/test/framework"
	_corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	_kubernetes "k8s.io/client-go/kubernetes/fake"
	watchtools "k8s.io/client-go/tools/watch"
)

type Environment struct {
	framework.Environment
	kubernetesClient *_kubernetes.Clientset
	coreClient       *_core.Clientset
}

func NewEnvironment() *Environment {
	kubernetesClient := _kubernetes.NewSimpleClientset()
	coreClient := _core.NewSimpleClientset()
	groupVersions := []*framework.GroupVersion{
		framework.NewGroupVersion(_kubernetes.AddToScheme, kubernetesClient),
		framework.NewGroupVersion(_core.AddToScheme, coreClient),
	}
	return &Environment{
		Environment:      framework.NewEnvironment(groupVersions),
		kubernetesClient: kubernetesClient,
		coreClient:       coreClient,
	}
}

type Must struct {
	env          *Environment
	errorHandler func(error)
}

func (must *Must) handleError(err error) {
	if err != nil {
		must.errorHandler(err)
	}
}

func (env *Environment) Must(errorHandler func(error)) *Must {
	return &Must{env: env, errorHandler: errorHandler}
}

func (env *Environment) MustError(t *testing.T) *Must {
	return env.Must(func(err error) { t.Error(err) })
}

func (env *Environment) MustFatal(t *testing.T) *Must {
	return env.Must(func(err error) { t.Fatal(err) })
}

// Client accessors

func (env *Environment) KubernetesClient() *_kubernetes.Clientset {
	return env.kubernetesClient
}

func (env *Environment) CoreClient() *_core.Clientset {
	return env.coreClient
}

// Typed methods for core/v1 Secret

func (env *Environment) LoadSecretFromFile(path string) *_corev1.Secret {
	return env.LoadObjectFromFile(path).(*_corev1.Secret)
}

func (env *Environment) AddSecret(obj *_corev1.Secret) {
	env.AddObject(obj)
}

func (env *Environment) AddSecretFromFile(path string) {
	env.AddObjectFromFile(path)
}

func (env *Environment) AddSecretsFromFiles(paths ...string) {
	env.AddObjectsFromFiles(paths...)
}

func (env *Environment) WithSecret(obj *_corev1.Secret) *Environment {
	return env.WithObject(obj).(*Environment)
}

func (env *Environment) WithSecretFromFile(path string) *Environment {
	return env.WithObjectFromFile(path).(*Environment)
}

func (env *Environment) WithSecretsFromFiles(paths ...string) *Environment {
	return env.WithObjectsFromFiles(paths...).(*Environment)
}

func (env *Environment) AssertSecret(obj *_corev1.Secret) error {
	return env.AssertObject(obj)
}

func (env *Environment) AssertSecretFromFile(path string) error {
	return env.AssertObjectFromFile(path)
}

func (env *Environment) AssertSecretCount(namespace string, labelSelector string, count int) error {
	return env.AssertObjectCount(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Secret"}, namespace, labelSelector, count)
}

func (env *Environment) GetSecret(namespace string, name string) (*_corev1.Secret, error) {
	retobj, err := env.GetObject(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Secret"}, namespace, name)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corev1.Secret), nil
}

func (env *Environment) ListSecrets(namespace string, labelSelector string) ([]*_corev1.Secret, error) {
	retobjs, err := env.ListObjects(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Secret"}, namespace, labelSelector)
	if err != nil {
		return nil, err
	}
	typedretobjs := make([]*_corev1.Secret, len(retobjs))
	for i, retobj := range retobjs {
		typedretobjs[i] = retobj.(*_corev1.Secret)
	}
	return typedretobjs, nil
}

func (env *Environment) CreateSecret(obj *_corev1.Secret) (*_corev1.Secret, error) {
	retobj, err := env.CreateObject(obj)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corev1.Secret), nil
}

func (env *Environment) CreateSecretFromFile(path string) (*_corev1.Secret, error) {
	retobj, err := env.CreateObjectFromFile(path)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corev1.Secret), nil
}

func (env *Environment) UpdateSecret(obj *_corev1.Secret) (*_corev1.Secret, error) {
	retobj, err := env.UpdateObject(obj)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corev1.Secret), nil
}

func (env *Environment) UpdateSecretFromFile(path string) (*_corev1.Secret, error) {
	retobj, err := env.UpdateObjectFromFile(path)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corev1.Secret), nil
}

func (env *Environment) PatchSecret(namespace string, name string, patchType types.PatchType, patch []byte) (*_corev1.Secret, error) {
	retobj, err := env.PatchObject(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Secret"}, namespace, name, patchType, patch)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corev1.Secret), nil
}

func (env *Environment) LabelSecret(namespace string, name string, key string, value string) (*_corev1.Secret, error) {
	retobj, err := env.LabelObject(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Secret"}, namespace, name, key, value)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corev1.Secret), nil
}

func (env *Environment) UnlabelSecret(namespace string, name string, key string) (*_corev1.Secret, error) {
	retobj, err := env.UnlabelObject(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Secret"}, namespace, name, key)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corev1.Secret), nil
}

func (env *Environment) DeleteSecret(namespace string, name string) error {
	return env.DeleteObject(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Secret"}, namespace, name)
}

func (env *Environment) WaitForSecret(obj *_corev1.Secret, conditions ...watchtools.ConditionFunc) (*_corev1.Secret, error) {
	retobj, err := env.WaitForObject(obj, conditions...)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corev1.Secret), nil
}

func (env *Environment) WaitForSecretFromFile(path string, conditions ...watchtools.ConditionFunc) (*_corev1.Secret, error) {
	retobj, err := env.WaitForObjectFromFile(path, conditions...)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corev1.Secret), nil
}

func (must *Must) AssertSecret(obj *_corev1.Secret) {
	err := must.env.AssertSecret(obj)
	must.handleError(err)
}

func (must *Must) AssertSecretFromFile(path string) {
	err := must.env.AssertSecretFromFile(path)
	must.handleError(err)
}

func (must *Must) AssertSecretCount(namespace string, labelSelector string, count int) {
	err := must.env.AssertSecretCount(namespace, labelSelector, count)
	must.handleError(err)
}

func (must *Must) GetSecret(namespace string, name string) *_corev1.Secret {
	retobj, err := must.env.GetSecret(namespace, name)
	must.handleError(err)
	return retobj
}

func (must *Must) ListSecrets(namespace string, labelSelector string) []*_corev1.Secret {
	retobjs, err := must.env.ListSecrets(namespace, labelSelector)
	must.handleError(err)
	return retobjs
}

func (must *Must) CreateSecret(obj *_corev1.Secret) *_corev1.Secret {
	retobj, err := must.env.CreateSecret(obj)
	must.handleError(err)
	return retobj
}

func (must *Must) CreateSecretFromFile(path string) *_corev1.Secret {
	retobj, err := must.env.CreateSecretFromFile(path)
	must.handleError(err)
	return retobj
}

func (must *Must) UpdateSecret(obj *_corev1.Secret) *_corev1.Secret {
	retobj, err := must.env.UpdateSecret(obj)
	must.handleError(err)
	return retobj
}

func (must *Must) UpdateSecretFromFile(path string) *_corev1.Secret {
	retobj, err := must.env.UpdateSecretFromFile(path)
	must.handleError(err)
	return retobj
}

func (must *Must) PatchSecret(namespace string, name string, patchType types.PatchType, patch []byte) *_corev1.Secret {
	retobj, err := must.env.PatchSecret(namespace, name, patchType, patch)
	must.handleError(err)
	return retobj
}

func (must *Must) LabelSecret(namespace string, name string, key string, value string) *_corev1.Secret {
	retobj, err := must.env.LabelSecret(namespace, name, key, value)
	must.handleError(err)
	return retobj
}

func (must *Must) UnlabelSecret(namespace string, name string, key string) *_corev1.Secret {
	retobj, err := must.env.UnlabelSecret(namespace, name, key)
	must.handleError(err)
	return retobj
}

func (must *Must) DeleteSecret(namespace string, name string) {
	err := must.env.DeleteSecret(namespace, name)
	must.handleError(err)
}

func (must *Must) WaitForSecret(obj *_corev1.Secret, conditions ...watchtools.ConditionFunc) *_corev1.Secret {
	retobj, err := must.env.WaitForSecret(obj, conditions...)
	must.handleError(err)
	return retobj
}

func (must *Must) WaitForSecretFromFile(path string, conditions ...watchtools.ConditionFunc) *_corev1.Secret {
	retobj, err := must.env.WaitForSecretFromFile(path, conditions...)
	must.handleError(err)
	return retobj
}

// Typed methods for core/v1 Namespace

func (env *Environment) LoadNamespaceFromFile(path string) *_corev1.Namespace {
	return env.LoadObjectFromFile(path).(*_corev1.Namespace)
}

func (env *Environment) AddNamespace(obj *_corev1.Namespace) {
	env.AddObject(obj)
}

func (env *Environment) AddNamespaceFromFile(path string) {
	env.AddObjectFromFile(path)
}

func (env *Environment) AddNamespacesFromFiles(paths ...string) {
	env.AddObjectsFromFiles(paths...)
}

func (env *Environment) WithNamespace(obj *_corev1.Namespace) *Environment {
	return env.WithObject(obj).(*Environment)
}

func (env *Environment) WithNamespaceFromFile(path string) *Environment {
	return env.WithObjectFromFile(path).(*Environment)
}

func (env *Environment) WithNamespacesFromFiles(paths ...string) *Environment {
	return env.WithObjectsFromFiles(paths...).(*Environment)
}

func (env *Environment) AssertNamespace(obj *_corev1.Namespace) error {
	return env.AssertObject(obj)
}

func (env *Environment) AssertNamespaceFromFile(path string) error {
	return env.AssertObjectFromFile(path)
}

func (env *Environment) AssertNamespaceCount(labelSelector string, count int) error {
	return env.AssertObjectCount(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Namespace"}, "", labelSelector, count)
}

func (env *Environment) GetNamespace(name string) (*_corev1.Namespace, error) {
	retobj, err := env.GetObject(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Namespace"}, "", name)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corev1.Namespace), nil
}

func (env *Environment) ListNamespaces(labelSelector string) ([]*_corev1.Namespace, error) {
	retobjs, err := env.ListObjects(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Namespace"}, "", labelSelector)
	if err != nil {
		return nil, err
	}
	typedretobjs := make([]*_corev1.Namespace, len(retobjs))
	for i, retobj := range retobjs {
		typedretobjs[i] = retobj.(*_corev1.Namespace)
	}
	return typedretobjs, nil
}

func (env *Environment) CreateNamespace(obj *_corev1.Namespace) (*_corev1.Namespace, error) {
	retobj, err := env.CreateObject(obj)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corev1.Namespace), nil
}

func (env *Environment) CreateNamespaceFromFile(path string) (*_corev1.Namespace, error) {
	retobj, err := env.CreateObjectFromFile(path)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corev1.Namespace), nil
}

func (env *Environment) UpdateNamespace(obj *_corev1.Namespace) (*_corev1.Namespace, error) {
	retobj, err := env.UpdateObject(obj)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corev1.Namespace), nil
}

func (env *Environment) UpdateNamespaceFromFile(path string) (*_corev1.Namespace, error) {
	retobj, err := env.UpdateObjectFromFile(path)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corev1.Namespace), nil
}

func (env *Environment) PatchNamespace(name string, patchType types.PatchType, patch []byte) (*_corev1.Namespace, error) {
	retobj, err := env.PatchObject(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Namespace"}, "", name, patchType, patch)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corev1.Namespace), nil
}

func (env *Environment) LabelNamespace(name string, key string, value string) (*_corev1.Namespace, error) {
	retobj, err := env.LabelObject(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Namespace"}, "", name, key, value)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corev1.Namespace), nil
}

func (env *Environment) UnlabelNamespace(name string, key string) (*_corev1.Namespace, error) {
	retobj, err := env.UnlabelObject(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Namespace"}, "", name, key)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corev1.Namespace), nil
}

func (env *Environment) DeleteNamespace(name string) error {
	return env.DeleteObject(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Namespace"}, "", name)
}

func (env *Environment) WaitForNamespace(obj *_corev1.Namespace, conditions ...watchtools.ConditionFunc) (*_corev1.Namespace, error) {
	retobj, err := env.WaitForObject(obj, conditions...)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corev1.Namespace), nil
}

func (env *Environment) WaitForNamespaceFromFile(path string, conditions ...watchtools.ConditionFunc) (*_corev1.Namespace, error) {
	retobj, err := env.WaitForObjectFromFile(path, conditions...)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corev1.Namespace), nil
}

func (must *Must) AssertNamespace(obj *_corev1.Namespace) {
	err := must.env.AssertNamespace(obj)
	must.handleError(err)
}

func (must *Must) AssertNamespaceFromFile(path string) {
	err := must.env.AssertNamespaceFromFile(path)
	must.handleError(err)
}

func (must *Must) AssertNamespaceCount(labelSelector string, count int) {
	err := must.env.AssertNamespaceCount(labelSelector, count)
	must.handleError(err)
}

func (must *Must) GetNamespace(name string) *_corev1.Namespace {
	retobj, err := must.env.GetNamespace(name)
	must.handleError(err)
	return retobj
}

func (must *Must) ListNamespaces(labelSelector string) []*_corev1.Namespace {
	retobjs, err := must.env.ListNamespaces(labelSelector)
	must.handleError(err)
	return retobjs
}

func (must *Must) CreateNamespace(obj *_corev1.Namespace) *_corev1.Namespace {
	retobj, err := must.env.CreateNamespace(obj)
	must.handleError(err)
	return retobj
}

func (must *Must) CreateNamespaceFromFile(path string) *_corev1.Namespace {
	retobj, err := must.env.CreateNamespaceFromFile(path)
	must.handleError(err)
	return retobj
}

func (must *Must) UpdateNamespace(obj *_corev1.Namespace) *_corev1.Namespace {
	retobj, err := must.env.UpdateNamespace(obj)
	must.handleError(err)
	return retobj
}

func (must *Must) UpdateNamespaceFromFile(path string) *_corev1.Namespace {
	retobj, err := must.env.UpdateNamespaceFromFile(path)
	must.handleError(err)
	return retobj
}

func (must *Must) PatchNamespace(name string, patchType types.PatchType, patch []byte) *_corev1.Namespace {
	retobj, err := must.env.PatchNamespace(name, patchType, patch)
	must.handleError(err)
	return retobj
}

func (must *Must) LabelNamespace(name string, key string, value string) *_corev1.Namespace {
	retobj, err := must.env.LabelNamespace(name, key, value)
	must.handleError(err)
	return retobj
}

func (must *Must) UnlabelNamespace(name string, key string) *_corev1.Namespace {
	retobj, err := must.env.UnlabelNamespace(name, key)
	must.handleError(err)
	return retobj
}

func (must *Must) DeleteNamespace(name string) {
	err := must.env.DeleteNamespace(name)
	must.handleError(err)
}

func (must *Must) WaitForNamespace(obj *_corev1.Namespace, conditions ...watchtools.ConditionFunc) *_corev1.Namespace {
	retobj, err := must.env.WaitForNamespace(obj, conditions...)
	must.handleError(err)
	return retobj
}

func (must *Must) WaitForNamespaceFromFile(path string, conditions ...watchtools.ConditionFunc) *_corev1.Namespace {
	retobj, err := must.env.WaitForNamespaceFromFile(path, conditions...)
	must.handleError(err)
	return retobj
}

// Typed methods for core.cs.sap.com/v1alpha1 ClusterSecret

func (env *Environment) LoadClusterSecretFromFile(path string) *_corecssapcomv1alpha1.ClusterSecret {
	return env.LoadObjectFromFile(path).(*_corecssapcomv1alpha1.ClusterSecret)
}

func (env *Environment) AddClusterSecret(obj *_corecssapcomv1alpha1.ClusterSecret) {
	env.AddObject(obj)
}

func (env *Environment) AddClusterSecretFromFile(path string) {
	env.AddObjectFromFile(path)
}

func (env *Environment) AddClusterSecretsFromFiles(paths ...string) {
	env.AddObjectsFromFiles(paths...)
}

func (env *Environment) WithClusterSecret(obj *_corecssapcomv1alpha1.ClusterSecret) *Environment {
	return env.WithObject(obj).(*Environment)
}

func (env *Environment) WithClusterSecretFromFile(path string) *Environment {
	return env.WithObjectFromFile(path).(*Environment)
}

func (env *Environment) WithClusterSecretsFromFiles(paths ...string) *Environment {
	return env.WithObjectsFromFiles(paths...).(*Environment)
}

func (env *Environment) AssertClusterSecret(obj *_corecssapcomv1alpha1.ClusterSecret) error {
	return env.AssertObject(obj)
}

func (env *Environment) AssertClusterSecretFromFile(path string) error {
	return env.AssertObjectFromFile(path)
}

func (env *Environment) AssertClusterSecretCount(labelSelector string, count int) error {
	return env.AssertObjectCount(schema.GroupVersionKind{Group: "core.cs.sap.com", Version: "v1alpha1", Kind: "ClusterSecret"}, "", labelSelector, count)
}

func (env *Environment) GetClusterSecret(name string) (*_corecssapcomv1alpha1.ClusterSecret, error) {
	retobj, err := env.GetObject(schema.GroupVersionKind{Group: "core.cs.sap.com", Version: "v1alpha1", Kind: "ClusterSecret"}, "", name)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corecssapcomv1alpha1.ClusterSecret), nil
}

func (env *Environment) ListClusterSecrets(labelSelector string) ([]*_corecssapcomv1alpha1.ClusterSecret, error) {
	retobjs, err := env.ListObjects(schema.GroupVersionKind{Group: "core.cs.sap.com", Version: "v1alpha1", Kind: "ClusterSecret"}, "", labelSelector)
	if err != nil {
		return nil, err
	}
	typedretobjs := make([]*_corecssapcomv1alpha1.ClusterSecret, len(retobjs))
	for i, retobj := range retobjs {
		typedretobjs[i] = retobj.(*_corecssapcomv1alpha1.ClusterSecret)
	}
	return typedretobjs, nil
}

func (env *Environment) CreateClusterSecret(obj *_corecssapcomv1alpha1.ClusterSecret) (*_corecssapcomv1alpha1.ClusterSecret, error) {
	retobj, err := env.CreateObject(obj)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corecssapcomv1alpha1.ClusterSecret), nil
}

func (env *Environment) CreateClusterSecretFromFile(path string) (*_corecssapcomv1alpha1.ClusterSecret, error) {
	retobj, err := env.CreateObjectFromFile(path)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corecssapcomv1alpha1.ClusterSecret), nil
}

func (env *Environment) UpdateClusterSecret(obj *_corecssapcomv1alpha1.ClusterSecret) (*_corecssapcomv1alpha1.ClusterSecret, error) {
	retobj, err := env.UpdateObject(obj)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corecssapcomv1alpha1.ClusterSecret), nil
}

func (env *Environment) UpdateClusterSecretFromFile(path string) (*_corecssapcomv1alpha1.ClusterSecret, error) {
	retobj, err := env.UpdateObjectFromFile(path)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corecssapcomv1alpha1.ClusterSecret), nil
}

func (env *Environment) PatchClusterSecret(name string, patchType types.PatchType, patch []byte) (*_corecssapcomv1alpha1.ClusterSecret, error) {
	retobj, err := env.PatchObject(schema.GroupVersionKind{Group: "core.cs.sap.com", Version: "v1alpha1", Kind: "ClusterSecret"}, "", name, patchType, patch)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corecssapcomv1alpha1.ClusterSecret), nil
}

func (env *Environment) LabelClusterSecret(name string, key string, value string) (*_corecssapcomv1alpha1.ClusterSecret, error) {
	retobj, err := env.LabelObject(schema.GroupVersionKind{Group: "core.cs.sap.com", Version: "v1alpha1", Kind: "ClusterSecret"}, "", name, key, value)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corecssapcomv1alpha1.ClusterSecret), nil
}

func (env *Environment) UnlabelClusterSecret(name string, key string) (*_corecssapcomv1alpha1.ClusterSecret, error) {
	retobj, err := env.UnlabelObject(schema.GroupVersionKind{Group: "core.cs.sap.com", Version: "v1alpha1", Kind: "ClusterSecret"}, "", name, key)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corecssapcomv1alpha1.ClusterSecret), nil
}

func (env *Environment) DeleteClusterSecret(name string) error {
	return env.DeleteObject(schema.GroupVersionKind{Group: "core.cs.sap.com", Version: "v1alpha1", Kind: "ClusterSecret"}, "", name)
}

func (env *Environment) WaitForClusterSecret(obj *_corecssapcomv1alpha1.ClusterSecret, conditions ...watchtools.ConditionFunc) (*_corecssapcomv1alpha1.ClusterSecret, error) {
	retobj, err := env.WaitForObject(obj, conditions...)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corecssapcomv1alpha1.ClusterSecret), nil
}

func (env *Environment) WaitForClusterSecretFromFile(path string, conditions ...watchtools.ConditionFunc) (*_corecssapcomv1alpha1.ClusterSecret, error) {
	retobj, err := env.WaitForObjectFromFile(path, conditions...)
	if err != nil {
		return nil, err
	}
	return retobj.(*_corecssapcomv1alpha1.ClusterSecret), nil
}

func (must *Must) AssertClusterSecret(obj *_corecssapcomv1alpha1.ClusterSecret) {
	err := must.env.AssertClusterSecret(obj)
	must.handleError(err)
}

func (must *Must) AssertClusterSecretFromFile(path string) {
	err := must.env.AssertClusterSecretFromFile(path)
	must.handleError(err)
}

func (must *Must) AssertClusterSecretCount(labelSelector string, count int) {
	err := must.env.AssertClusterSecretCount(labelSelector, count)
	must.handleError(err)
}

func (must *Must) GetClusterSecret(name string) *_corecssapcomv1alpha1.ClusterSecret {
	retobj, err := must.env.GetClusterSecret(name)
	must.handleError(err)
	return retobj
}

func (must *Must) ListClusterSecrets(labelSelector string) []*_corecssapcomv1alpha1.ClusterSecret {
	retobjs, err := must.env.ListClusterSecrets(labelSelector)
	must.handleError(err)
	return retobjs
}

func (must *Must) CreateClusterSecret(obj *_corecssapcomv1alpha1.ClusterSecret) *_corecssapcomv1alpha1.ClusterSecret {
	retobj, err := must.env.CreateClusterSecret(obj)
	must.handleError(err)
	return retobj
}

func (must *Must) CreateClusterSecretFromFile(path string) *_corecssapcomv1alpha1.ClusterSecret {
	retobj, err := must.env.CreateClusterSecretFromFile(path)
	must.handleError(err)
	return retobj
}

func (must *Must) UpdateClusterSecret(obj *_corecssapcomv1alpha1.ClusterSecret) *_corecssapcomv1alpha1.ClusterSecret {
	retobj, err := must.env.UpdateClusterSecret(obj)
	must.handleError(err)
	return retobj
}

func (must *Must) UpdateClusterSecretFromFile(path string) *_corecssapcomv1alpha1.ClusterSecret {
	retobj, err := must.env.UpdateClusterSecretFromFile(path)
	must.handleError(err)
	return retobj
}

func (must *Must) PatchClusterSecret(name string, patchType types.PatchType, patch []byte) *_corecssapcomv1alpha1.ClusterSecret {
	retobj, err := must.env.PatchClusterSecret(name, patchType, patch)
	must.handleError(err)
	return retobj
}

func (must *Must) LabelClusterSecret(name string, key string, value string) *_corecssapcomv1alpha1.ClusterSecret {
	retobj, err := must.env.LabelClusterSecret(name, key, value)
	must.handleError(err)
	return retobj
}

func (must *Must) UnlabelClusterSecret(name string, key string) *_corecssapcomv1alpha1.ClusterSecret {
	retobj, err := must.env.UnlabelClusterSecret(name, key)
	must.handleError(err)
	return retobj
}

func (must *Must) DeleteClusterSecret(name string) {
	err := must.env.DeleteClusterSecret(name)
	must.handleError(err)
}

func (must *Must) WaitForClusterSecret(obj *_corecssapcomv1alpha1.ClusterSecret, conditions ...watchtools.ConditionFunc) *_corecssapcomv1alpha1.ClusterSecret {
	retobj, err := must.env.WaitForClusterSecret(obj, conditions...)
	must.handleError(err)
	return retobj
}

func (must *Must) WaitForClusterSecretFromFile(path string, conditions ...watchtools.ConditionFunc) *_corecssapcomv1alpha1.ClusterSecret {
	retobj, err := must.env.WaitForClusterSecretFromFile(path, conditions...)
	must.handleError(err)
	return retobj
}
