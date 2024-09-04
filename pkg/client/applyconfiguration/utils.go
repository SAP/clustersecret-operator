/*
SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and clustersecret-operator contributors
SPDX-License-Identifier: Apache-2.0
*/

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package applyconfiguration

import (
	v1alpha1 "github.com/sap/clustersecret-operator/pkg/apis/core.cs.sap.com/v1alpha1"
	corecssapcomv1alpha1 "github.com/sap/clustersecret-operator/pkg/client/applyconfiguration/core.cs.sap.com/v1alpha1"
	internal "github.com/sap/clustersecret-operator/pkg/client/applyconfiguration/internal"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	testing "k8s.io/client-go/testing"
)

// ForKind returns an apply configuration type for the given GroupVersionKind, or nil if no
// apply configuration type exists for the given GroupVersionKind.
func ForKind(kind schema.GroupVersionKind) interface{} {
	switch kind {
	// Group=core.cs.sap.com, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithKind("ClusterSecret"):
		return &corecssapcomv1alpha1.ClusterSecretApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ClusterSecretCondition"):
		return &corecssapcomv1alpha1.ClusterSecretConditionApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ClusterSecretSpec"):
		return &corecssapcomv1alpha1.ClusterSecretSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ClusterSecretStatus"):
		return &corecssapcomv1alpha1.ClusterSecretStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("SecretTemplateSpec"):
		return &corecssapcomv1alpha1.SecretTemplateSpecApplyConfiguration{}

	}
	return nil
}

func NewTypeConverter(scheme *runtime.Scheme) *testing.TypeConverter {
	return &testing.TypeConverter{Scheme: scheme, TypeResolver: internal.Parser()}
}
