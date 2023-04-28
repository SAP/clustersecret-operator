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
