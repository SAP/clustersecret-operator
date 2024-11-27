/*
SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and clustersecret-operator contributors
SPDX-License-Identifier: Apache-2.0
*/

package framework

import "k8s.io/apimachinery/pkg/runtime"

type CreateCallbackFunc func(runtime.Object)
type UpdateCallbackFunc func(runtime.Object, runtime.Object)
type DeleteCallbackFunc func(runtime.Object)
