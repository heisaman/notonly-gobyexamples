
```shell
$ mkdir cnat-operator && cd cnat-operator
$ operator-sdk init --domain programming-kubernetes.info
$ operator-sdk create api --group cnat  --version v1alpha1 --kind At  --resource --controller
$ make generate
$ make manifests  // create crd manifest
$ kubectl apply -f config/crd/bases/cnat.programming-kubernetes.info_ats.yaml  // apply crd manifest
$ make install run  // start the operator
```