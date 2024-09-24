#!/usr/bin/env bash

set -e
set -o pipefail

KEPLER_NS=kepler
PROMETHEUS_NS=monitoring

echo "> Creating Kind cluster"
echo "> Configure cluster to run Kepler"
echo "> Mount /proc - to expose information about processes running on the host"
echo "> Mount /usr/src - to expose kernel headers allowing dynamic eBPF program compilation"

export CLUSTER_NAME="local-cluster"

# Check if the cluster already exists
if kind get clusters | grep -q "$CLUSTER_NAME"; then
    echo "> Kind cluster $CLUSTER_NAME already exists"
else
    kind create cluster --name="$CLUSTER_NAME" --config=./deployment/local-cluster-config.yaml
fi

# Install Prometheus via Helm
echo "> Install Prometheus"
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus prometheus-community/kube-prometheus-stack \
    --namespace="$PROMETHEUS_NS" \
    --create-namespace \
    --wait

# Install Kepler via Helm
echo "> Install Kepler"
helm repo add kepler https://sustainable-computing-io.github.io/kepler-helm-chart
helm repo update
helm install kepler kepler/kepler \
    --namespace="$KEPLER_NS" \
    --create-namespace \
    --set serviceMonitor.enabled=true \
    --set serviceMonitor.labels.release=prometheus

# Get Kepler pod and wait for it to be ready
echo "> Wait for kepler pod to be ready"
KPLR_POD=$(kubectl get pod -l app.kubernetes.io/name=kepler -o jsonpath="{.items[0].metadata.name}" -n kepler)
kubectl wait --for=condition=Ready pod $KPLR_POD --timeout=-1s -n kepler


 # Add kepler dashboard to Grafana
echo "> Add kepler dashboard to Grafana"
GF_POD=$(
    kubectl get pod \
        -n monitoring \
        -l app.kubernetes.io/name=grafana \
        -o jsonpath="{.items[0].metadata.name}"
)
kubectl cp deployment/kepler_dashboard.json monitoring/$GF_POD:/tmp/dashboards/kepler_dashboard.json


# echo "> Install OLM"
# curl -sL https://github.com/operator-framework/operator-lifecycle-manager/releases/download/v0.28.0/install.sh | bash -s v0.28.0
# operator-sdk olm install


# echo "> Install SUSQL"
# kubectl create -f https://operatorhub.io/install/susql-operator.yaml

# echo "> Wait for susql to be ready"
# kubectl wait --for=condition=Installed csv/susql-operator.v0.0.30 -n operators --timeout=300s


# Optional: Delete Kind cluster after use
# kind delete cluster --name="$CLUSTER_NAME"

 echo "> Done"