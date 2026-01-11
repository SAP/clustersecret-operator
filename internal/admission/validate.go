/*
SPDX-FileCopyrightText: 2026 SAP SE or an SAP affiliate company and clustersecret-operator contributors
SPDX-License-Identifier: Apache-2.0
*/

package admission

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/go-multierror"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation"

	corev1alpha1 "github.com/sap/clustersecret-operator/pkg/apis/core.cs.sap.com/v1alpha1"
)

func (h *Handler) validate(request *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse {
	// check that we are called with the right resources only
	if request.Resource != metav1.GroupVersionResource(corev1alpha1.ClusterSecretGroupVersionResource) {
		return admissionError(http.StatusBadRequest, fmt.Errorf("admission error: this webhook must not be called for resources of type '%s'", &request.Resource))
	}

	// do nothing except for creations or updates
	if request.Operation != admissionv1.Create && request.Operation != admissionv1.Update {
		return &admissionv1.AdmissionResponse{Allowed: true}
	}

	// deserialize clustersecret
	deserializer := codecs.UniversalDeserializer()
	var clusterSecret corev1alpha1.ClusterSecret
	if _, _, err := deserializer.Decode(request.Object.Raw, nil, &clusterSecret); err != nil {
		return admissionError(http.StatusInternalServerError, fmt.Errorf("admission error: %s", err))
	}

	// perform validations ...

	// ... check that stringData is empty
	if clusterSecret.Spec.Template.StringData != nil {
		return admissionError(http.StatusBadRequest, fmt.Errorf("admission error: unexpected field stringData"))
	}

	// ... check namespace selector
	if clusterSecret.Spec.NamespaceSelector != nil {
		if err := validateLabelSelector(clusterSecret.Spec.NamespaceSelector); err != nil {
			return admissionError(http.StatusBadRequest, fmt.Errorf("admission error: %s", err))
		}
	}

	// ... check data keys
	for key := range clusterSecret.Spec.Template.Data {
		if err := validateSecretKey(key); err != nil {
			return admissionError(http.StatusBadRequest, fmt.Errorf("admission error: %s", err))
		}
	}

	// assemble response and return
	response := admissionv1.AdmissionResponse{Allowed: true}
	return &response
}

func validateLabelSelector(selector *metav1.LabelSelector) error {
	for key, value := range selector.MatchLabels {
		if err := validateLabelKey(key); err != nil {
			return fmt.Errorf("invalid label key: %s (%s)", key, err)
		}
		if err := validateLabelValue(value); err != nil {
			return fmt.Errorf("invalid label value: %s: %s (%s)", key, value, err)
		}
	}
	for _, expr := range selector.MatchExpressions {
		if err := validateLabelKey(expr.Key); err != nil {
			return fmt.Errorf("invalid label key: %s (%s)", expr.Key, err)
		}
		switch expr.Operator {
		case metav1.LabelSelectorOpIn, metav1.LabelSelectorOpNotIn:
			if len(expr.Values) == 0 {
				return fmt.Errorf("invalid label expression value set (must not be empty): %s", expr.Key)
			}
		case metav1.LabelSelectorOpExists, metav1.LabelSelectorOpDoesNotExist:
			if len(expr.Values) > 0 {
				return fmt.Errorf("invalid label expression value set (must be empty): %s", expr.Key)
			}
		default:
			return fmt.Errorf("invalid label expression operator: %s %s", expr.Key, expr.Operator)
		}
		for _, value := range expr.Values {
			if err := validateLabelValue(value); err != nil {
				return fmt.Errorf("invalid label value: %s: %s (%s)", expr.Key, value, err)
			}
		}
	}
	return nil
}

func validateLabelKey(key string) error {
	var merr *multierror.Error
	for _, msg := range validation.IsQualifiedName(key) {
		merr = multierror.Append(merr, errors.New(msg))
	}
	return merr.ErrorOrNil()
}

func validateLabelValue(value string) error {
	var merr *multierror.Error
	for _, msg := range validation.IsValidLabelValue(value) {
		merr = multierror.Append(merr, errors.New(msg))
	}
	return merr.ErrorOrNil()
}

func validateSecretKey(key string) error {
	if !regexp.MustCompile(`^[A-Za-z0-9_\-.]*$`).MatchString(key) {
		return fmt.Errorf("invalid secret key: %s", key)
	}
	return nil
}
