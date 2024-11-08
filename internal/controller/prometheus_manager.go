/*
Copyright 2023, 2024.

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
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	maxQueryTime                         = "1y" // Look back 'maxQueryTime' for the most recent value
	sourceRoundTripper http.RoundTripper = nil
	susqlRoundTripper  http.RoundTripper = nil
)

// Functions to get data from the cluster
func (r *LabelGroupReconciler) GetMostRecentValue(susqlPrometheusQuery string) (float64, error) {
	// Return the most recent value found in the table
	if susqlRoundTripper == nil {
		if strings.HasPrefix(r.SusQLPrometheusDatabaseUrl, "https://") {
			rttls := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
			susqlRoundTripper = config.NewAuthorizationCredentialsRoundTripper("Bearer", config.NewFileSecret("/var/run/secrets/kubernetes.io/serviceaccount/token"), rttls)
		}
	}
	client, err := api.NewClient(api.Config{
		Address:      r.SusQLPrometheusDatabaseUrl,
		RoundTripper: susqlRoundTripper,
	})

	if err != nil {
		r.Logger.V(0).Error(err, fmt.Sprintf("[GetMostRecentValue] Couldn't create HTTP client.\n")+
			fmt.Sprintf("\tQuery:  %s\n", susqlPrometheusQuery)+
			fmt.Sprintf("\tSusQLPrometheusDatabaseUrl:  %s\n", r.SusQLPrometheusDatabaseUrl))
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	queryString := fmt.Sprintf("max_over_time(%s[%s])", susqlPrometheusQuery, maxQueryTime)
	results, warnings, err := v1api.Query(ctx, queryString, time.Now(), v1.WithTimeout(0*time.Second))

	r.Logger.V(5).Info(fmt.Sprintf("[GetMostRecentValue] Query: %s", queryString)) // trace
	r.Logger.V(5).Info(fmt.Sprintf("[GetMostRecentValue] Results: '%v'", results)) // trace

	if len(warnings) > 0 {
		r.Logger.V(0).Info(fmt.Sprintf("WARNING [GetMostRecentValue] %v\n", warnings) +
			fmt.Sprintf("\tQuery:  %s\n", queryString) +
			fmt.Sprintf("\tSusQLPrometheusDatabaseUrl:  %s", r.SusQLPrometheusDatabaseUrl))
	}

	if err != nil {
		r.Logger.V(0).Error(err, fmt.Sprintf("[GetMostRecentValue] Querying Prometheus didn't work.\n")+
			fmt.Sprintf("\tQuery:  %s\n", queryString)+
			fmt.Sprintf("\tSusQLPrometheusDatabaseUrl:  %s\n", r.SusQLPrometheusDatabaseUrl))
		return 0.0, err
	}

	if len(results.(model.Vector)) > 0 {
		return float64(results.(model.Vector)[0].Value), err
	} else {
		return 0.0, err
	}
}

func (r *LabelGroupReconciler) GetMetricValuesForPodNames(metricName string, podNames []string, namespaceName string) (map[string]float64, error) {
	if sourceRoundTripper == nil {
		if strings.HasPrefix(r.SourcePrometheusUrl, "https://") {
			rttls := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
			sourceRoundTripper = config.NewAuthorizationCredentialsRoundTripper("Bearer", config.NewFileSecret("/var/run/secrets/kubernetes.io/serviceaccount/token"), rttls)
		}
	}
	client, err := api.NewClient(api.Config{
		Address:      r.SourcePrometheusUrl,
		RoundTripper: sourceRoundTripper,
	})

	if err != nil {
		r.Logger.V(0).Error(err, "[GetMetricValuesForPodNames] Couldn't create an HTTP client.\n"+
			fmt.Sprintf("\tmetricName: %s\n", metricName)+
			fmt.Sprintf("\tSourcePrometheusUrl: %s\n", r.SourcePrometheusUrl))
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//	queryString := fmt.Sprintf("%s{pod_name=~\"%s\",mode=\"dynamic\"}", metricName, strings.Join(podNames, "|"))

	// new query for issue 2: can improve runtime efficiency...
	queryString := fmt.Sprintf("sum(%s{pod_name=\"%s\",container_namespace=\"%s\",mode=\"dynamic\"})", metricName, podNames[0], namespaceName)
	queryString = queryString + "+" + fmt.Sprintf("sum(%s{pod_name=\"%s\",container_namespace=\"%s\",mode=\"idle\"})", metricName, podNames[0], namespaceName)
	for i := 1; i < len(podNames); i++ {
		queryString = queryString + "+" + fmt.Sprintf("sum(%s{pod_name=\"%s\",container_namespace=\"%s\",mode=\"dynamic\"})", metricName, podNames[i], namespaceName)
		queryString = queryString + "+" + fmt.Sprintf("sum(%s{pod_name=\"%s\",container_namespace=\"%s\",mode=\"idle\"})", metricName, podNames[i], namespaceName)
	}

	results, warnings, err := v1api.Query(ctx, queryString, time.Now(), v1.WithTimeout(0*time.Second))

	r.Logger.V(5).Info(fmt.Sprintf("[GetMetricValuesForPodNames] Query: %s", queryString)) // trace

	if err != nil {
		r.Logger.V(0).Error(err, "[GetMetricValuesForPodNames] Querying Prometheus didn't work.\n"+
			fmt.Sprintf("\tmetricName: %s\n", metricName)+
			fmt.Sprintf("\tSourcePrometheusUrl: %s\n", r.SourcePrometheusUrl)+
			fmt.Sprintf("\tqueryString: %s\n", queryString))
		return nil, err
	}

	if len(warnings) > 0 {
		r.Logger.V(0).Info(fmt.Sprintf("WARNING [GetMetricValuesForPodNames] %v\n", warnings) +
			fmt.Sprintf("\tmetricName: %s\n", metricName) +
			fmt.Sprintf("\tSourcePrometheusUrl: %s\n", r.SourcePrometheusUrl) +
			fmt.Sprintf("\tqueryString: %s", queryString))
	}

	metricValues := make(map[string]float64, len(results.(model.Vector)))

	for _, result := range results.(model.Vector) {
		r.Logger.V(5).Info(fmt.Sprintf("[GetMetricValuesForPodNames] Container id %s value is %f.", string(result.Metric["container_id"]), float64(result.Value))) // trace
		metricValues[string(result.Metric["container_id"])] = float64(result.Value)
	}

	return metricValues, nil
}

type SusqlMetrics struct {
	totalEnergy *prometheus.GaugeVec
	totalCarbon *prometheus.GaugeVec
}

var (
	susqlMetrics = &SusqlMetrics{
		totalEnergy: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "susql",
			Name:      "total_energy_joules",
			Help:      "Accumulated energy over time for set of labels",
		}, susqlPrometheusLabelNames),
		totalCarbon: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "susql",
			Name:      "total_carbon_dioxide_grams",
			Help:      "Accumulated carbon dioxide grams over time for set of labels",
		}, susqlPrometheusLabelNames),
	}

	prometheusRegistry *prometheus.Registry
	prometheusHandler  http.Handler
)

func (r *LabelGroupReconciler) InitializeMetricsExporter() {
	// Initiate the exporting of prometheus metrics for the energy
	r.Logger.V(5).Info("Entering InitializeMetricsExporter().")
	if prometheusRegistry == nil {
		prometheusRegistry = prometheus.NewRegistry()
		prometheusRegistry.MustRegister(susqlMetrics.totalEnergy, susqlMetrics.totalCarbon)

		prometheusHandler = promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{Registry: prometheusRegistry})
		http.Handle("/metrics", prometheusHandler)

		if metricsUrl, parseErr := url.Parse(r.SusQLPrometheusMetricsUrl); parseErr == nil {
			r.Logger.V(2).Info(fmt.Sprintf("[InitializeMetricsExporter] Serving metrics at '%s:%s'...\n", metricsUrl.Hostname(), metricsUrl.Port()))

			go func() {
				err := http.ListenAndServe(metricsUrl.Hostname()+":"+metricsUrl.Port(), nil)

				if err != nil {
					r.Logger.V(0).Error(err, "PANIC [InitializeMetricsExporter] ListenAndServe")
					panic("PANIC [InitializeMetricsExporter]: ListenAndServe: " + err.Error())
				}
			}()
		} else {
			r.Logger.V(0).Error(parseErr, fmt.Sprintf("PANIC [InitializeMetricsExporter] Parsing the URL '%s' to set the metrics address didn't work.", r.SusQLPrometheusMetricsUrl))
			panic(fmt.Sprintf("PANIC [InitializeMetricsExporter]: Parsing the URL '%s' to set the metrics address didn't work (%v)", r.SusQLPrometheusMetricsUrl, parseErr))
		}
	}
}

func (r *LabelGroupReconciler) SetAggregatedEnergyForLabels(totalEnergy float64, prometheusLabels map[string]string) error {
	// Save aggregated energy to Prometheus table
	susqlMetrics.totalEnergy.With(prometheusLabels).Set(totalEnergy)

	r.Logger.V(5).Info(fmt.Sprintf("[SetAggregatedEnergyForLabels] Setting energy %f for %v.", totalEnergy, prometheusLabels)) // trace

	return nil
}

func (r *LabelGroupReconciler) SetAggregatedCarbonForLabels(totalCarbon float64, prometheusLabels map[string]string) error {
	// Save aggregated carbon to Prometheus table
	susqlMetrics.totalCarbon.With(prometheusLabels).Set(totalCarbon)

	r.Logger.V(5).Info(fmt.Sprintf("[SetAggregatedEnergyForLabels] Setting carbon %f for %v.", totalCarbon, prometheusLabels)) // trace

	return nil
}
