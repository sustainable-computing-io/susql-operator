# SusQL Operator Installation via OpenShift Community Operator Catalog

Installation of the SusQL Operator on OpenShift is very easy.

## Prerequisites

- [Kepler](https://sustainable-computing.io/) is assumed to be installed in the cluster. `oc` is assumed to be installed, and you are assumed to be logged in to the cluster in your CLI environment.
- Installation of the SusQL Operator has been extensively used with Red Hat OpenShift clusters 4.14 and higher.
- Consider enabling User Project Monitoring when using on OpenShift.  https://docs.openshift.com/container-platform/latest/monitoring/enabling-monitoring-for-user-defined-projects.html
- Other versions of OpenShift clusters are also expected to function.
- Ensure that the following command is run to install the service monitor. Replace`<NAMESPACE>` with the namespace SusQL is installed into. (Usually `openshift-operators`)
```
oc apply -n <NAMESPACE> -f oc apply -f https://raw.githubusercontent.com/sustainable-computing-io/susql-operator/main/hack/susql-servicemonitor.yaml
```

## Installation

Use the OpenShift web console to install the SusQL Operator:
- Click on "Operator->OperatorHub"
- Search for "SusQL"
- Click on it, and follow the GUI prompts to install.

# Customization

Before deploying the SusQL Operator create a `ConfigMap` called `susql-config` in
the same namespace that the operator will run in.
[susql-config.yaml](susql-config.yaml) is a good starting point. If you download it first, you
could create the ConfigMap with `oc apply -n <NAMESPACE> -f susql-config.yaml`.
If you update (or create) the ConfigMap after the SusQL Operator has been installed, then restarting the SusQL Operator controller pod will
enable the changes. (e.g., Delete the pod, and allow it to be recreated automatically.)

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
