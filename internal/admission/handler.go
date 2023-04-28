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

package admission

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/klog/v2"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Validate(w http.ResponseWriter, r *http.Request) {
	handle(w, r, h.validate)
}

func (h *Handler) Mutate(w http.ResponseWriter, r *http.Request) {
	handle(w, r, h.mutate)
}

type admitFunc func(*admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse

func handle(w http.ResponseWriter, r *http.Request, admit admitFunc) {
	switch r.Method {
	case http.MethodGet, http.MethodHead:
		// GET, HEAD should be supported by all web servers, so we return 400 here instead of 405
		httpError(w, http.StatusBadRequest, fmt.Errorf("admission error: bad method, expect POST"))
		return
	case http.MethodPost:
		// ok
	default:
		// other methods are rejected with 405
		httpError(w, http.StatusMethodNotAllowed, fmt.Errorf("admission error: bad method, expect POST"))
		return
	}

	if r.Body == nil {
		httpError(w, http.StatusBadRequest, fmt.Errorf("admission error: empty reqeuest"))
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpError(w, http.StatusInternalServerError, fmt.Errorf("admission error: %s", err))
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		httpError(w, http.StatusUnsupportedMediaType, fmt.Errorf("admission error: got content-type '%s', expect 'application/json'", contentType))
		return
	}

	requestAdmissionReview := admissionv1.AdmissionReview{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(reqBody, nil, &requestAdmissionReview); err != nil {
		httpError(w, http.StatusBadRequest, fmt.Errorf("admission error: %s", err))
		return
	}
	if requestAdmissionReview.APIVersion != admissionv1.GroupName+"/v1" || requestAdmissionReview.Kind != "AdmissionReview" {
		httpError(w, http.StatusBadRequest, fmt.Errorf("admission error: got '%s' '%s', expect '%s' '%s'", requestAdmissionReview.APIVersion, requestAdmissionReview.Kind, admissionv1.GroupName+"/v1", "AdmissionReview"))
		return
	}
	if requestAdmissionReview.Request == nil || requestAdmissionReview.Request.UID == "" {
		httpError(w, http.StatusBadRequest, fmt.Errorf("admission error: empty or incomplete review request"))
		return
	}

	responseAdmissionReview := admissionv1.AdmissionReview{}
	responseAdmissionReview.Response = admit(requestAdmissionReview.Request)
	responseAdmissionReview.Kind = requestAdmissionReview.Kind
	responseAdmissionReview.APIVersion = requestAdmissionReview.APIVersion
	responseAdmissionReview.Response.UID = requestAdmissionReview.Request.UID

	respBody, err := json.Marshal(responseAdmissionReview)
	if err != nil {
		httpError(w, http.StatusInternalServerError, fmt.Errorf("admission error: %s", err))
		return
	}
	if _, err := w.Write(respBody); err != nil {
		panic(err)
	}
}

func httpError(w http.ResponseWriter, code int, err error) {
	klog.Error(err)
	http.Error(w, err.Error(), code)
}
