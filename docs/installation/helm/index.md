---
title: "Helm Installation"
weight: 220
---

## Installation

This is a OCI helm chart, helm started to support OCI in version 3.8.0.

```shell
helm upgrade -i doris-operator oci://ghcr.io/linsoss/helm/doris-operator --version 1.0.2
```

## Values

| **Key**             | **Type** | **Default**                          | **Description**                            |
|---------------------|----------|--------------------------------------|--------------------------------------------|
| manager.image       | string   | ghcr.io/linsoss/doris-operator:1.0.2 | Controller container image tag             |
| manager.resources   | object   | {}                                   | Controller container resource requirement  |
| rbacProxy.image     | string   | bitnami/kube-rbac-proxy:0.14.1       | rbac-proxy container image tag             |
| rbacProxy.resources | object   | {}                                   | rbac-proxy container resource requirements |
| imagePullPolicy     | string   | IfNotPresent                         | image pull policy                          |
| imagePullSecrets    | list     | []                                   | image pull secrets                         |

