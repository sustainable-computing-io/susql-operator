#!/bin/bash

#
# Deploy SusQL controller using helm charts
#

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo "Error: kubectl is not installed. Please install kubectl to interact with Kubernetes clusters."
    exit 1
fi

if [[ ! -z ${SUSQL_ENHANCED} ]]; then
    echo "->Using enhanced SusQL configuration."
fi

# Check if there is a current context
current_context=$(kubectl config current-context)
if [[ -z $current_context ]]; then # exit
    echo "Error: You are not currently signed into any Kubernetes cluster."
    exit 1
else # continue
    echo "You are currently signed into the following Kubernetes cluster:"
    echo "$current_context"
fi

export SUSQL_NAMESPACE="openshift-kepler-operator"
export KEPLER_PROMETHEUS_NAMESPACE="openshift-monitoring"

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

if [[ -z ${KEPLER_METRIC_NAME} ]]; then
    KEPLER_METRIC_NAME="kepler_container_joules_total"
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

if [[ -z ${SUSQL_SAMPLING_RATE} ]]; then
    SUSQL_SAMPLING_RATE="2"
fi

if [[ -z ${SUSQL_LOG_LEVEL} ]]; then
    SUSQL_LOG_LEVEL="-5"
fi

if [[ -z ${CARBON_METHOD} ]]; then
    CARBON_METHOD="static"
fi

if [[ -z ${CARBON_INTENSITY} ]]; then
    CARBON_INTENSITY="0.0001158333333333"
fi

if [[ -z ${CARBON_INTENSITY_URL} ]]; then
    CARBON_INTENSITY_URL="https://api.electricitymap.org/v3/carbon-intensity/latest?zone=%s"
fi

if [[ -z ${CARBON_LOCATION} ]]; then
    CARBON_LOCATION="JP-TK"
fi

if [[ -z ${CARBON_QUERY_RATE} ]]; then
    CARBON_QUERY_RATE="7200"
fi

if [[ -z ${CARBON_QUERY_FILTER} ]]; then
    CARBON_QUERY_FILTER="carbonIntensity"
fi

if [[ -z ${CARBON_QAUERY_CONV_2J} ]]; then
    CARBON_QAUERY_CONV_2J="0.0000002777777778"
fi

# Check if namespace exists
if [[ -z $(kubectl get namespaces --no-headers -o custom-columns=':{.metadata.name}' | grep ${SUSQL_NAMESPACE}) ]]; then
    echo "Namespace '${SUSQL_NAMESPACE}' doesn't exist. Creating it."
    kubectl create namespace ${SUSQL_NAMESPACE}
else
    echo "Namespace '${SUSQL_NAMESPACE}' found. Using it."
fi

# Set SusQL installation variables
SUSQL_DIR=".."
if [[ -z ${SUSQL_REGISTRY} ]]; then
    SUSQL_REGISTRY="quay.io/sustainable_computing_io"
fi
if [[ -z ${SUSQL_IMAGE_NAME} ]]; then
    SUSQL_IMAGE_NAME="susql_operator"
fi
if [[ -z ${SUSQL_IMAGE_TAG} ]]; then
    SUSQL_IMAGE_TAG="latest"
fi

# display configured variables 
echo "==================================================================================================="
echo "SUSQL_NAMESPACE - '${SUSQL_NAMESPACE}'"
echo "KEPLER_PROMETHEUS_NAMESPACE - '${KEPLER_PROMETHEUS_NAMESPACE}'"
echo "PROMETHEUS_PROTOCOL - '${PROMETHEUS_PROTOCOL}'"
echo "PROMETHEUS_SERVICE - '${PROMETHEUS_SERVICE}'"
echo "PROMETHEUS_NAMESPACE - '${PROMETHEUS_NAMESPACE}'"
echo "PROMETHEUS_DOMAIN - '${PROMETHEUS_DOMAIN}'"
echo "PROMETHEUS_PORT - '${PROMETHEUS_PORT}'"
echo "KEPLER_PROMETHEUS_URL - '${KEPLER_PROMETHEUS_URL}'"
echo "KEPLER_METRIC_NAME - '${KEPLER_METRIC_NAME}'"
echo "CARBON_METHOD - '${CARBON_METHOD}'"
echo "CARBON_INTENSITY - '${CARBON_INTENSITY}'"
echo "CARBON_INTENSITY_URL - '${CARBON_INTENSITY_URL}'"
echo "CARBON_LOCATION - '${CARBON_LOCATION}'"
echo "CARBON_QUERY_RATE - '${CARBON_QUERY_RATE}'"
echo "CARBON_QUERY_FILTER - '${CARBON_QUERY_FILTER}'"
echo "CARBON_QUERY_CONV_2J - '${CARBON_QUERY_CONV_2J}'"
echo "SUSQL_PROMETHEUS_URL - '${SUSQL_PROMETHEUS_URL}'"
echo "SUSQL_SAMPLING_RATE - '${SUSQL_SAMPLING_RATE}'"
echo "SUSQL_ENHANCED - '${SUSQL_ENHANCED}'"
echo "SUSQL_REGISTRY - '${SUSQL_REGISTRY}'"
echo "SUSQL_IMAGE_NAME - '${SUSQL_IMAGE_NAME}'"
echo "SUSQL_IMAGE_TAG - '${SUSQL_IMAGE_TAG}'"
echo "SUSQL_LOG_LEVEL - '${SUSQL_LOG_LEVEL}'"
echo "==================================================================================================="
# Actions to perform, separated by comma
actions=${1:-"kepler-check,prometheus-undeploy,prometheus-deploy,susql-undeploy,susql-deploy"}

# output deploy information
LOGFILE=.susql-deploy-info.txt
LASTLOGFILE=.susql-deploy-info-last.txt
if [[ -f ${LOGFILE} ]]
then
    mv -f ${LOGFILE} ${LASTLOGFILE}
fi

echo "# SusQL Deployment on "$(hostname)" at "$(date) > ${LOGFILE}
echo "# deploy action was: '${action}'" >> ${LOGFILE}
echo "export SUSQL_NAMESPACE=${SUSQL_NAMESPACE}" >> ${LOGFILE}
echo "export KEPLER_PROMETHEUS_NAMESPACE=${KEPLER_PROMETHEUS_NAMESPACE}" >> ${LOGFILE}
echo "export PROMETHEUS_PROTOCOL=${PROMETHEUS_PROTOCOL}" >> ${LOGFILE}
echo "export PROMETHEUS_SERVICE=${PROMETHEUS_SERVICE}" >> ${LOGFILE}
echo "export PROMETHEUS_NAMESPACE=${PROMETHEUS_NAMESPACE}" >> ${LOGFILE}
echo "export PROMETHEUS_DOMAIN=${PROMETHEUS_DOMAIN}" >> ${LOGFILE}
echo "export PROMETHEUS_PORT=${PROMETHEUS_PORT}" >> ${LOGFILE}
echo "export KEPLER_PROMETHEUS_URL=${KEPLER_PROMETHEUS_URL}" >> ${LOGFILE}
echo "export KEPLER_METRIC_NAME=${KEPLER_METRIC_NAME}" >> ${LOGFILE}
echo "export CARBON_METHOD=${CARBON_METHOD}" >> ${LOGFILE}
echo "export CARBON_INTENSITY=${CARBON_INTENSITY}" >> ${LOGFILE}
echo "export CARBON_INTENSITY_URL=${CARBON_INTENSITY_URL}" >> ${LOGFILE}
echo "export CARBON_LOCATION=${CARBON_LOCATION}" >> ${LOGFILE}
echo "export CARBON_QUERY_RATE=${CARBON_QUERY_RATE}" >> ${LOGFILE}
echo "export CARBON_QUERY_FILTER=${CARBON_QUERY_FILTER}" >> ${LOGFILE}
echo "export CARBON_QUERY_CONV_2J=${CARBON_QUERY_CONV_2J}" >> ${LOGFILE}
echo "export SUSQL_PROMETHEUS_URL=${SUSQL_PROMETHEUS_URL}" >> ${LOGFILE}
echo "export SUSQL_SAMPLING_RATE=${SUSQL_SAMPLING_RATE}" >> ${LOGFILE}
echo "export SUSQL_ENHANCED=${SUSQL_ENHANCED}" >> ${LOGFILE}
echo "export SUSQL_REGISTRY=${SUSQL_REGISTRY}" >> ${LOGFILE}
echo "export SUSQL_IMAGE_NAME=${SUSQL_IMAGE_NAME}" >> ${LOGFILE}
echo "export SUSQL_IMAGE_TAG=${SUSQL_IMAGE_TAG}" >> ${LOGFILE}
echo "export SUSQL_LOG_LEVEL=${SUSQL_LOG_LEVEL}" >> ${LOGFILE}

for action in $(echo ${actions} | tr ',' '\n')
do
    if [[ ${action} = "kepler-check" ]]; then
        # Check if Kepler is serving metrics through prometheus
        echo "->Checking if Kepler is deployed..."

        sed "s|KEPLER_PROMETHEUS_URL|${KEPLER_PROMETHEUS_URL}|g" kepler-check.yaml | kubectl apply --namespace ${KEPLER_PROMETHEUS_NAMESPACE} -f -

        echo "->Checking kepler-check status";
        while true
        do
            echo -n "." && sleep 1;
            phase=$(kubectl get pod kepler-check -o jsonpath='{.status.phase}' --namespace ${KEPLER_PROMETHEUS_NAMESPACE})

            if [[ ${phase} != "Pending" ]] && [[ ${phase} != "Running" ]]; then
                break
            fi
        done

        echo ""
        echo "Kepler check '${phase}'"

        logs=$(kubectl logs -f kepler-check --namespace ${KEPLER_PROMETHEUS_NAMESPACE})
        echo "->Deleting kepler-check pod..."
        kubectl delete -f kepler-check.yaml --namespace ${KEPLER_PROMETHEUS_NAMESPACE}

        echo "->Kepler service"
        if [[ ${phase} == "Failed" ]]; then
            echo "-----------------"
            echo "Kepler check logs"
            echo "-----------------"
            echo ${logs}
            echo "Kepler service at '${KEPLER_PROMETHEUS_URL}' was not found. Check values and try again."
            exit 1
        else
            echo "Kepler service at '${KEPLER_PROMETHEUS_URL}' was found. Using it to collect SusQL data."
        fi

    elif [[ ${action} = "prometheus-deploy" ]]; then
        # Install prometheus from community helm charts
        echo "->Deploying Prometheus controller to store SusQL data..."

        helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
        helm repo update
        helm upgrade --install prometheus -f ${PROMETHEUS_YAML} --namespace ${SUSQL_NAMESPACE} prometheus-community/prometheus

    elif [[ ${action} = "prometheus-undeploy" ]]; then
        echo "->Undeploying Prometheus controller..."
        echo "All data will be lost if the Prometheus deployment was not using a volume for permanent storage."
        echo -n "Would you like to procceed? [y/n]: "
        read response

        if [[ ${response} == "Y" || ${response} == "y" ]]; then
            output=$(helm uninstall prometheus --namespace ${SUSQL_NAMESPACE} 2>&1)
            if [ $? -ne 0 ]; then
                # Extract the error message
                error_message=$(echo "$output" | tail -n 1)

                # Check if the error message contains "not found"
                if [[ "$error_message" == *"not found"* ]]; then
                    # Custom message if release not found
                    echo "Prometheus: release: not found"
                else
                    # Display the original error message
                    echo "$error_message"
                fi
            else
                # Success message if the command was successful
                echo "Prometheus uninstall successful."
            fi
        fi

    elif [[ ${action} = "susql-deploy" ]]; then
        echo "->Deploying SusQL controller..."
        if [[ ! -z ${SUSQL_ENHANCED} ]]; then
            kubectl apply -f ../config/rbac/susql-rbac.yaml
        fi

        cd ${SUSQL_DIR} && make manifests && make install
        cd -
        helm upgrade --install --wait susql-controller ${SUSQL_DIR}/deployment/susql-controller --namespace ${SUSQL_NAMESPACE} \
            --set keplerPrometheusUrl="${KEPLER_PROMETHEUS_URL}" \
            --set keplerMetricName="${KEPLER_METRIC_NAME}" \
            --set susqlPrometheusDatabaseUrl="${SUSQL_PROMETHEUS_URL}" \
            --set susqlPrometheusMetricsUrl="http://0.0.0.0:8082" \
            --set CarbonMethod="${CARBON_METHOD}" \
            --set CarbonIntensity="${CARBON_INTENSITY}" \
            --set CarbonIntensityUrl="${CARBON_INTENSITY_URL}" \
            --set CarbonLocation="${CARBON_LOCATION}" \
            --set CarbonQueryRate="${CARBON_QUERY_RATE}" \
            --set CarbonQueryFilter="${CARBON_QUERY_FILTER}" \
            --set CarbonQueryConv2J="${CARBON_QUERY_CONV_2J}" \
            --set samplingRate="${SUSQL_SAMPLING_RATE}" \
            --set susqlLogLevel="${SUSQL_LOG_LEVEL}" \
            --set imagePullPolicy="Always" \
            --set containerImage="${SUSQL_REGISTRY}/${SUSQL_IMAGE_NAME}:${SUSQL_IMAGE_TAG}"
        if [[ ! -z ${SUSQL_ENHANCED} ]]; then
            kubectl apply -f ../config/servicemonitor/susql-smon.yaml
        fi

    elif [[ ${action} = "susql-undeploy" ]]; then
        echo "->Undeploying SusQL controller..."
        echo -n "Would you like to uninstall the CRD? WARNING: Doing this would remove ALL deployments for this CRD already in the cluster. [y/n]: "
        read response

        if [[ ${response} == "Y" || ${response} == "y" ]]; then
            output=$(cd ${SUSQL_DIR} && make uninstall 2>&1)
            if [ $? -ne 0 ]; then
                # Extract the error message
                error_message=$(echo "$output" | tail -n 2)

                # Check if the error message contains "not found"
                if [[ "$error_message" == *"not found"* ]]; then
                    # Custom message if release not found
                    echo "SusQL Controller not found"
                else
                    # Display the original error message
                    echo "$error_message"
                fi
            else
                # Success message if the command was successful
                echo "SusQL uninstall successful."
                helm -n ${SUSQL_NAMESPACE} uninstall susql-controller
            fi
        fi

        if [[ ! -z ${SUSQL_ENHANCED} ]]; then
            kubectl delete -f ../config/rbac/susql-rbac.yaml
            kubectl delete -f ../config/servicemonitor/susql-smon.yaml
        fi

    else
        echo "Nothing to do"

    fi
done
