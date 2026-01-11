/*
SPDX-FileCopyrightText: 2026 SAP SE or an SAP affiliate company and clustersecret-operator contributors
SPDX-License-Identifier: Apache-2.0
*/

package admission

import (
	"fmt"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/sap/clustersecret-operator/internal/controller"
	encodingutils "github.com/sap/clustersecret-operator/internal/utils/encoding"

	corev1alpha1 "github.com/sap/clustersecret-operator/pkg/apis/core.cs.sap.com/v1alpha1"
)

func (h *Handler) mutate(request *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse {
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

	// perform mutations ...
	var patches []map[string]interface{}

	// ... inject finalizer if missing (on creation)
	if request.Operation == admissionv1.Create {
		exists := false
		for _, finalizer := range clusterSecret.Finalizers {
			if finalizer == controller.ControllerName {
				exists = true
				break
			}
		}
		if !exists {
			clusterSecret.Finalizers = append(clusterSecret.Finalizers, controller.ControllerName)
			patches = append(patches, map[string]interface{}{"op": "add", "path": "/metadata/finalizers", "value": clusterSecret.Finalizers})
		}
	}

	// ... rewrite stringData to data
	if request.Operation == admissionv1.Create || request.Operation == admissionv1.Update {
		if len(clusterSecret.Spec.Template.StringData) > 0 {
			if clusterSecret.Spec.Template.Data == nil {
				clusterSecret.Spec.Template.Data = make(map[string][]byte)
			}
			for key, value := range clusterSecret.Spec.Template.StringData {
				clusterSecret.Spec.Template.Data[key] = []byte(value)
			}
			patches = append(patches, map[string]interface{}{"op": "add", "path": "/spec/template/data", "value": clusterSecret.Spec.Template.Data})
		}
		if clusterSecret.Spec.Template.StringData != nil {
			patches = append(patches, map[string]interface{}{"op": "remove", "path": "/spec/template/stringData"})
		}
	}

	// assemble response and return
	response := admissionv1.AdmissionResponse{Allowed: true}
	if len(patches) > 0 {
		response.PatchType = &[]admissionv1.PatchType{admissionv1.PatchTypeJSONPatch}[0]
		response.Patch = encodingutils.ToJson(patches)
	}
	return &response
}
