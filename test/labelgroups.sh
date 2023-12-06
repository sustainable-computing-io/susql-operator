#!/bin/bash

namespace=default

t_start=$(date +%s.%N)

for labelgroup in $(kubectl -n ${namespace} get labelgroups -o custom-columns=':{.metadata.name}')
do
    totalEnergy=$(kubectl -n ${namespace} get labelgroup ${labelgroup} -o jsonpath='{.status.totalEnergy}')
    susqlPrometheusQuery=$(kubectl -n ${namespace} get labelgroup ${labelgroup} -o jsonpath='{.status.susqlPrometheusQuery}')
    phase=$(kubectl -n ${namespace} get labelgroup ${labelgroup} -o jsonpath='{.status.phase}')
    labels=$(kubectl -n ${namespace} get labelgroup ${labelgroup} -o jsonpath='{.spec.labels}')

    echo "LabelGroup: ${labelgroup}"
    echo "    - Labels: ${labels}"
    echo "    - Total Energy: ${totalEnergy}"
    echo "    - SusQL Query: ${susqlPrometheusQuery}"
    echo "    - Phase: ${phase}"
    echo 
done

t_stop=$(date +%s.%N)
t_total=$(echo "${t_stop} - ${t_start}" | bc -l)
echo "Total query time: ${t_total}"
