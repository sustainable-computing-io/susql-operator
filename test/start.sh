oc apply -f labelgroups.yaml
sleep 10
oc apply -f energy-consumer-job.yaml
oc apply -f training-job-1.yaml
oc apply -f training-job-2.yaml
