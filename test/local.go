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

package test

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/watch"

	corev1alpha1 "github.com/sap/clustersecret-operator/pkg/apis/core.cs.sap.com/v1alpha1"
)

func (env *Environment) WaitForClusterSecretReady(obj *corev1alpha1.ClusterSecret) (*corev1alpha1.ClusterSecret, error) {
	// todo: review
	isReady := func(event watch.Event) (bool, error) {
		return (event.Type == watch.Added || event.Type == watch.Modified) && event.Object.(*corev1alpha1.ClusterSecret).Status.State == corev1alpha1.StateReady, nil
	}
	return env.WaitForClusterSecret(obj, isReady)
}

func (env *Environment) WaitForClusterSecretDeleted(obj *corev1alpha1.ClusterSecret) error {
	// todo: review
	isDeleted := func(event watch.Event) (bool, error) {
		return event.Type == watch.Deleted, nil
	}
	_, err := env.WaitForClusterSecret(obj, isDeleted)
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}

func (must *Must) WaitForClusterSecretReady(obj *corev1alpha1.ClusterSecret) *corev1alpha1.ClusterSecret {
	retobj, err := must.env.WaitForClusterSecretReady(obj)
	must.handleError(err)
	return retobj
}

func (must *Must) WaitForClusterSecretDeleted(obj *corev1alpha1.ClusterSecret) {
	err := must.env.WaitForClusterSecretDeleted(obj)
	must.handleError(err)
}
