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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LabelGroupSpec defines the desired state of LabelGroup
type LabelGroupSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Do not use the most recent value stored in the database
	DisableUsingMostRecentValue bool `json:"disableUsingMostRecentValue,omitempty"`

	// List of labels to be tracked for energy measurments (up to 5)
	Labels []string `json:"labels,omitempty"`

	// +kubebuilder:default:="0.00000000011583333"
	// Static Carbon Intensity Factor in Grams CO2 / Joule
	StaticCarbonIntensity string `json:"staticcarbonintensity,omitempty"`
}

// LabelGroupStatus defines the observed state of LabelGroup
type LabelGroupStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Transition phase of the label group
	Phase LabelGroupPhase `json:"phase,omitempty"`

	// SusQL Kubernetes labels constructed from the spec
	KubernetesLabels map[string]string `json:"kubernetesLabels,omitempty"`

	// SusQL Prometheus labels constructed from the spec
	PrometheusLabels map[string]string `json:"prometheusLabels,omitempty"`

	// TotalEnergy keeps track of the accumulated energy over time
	TotalEnergy string `json:"totalEnergy,omitempty"`

	// TotalGCO2 keeps track of the accumulated grams of carbon emission over time
	TotalGCO2 string `json:"totalgco2,omitempty"`

	// Prometheus query to get the total energy for this label group
	SusQLPrometheusQuery string `json:"susqlPrometheusQuery,omitempty"`

	// Active containers associated with these set of labels
	ActiveContainerIds map[string]float64 `json:"activeContainerIds,omitempty"`
}

// LabelGroupPhase defines the label for the LabelGroupStatus
type LabelGroupPhase string

const (
	// Initializing: The label group is picked up for the first time and setup
	Initializing LabelGroupPhase = "Initializing"

	// Reloading: Use most recent value in the database if requested
	Reloading LabelGroupPhase = "Reloading"

	// Aggregating: The label group is aggregating the energy for the registered labels
	Aggregating LabelGroupPhase = "Aggregating"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LabelGroup is the Schema for the labelgroups API
type LabelGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LabelGroupSpec   `json:"spec,omitempty"`
	Status LabelGroupStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LabelGroupList contains a list of LabelGroup
type LabelGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LabelGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LabelGroup{}, &LabelGroupList{})
}
