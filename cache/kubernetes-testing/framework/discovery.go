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
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (env *environmentImpl) groupVersionKind(gvr schema.GroupVersionResource) schema.GroupVersionKind {
	if gvk, ok := env.groupVersionKinds[gvr]; ok {
		return gvk
	}
	panic(fmt.Sprintf("unable to determine group version kind for %s", gvr))
}

func (env *environmentImpl) groupVersionResource(gvk schema.GroupVersionKind) schema.GroupVersionResource {
	if gvr, ok := env.groupVersionResources[gvk]; ok {
		return gvr
	}
	panic(fmt.Sprintf("unable to determine group version resource for %s", gvk))
}
