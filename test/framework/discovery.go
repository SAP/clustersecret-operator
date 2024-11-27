/*
SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and clustersecret-operator contributors
SPDX-License-Identifier: Apache-2.0
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
