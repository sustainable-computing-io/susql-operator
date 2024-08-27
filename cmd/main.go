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

package main

import (
	"crypto/tls"
	"flag"
	"os"
	"strconv"
	"time"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	susqlv1 "github.com/sustainable-computing-io/susql-operator/api/v1"
	"github.com/sustainable-computing-io/susql-operator/internal/controller"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	susqlLog = ctrl.Log.WithName("susql")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(susqlv1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func getEnv(key, defval string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defval
}

func main() {
	var enableLeaderElection bool = true
	var probeAddr string = ":8081"
	var keplerPrometheusUrl string = "https://thanos-querier.openshift-monitoring.svc.cluster.local:9091"
	var keplerMetricName string = "kepler_container_joules_total"
	var susqlPrometheusMetricsUrl string = "http://0.0.0.0:8082"
	var susqlPrometheusDatabaseUrl string = "https://thanos-querier.openshift-monitoring.svc.cluster.local:9091"
	var samplingRate string = "2"
	var susqlLogLevel string = "-5"
	// Carbon Intensity Factor in grams CO2 / Joule
	var carbonIntensity string = "0.00000000011583333"
	var carbonMethod string = "static" // options: static, simpledynamic, sdk
	var carbonIntensityUrl string = "https://api.electricitymap.org/v3/carbon-intensity/latest?zone=%s"
	var carbonLocation string = "JP-TK"
	var carbonQueryRate string = "60"
	var carbonQueryFilter string = ".carbonIntensity"

	// NOTE: these can be set as env or flag, flag takes precedence over env
	keplerPrometheusUrlEnv := getEnv("KEPLER-PROMETHEUS-URL", keplerPrometheusUrl)
	keplerMetricNameEnv := getEnv("KEPLER-METRIC-NAME", keplerMetricName)
	susqlPrometheusDatabaseUrlEnv := getEnv("SUSQL-PROMETHEUS-DATABASE-URL", susqlPrometheusDatabaseUrl)
	susqlPrometheusMetricsUrlEnv := getEnv("SUSQL-PROMETHEUS-METRICS-URL", susqlPrometheusMetricsUrl)
	samplingRateEnv := getEnv("SAMPLING-RATE", samplingRate)
	probeAddrEnv := getEnv("HEALTH-PROBE-BIND-ADDRESS", probeAddr)
	susqlLogLevelEnv := getEnv("SUSQL-LOG-LEVEL", susqlLogLevel)
	carbonIntensityEnv := getEnv("CARBON-INTENSITY", carbonIntensity)
	carbonMethodEnv := getEnv("CARBON-METHOD", carbonMethod)
	carbonIntensityUrlEnv := getEnv("CARBON-INTENSITY-URL", carbonIntensityUrl)
	carbonLocationEnv := getEnv("CARBON-LOCATION", carbonLocation)
	carbonQueryRateEnv := getEnv("CARBON-QUERY-RATE", carbonQueryRate)
	carbonQueryFilterEnv := getEnv("CARBON-QUERY-FILTER", carbonQueryFilter)
	enableLeaderElectionEnv, err := strconv.ParseBool(getEnv("LEADER-ELECT", strconv.FormatBool(enableLeaderElection)))
	if err != nil {
		enableLeaderElectionEnv = false
	}

	flag.StringVar(&keplerPrometheusUrl, "kepler-prometheus-url", keplerPrometheusUrlEnv, "The URL for the Prometheus server where Kepler stores the energy data")
	flag.StringVar(&keplerMetricName, "kepler-metric-name", keplerMetricNameEnv, "The metric name to be queried in the kepler Prometheus server")
	flag.StringVar(&susqlPrometheusDatabaseUrl, "susql-prometheus-database-url", susqlPrometheusDatabaseUrlEnv, "The URL for the Prometheus database where SusQL stores the energy data")
	flag.StringVar(&susqlPrometheusMetricsUrl, "susql-prometheus-metrics-url", susqlPrometheusMetricsUrlEnv, "The URL for the Prometheus metrics where SusQL exposes the energy data")
	flag.StringVar(&samplingRate, "sampling-rate", samplingRateEnv, "Sampling rate in seconds")
	flag.StringVar(&probeAddr, "health-probe-bind-address", probeAddrEnv, "The address the probe endpoint binds to.")
	flag.StringVar(&susqlLogLevel, "susql-log-level", susqlLogLevelEnv, "SusQL log level")
	flag.StringVar(&carbonIntensity, "carbon-intensity", carbonIntensityEnv, "Carbon Intensity Factor in grams CO2 / Joule")
	flag.StringVar(&carbonMethod, "carbon-method", carbonMethodEnv, "Method used to calculate CO2 emissions")
	flag.StringVar(&carbonIntensityUrl, "carbon-intensity-url", carbonIntensityUrlEnv, "URL used to query calculate carbon intensity")
	flag.StringVar(&carbonLocation, "carbon-location", carbonLocationEnv, "Location identfier used in carbon intensity query")
	flag.StringVar(&carbonQueryRate, "carbon-query-rate", carbonQueryRateEnv, "How often to query carbon intensity query (minutes)")
	flag.StringVar(&carbonQueryFilter, "carbon-query-filter", carbonQueryFilterEnv, "jq parameter to extract carbon intensity from JSON returned by query")
	flag.BoolVar(&enableLeaderElection, "leader-elect", enableLeaderElectionEnv,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	susqlLogLevelInt, err := strconv.Atoi(susqlLogLevel)
	if err != nil {
		susqlLogLevelInt = -5
	}

	opts := zap.Options{
		Development: true,
		Level:       zapcore.Level(susqlLogLevelInt),
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	susqlLog.Info("SusQL configuration values at runtime")
	susqlLog.Info("enableLeaderElection=" + strconv.FormatBool(enableLeaderElection))
	susqlLog.Info("probeAddr=" + probeAddr)
	susqlLog.Info("keplerPrometheusUrl=" + keplerPrometheusUrl)
	susqlLog.Info("keplerMetricName=" + keplerMetricName)
	susqlLog.Info("susqlPrometheusMetricsUrl=" + susqlPrometheusMetricsUrl)
	susqlLog.Info("susqlPrometheusDatabaseUrl=" + susqlPrometheusDatabaseUrl)
	susqlLog.Info("samplingRate=" + samplingRate)
	susqlLog.Info("susqlLogLevel=" + susqlLogLevel)
	susqlLog.Info("carbonMethod=" + carbonMethod)
	susqlLog.Info("carbonIntensity=" + carbonIntensity)
	susqlLog.Info("carbonIntensityUrl=" + carbonIntensityUrl)
	susqlLog.Info("carbonLocation=" + carbonLocation)
	susqlLog.Info("carbonQueryRate=" + carbonQueryRate)
	susqlLog.Info("carbonQueryFilter=" + carbonQueryFilter)

	// if the enable-http2 flag is false (the default), http/2 should be disabled
	// due to its vulnerabilities. More specifically, disabling http/2 will
	// prevent from being vulnerable to the HTTP/2 Stream Cancellation and
	// Rapid Reset CVEs. For more information see:
	// - https://github.com/advisories/GHSA-qppj-fm5r-hxr3
	// - https://github.com/advisories/GHSA-4374-p667-p6c8
	disableHTTP2 := func(c *tls.Config) {
		susqlLog.Info("disabling http/2")
		c.NextProtos = []string{"http/1.1"}
	}

	tlsOpts := []func(*tls.Config){}
	tlsOpts = append(tlsOpts, disableHTTP2)
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress:   "0", // was: tunable metricsAddr
			SecureServing: false,
			TLSOpts:       tlsOpts,
		},
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "cac735ee.ibm.com",
		PprofBindAddress:       "127.0.0.1:6060",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		susqlLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// TODO: verify that carbonMethod is an expected value. If not log warning and set to default value.
	// (static, simpledynamic, sdk)
	// Note: assume that carbonIntensityUrl, carbonLocation, and carbonQueryFilter are OK. If not, we will log errors at runtime.

	samplingRateInteger, err := strconv.Atoi(samplingRate)
	if err != nil {
		samplingRateInteger = 2
	}

	carbonQueryRateInteger, err := strconv.Atoi(carbonQueryRate)
	if err != nil {
		carbonQueryRateInteger = 60
	}

	carbonIntensityFloat, err := strconv.ParseFloat(carbonIntensity, 64)
	if err != nil {
		susqlLog.Error(err, "Unable to obtain initial carbon intensity value. Using 0.0.")
		carbonIntensityFloat = 0.0
	}

	susqlLog.Info("Setting up labelGroupReconciler.")

	if err = (&controller.LabelGroupReconciler{
		Client:                     mgr.GetClient(),
		Scheme:                     mgr.GetScheme(),
		KeplerPrometheusUrl:        keplerPrometheusUrl,
		KeplerMetricName:           keplerMetricName,
		SusQLPrometheusDatabaseUrl: susqlPrometheusDatabaseUrl,
		SusQLPrometheusMetricsUrl:  susqlPrometheusMetricsUrl,
		SamplingRate:               time.Duration(samplingRateInteger) * time.Second,
		CarbonMethod:               carbonMethod,
		CarbonIntensity:            carbonIntensityFloat,
		CarbonIntensityUrl:         carbonIntensityUrl,
		CarbonLocation:             carbonLocation,
		CarbonQueryRate:            time.Duration(carbonQueryRateInteger) * time.Minute,
		CarbonQueryFilter:          carbonQueryFilter,
		Logger:                     susqlLog,
	}).SetupWithManager(mgr); err != nil {
		susqlLog.Error(err, "unable to create controller", "controller", "LabelGroup")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	susqlLog.Info("Adding healthz check.")

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		susqlLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		susqlLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	susqlLog.Info("Starting manager.")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		susqlLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
