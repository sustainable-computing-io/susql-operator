#!/bin/bash

#
# Deploy SusQL controller using helm charts
#

export SUSQL_NAMESPACE="openshift-kepler-operator"

if [[ -z ${PROMETHEUS_PROTOCOL} ]]; then
    PROMETHEUS_PROTOCOL="http"
fi

if [[ -z ${PROMETHEUS_PORT} ]]; then
    PROMETHEUS_PORT="9090"
fi

# Get Kepler metrics variables
if [[ -z ${PROMETHEUS_NAMESPACE} ]]; then
    PROMETHEUS_NAMESPACE="monitoring"
fi

if [[ -z ${PROMETHEUS_SERVICE} ]]; then
    PROMETHEUS_SERVICE="prometheus-k8s"
fi

if [[ -z ${PROMETHEUS_DOMAIN} ]]; then
    PROMETHEUS_DOMAIN="svc.cluster.local"
    PROMETHEUS_YAML="prometheus.yaml"
else
    PROMETHEUS_YAML="temp-prometheus.yaml"
    sed -e 's/svc.cluster.local/'${PROMETHEUS_DOMAIN}'/g' prometheus.yaml >${PROMETHEUS_YAML}
    sed -i -e 's/svc.cluster.local/'${PROMETHEUS_DOMAIN}'/g' susql-controller/values.yaml
fi

if [[ -z ${KEPLER_PROMETHEUS_URL} ]]; then
    KEPLER_PROMETHEUS_URL="${PROMETHEUS_PROTOCOL}://${PROMETHEUS_SERVICE}.${PROMETHEUS_NAMESPACE}.${PROMETHEUS_DOMAIN}:${PROMETHEUS_PORT}"
fi

if [[ -z ${SUSQL_PROMETHEUS_URL} ]]; then
    if [[ -z ${SHARED_PROMETHEUS} ]]; then
        # using separate prometheus instance
        SUSQL_PROMETHEUS_URL="http://prometheus-susql.${SUSQL_NAMESPACE}.${PROMETHEUS_DOMAIN}:9090"
    else
        # using shared prometheus instance
        SUSQL_PROMETHEUS_URL=${KEPLER_PROMETHEUS_URL}
    fi
fi

# Check if namespace exists
if [[ -z $(kubectl get namespaces --no-headers -o custom-columns=':{.metadata.name}' | grep ${SUSQL_NAMESPACE}) ]]; then
    echo "Namespace '${SUSQL_NAMESPACE}' doesn't exist. Creating it."
    kubectl create namespace ${SUSQL_NAMESPACE}
else
    echo "Namespace '${SUSQL_NAMESPACE}' found. Deploying using it."
fi

# Set SusQL installation variables
SUSQL_DIR=".."
SUSQL_REGISTRY="quay.io/sustainable_computing_io"
SUSQL_IMAGE_NAME="susql_operator"
SUSQL_IMAGE_TAG="latest"

# Actions to perform, separated by comma
actions=${1:-"kepler-check,prometheus-undeploy,prometheus-deploy,susql-undeploy,susql-deploy"}

for action in $(echo ${actions} | tr ',' '\n')
do
    if [[ ${action} = "kepler-check" ]]; then
        # Check if Kepler is serving metrics through prometheus
        echo "Checking if Kepler is deployed..."

        sed "s|KEPLER_PROMETHEUS_URL|${KEPLER_PROMETHEUS_URL}|g" kepler-check.yaml | kubectl apply --namespace ${SUSQL_NAMESPACE} -f -

        while true
        do
            phase=$(kubectl get pod kepler-check -o jsonpath='{.status.phase}' --namespace ${SUSQL_NAMESPACE})

            if [[ ${phase} != "Pending" ]] && [[ ${phase} != "Running" ]]; then
                break
            fi
        done

        logs=$(kubectl logs -f kepler-check --namespace ${SUSQL_NAMESPACE})
        kubectl delete -f kepler-check.yaml --namespace ${SUSQL_NAMESPACE}

        if [[ ${phase} == "Failed" ]]; then
            echo ${logs}
            echo "Kepler service at '${KEPLER_PROMETHEUS_URL}' was not found. Check values and try again."
            exit 1
        else
            echo "Kepler service at '${KEPLER_PROMETHEUS_URL}' was found. Using it to collect SusQL data."
        fi

    elif [[ ${action} = "prometheus-deploy" ]]; then
        # Install prometheus from community helm charts
        echo "Deploying Prometheus controller to store susql data..."

        helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
        helm repo update
        helm upgrade --install prometheus -f ${PROMETHEUS_YAML} --namespace ${SUSQL_NAMESPACE} prometheus-community/prometheus

    elif [[ ${action} = "prometheus-undeploy" ]]; then
        echo "Undeploying Prometheus controller..."
        echo "All data will be lost if the Prometheus deployment was not using a volume for permanent storage."
        echo -n "Would you like to procceed? [y/n]: "
        read response

        if [[ ${response} == "Y" || ${response} == "y" ]]; then
            helm uninstall prometheus --namespace ${SUSQL_NAMESPACE}
        fi

    elif [[ ${action} = "susql-deploy" ]]; then
        kubectl apply -f ../config/rbac/susql-rbac.yaml

        cd ${SUSQL_DIR} && make manifests && make install
        cd -
        helm upgrade --install --wait susql-controller ${SUSQL_DIR}/deployment/susql-controller --namespace ${SUSQL_NAMESPACE} \
            --set keplerPrometheusUrl="${KEPLER_PROMETHEUS_URL}" \
            --set susqlPrometheusDatabaseUrl="${SUSQL_PROMETHEUS_URL}" \
            --set susqlPrometheusMetricsUrl="http://0.0.0.0:8082" \
            --set imagePullPolicy="Always" \
            --set containerImage="${SUSQL_REGISTRY}/${SUSQL_IMAGE_NAME}:${SUSQL_IMAGE_TAG}"
        kubectl apply -f ../config/servicemonitor/susql-smon.yaml

    elif [[ ${action} = "susql-undeploy" ]]; then
        echo "Undeploying SusQL controller..."
        echo -n "Would you like to uninstall the CRD? WARNING: Doing this would remove ALL deployments for this CRD already in the cluster. [y/n]: "
        read response

        if [[ ${response} == "Y" || ${response} == "y" ]]; then
            cd ${SUSQL_DIR} && make uninstall
            cd -
        fi

        helm -n ${SUSQL_NAMESPACE} uninstall susql-controller
        kubectl delete -f ../config/rbac/susql-rbac.yaml
        kubectl delete -f ../config/servicemonitor/susql-smon.yaml

    else
        echo "Nothing to do"

    fi
done
