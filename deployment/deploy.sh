#!/bin/bash

#
# Deploy SusQL controller using helm charts
#

# Get Kepler installation variables
if [[ -z ${KEPLER_NAMESPACE} ]]; then
    KEPLER_NAMESPACE="monitoring"
fi

if [[ -z ${KEPLER_SERVICE} ]]; then
    KEPLER_SERVICE="prometheus-k8s"
fi

KEPLER_PROMETHEUS_URL="http://${KEPLER_SERVICE}.${KEPLER_NAMESPACE}.svc.cluster.local:9090"

# Check if namespace exists
if [[ -z $(kubectl get namespaces --no-headers -o custom-columns=':{.metadata.name}' | grep susql) ]]; then
    echo "Namespace 'susql' doesn't exist. Creating it."
    kubectl create namespace susql
else
    echo "Namespace 'susql' found. Deploying using it."
fi

# Set SusQL installation variables
SUSQL_DIR=".."
SUSQL_REGISTRY="docker.io/metalcycling"
SUSQL_IMAGE_NAME="susql-controller"
SUSQL_IMAGE_TAG="latest"

# Actions to perform, separated by comma
actions=${1:-"kepler-check,prometheus-undeploy,prometheus-deploy,susql-undeploy,susql-deploy"}

for action in $(echo ${actions} | tr ',' '\n')
do
    if [[ ${action} = "kepler-check" ]]; then
        # Check if Kepler is serving metrics through prometheus
        echo "Checking if Kepler is deployed..."

        sed "s|KEPLER_PROMETHEUS_URL|${KEPLER_PROMETHEUS_URL}|g" kepler-check.yaml | kubectl apply -f -

        while true
        do
            phase=$(kubectl get pod kepler-check -o jsonpath='{.status.phase}')

            if [[ ${phase} != "Pending" ]] && [[ ${phase} != "Running" ]]; then
                break
            fi
        done

        kubectl delete -f kepler-check.yaml

        if [[ ${phase} == "Failed" ]]; then
            echo "Kepler service at '${KEPLER_PROMETHEUS_URL}' was not found. Check values and try again."
            exit 1
        else
            echo "Kepler service at '${KEPLER_PROMETHEUS_URL}' was found. Using it to collect SusQL data."
        fi

    elif [[ ${action} = "prometheus-deploy" ]]; then
        # Install prometheus from community helm charts
        helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
        helm repo update
        helm upgrade --install prometheus -f prometheus.yaml --namespace susql prometheus-community/prometheus

    elif [[ ${action} = "prometheus-undeploy" ]]; then
        helm uninstall prometheus --namespace susql

    elif [[ ${action} = "susql-undeploy" ]]; then
        cd ${SUSQL_DIR} && make uninstall && cd -
        helm -n susql uninstall susql-controller

    elif [[ ${action} = "susql-deploy" ]]; then
        cd ${SUSQL_DIR} && make manifests && make install && cd -
        helm upgrade --install --wait susql-controller ${SUSQL_DIR}/deployment/susql-controller --namespace susql \
            --set keplerPrometheusUrl="${KEPLER_PROMETHEUS_URL}" \
            --set susqlPrometheusDatabaseUrl="http://prometheus-susql.susql.svc.cluster.local:9090" \
            --set susqlPrometheusMetricsUrl="http://0.0.0.0:8082" \
            --set imagePullPolicy="IfNotPresent" \
            --set containerImage="${SUSQL_REGISTRY}/${SUSQL_IMAGE_NAME}:${SUSQL_IMAGE_TAG}"

    else
        echo "Nothing to do"

    fi
done
