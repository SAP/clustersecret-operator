---
title: "clustersecret-operator"
linkTitle: "clustersecret-operator"
weight: 10
type: "docs"
---

clustersecret-operator adds a new resource type `clustersecrets.core.cs.sap.com`, with kind `ClusterSecret`, to Kubernetes clusters.
It allows to define secrets at cluster scope, along with an optional selector defining in which
namespaces the according Kubernetes secrets shall exist. The controlller provided by this repository
takes care of distributing the secrets, and keeping everything in sync.

This website provides the full technical documentation for the project, and can be
used as a reference; if you feel that there's anything missing, please let us know
or [raise a PR](https://github.com/sap/clustersecret-operator/pulls) to add it.