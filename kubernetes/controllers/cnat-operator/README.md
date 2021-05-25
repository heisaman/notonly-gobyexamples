
```shell
$ mkdir cnat-operator && cd cnat-operator
$ operator-sdk init --domain programming-kubernetes.info
$ operator-sdk create api --group cnat.programming-kubernetes.info  --version v1alpha1 --kind At  --resource --controller
```
