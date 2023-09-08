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
	//"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/common/model"

	"github.com/prometheus/client_golang/prometheus"
	//"github.com/prometheus/client_golang/prometheus/promhttp"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

// Functions to get data from the cluster
func (r *LabelGroupReconciler) GetMetricValuesForPodNames(metricName string, podNames []string) (map[string]float64, error) {
	client, err := api.NewClient(api.Config{
		Address: r.KeplerPrometheusUrl,
	})

	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	queryString := fmt.Sprintf("%s{pod_name=~\"%s\",mode=\"dynamic\"}", metricName, strings.Join(podNames, "|"))
	results, warnings, err := v1api.Query(ctx, queryString, time.Now(), v1.WithTimeout(0*time.Second))

	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		return nil, err
	}

	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}

	metricValues := make(map[string]float64, len(results.(model.Vector)))

	for _, result := range results.(model.Vector) {
		metricValues[string(result.Metric["container_id"])] = float64(result.Value)
	}

	return metricValues, nil
}

type SusqlMetrics struct {
	totalEnergy *prometheus.GaugeVec
}

func (r *LabelGroupReconciler) SetAggregatedEnergyForLabels(totalEnergy float64, prometheusLabels map[string]string) (error) {
	/*client, err := api.NewClient(api.Config{
		Address: r.SusQLPrometheusUrl,
	})

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	v1api.Series()

	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}*/

	/*susqlMetrics := &SusqlMetrics{
		totalEnergy: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "susql",
			Name: "total_energy_joules",
			Help: "Accumulated energy over time for set of labels",
		},  susqlPrometheusLabelNames),
	}

	fmt.Println(susqlMetrics)

	http.Handle("/metrics", promhttp.Handler())*/

	return nil
}
