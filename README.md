# Cluster Scope Secrets For Kubernetes

[![REUSE status](https://api.reuse.software/badge/github.com/SAP/clustersecret-operator)](https://api.reuse.software/info/github.com/SAP/clustersecret-operator)

## About this project

This repository adds cluster-scope secrets (as kind `ClusterSecret`) to Kubernetes clusters.
It contains a cluster-scope custom resource definition `clustersecrets.core.cs.sap.com` 
and an according operator reconciling resources of this type.

A typical `ClusterSecret` could look as follows:

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

When reconciling this object, the operator will ensure that an according regular `Secret` with the same name (`my-secret`) exists
in all namespaces matching the provided label selector.

## Requirements and Setup

The recommended deployment method is to use the [Helm chart](https://github.com/sap/clustersecret-operator-helm):

```bash
helm upgrade -i clustersecret-operator oci://ghcr.io/sap/clustersecret-operator-helm/clustersecret-operator
```

## Documentation

The project's documentation can be found here: [https://sap.github.io/clustersecret-operator](https://sap.github.io/clustersecret-operator).  
The API reference is here: [https://pkg.go.dev/github.com/sap/clustersecret-operator](https://pkg.go.dev/github.com/sap/clustersecret-operator).

## Support, Feedback, Contributing

This project is open to feature requests/suggestions, bug reports etc. via [GitHub issues](https://github.com/SAP/clustersecret-operator/issues). Contribution and feedback are encouraged and always welcome. For more information about how to contribute, the project structure, as well as additional contribution information, see our [Contribution Guidelines](CONTRIBUTING.md).

## Code of Conduct

We as members, contributors, and leaders pledge to make participation in our community a harassment-free experience for everyone. By participating in this project, you agree to abide by its [Code of Conduct](https://github.com/SAP/.github/blob/main/CODE_OF_CONDUCT.md) at all times.

## Licensing

Copyright 2023 SAP SE or an SAP affiliate company and clustersecret-operator contributors. Please see our [LICENSE](LICENSE) for copyright and license information. Detailed information including third-party components and their licensing/copyright information is available [via the REUSE tool](https://api.reuse.software/info/github.com/SAP/clustersecret-operator).
