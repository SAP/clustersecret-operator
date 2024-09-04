/*
SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and clustersecret-operator contributors
SPDX-License-Identifier: Apache-2.0
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/sap/clustersecret-operator/pkg/apis/core.cs.sap.com/v1alpha1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/listers"
	"k8s.io/client-go/tools/cache"
)

// ClusterSecretLister helps list ClusterSecrets.
// All objects returned here must be treated as read-only.
type ClusterSecretLister interface {
	// List lists all ClusterSecrets in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ClusterSecret, err error)
	// Get retrieves the ClusterSecret from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.ClusterSecret, error)
	ClusterSecretListerExpansion
}

// clusterSecretLister implements the ClusterSecretLister interface.
type clusterSecretLister struct {
	listers.ResourceIndexer[*v1alpha1.ClusterSecret]
}

// NewClusterSecretLister returns a new ClusterSecretLister.
func NewClusterSecretLister(indexer cache.Indexer) ClusterSecretLister {
	return &clusterSecretLister{listers.New[*v1alpha1.ClusterSecret](indexer, v1alpha1.Resource("clustersecret"))}
}
