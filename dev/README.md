# Random notes on developing for SusQL Operator

- Operator build has been verified on RHEL and Ubuntu.
- docker.io and quay.io have been used.
- Sample steps to build and push image.

```
export IMG=REGISTRYURL/REPOSITORYNAME/susql-controller:latest
podman login 
make build && make docker-build 
podman push ${IMG}
```


