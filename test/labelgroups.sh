#!/bin/bash

namespace=default

t_start=$(date +%s.%N)

alldata=$(kubectl -n ${namespace} get labelgroups -o json)

for labelgroup in $(echo ${alldata} | jq -cr '.items[].metadata.name')
do
    newdata=$(echo ${alldata} | jq '.items[] | select(.metadata.name=="'${labelgroup}'")')
    totalEnergy=$(echo ${newdata} | jq -cr '.status.totalEnergy')
    totalCarbon=$(echo ${newdata} | jq -cr '.status.totalCarbon')
    susqlPrometheusQuery=$(echo ${newdata} | jq -cr '.status.susqlPrometheusQuery')
    phase=$(echo ${newdata} | jq -cr '.status.phase')
    labels=$(echo ${newdata} | jq -cr '.spec.labels')

    echo "LabelGroup: ${labelgroup}"
    echo "    - Labels: ${labels}"
    echo "    - Total Energy (J): ${totalEnergy}"
    echo "    - Total CO2 (g): ${totalCarbon}"
    echo "    - SusQL Query: ${susqlPrometheusQuery}"
    echo "    - Phase: ${phase}"
    echo 
done

t_stop=$(date +%s.%N)
t_total=$(echo "${t_stop} - ${t_start}" | bc -l)
echo "Total query time: ${t_total}"
