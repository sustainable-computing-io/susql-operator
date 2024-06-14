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

// EnergyDataSpec defines the desired state of EnergyData
type EnergyDataSpec struct {
	LabelGroupName string `json:"labelGroupName"`
}

// EnergyDataStatus defines the observed state of EnergyData
type EnergyDataStatus struct {
	TotalEnergy string `json:"totalEnergy"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// EnergyData is the Schema for the energydatas API
type EnergyData struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EnergyDataSpec   `json:"spec,omitempty"`
	Status EnergyDataStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EnergyDataList contains a list of EnergyData
type EnergyDataList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EnergyData `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EnergyData{}, &EnergyDataList{})
}
