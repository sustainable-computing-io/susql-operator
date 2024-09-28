# SusQL Operator Installation via Helm

Instructions to install SusQL via Helm.

**NOTE:** Your controller will automatically use the context in your kubeconfig file at the time of installation (i.e. whatever cluster `kubectl cluster-info` shows).

### Prerequisites

* Kepler is assumed to be installed in the cluster.
* `helm`, `kubectl`, and `go` are also required to deploy with Helm.
* Consider enabling User Project Monitoring when using on OpenShift.
    * https://docs.openshift.com/container-platform/latest/monitoring/enabling-monitoring-for-user-defined-projects.html

### Installation

To install SusQL via Helm go to the `deployment` directory and run the command `bash deploy.sh`. This script does a few actions:

* Check if Kepler is installed and exposing metrics through Prometheus
    * In general, Kepler metrics are exposed, clusterwide, at:

      ```
      http://<PROMETHEUS_SERVICE>.<PROMETHEUS_NAMESPACE>.<PROMETHEUS_DOMAIN>:9090
      ```

      The deployment script assumes `PROMETHEUS_SERVICE=prometheus-k8s`, `PROMETHEUS_NAMESPACE=monitoring`, and `PROMETHEUS_DOMAIN=svc.cluster.local`. If this is not the case, use the deployment script as:

      ```
      $ PROMETHEUS_SERVICE=<prometheus-service> PROMETHEUS_NAMESPACE=<prometheus-namespace> PROMETHEUS_DOMAIN=<prometheus-domain> bash deploy.sh
      ```

      Or for example, the following works for a 4.12 single node OpenShift local cluster:

      ```
      $ PROMETHEUS_PROTOCOL=https PROMETHEUS_PORT=9091 PROMETHEUS_NAMESPACE=openshift-monitoring PROMETHEUS_SERVICE=thanos-querier SUSQL_ENHANCED=true bash deploy.sh
      ```

* Create the namespace `openshift-kepler-operator`

* Install the SusQL operator in the namespace `openshift-kepler-operator`
    * This installation also deploys the CRD and sets the cluster permissions

* Install a separate Prometheus instance in the namespace `openshift-kepler-operator`

**NOTE**: This set of ***actions*** can be controlled by calling `bash deploy.sh susql-deploy`, for example, if only the SusQL deployment is needed. Check the script for all possible options or run it with the default set of actions.

### Installation Configuration Options

The following environment variables will influence the way that the SusQL Operator is installed via Helm.

| Environment Variable        | Default Value                 | Description                                      |
|:---------------------------:|:-----------------------------:|:------------------------------------------------:|
| SUSQL_NAMESPACE             | openshift-kepler-operator     | namespace that SusQL resources run in            |
| KEPLER_PROMETHEUS_NAMESPACE | openshift-monitoring          | namespace that Kepler Prometheus runs in         |
| PROMETHEUS_PROTOCOL         | http                          | Either http or https for Kepler Prometheus access|
| PROMETHEUS_SERVICE          | prometheus-k8s                | service name for the Kepler Prometheus           |
| PROMETHEUS_NAMESPACE        | monitoring                    | namespace used by the Kepler Prometheus          |
| PROMETHEUS_DOMAIN           | svc.cluster.local             | Domain used by the Kepler Prometheus             |
| PROMETHEUS_PORT             | 9090                          | Port used by the Kepler Prometheus               |
| KEPLER_PROMETHEUS_URL       | http://prometheus-k8s.monitoring.svc.cluster.local:9090 | A shortcut to specify final Kepler Prometheus URL |
| KEPLER_METRIC_NAME          | kepler_container_joules_total | Metric queried in the Kepler Prometheus          |
| SUSQL_PROMETHEUS_URL        | http://prometheus-susql.openshift-kepler-operator.svc.cluster.local:9090 | SusQL Prometheus URL |
| SUSQL_SAMPLING_RATE         | 2                             | Sampling rate in seconds                         |
| SUSQL_LOG_LEVEL             | -5                            | Log level                                        |
| SUSQL_ENHANCED              |                               | If set to any string, then use enhanced RBAC and SMON configuration |
| SUSQL_REGISTRY              | quay.io/sustainable_computing_io | Container registry that SusQL is stored in    |
| SUSQL_IMAGE_NAME            | susql_operator                | Image name used on SusQL container registry      |
| SUSQL_IMAGE_TAG             | latest                        | Tag for SusQL container                          |
| CARBON_METHOD               | static                        | "static", "simpledynamic", "scadk"               |
| CARBON_INTENSITY            | "0.00011583333"               | Carbon intensity in grams CO2 / Joule            |
| CARBON_INTENSITY_URL        |                               | Web API to query carbon intensity                |
| CARBON_LOCATION             |                               | Location for carbon intensity query              |
| CARBON_QUERY_RATE           | 7200                          | # of seconds between carbon intensity queries    |
| CARBON_QUERY_FILTER         | carbonIntensity               | JSON identifer of carbon intensity value         |
| CARBON_QUERY_CONV_2J        |                               | Convert to grams CO2/J                           |



## License

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

