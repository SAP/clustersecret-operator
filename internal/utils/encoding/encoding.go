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

package encoding

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// decode JSON byte string into specified variable; panic in case of any error
func FromJson(raw []byte, any interface{}) interface{} {
	err := json.Unmarshal(raw, any)
	if err != nil {
		panic(err)
	}
	return any
}

// decode JSON string into specified variable; panic in case of any error
func FromJsonString(str string, any interface{}) interface{} {
	return FromJson([]byte(str), any)
}

// encode specified variable into JSON byte string; panic in case of any error
func ToJson(any interface{}) []byte {
	raw, err := json.Marshal(any)
	if err != nil {
		panic(err)
	}
	return raw
}

// encode specified variable into JSON string; panic in case of any error
func ToJsonString(any interface{}) string {
	return string(ToJson(any))
}

// encode specified variable into prettified JSON byte string; panic in case of any error
func ToPrettyJson(any interface{}) []byte {
	raw, err := json.MarshalIndent(any, "", "    ")
	if err != nil {
		panic(err)
	}
	return raw
}

// encode specified variable into prettified JSON string; panic in case of any error
func ToPrettyJsonString(any interface{}) string {
	return string(ToPrettyJson(any))
}

// decode YAML byte string into specified variable; panic in case of any error
func FromYaml(raw []byte, any interface{}) interface{} {
	err := yaml.Unmarshal(raw, any)
	if err != nil {
		panic(err)
	}
	return any
}

// decode YAML string into specified variable; panic in case of any error
func FromYamlString(str string, any interface{}) interface{} {
	return FromYaml([]byte(str), any)
}

// encode specified variable into YAML byte string; panic in case of any error
func ToYaml(any interface{}) []byte {
	raw, err := yaml.Marshal(any)
	if err != nil {
		panic(err)
	}
	return raw
}

// encode specified variable into YAML string; panic in case of any error
func ToYamlString(any interface{}) string {
	return string(ToYaml(any))
}
