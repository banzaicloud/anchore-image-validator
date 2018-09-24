
# Anchore Image Validator

This admission-server that is used as a ValidatingWebhook in a k8s cluster. If it's working, kubernetes will send requests to the admission server when a Pod creation is initiated. The server checks the image, which is defined in PodSpec, against configured Anchore-engine API. If the API responds with an error, that the image is not valid according to defined policy, k8s will reject the Pod creation request.

- If an image is not valid, the release will be able to put into whitelist using CRD (whitelists).
- Every image check results are logged as a sructured format to specified CRD (audits).

### Accessing banzaicloud security via k8s api:

```shell
$ curl http://<k8s apiserver>/apis/security.banzaicloud.com/v1alpha1
```

```json
{
  "kind": "APIResourceList",
  "apiVersion": "v1",
  "groupVersion": "security.banzaicloud.com/v1alpha1",
  "resources": [
    {
      "name": "whitelists",
      "singularName": "whitelist",
      "namespaced": false,
      "kind": "WhiteList",
      "verbs": [ ... ],
      "shortNames": [
        "wl"
      ]
    },
    {
      "name": "audits",
      "singularName": "audit",
      "namespaced": false,
      "kind": "Audit",
      "verbs": [ ... ]
    }
  ]
}
```

#### Resources accessible via `kubectl` command:

```shell
$ kubectl get whitelist
$ kubectl get audit
```


#### Example whitelist:

```yaml
apiVersion: security.banzaicloud.com/v1alpha1
kind:  WhiteList
metadata:
  name: <name of whiltelist>
spec:
  releaseName: <helm release name>
  reason: <whitelisting reason>
  creator: <createor>
```

#### Examle audit:

```yaml
apiVersion: security.banzaicloud.com/v1alpha1
kind:  Audit
metadata:
  name: <name of audit (generated from Pod OwnerReference)>
  ownerReferences:
  - <scanned Pod OwnerReference>
spec:
  releaseName: <helm release name>
  resource: pod
  image:
    - <coinatinre image1>
    - <coinatinre image2>
  result:
    - <image1 scan result>
    - <image2 scan result>
  action: <allow or reject>
status:
  state: <optional>
```

### Some environment variables have to be defined:

|           ENV          |       Descripton      |
|------------------------|-----------------------|
|ANCHORE_ENGINE_USERNAME |Anchore-engine username|
|ANCHORE_ENGINE_PASSWORD |Anchore-engine password|
|ANCHORE_ENGINE_URL      |Anchore-engine URL     |
