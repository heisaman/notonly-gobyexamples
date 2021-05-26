
```shell
$ mkdir cnat-operator && cd cnat-operator
$ operator-sdk init --domain programming-kubernetes.info
$ operator-sdk create api --group cnat  --version v1alpha1 --kind At  --resource --controller
$ make generate
$ make manifests
```