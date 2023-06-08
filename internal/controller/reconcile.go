/*
SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and clustersecret-operator contributors
SPDX-License-Identifier: Apache-2.0
*/

package controller

import (
	"context"
	"fmt"

	multierror "github.com/hashicorp/go-multierror"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels" // could also be aliased 'kubeclients' but we keep it as 'kubernetes' since most people do
	"k8s.io/klog/v2"

	conversionutils "github.com/sap/clustersecret-operator/internal/utils/conversion"

	corev1alpha1 "github.com/sap/clustersecret-operator/pkg/apis/core.cs.sap.com/v1alpha1"
)

type secretKey struct {
	namespace string
	name      string
}

type secretOperation struct {
	old *corev1.Secret
	new *corev1.Secret
}

const (
	LabelKeyName            = "clustersecrets.core.cs.sap.com/name"
	AnnotationKeyGeneration = "clustersecrets.core.cs.sap.com/generation"
)

func (c *Controller) reconcileNamespace(namespaceName string) error {
	// note: due to the implementation details of the workqueue it is guaranteed that this function will not run concurrently for the same namespace

	klog.V(2).Infof("reconciling namespace %s", namespaceName)

	// wait for caches to be synchronized
	if c.synchronizer != nil {
		c.synchronizer.WaitUntilSynced()
	}

	// fetch namespace
	namespace, err := c.namespaceLister.Get(namespaceName)
	if err != nil {
		if errors.IsNotFound(err) {
			// that is a bit strange, since this should be only triggered after namespace create/update events
			// however it may happen if the  namespace has been deleted concurrently
			klog.Warning("namespace %s does not exist; skipping reconcile", namespaceName)
			return nil
		} else {
			return err
		}
	}

	// if namespace is being deleted (has a deletionTimestamp), no action is required
	if !namespace.DeletionTimestamp.IsZero() {
		return nil
	}

	// determine all clustersecrets that potentially need reconciliation ...
	clusterSecretNames := make(map[string]struct{})

	// ... first, find all managed secrets in specified namespace
	secretSelector, err := labels.Parse(LabelKeyName)
	if err != nil {
		panic("this cannot happen")
	}
	existingSecrets, err := c.secretLister.Secrets(namespaceName).List(secretSelector)
	if err != nil {
		c.eventRecorder.Event(namespace, corev1.EventTypeWarning, "Error", err.Error())
		return err
	}
	for _, secret := range existingSecrets {
		clusterSecretNames[secret.Name] = struct{}{}
	}

	// ... then, find all clustersecrets selecting the specified namespace
	clusterSecrets, err := c.clusterSecretLister.List(labels.Everything())
	if err != nil {
		c.eventRecorder.Event(namespace, corev1.EventTypeWarning, "Error", err.Error())
		return err
	}
	for _, clusterSecret := range clusterSecrets {
		namespaceSelector := buildNamespaceSelectorFromClusterSecret(clusterSecret)
		if namespaceSelector.Matches(labels.Set(namespace.Labels)) {
			clusterSecretNames[clusterSecret.Name] = struct{}{}
		}
	}

	// schedule a reconciliation for all these determined clustersecrets
	for clusterSecretName := range clusterSecretNames {
		c.eventRecorder.Eventf(namespace, corev1.EventTypeNormal, "TriggerClusterSecretReconcile", "Successfully triggered reconciliation of clustersecret %s", clusterSecretName)
		c.workqueue.Add(workqueueItem{key: workqueueItemKeyClusterSecret, name: clusterSecretName})
	}

	// return
	return nil
}

func (c *Controller) reconcileClusterSecret(clusterSecretName string) error {
	// note: due to the implementation details of the workqueue it is guaranteed that this function will not run concurrently for the same clustersecret

	klog.V(2).Infof("reconciling clustersecret %s", clusterSecretName)

	// wait for caches to be synchronized
	if c.synchronizer != nil {
		c.synchronizer.WaitUntilSynced()
	}

	// fetch clustersecret (if existing)
	// note: we cannot fetch it from the lister because the cached state might not yet reflect updates done by (very recent) previous invocations of this function
	clusterSecret, err := c.coreclient.CoreV1alpha1().ClusterSecrets().Get(context.TODO(), clusterSecretName, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		clusterSecret = nil
	}

	// set finalizer
	if clusterSecret != nil && clusterSecret.DeletionTimestamp.IsZero() {
		if err := c.setClusterSecretFinalizer(clusterSecret); err != nil {
			c.eventRecorder.Event(clusterSecret, corev1.EventTypeWarning, "Error", err.Error())
			return err
		}
	}

	// check that stringData is not set (should have been rewritten to data by admission webhook)
	if clusterSecret != nil && clusterSecret.Spec.Template.StringData != nil {
		err := fmt.Errorf("unexpected stringData in clustersecret %s", clusterSecret.Name)
		c.eventRecorder.Event(clusterSecret, corev1.EventTypeWarning, "Error", err.Error())
		return err
	}

	// fetch all secrets managed by this clustersecret in all namespaces
	secretSelector := labels.SelectorFromSet(map[string]string{LabelKeyName: clusterSecretName})
	existingSecrets, err := c.secretLister.List(secretSelector)
	if err != nil {
		if clusterSecret != nil {
			c.eventRecorder.Event(clusterSecret, corev1.EventTypeWarning, "Error", err.Error())
		}
		return err
	}

	// determine set of secrets to reconcile ...
	operations := make(map[secretKey]*secretOperation)
	// ... first, consider all existing managed secrets
	for _, secret := range existingSecrets {
		key := secretKey{secret.Namespace, secret.Name}
		operations[key] = &secretOperation{old: secret}
	}
	// ... then (if clustersecret is not deleted or in deletion), consider the wanted generated secret in all selected namespaces
	if clusterSecret != nil && clusterSecret.DeletionTimestamp.IsZero() {
		namespaceSelector := buildNamespaceSelectorFromClusterSecret(clusterSecret)
		matchingNamespaces, err := c.namespaceLister.List(namespaceSelector)
		if err != nil {
			c.eventRecorder.Event(clusterSecret, corev1.EventTypeWarning, "Error", err.Error())
			return err
		}
		for _, namespace := range matchingNamespaces {
			// skip if namespace is in deletion; if a secret exists in that namespace, it will be deleted through the operations entry
			// (which is not necessary, but does not harm)
			if !namespace.DeletionTimestamp.IsZero() {
				continue
			}
			key := secretKey{namespace.Name, clusterSecret.Name}
			if operation, ok := operations[key]; ok {
				operation.new = buildSecretFromClusterSecret(namespace.Name, clusterSecret)
			} else {
				operations[key] = &secretOperation{new: buildSecretFromClusterSecret(namespace.Name, clusterSecret)}
			}
		}
		for key, operation := range operations {
			// if secret is going to be udpated, set the resourceVersion to enable/allow optimistic locking on update
			if operation.old != nil && operation.new != nil {
				operation.new.ResourceVersion = operation.old.ResourceVersion
				// skip/remove all secrets which are already up-to-date
				if operation.old != nil && conversionutils.Atoi(operation.old.Annotations[AnnotationKeyGeneration]) >= clusterSecret.Generation {
					delete(operations, key)
				}
			}
		}
	}

	// update status (if applicable); set to Processing or Deleting respectively (unless it's already in Error state; in that case it stays Error)
	if clusterSecret != nil && clusterSecret.Status.State != corev1alpha1.StateError {
		if clusterSecret.DeletionTimestamp.IsZero() {
			if clusterSecret.Generation > clusterSecret.Status.ObservedGeneration || len(operations) > 0 {
				if err := c.updateClusterSecretStatus(clusterSecret, corev1alpha1.StateProcessing); err != nil {
					c.eventRecorder.Event(clusterSecret, corev1.EventTypeWarning, "Error", err.Error())
					return err
				}
			}
		} else {
			if err := c.updateClusterSecretStatus(clusterSecret, corev1alpha1.StateDeleting); err != nil {
				c.eventRecorder.Event(clusterSecret, corev1.EventTypeWarning, "Error", err.Error())
				return err
			}
		}
	}

	// reconcile all determined secrets (as determined in operations), and update status (if applicable) to Ready or Error, respectively
	var merr *multierror.Error
	for key, operation := range operations {
		if operation.new == nil {
			// this is a deletion
			// note: we can assume that operation.old is not nil because of the way how operations was defined
			klog.V(2).Infof("deleting secret %s/%s (if existing)", key.namespace, key.name)
			err := c.kubeclient.CoreV1().Secrets(key.namespace).Delete(
				context.TODO(),
				key.name,
				metav1.DeleteOptions{Preconditions: &metav1.Preconditions{ResourceVersion: &operation.old.ResourceVersion}},
			)
			if err != nil {
				if !errors.IsNotFound(err) {
					merr = multierror.Append(
						merr, fmt.Errorf("error deleting secret %s/%s", key.namespace, key.name), err,
					)
				}
			}
			if recorder, ok := c.synchronizer.(Recorder); ok {
				recorder.RecordDeletion(operation.old)
			}
		} else if operation.old == nil {
			// this is a creation
			// note: this can fail in particular if the secret already exists, but is not managed by us
			klog.V(2).Infof("create secret %s/%s", key.namespace, key.name)
			secret, err := c.kubeclient.CoreV1().Secrets(key.namespace).Create(
				context.TODO(),
				operation.new,
				metav1.CreateOptions{FieldManager: ControllerName},
			)
			if err != nil {
				merr = multierror.Append(
					merr, fmt.Errorf("error creating secret %s/%s", key.namespace, key.name), err,
				)
			}
			if recorder, ok := c.synchronizer.(Recorder); ok {
				recorder.RecordCreation(secret)
			}
		} else {
			// this is an update
			klog.V(2).Infof("update secret %s/%s", key.namespace, key.name)
			secret, err := c.kubeclient.CoreV1().Secrets(key.namespace).Update(
				context.TODO(),
				operation.new,
				metav1.UpdateOptions{FieldManager: ControllerName},
			)
			if err != nil {
				merr = multierror.Append(merr, fmt.Errorf("error updating secret %s/%s", key.namespace, key.name), err)
			}
			if recorder, ok := c.synchronizer.(Recorder); ok {
				recorder.RecordUpdate(operation.old, secret)
			}
		}
	}
	if merr.ErrorOrNil() != nil {
		if clusterSecret != nil {
			if err := c.updateClusterSecretStatus(clusterSecret, corev1alpha1.StateError); err != nil {
				multierror.Append(merr, err)
			}
		}
		if clusterSecret != nil {
			c.eventRecorder.Event(clusterSecret, corev1.EventTypeWarning, "Error", merr.Error())
		}
		return merr
	}
	if clusterSecret != nil {
		c.eventRecorder.Eventf(clusterSecret, corev1.EventTypeNormal, "ClusterSecretReconcile", "Successfully reconciled clustersecret %s", clusterSecret.Name)
	}
	if clusterSecret != nil && clusterSecret.DeletionTimestamp.IsZero() {
		if err := c.updateClusterSecretStatus(clusterSecret, corev1alpha1.StateReady); err != nil {
			c.eventRecorder.Event(clusterSecret, corev1.EventTypeWarning, "Error", err.Error())
			return err
		}
	}

	// unset finalizer
	if clusterSecret != nil && !clusterSecret.DeletionTimestamp.IsZero() {
		if err := c.unsetClusterSecretFinalizer(clusterSecret); err != nil {
			c.eventRecorder.Event(clusterSecret, corev1.EventTypeWarning, "Error", err.Error())
			return err
		}
	}

	// return
	return nil
}
