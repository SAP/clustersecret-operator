/*
SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and clustersecret-operator contributors
SPDX-License-Identifier: Apache-2.0
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
