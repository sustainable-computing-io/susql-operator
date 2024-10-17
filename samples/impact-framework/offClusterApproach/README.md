# Using SusQL with the Green Software Foundation's Impact Framework

This particular experiment was performed using a recent RHEL 9 x86 machine,
however, any OS that supports a sufficiently recent version of nodejs should
work in principle.

### steps
- Ensure that you are logged into your cluster and can use the `oc` (or `kubectl` command).
- Make sure that a recent version of `node` is installed.  (This test used v22.9.0)
- Install Impact Framework
  - `npm install -g @grnsft/if`
- Install Prometheus Importer
  - `npm install -g "https://github.com/trent-s/prometheus-importer"`
- Update, just to be sure:
  - `cd test; npm update; cd -`

- Create required credential file, and attempt to use prometheus-importer:
```
cd test
echo "BEARER_TOKEN="$(oc whoami -t) >.env
echo HOST=https://$(oc get routes -n openshift-monitoring thanos-querier -o jsonpath='{.status.ingress[0].host}') >>.env

if-run --manifest  ../ptest.yaml --debug
cd -
```


Note:
- `oc` and `kubectl` are (should be) interchangable in this context.
