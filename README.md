# SusQL Operator

SusQL is a Kubernetes operator that aggregates energy and estimated carbon dioxide emission data for pods tagged with SusQL specific labels. The energy measurements are taken from [Kepler](https://sustainable-computing.io/) which should be installed/deployed in the cluster before using SusQL. Watch a video with a demonstration by clicking on the image bellow.

[![Watch the video](https://img.youtube.com/vi/NRVD7gJECfA/maxresdefault.jpg)](https://youtu.be/NRVD7gJECfA)

## Getting Started

SusQL is an operator that can be deployed in a Kubernetes/OpenShift cluster. You can use [kind](https://sigs.k8s.io/kind) or [minikube](https://minikube.sigs.k8s.io/) to get a local cluster for testing, or run against a remote cluster.

## Carbon Dioxide Emission Calculation

CO2 emission calculation is currently a work in progress. At the moment we use
a static conversion factor based on
[US EPA](https://www.epa.gov/energy/greenhouse-gases-equivalencies-calculator-calculations-and-references)
data. Although users can configure this value to match their runtime environment, we are working on providing
automatic up to date conversion functionality.

### Prerequisites

Kepler is assumed to be installed in the cluster.

### Installation

Follow these instructions to install the SusQL Operator with Helm.
- [Installation with Helm](doc/helm-installation.md)

Follow these instructions to install the SusQL Operator from [OperatorHub](https://operatorhub.io).
- [Installation with OperatorHub](doc/operatorhub-installation.md)


## Using SusQL

To begin using SusQL, a `LabelGroup` is used to specify the set of labels that the controller uses to identify pods that belong to the same energy aggregation. An example of a `LabelGroup` could be:

```
apiVersion: susql.ibm.com/v1
kind: LabelGroup
metadata:
    name: labelgroup-name
    namespace: default
spec:
    labels:
        - my-label-1
        - my-label-2
```

A pod that would be part of the ***group of pods*** belonging to the same energy aggregation would specify the `LabelGroup` labels as:

```
apiVersion: v1
kind: Pod
metadata:
    name: pod-name
    labels:
        susql.label/1: my-label-1
        susql.label/2: my-label-2
spec:
    containers:
        - name: container
          image: ubuntu
          command: ["sleep"]
          args: ["infinity"]
```

Energy of the group of pods is exposed in 2 ways:

* Through Prometheus at `http://prometheus-susql.openshift-kepler-operator.svc.cluster.local:9090` using the query `susql_total_energy_joules{susql_label_1=my-label-1,susql_label_2=my-label-2}`
* From `status` of the `LabelGroup` CRD given as `labelgroup.status.totalEnergy`

## Other Examples
- A step by step explanation of how to aggregate a [GPU based Jupyter Notebook workload on OpenShift AI](doc/openshift-ai-example-notebook.md).


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

