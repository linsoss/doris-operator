---
title: "Helm 方式"
weight: 220
---

## 安装

以下通过 OCI 安装 Doris operator helm chart 的方式，helm 在 3.8.0 开始支持 OCI。

```shell
helm upgrade -i doris-operator oci://ghcr.io/linsoss/helm/doris-operator --version {{< param last_doris_operator_version >}}
```

## Values

| **Key**             | **Type** | **Default**                                                              | **Description**                            |
|---------------------|----------|--------------------------------------------------------------------------|--------------------------------------------|
| manager.image       | string   | ghcr.io/linsoss/doris-operator:{{< param last_doris_operator_version >}} | Controller container image tag             |
| manager.resources   | object   | {}                                                                       | Controller container resource requirement  |
| rbacProxy.image     | string   | bitnami/kube-rbac-proxy:0.14.1                                           | rbac-proxy container image tag             |
| rbacProxy.resources | object   | {}                                                                       | rbac-proxy container resource requirements |
| imagePullPolicy     | string   | IfNotPresent                                                             | image pull policy                          |
| imagePullSecrets    | list     | []                                                                       | image pull secrets                         |

