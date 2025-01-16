# SusQL Operator

SusQL is a Kubernetes operator that aggregates energy and estimated carbon dioxide emission data for pods tagged with SusQL specific labels. The energy measurements are obtained from [Kepler](https://sustainable-computing.io/) which should be deployed on the cluster before using SusQL. Click the picture below to watch the demo video.

[![SusQL Demo](https://github.com/sustainable-computing-io/susql-operator/wiki/files/SusQL-Demo-2024-10-Thumbnail.png)](https://youtu.be/9CwuhOfVtjE)

## Getting Started

SusQL is an operator that can be deployed in a Kubernetes/OpenShift cluster. You can also use [kind](https://sigs.k8s.io/kind) or [minikube](https://minikube.sigs.k8s.io/) as a local cluster for testing, or run against a remote cluster.

## Carbon Dioxide Emission Calculation

By default SusQL calculates carbon dioxide emission in grams of CO2 using a carbon intensity value from 
[US EPA](https://www.epa.gov/energy/greenhouse-gases-equivalencies-calculator-calculations-and-references).
SusQL can be configured to use other static carbon intensity values or query carbon intensity values for a
given location from web API's such as those provided by
the Green Software Foundation's [Carbon Aware SDK](https://github.com/Green-Software-Foundation/carbon-aware-sdk).

Detailed information on configuration of CO2 emission calculation in SusQL is available in the [SusQL carbon
calculation documentation.](doc/carbon.md)

## Prerequisites

Kepler is assumed to be installed in the cluster.

## Installation

- Follow these instructions for easy SusQL installation from the Red Hat Community Operator catalog on an OpenShift cluster.
  - [Installation on OpenShift](doc/openshift-installation.md)

- Follow these instructions to install the SusQL Operator from [OperatorHub.io](https://operatorhub.io) on a Kubernetes cluster including OpenShift.
  - [Installation from OperatorHub.io](doc/operatorhub-installation.md)

- Follow these instructions to install the SusQL Operator from a Helm chart on a Kubernetes cluster, including OpenShift.
  - [Installation with Helm](doc/helm-installation.md)
 

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

Energy of the group of pods is exposed in two ways:

* Through Prometheus at `http://prometheus-susql.openshift-kepler-operator.svc.cluster.local:9090` using the query `susql_total_energy_joules{susql_label_1=my-label-1,susql_label_2=my-label-2}`
* From `status` of the `LabelGroup` CRD given as `labelgroup.status.totalEnergy`

## Other Examples
- A step by step explanation of how to aggregate a [GPU based Jupyter Notebook workload on OpenShift AI](doc/openshift-ai-example-notebook.md).


## License

Copyright 2023, 2024, 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

