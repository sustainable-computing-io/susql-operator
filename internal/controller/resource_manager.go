/*
Copyright 2023.

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
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	susqlv1 "github.com/sustainable-computing-io/susql-operator/api/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// Function to filter pods with matching labels from namespace label is defined
func (r *LabelGroupReconciler) filterPodsInNamespace(ctx context.Context, namespace string, labelSelector map[string]string) ([]string, error) {
	// Initialize list options with label selector
	listOptions := &client.ListOptions{
		Namespace:     namespace,
		LabelSelector: labels.SelectorFromSet(labels.Set(labelSelector)),
	}

	// List pods in the specified namespace with label selector applied
	var podList v1.PodList
	if err := r.Client.List(ctx, &podList, listOptions); err != nil {
		r.Logger.V(5).Info(fmt.Sprintf("[filterPodsInNamespace] labelSelector: %#v", labelSelector))
		r.Logger.V(5).Info(fmt.Sprintf("[filterPodsInNamespace] ctx: %#v", ctx))
		r.Logger.V(5).Info(fmt.Sprintf("[filterPodsInNamespace] podList: %#v", podList))
		r.Logger.V(5).Info(fmt.Sprintf("[filterPodsInNamespace] listOptions: %#v", listOptions))
		r.Logger.V(0).Error(err, "[filterPodsInNamespace] List Error:")
		return nil, err
	}

	var podNames []string
	for _, pod := range podList.Items {
		podNames = append(podNames, pod.Name)
	}

	return podNames, nil
}

// Functions to get data from the cluster
func (r *LabelGroupReconciler) GetPodNamesMatchingLabels(ctx context.Context, labelGroup *susqlv1.LabelGroup) ([]string, []string, error) {
	pods := &v1.PodList{}

	if err := r.List(ctx, pods, client.UnsafeDisableDeepCopy, (client.MatchingLabels)(labelGroup.Status.KubernetesLabels)); err != nil {
		r.Logger.V(5).Info(fmt.Sprintf("[GetPodNamesMatchingLabels] pods: %#v", pods))
		r.Logger.V(5).Info(fmt.Sprintf("[GetPodNamesMatchingLabels] labelgroup: %#v", labelGroup))
		r.Logger.V(0).Error(err, "[GetPodNamesMatchingLabels] List Error:")
		return nil, nil, err
	}

	var podNames []string
	var namespaceNames []string

	for _, pod := range pods.Items {
		podNames = append(podNames, pod.Name)
		namespaceNames = append(namespaceNames, pod.Namespace)
	}

	return podNames, namespaceNames, nil
}
