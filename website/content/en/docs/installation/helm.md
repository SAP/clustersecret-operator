---
title: "Helm"
linkTitle: "Helm"
weight: 10
type: "docs"
description: >
  Installation by Helm
---

## Helm deployment

The recommended way to deploy clustersecret-operator is to use the [Helm chart](https://github.com/sap/clustersecret-operator-helm),
also available in packaged form:
- as helm package: [https://sap.github.io/clustersecret-operator-helm](https://sap.github.io/clustersecret-operator-helm)
- as OCI package: [oci://ghcr.io/sap/clustersecret-operator-helm](oci://ghcr.io/sap/clustersecret-operator-helm)

The chart does not require any mandatory parameters, so deploying clustersecret-operator is as easy as

```bash
helm repo add clustersecret-operator https://sap.github.io/clustersecret-operator-helm
helm -n clustersecret-operator upgrade -i clustersecret-operator clustersecret-operator/clustersecret-operator
```