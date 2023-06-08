/*
SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and clustersecret-operator contributors
SPDX-License-Identifier: Apache-2.0
*/

package framework

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"k8s.io/apimachinery/pkg/runtime"
)

func contentHash(obj interface{}) string {
	unstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		panic("this cannot happen")
	}
	unstructured = runtime.DeepCopyJSON(unstructured)
	delete(unstructured, "status")
	delete(unstructured, "metadata")
	raw, err := json.Marshal(unstructured)
	if err != nil {
		panic("this cannot happen")
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}
