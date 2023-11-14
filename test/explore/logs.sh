echo susql-controller logs
p=$(oc get po -n openshift-kepler-operator | grep susql-controller| cut -d' ' -f 1)
oc logs $p -n openshift-kepler-operator

echo
echo prometheus-susql logs
p=$(oc get po -n openshift-kepler-operator | grep prometheus-susql | cut -d' ' -f 1)
oc logs $p -n openshift-kepler-operator
