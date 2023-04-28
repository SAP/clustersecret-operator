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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	Group                 = "core.cs.sap.com"
	Version               = "v1alpha1"
	ClusterSecretKind     = "ClusterSecret"
	ClusterSecretResource = "clustersecrets"
)

var (
	GroupVersion = schema.GroupVersion{
		Group:   Group,
		Version: Version,
	}
	ClusterSecretGroupKind = schema.GroupKind{
		Group: Group,
		Kind:  ClusterSecretKind,
	}
	ClusterSecretGroupVersionKind = schema.GroupVersionKind{
		Group:   Group,
		Version: Version,
		Kind:    ClusterSecretKind,
	}
	ClusterSecretGroupResource = schema.GroupResource{
		Group:    Group,
		Resource: ClusterSecretResource,
	}
	ClusterSecretGroupVersionResource = schema.GroupVersionResource{
		Group:    Group,
		Version:  Version,
		Resource: ClusterSecretResource,
	}
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterSecret is the Schema for the clustersecrets API
type ClusterSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// ClusterSecret spec
	Spec ClusterSecretSpec `json:"spec"`
	// ClusterSecret status
	Status ClusterSecretStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterSecretList contains a list of ClusterSecret
type ClusterSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ClusterSecret `json:"items"`
}

// ClusterSecretSpec defines the desired state of ClusterSecret
type ClusterSecretSpec struct {
	// Namespace selector; defines to which namespaces the secrets will be distributed
	NamespaceSelector *metav1.LabelSelector `json:"namespaceSelector,omitempty"`
	// Secret template; defines how the distributed secrets shall look like
	Template SecretTemplateSpec `json:"template"`
}

// ClusterSecretStatus reflects the actual state of ClusterSecret
type ClusterSecretStatus struct {
	// Observed generation
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// State in a short human readable form
	State string `json:"state,omitempty"`
	// State expressed as conditions (for usage with kubectl wait et al.)
	Conditions []ClusterSecretCondition `json:"conditions,omitempty"`
}

// SecretTemplateSpec defines how the managed secrets should look like
type SecretTemplateSpec struct {
	// Secret type
	Type corev1.SecretType `json:"type"`
	// Secret data as base64 encoded raw data
	Data map[string][]byte `json:"data,omitempty"`
	// Secret data as string
	StringData map[string]string `json:"stringData,omitempty"`
}

const (
	StateProcessing = "Processing"
	StateDeleting   = "Deleting"
	StateError      = "Error"
	StateReady      = "Ready"
)

// Type of a condition
type ClusterSecretConditionType string

const (
	ClusterSecretConditionTypeReady = "Ready"
)

// Condition represents a certain aspect of the overall state of a ClusterSecret object
type ClusterSecretCondition struct {
	// Type of the condition, known values are ('Ready').
	Type ClusterSecretConditionType `json:"type"`
	// Status of the condition, one of ('True', 'False', 'Unknown').
	Status corev1.ConditionStatus `json:"status"`
	// LastUpdateTime is the timestamp corresponding to the last status
	// update to this condition.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// LastTransitionTime is the timestamp corresponding to the last status
	// change of this condition.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// Reason is a brief machine readable explanation for the condition's last
	// transition.
	Reason string `json:"reason,omitempty"`
	// Message is a human readable description of the details of the last
	// transition, complementing reason.
	Message string `json:"message,omitempty"`
}
