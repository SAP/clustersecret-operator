/*
SPDX-FileCopyrightText: 2026 SAP SE or an SAP affiliate company and clustersecret-operator contributors
SPDX-License-Identifier: Apache-2.0
*/

package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	conversionutils "github.com/sap/clustersecret-operator/internal/utils/conversion"
	stringutils "github.com/sap/clustersecret-operator/internal/utils/strings"

	corev1alpha1 "github.com/sap/clustersecret-operator/pkg/apis/core.cs.sap.com/v1alpha1"
)

func (c *Controller) setClusterSecretFinalizer(clusterSecret *corev1alpha1.ClusterSecret) error {
	if stringutils.ContainsString(clusterSecret.Finalizers, ControllerName) {
		return nil
	}
	newClusterSecret := clusterSecret.DeepCopy()
	newClusterSecret.Finalizers = append(newClusterSecret.Finalizers, ControllerName)
	updatedClusterSecret, err := c.coreclient.CoreV1alpha1().ClusterSecrets().Update(context.TODO(), newClusterSecret, metav1.UpdateOptions{FieldManager: ControllerName})
	if err != nil {
		return err
	}
	if recorder, ok := c.synchronizer.(Recorder); ok {
		recorder.RecordUpdate(clusterSecret, updatedClusterSecret)
	}
	*clusterSecret = *updatedClusterSecret
	return nil
}

func (c *Controller) unsetClusterSecretFinalizer(clusterSecret *corev1alpha1.ClusterSecret) error {
	if !stringutils.ContainsString(clusterSecret.Finalizers, ControllerName) {
		return nil
	}
	newClusterSecret := clusterSecret.DeepCopy()
	newClusterSecret.Finalizers = stringutils.RemoveString(newClusterSecret.Finalizers, ControllerName)
	updatedClusterSecret, err := c.coreclient.CoreV1alpha1().ClusterSecrets().Update(context.TODO(), newClusterSecret, metav1.UpdateOptions{FieldManager: ControllerName})
	if err != nil {
		return err
	}
	if recorder, ok := c.synchronizer.(Recorder); ok {
		recorder.RecordUpdate(clusterSecret, updatedClusterSecret)
	}
	*clusterSecret = *updatedClusterSecret
	return nil
}

func (c *Controller) updateClusterSecretStatus(clusterSecret *corev1alpha1.ClusterSecret, state string) error {
	// return immediately if status is already up-to-date
	if clusterSecret.Status.ObservedGeneration == clusterSecret.Generation && clusterSecret.Status.State == state {
		return nil
	}

	// store current time for consistent later use
	now := metav1.Now()

	// build new ready condition
	var readyCondition corev1alpha1.ClusterSecretCondition
	for _, cond := range clusterSecret.Status.Conditions {
		if cond.Type == corev1alpha1.ClusterSecretConditionTypeReady {
			readyCondition = cond
			break
		}
	}
	newReadyCondition := corev1alpha1.ClusterSecretCondition{
		Type: corev1alpha1.ClusterSecretConditionTypeReady,
	}
	if state == corev1alpha1.StateReady {
		newReadyCondition.Status = corev1.ConditionTrue
	} else {
		newReadyCondition.Status = corev1.ConditionFalse
	}
	newReadyCondition.LastUpdateTime = now
	if newReadyCondition.Status == readyCondition.Status {
		newReadyCondition.LastTransitionTime = readyCondition.LastTransitionTime
	} else {
		newReadyCondition.LastTransitionTime = now
	}
	newReadyCondition.Reason = "ClusterSecret" + state
	newReadyCondition.Message = ""

	// prepare new clustersecret (with new status)
	newClusterSecret := clusterSecret.DeepCopy()
	newClusterSecret.Status = corev1alpha1.ClusterSecretStatus{
		ObservedGeneration: newClusterSecret.Generation,
		State:              state,
		Conditions:         []corev1alpha1.ClusterSecretCondition{newReadyCondition},
	}

	// update status
	updatedClusterSecret, err := c.coreclient.CoreV1alpha1().ClusterSecrets().UpdateStatus(context.TODO(), newClusterSecret, metav1.UpdateOptions{FieldManager: ControllerName})
	if err != nil {
		return err
	}
	if recorder, ok := c.synchronizer.(Recorder); ok {
		recorder.RecordUpdate(clusterSecret, updatedClusterSecret)
	}
	*clusterSecret = *updatedClusterSecret
	return nil
}

func buildNamespaceSelectorFromClusterSecret(clusterSecret *corev1alpha1.ClusterSecret) labels.Selector {
	if clusterSecret.Spec.NamespaceSelector == nil {
		return labels.Everything()
	}
	namespaceSelector, err := metav1.LabelSelectorAsSelector(clusterSecret.Spec.NamespaceSelector)
	if err != nil {
		panic("this cannot happen")
	}
	return namespaceSelector
}

func buildSecretFromClusterSecret(namespace string, clusterSecret *corev1alpha1.ClusterSecret) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      clusterSecret.Name,
			Labels: map[string]string{
				LabelKeyName: clusterSecret.Name,
			},
			Annotations: map[string]string{
				AnnotationKeyGeneration: conversionutils.Itoa(clusterSecret.Generation),
			},
		},
		Type: clusterSecret.Spec.Template.Type,
		Data: clusterSecret.Spec.Template.Data,
	}
}
