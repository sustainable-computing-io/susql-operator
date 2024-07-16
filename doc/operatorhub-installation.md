# SusQL Operator Installation via OperatorHub

Installation of the SusQL Operator via OperatorHub is fairly easy.

## Prerequisites

Kepler is assumed to be installed in the cluster. `oc` is assumed to be installed, and you are assumed to be logged in to the cluster in your CLI environment.
Installation of the SusQL Operator has been extensively used with Red Hat OpenShift clusters 4.14 and higher.
Other flavors and versions of Kubernetes clusters are also expected to function.

## Installation

Create a file called `catalogsource.yaml` as follows:

```
apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: operatorhubio-catalog
  namespace: openshift-marketplace
spec:
  sourceType: grpc
  image: quay.io/operatorhubio/catalog:latest
  displayName: Community Operators
  publisher: OperatorHub.io
  updateStrategy:
    registryPoll:
      interval: 10m
```

Apply this file with these commands:
```
oc apply -f catalogsource.yaml
oc patch -n openshift-marketplace OperatorHub cluster --type merge -p '{"spec":{"sources":[{"name":"operatorhubio-catalog","disabled":false},{"name":"community-operators","disabled":false}]}}'
```

Next, use the OpenShift web console to install the SusQL Operator:
- Click on "Operator->OperatorHub"
- Search for "SusQL"
- Click on it, and follow the GUI prompts to install.

# Post Install Customization

After installation, click "Operators->Installed Operators->SusQL->YAML", then edit one of the values under =spec.install.spec.deployments.labels.spec.template.spec.containers.resources.env.name.values` and save the YAML file when finished.
The operator will automatically restart with the newly specified value(s).

Currently the following variables can be modified. Searching directly for these variables may be more efficient than perusing the entire YAML.

- `KEPLER-PROMETHEUS-URL`
- `KEPLER-METRIC-NAME`
- `SUSQL-PROMETHEUS-DATABASE-URL`
- `SUSQL-PROMETHEUS-METRICS-URL`
- `SAMPLING-RATE`
- `LEADER-ELECT`
- `HEALTH-PROBE-BIND-ADDRESS`
- `METRICS-BIND-ADDRESS`

## License

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

