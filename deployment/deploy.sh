#!/bin/bash

#
# Deploy SusQL controller using helm charts
#

# Actions to perform, separated by comma
#actions="uninstall,install"
#actions="uninstall"
#actions="install"
#actions="prometheus-undeploy,prometheus-deploy"
#actions="susql-undeploy,susql-deploy"
#actions="susql-deploy"
actions="susql-undeploy"

# Check if namespace exists
if [[ -z $(kubectl get namespaces --no-headers -o custom-columns=':{.metadata.name}' | grep susql) ]]; then
    kubectl create namespace susql
fi

# Variables
SUSQL_DIR=".."
REGISTRY="quay.io/metalcycling"
IMAGE_NAME="susql-controller"
IMAGE_TAG="latest"

# Run actions
for action in $(echo ${actions} | tr ',' '\n')
do
    if [[ ${action} = "prometheus-deploy" ]]; then
        helm upgrade --install prometheus -f prometheus.yaml --namespace susql prometheus-community/prometheus

    elif [[ ${action} = "prometheus-undeploy" ]]; then
        helm uninstall prometheus --namespace susql

    elif [[ ${action} = "susql-undeploy" ]]; then
        #cd ${SUSQL_DIR} && make uninstall && cd -
        #kubectl delete -f ${SUSQL_DIR}/config/rbac/labelgroup_editor_role.yaml
        helm -n susql uninstall susql-controller

    elif [[ ${action} = "susql-deploy" ]]; then
        #cd ${SUSQL_DIR} && make install && cd -
        #kubectl apply -f ${SUSQL_DIR}/config/rbac/labelgroup_editor_role.yaml
        helm upgrade --install --wait susql-controller ${SUSQL_DIR}/deployment/susql-controller --namespace susql \
            --set imagePullPolicy="IfNotPresent"
        #helm template --debug ${SUSQL_DIR}/deployment/susql-controller --namespace susql

    else
        echo "Nothing to do"

    fi
done
