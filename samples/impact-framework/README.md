# Using SusQL with the Green Software Foundation's Impact Framework

### Using prebuilt container image:
- Ensure that you are logged into your cluster, then start the Impact Framework container: `oc apply -f impact-framework.yaml --wait`
- Log into the newly created container: `oc rsh impact-framework bash`
- Clone the Impact Framework repository, install it, and run a sample manifest:
```
cd if
npm install
npm run if-run -- --manifest  manifests/examples/builtins/sum/success.yml
```

### Building your own container image
For those who wish to build their own Impact Framework, log into both your 
cluster and your image repository, and run the following commands:
```
export IMG=<YOURREPO>/<YOURID>/<YOURIMAGENAME
podman build --tag ${IMG} .
podman push ${IMG}
```
Then edit `impact-framework.yaml` so that `image` points to your newly built image.

Notes:
- `oc` and `kubectl` are (should be) interchangable in this context.
- `podman` and `docker` are (should be) interchangable in this context.
