/*
SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and clustersecret-operator contributors
SPDX-License-Identifier: Apache-2.0
*/

package admission

import (
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func admissionError(code int, err error) *admissionv1.AdmissionResponse {
	klog.Error(err)
	return &admissionv1.AdmissionResponse{
		Allowed: false,
		Result: &metav1.Status{
			Status:  http.StatusText(code),
			Message: err.Error(),
			Reason:  metav1.StatusReasonForbidden,
			Code:    int32(code),
		},
		Warnings: []string{err.Error()},
	}
}
