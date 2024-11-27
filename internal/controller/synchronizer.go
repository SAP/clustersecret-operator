/*
SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and clustersecret-operator contributors
SPDX-License-Identifier: Apache-2.0
*/

package controller

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"
)

type Synchronizer interface {
	Init(map[schema.GroupVersionKind]cache.SharedIndexInformer)
	WaitUntilSynced()
}

type Recorder interface {
	RecordCreation(runtime.Object)
	RecordUpdate(runtime.Object, runtime.Object)
	RecordDeletion(runtime.Object)
}
