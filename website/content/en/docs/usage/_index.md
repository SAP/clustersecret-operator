---
title: "Usage"
linkTitle: "Usage"
weight: 30
type: "docs"
description: >
  How to use ClusterSecrets
---

A typical ClusterSecret resource looks like this:

```yaml
apiVersion: core.cs.sap.com/v1alpha1
kind: ClusterSecret
metadata:
  name: my-secret
spec:
  namespaceSelector:
    matchLabels:
      mylabel: myvalue
  template:
    type: Opaque
    data:
      mykey: bXl2YWx1ZQ==
```

The ClusterSecret `spec` consists of two parts:
- `spec.namespaceSeletor` follows the [usual syntax](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#resources-that-support-set-based-requirements)
- `spec.template` mirrors the usual secret spec, at least partially, allowing to specify `type` (mandatory), and at least one of `data` or `stringData`; if `stringData` is provided, it will be rewritten to `data` by the mutating admission webhook.

The controller will then ensure that an according secret (having the same name as the ClusterSecret) exists in all selected namespaces; in addition to ClusterSecret resources, the controller watches namespaces, and immediately reacts to creation of namespaces, or label changes.

