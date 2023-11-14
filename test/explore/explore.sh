oc apply -f explore.yaml
sleep 2

echo kepler metrics via thanos querier
oc exec -it curl-pod -n openshift-kepler-operator -- sh -c 'curl -k https://thanos-querier.openshift-monitoring.svc.cluster.local:9091/metrics -H "Authorization: Bearer $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)"' 2>&1 | tee keplerMetrics.txt
echo
echo "susql metrics directly from susql prometheus"
oc exec -it curl-pod -n openshift-kepler-operator -- sh -c 'curl -k http://prometheus-susql.openshift-kepler-operator.svc.cluster.local:9090/metrics' 2>&1| tee susqlMetrics.txt
oc delete -f explore.yaml &
