# Using SusQL with the Green Software Foundation's Impact Framework

This particular experiment was performed using a recent RHEL 9 x86 machine,
however, any OS that supports a sufficiently recent version of `nodejs` should
work in principle.

The key to importing data from SusQL is a Prometheus Importer Plugin for the GSF Impact Framework.
Currently there are two such plugins described in the (IF Explorer)[https://explorer.if.greensoftware.foundation]:
- (`Prometheus Importer`)[https://github.com/andreic94/if-prometheus-importer/blob/main/README.md] by `andreic94`, et al.
- (`prometheus-importer`)[https://github.com/Shivani-G/prometheus-importer/blob/main/README.md] by `Shibani-G`.

The following is an approach using the later plugin:

### steps
- Ensure that you are logged into your cluster and can use the `oc` (or `kubectl` command).
- Make sure that a recent version of `node` is installed.  (This test used v22.9.0)
- Install Impact Framework
  - `npm install -g @grnsft/if`
- Install Prometheus Importer
  - `npm install -g "https://github.com/Shivani-G/prometheus-importer"`
- Update, just to be sure: (Starting from the directory that contains this README...)
  - `cd test; npm update; cd -`

- Create required credential file, and attempt to use prometheus-importer: (Starting from the directory that contains this README...)
```
cd test
echo "BEARER_TOKEN="$(oc whoami -t) >.env
echo HOST=https://$(oc get routes -n openshift-monitoring thanos-querier -o jsonpath='{.status.ingress[0].host}') >>.env

if-run --manifest  ../ptest.yaml --debug
cd -
```


Note:
- `oc` and `kubectl` are (should be) interchangable in this context.
