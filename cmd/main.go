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

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var keplerPrometheusUrl string
	var keplerMetricName string
	var susqlPrometheusMetricsUrl string
	var susqlPrometheusDatabaseUrl string
	var samplingRate string

	// NOTE: these can be set as env or flag, flag takes precedence over env
	keplerPrometheusUrlEnv := os.Getenv("KEPLER-PROMETHEUS-URL")
	keplerMetricNameEnv := os.Getenv("KEPLER-METRIC-NAME")
	susqlPrometheusDatabaseUrlEnv := os.Getenv("SUSQL-PROMETHEUS-DATABASE-URL")
	susqlPrometheusMetricsUrlEnv := os.Getenv("SUSQL-PROMETHEUS-METRICS-URL")
	samplingRateEnv := os.Getenv("SAMPLING-RATE")
	metricsAddrEnv := os.Getenv("METRICS-BIND-ADDRESS")
	probeAddrEnv := os.Getenv("HEALTH-PROBE-BIND-ADDRESS")
	enableLeaderElectionEnv, err := strconv.ParseBool(os.Getenv("LEADER-ELECT"))
	if err != nil {
		enableLeaderElectionEnv = false
	}

	flag.StringVar(&keplerPrometheusUrl, "kepler-prometheus-url", keplerPrometheusUrlEnv, "The URL for the Prometheus server where Kepler stores the energy data")
	flag.StringVar(&keplerMetricName, "kepler-metric-name", keplerMetricNameEnv, "The metric name to be queried in the kepler Prometheus server")
	flag.StringVar(&susqlPrometheusDatabaseUrl, "susql-prometheus-database-url", susqlPrometheusDatabaseUrlEnv, "The URL for the Prometheus database where SusQL stores the energy data")
	flag.StringVar(&susqlPrometheusMetricsUrl, "susql-prometheus-metrics-url", susqlPrometheusMetricsUrlEnv, "The URL for the Prometheus metrics where SusQL exposes the energy data")
	flag.StringVar(&samplingRate, "sampling-rate", samplingRateEnv, "Sampling rate in seconds")
	flag.StringVar(&metricsAddr, "metrics-bind-address", metricsAddrEnv, "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", probeAddrEnv, "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", enableLeaderElectionEnv,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	opts := zap.Options{
		Development: true,
		Level:       zapcore.Level(-5),
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	susqlLog.Info("SusQL configuration values at runtime")
	susqlLog.Info("metricsAddr=" + metricsAddr)
	susqlLog.Info("enableLeaderElection=" + strconv.FormatBool(enableLeaderElection))
	susqlLog.Info("probeAddr=" + probeAddr)
	susqlLog.Info("keplerPrometheusUrl=" + keplerPrometheusUrl)
	susqlLog.Info("keplerMetricName=" + keplerMetricName)
	susqlLog.Info("susqlPrometheusMetricsUrl=" + susqlPrometheusMetricsUrl)
	susqlLog.Info("susqlPrometheusDatabaseUrl=" + susqlPrometheusDatabaseUrl)
	susqlLog.Info("samplingRate=" + samplingRate)

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "cac735ee.ibm.com",
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

	samplingRateInteger, err := strconv.Atoi(samplingRate)
	if err != nil {
		samplingRateInteger = 2
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
