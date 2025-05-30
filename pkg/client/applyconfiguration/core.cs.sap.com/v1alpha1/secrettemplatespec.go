/*
SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and clustersecret-operator contributors
SPDX-License-Identifier: Apache-2.0
*/

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
)

// SecretTemplateSpecApplyConfiguration represents a declarative configuration of the SecretTemplateSpec type for use
// with apply.
type SecretTemplateSpecApplyConfiguration struct {
	Type       *v1.SecretType    `json:"type,omitempty"`
	Data       map[string][]byte `json:"data,omitempty"`
	StringData map[string]string `json:"stringData,omitempty"`
}

// SecretTemplateSpecApplyConfiguration constructs a declarative configuration of the SecretTemplateSpec type for use with
// apply.
func SecretTemplateSpec() *SecretTemplateSpecApplyConfiguration {
	return &SecretTemplateSpecApplyConfiguration{}
}

// WithType sets the Type field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Type field is set to the value of the last call.
func (b *SecretTemplateSpecApplyConfiguration) WithType(value v1.SecretType) *SecretTemplateSpecApplyConfiguration {
	b.Type = &value
	return b
}

// WithData puts the entries into the Data field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the Data field,
// overwriting an existing map entries in Data field with the same key.
func (b *SecretTemplateSpecApplyConfiguration) WithData(entries map[string][]byte) *SecretTemplateSpecApplyConfiguration {
	if b.Data == nil && len(entries) > 0 {
		b.Data = make(map[string][]byte, len(entries))
	}
	for k, v := range entries {
		b.Data[k] = v
	}
	return b
}

// WithStringData puts the entries into the StringData field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the StringData field,
// overwriting an existing map entries in StringData field with the same key.
func (b *SecretTemplateSpecApplyConfiguration) WithStringData(entries map[string]string) *SecretTemplateSpecApplyConfiguration {
	if b.StringData == nil && len(entries) > 0 {
		b.StringData = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.StringData[k] = v
	}
	return b
}
