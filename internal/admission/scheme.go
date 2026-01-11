/*
SPDX-FileCopyrightText: 2026 SAP SE or an SAP affiliate company and clustersecret-operator contributors
SPDX-License-Identifier: Apache-2.0
*/

package admission

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var scheme = runtime.NewScheme()
var codecs = serializer.NewCodecFactory(scheme)
