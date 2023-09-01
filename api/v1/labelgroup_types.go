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

	// Set of labels to be tracked for energy measurments
	Labels map[string]string `json:"labels,omitempty"`
}

// LabelGroupStatus defines the observed state of LabelGroup
type LabelGroupStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// TotalEnergy keeps track of the accumulated energy over time
	TotalEnergy string `json:"totalEnergy"`

	// Active containers associated with these set of labels
	ActiveContainerIds map[string]float64 `json:"activeContainerIds"`
}

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
