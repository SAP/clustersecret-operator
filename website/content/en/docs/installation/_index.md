---
title: "Installation"
linkTitle: "Installation"
weight: 10
type: "docs"
description: >
  Overview on available installation methods
---

clustersecret-operator introduces one custom resource type, `clustersecrets.core.cs.sap.com`, with kind `ClusterSecret`.
The according definitions can be found 
[here](https://github.com/sap/clustersecret-operator/blob/main/crds/clustersecrets.yaml).
This definition must be deployed before the executables provided by this repository can be started.
The core of the clustersecret-operator installation are the controller and webhook executables built from this repository.
Docker images are available here:
- controller: `ghcr.io/sap/clustersecret-operator/controller`
- webhook: `ghcr.io/sap/clustersecret-operator/webhook`

A complete deployment consists of:
- the custom resource definition
- the controller deployment
- the webhook deployment
- rbac objects for controller and webhook (service accounts, (cluster) roles, according (cluster) role bindings)
- a service for the webhooks
- webhook configurations.

Note that it is highly recommended to always activate the webhooks, as they are not only validating, but
also adding essential defaulting logic. Running without this mutating functionality
might lead to unexpected behavior.

The following deployment methods are available (recommended is Helm).