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

	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	susqlv1 "github.com/sustainable-computing-io/susql-operator/api/v1"
)

// Load Data from new CRD to persist
func (r *LabelGroupReconciler) loadDataFromCRD(ctx context.Context, labelGroupName string) (float64, error) {
	energyData := &susqlv1.EnergyData{}
	err := r.Get(ctx, client.ObjectKey{Name: labelGroupName}, energyData)
	if err != nil {
		if errors.IsNotFound(err) {
			return 0, nil // No data found
		}
		return 0, err
	}

	totalEnergy, err := strconv.ParseFloat(energyData.Status.TotalEnergy, 64)
	if err != nil {
		return 0, err
	}

	return totalEnergy, nil
}

// Save Data to CRD to persist
func (r *LabelGroupReconciler) saveDataToCRD(ctx context.Context, labelGroupName string, totalEnergy float64) error {
	energyData := &susqlv1.EnergyData{}
	err := r.Get(ctx, client.ObjectKey{Name: labelGroupName}, energyData)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if errors.IsNotFound(err) {
		// Create new CRD instance if not found
		energyData = &susqlv1.EnergyData{
			ObjectMeta: metav1.ObjectMeta{
				Name: labelGroupName,
			},
			Spec: susqlv1.EnergyDataSpec{
				LabelGroupName: labelGroupName,
			},
		}
	}

	energyData.Status.TotalEnergy = fmt.Sprintf("%f", totalEnergy)
	return r.Status().Update(ctx, energyData)
}
