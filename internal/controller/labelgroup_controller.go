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
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	susql "github.com/metalcycling/susql/api/v1"
)

// LabelGroupReconciler reconciles a LabelGroup object
type LabelGroupReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	KeplerPrometheusUrl string
}

const (
	energyMetricName = "kepler_container_joules_total" // Kepler metric to query
	samplingRate = 2 * time.Second // Sampling rate for all the label groups
)

//+kubebuilder:rbac:groups=susql.ibm.com,resources=labelgroups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=susql.ibm.com,resources=labelgroups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=susql.ibm.com,resources=labelgroups/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LabelGroup object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *LabelGroupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Get label group to process
	labelGroup := &susql.LabelGroup{}

	// Get deep copy of LabelGroup object in reconciler cache
	if err := r.Get(ctx, req.NamespacedName, labelGroup); err != nil {
		// LabelGroup not found
		return ctrl.Result{}, nil
	}
	
	// Get list of pods matching the label group
	podNames, err := r.GetPodNamesMatchingLabels(ctx, labelGroup)

	if err != nil {
		fmt.Println("Error getting pods")
		return ctrl.Result{}, err
	}

	// Aggregate Kepler measurements for these set of pods
	metricValues, err := r.GetMetricValuesForPodNames(energyMetricName, podNames)

	// Compute total energy
	// 1) Get the current total energy from ETCD
	var totalEnergy float64

	if value, err := strconv.ParseFloat(labelGroup.Status.TotalEnergy, 64); err == nil {
		totalEnergy = value
	} else {
		totalEnergy = 0.0
	}

	if labelGroup.Status.ActiveContainerIds == nil {
		// First pass with this pod group
		labelGroup.Status.ActiveContainerIds = make(map[string]float64)
	}

	// 2) Check if the active containers are still active by comparing them to the current ones
	//    - In the set of new containers, remove all containers that are active
	for containerId, oldValue := range labelGroup.Status.ActiveContainerIds {
		if newValue, found := metricValues[containerId]; found {
			totalEnergy += (newValue - oldValue)
			labelGroup.Status.ActiveContainerIds[containerId] = newValue
			delete(metricValues, containerId)
		} else {
			// Delete inactive container since it doesn't appear in queried containers
			delete(labelGroup.Status.ActiveContainerIds, containerId)
		}
	}

	// 3) Add the values of the remaining new containers to the total energy and to the active containers
	for containerId, newValue := range metricValues {
		totalEnergy += newValue
		labelGroup.Status.ActiveContainerIds[containerId] = newValue
	}

	// 4) Update ETCD with the values
	labelGroup.Status.TotalEnergy = fmt.Sprintf("%.2f", totalEnergy)

	if err := r.Status().Update(ctx, labelGroup); err != nil {
		return ctrl.Result{}, err
	}

	// Requeue
	return ctrl.Result{RequeueAfter: samplingRate}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LabelGroupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&susql.LabelGroup{}).
		Complete(r)
}