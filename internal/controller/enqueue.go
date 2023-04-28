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

package controller

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"

	corev1alpha1 "github.com/sap/clustersecret-operator/pkg/apis/core.cs.sap.com/v1alpha1"
)

func (c *Controller) enqueueNamespace(eventType string, obj interface{}) {
	// no need to handle cache.DeletedFinalStateUnknown (i.e. recover objects from tombstone), since we are not handling delete events for namespaces
	namespace, ok := obj.(*corev1.Namespace)
	if !ok {
		panic("this cannot happen")
	}
	klog.V(2).Infof("enqueuing namespace %s (%s)", namespace.Name, eventType)
	c.workqueue.Add(workqueueItem{key: workqueueItemKeyNamespace, name: namespace.Name})
}

func (c *Controller) enqueueClusterSecret(eventType string, obj interface{}) {
	clusterSecret, ok := obj.(*corev1alpha1.ClusterSecret)
	if !ok {
		// try to recover from tombstone (can only happen in case of delete events, see https://pkg.go.dev/k8s.io/client-go/tools/cache#ResourceEventHandler)
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			// that is now really strange but we don't know if it's safe to panic here; so we just silently return, i.e. ignore the object
			return
		}
		klog.V(2).Infof("recovered deleted object %s from tombstone", tombstone.Key)
		clusterSecret, ok = tombstone.Obj.(*corev1alpha1.ClusterSecret)
		if !ok {
			panic("this cannot happen")
		}
	}
	klog.V(2).Infof("enqueuing clustersecret %s (%s)", clusterSecret.Name, eventType)
	c.workqueue.Add(workqueueItem{key: workqueueItemKeyClusterSecret, name: clusterSecret.Name})
}
