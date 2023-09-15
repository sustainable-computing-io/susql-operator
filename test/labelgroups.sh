#!/bin/bash

t_start=$(date +%s.%N)

for labelgroup in $(kubectl get labelgroups -o custom-columns=':{.metadata.name}')
do
    totalEnergy=$(kubectl get labelgroup ${labelgroup} -o jsonpath='{.status.totalEnergy}')
    susqlPrometheusQuery=$(kubectl get labelgroup ${labelgroup} -o jsonpath='{.status.susqlPrometheusQuery}')
    phase=$(kubectl get labelgroup ${labelgroup} -o jsonpath='{.status.phase}')
    labels=$(kubectl get labelgroup ${labelgroup} -o jsonpath='{.spec.labels}')

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
