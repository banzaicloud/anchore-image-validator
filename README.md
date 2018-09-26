
# Anchore Image Validator

Anchore Image Validator lets you automatically detect or block security issues just before a Kubernetes pod starts. 

This repository contains an [admission webhook](https://banzaicloud.com/blog/k8s-admission-webhooks/) server that can be configured as a ValidatingWebhook in a k8s cluster. Kubernetes will send requests to the admission server when a Pod creation is initiated. The server checks the image defined in the pod specification using the configured Anchore-engine API. If the result indicates that the image does not comply with the defined policy, k8s will reject the Pod creation request.

- If an image is not valid, the release can be added to a *whitelist* resource (CRD) to bypass the blocking.
- The results of image checks are stored as an *audit* resource (CRD) in a sructured format.

### Accessing banzaicloud security features via k8s api:

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
  name: <name of whitelist>
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
