/*
Copyright 2024.

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
	coreruntime "runtime"
	"strconv"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	susqlv1 "github.com/sustainable-computing-io/susql-operator/api/v1"
)

// LabelGroupReconciler reconciles a LabelGroup object
type LabelGroupReconciler struct {
	client.Client
	Scheme                     *runtime.Scheme
	KeplerPrometheusUrl        string
	KeplerMetricName           string
	SusQLPrometheusDatabaseUrl string
	SusQLPrometheusMetricsUrl  string
	SamplingRate               time.Duration // Sampling rate for all label groups
	StaticCarbonIntensity      float64
	Logger                     logr.Logger
}

const (
	susqlMetricName = "susql_total_energy_joules" // SusQL metric to query
	fixingDelay     = 15 * time.Second            // Time to wait in the event the label group was badly constructed
	nopodDelay      = 15 * time.Second            // Time to wait in the event no pods are found
	errorDelay      = 1 * time.Second             // Time to wait when an error happens due to network connectivity issues
)

var (
	susqlKubernetesLabelNames = []string{"susql.label/1", "susql.label/2", "susql.label/3", "susql.label/4", "susql.label/5", "susql.label/6"} // Names of the SusQL Kubernetes labels
	susqlPrometheusLabelNames = []string{"susql_label_1", "susql_label_2", "susql_label_3", "susql_label_4", "susql_label_5", "susql_label_6"} // Names of the SusQL Prometheus labels
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

	r.Logger.V(5).Info("[Reconcile] Entered Reconcile().")

	var m coreruntime.MemStats
	coreruntime.ReadMemStats(&m)
	r.Logger.V(5).Info(fmt.Sprintf("Memory: Alloc=%.2f MB  TotalAlloc=%.2f MB  Sys= %.2f MB  NumGC=%v", float32(m.Alloc)/1024.0/1024.0, float32(m.TotalAlloc)/1024.0/1024.0, float32(m.Sys)/1024.0/1024.0, m.NumGC))

	// Get label group object to process if it exists
	labelGroup := &susqlv1.LabelGroup{}

	// Get deep copy of LabelGroup object in reconciler cache
	err := r.Get(ctx, req.NamespacedName, labelGroup)
	if err != nil {
		// LabelGroup not found
		return ctrl.Result{}, nil
	}

	// Check that the susql prometheus labels are created
	if len(labelGroup.Status.PrometheusLabels) == 0 && labelGroup.Status.Phase != susqlv1.Initializing {
		r.Logger.V(1).Info(fmt.Sprintf("[Reconcile] The SusQL prometheus labels for LabelGroup '%s' in namespace '%s' have not been created. Reinitializing this LabelGroup.", labelGroup.Name, labelGroup.Namespace))

		labelGroup.Status.Phase = susqlv1.Initializing

		if err := r.Status().Update(ctx, labelGroup); err != nil {
			r.Logger.V(0).Error(err, "[Reconcile] Couldn't update the phase.")
		}

		return ctrl.Result{}, nil
	}

	// Decide what action to take based on the state of the labelGroup
	switch labelGroup.Status.Phase {
	case susqlv1.Initializing:
		r.Logger.V(5).Info("[Reconcile-Initializing] Entered initializing case.")
		if len(labelGroup.Spec.Labels) > len(susqlPrometheusLabelNames) {
			r.Logger.V(0).Error(fmt.Errorf("[Reconcile-Initializing] The number of provided labels is greater than the maximum number of supported labels (e.g., up to %d labels).", len(susqlPrometheusLabelNames)), "")
			return ctrl.Result{RequeueAfter: fixingDelay}, nil
		}

		susqlKubernetesLabels := make(map[string]string)

		for ldx := 0; ldx < len(labelGroup.Spec.Labels); ldx++ {
			susqlKubernetesLabels[susqlKubernetesLabelNames[ldx]] = labelGroup.Spec.Labels[ldx]
		}

		susqlPrometheusLabels := make(map[string]string)

		for ldx := 0; ldx < len(susqlKubernetesLabelNames); ldx++ {
			if ldx < len(labelGroup.Spec.Labels) {
				susqlPrometheusLabels[susqlPrometheusLabelNames[ldx]] = labelGroup.Spec.Labels[ldx]
			} else {
				susqlPrometheusLabels[susqlPrometheusLabelNames[ldx]] = ""
			}
		}

		var susqlPrometheusQuery string
		susqlPrometheusQuery = susqlMetricName
		susqlPrometheusQuery += "{"
		for ldx := 0; ldx < len(susqlKubernetesLabelNames); ldx++ {
			if ldx < len(labelGroup.Spec.Labels) {
				susqlPrometheusQuery += fmt.Sprintf("%s=\"%s\"", susqlPrometheusLabelNames[ldx], labelGroup.Spec.Labels[ldx])
			} else {
				susqlPrometheusQuery += fmt.Sprintf("%s=\"\"", susqlPrometheusLabelNames[ldx])
			}
			if ldx < len(susqlKubernetesLabelNames)-1 {
				susqlPrometheusQuery += ","
			}
		}
		susqlPrometheusQuery += "}"

		labelGroup.Status.KubernetesLabels = susqlKubernetesLabels
		labelGroup.Status.PrometheusLabels = susqlPrometheusLabels
		labelGroup.Status.SusQLPrometheusQuery = susqlPrometheusQuery
		labelGroup.Status.Phase = susqlv1.Reloading

		if err := r.Status().Update(ctx, labelGroup); err != nil {
			r.Logger.V(0).Error(err, "[Reconcile-Initializing] Couldn't update status of the LabelGroup.")
			return ctrl.Result{RequeueAfter: fixingDelay}, nil
		}

		// Requeue
		return ctrl.Result{}, nil

	case susqlv1.Reloading:
		r.Logger.V(5).Info("[Reconcile-Reloading] Entered reloading case.")
		// Reload data from existing database
		if !labelGroup.Spec.DisableUsingMostRecentValue {
			totalEnergy, err := r.GetMostRecentValue(labelGroup.Status.SusQLPrometheusQuery)

			if err != nil {
				r.Logger.V(0).Error(err, "[Reconcile-Reloading] Couldn't retrieve most recent value.")
				return ctrl.Result{RequeueAfter: fixingDelay}, nil
			}

			labelGroup.Status.TotalEnergy = fmt.Sprintf("%f", totalEnergy)

			labelGroup.Status.TotalGCO2 = fmt.Sprintf("%.15f", float64(totalEnergy)*r.StaticCarbonIntensity)
		}

		labelGroup.Status.Phase = susqlv1.Aggregating

		if err := r.Status().Update(ctx, labelGroup); err != nil {
			r.Logger.V(0).Error(err, "[Reconcile-Reloading] Couldn't update status of the LabelGroup.")
			return ctrl.Result{RequeueAfter: fixingDelay}, nil
		}

		// Requeue
		return ctrl.Result{}, nil

	case susqlv1.Aggregating:
		r.Logger.V(5).Info("[Reconcile-Aggregating] Entered aggregating case.") // trace

		// Get list of pods matching the label group and namespace
		podsInNamespace, err := r.filterPodsInNamespace(ctx, labelGroup.Namespace, labelGroup.Status.KubernetesLabels)

		if err != nil || len(podsInNamespace) == 0 {
			r.Logger.V(5).Info("[Reconcile-Aggregating] Unable to get podlist.")                                                 // trace
			r.Logger.V(5).Info(fmt.Sprintf("[Reconcile-Aggregating] LabelName: %s", labelGroup.Name))                            // trace
			r.Logger.V(5).Info(fmt.Sprintf("[Reconcile-Aggregating] Namespace: %s", labelGroup.Namespace))                       // trace
			r.Logger.V(5).Info(fmt.Sprintf("[Reconcile-Aggregating] KubernetesLabels: %#v", labelGroup.Status.KubernetesLabels)) // trace
			r.Logger.V(5).Info(fmt.Sprintf("[Reconcile-Aggregating] podNamesinNamespace: %s", podsInNamespace))                  // trace
			r.Logger.V(5).Info(fmt.Sprintf("[Reconcile] ctx: %#v", ctx))                                                         // trace
			if err != nil {
				r.Logger.V(0).Error(err, "[Reconcile-Aggregating] ERROR: Couldn't get pods for the labels provided.")
			}

			return ctrl.Result{RequeueAfter: nopodDelay}, nil
		}

		// Aggregate Kepler measurements for these set of pods
		metricValues, err := r.GetMetricValuesForPodNames(r.KeplerMetricName, podsInNamespace, labelGroup.Namespace)

		if err != nil {
			r.Logger.V(0).Error(err, "[Reconcile-Aggregating] Querying Prometheus didn't work.")
			return ctrl.Result{RequeueAfter: errorDelay}, nil
		}

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
		r.Logger.V(5).Info(fmt.Sprintf("[Reconcile-Aggregating] metricValues: %#v", metricValues))                               // trace
		r.Logger.V(5).Info(fmt.Sprintf("[Reconcile-Aggregating] ActiveContainerIds: %#v", labelGroup.Status.ActiveContainerIds)) // trace

		// 3) Add the values of the remaining new containers to the total energy and update the list of active containers
		for containerId, newValue := range metricValues {
			totalEnergy += newValue
			labelGroup.Status.ActiveContainerIds[containerId] = newValue
		}

		// 4) Update ETCD with the values
		labelGroup.Status.TotalEnergy = fmt.Sprintf("%.2f", totalEnergy)

		labelGroup.Status.TotalGCO2 = fmt.Sprintf("%.15f", float64(totalEnergy)*r.StaticCarbonIntensity)

		if err := r.Status().Update(ctx, labelGroup); err != nil {
			return ctrl.Result{}, err
		}

		// 5) Add energy aggregation to Prometheus table
		r.SetAggregatedEnergyForLabels(totalEnergy, labelGroup.Status.PrometheusLabels)

		// Requeue
		return ctrl.Result{RequeueAfter: r.SamplingRate}, nil

	default:
		r.Logger.V(5).Info("[Reconcile-default] Entered default case.")
		// First time seeing this object
		labelGroup.Status.Phase = susqlv1.Initializing

		if err := r.Status().Update(ctx, labelGroup); err != nil {
			r.Logger.V(0).Error(err, "[Reconcile-default] Couldn't set object to 'Initializing'.")
		}

		return ctrl.Result{}, nil
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *LabelGroupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	controllerManager := ctrl.NewControllerManagedBy(mgr).
		For(&susqlv1.LabelGroup{}).
		// Watch for changes to Pods and enqueue requests for LabelGroup owners
		Owns(&corev1.Pod{}).
		Complete(r)

	r.Logger.V(5).Info("[SetupWithManager] Initializing Metrics Exporter.")

	// Start server to export metrics
	r.InitializeMetricsExporter()

	return controllerManager
}
