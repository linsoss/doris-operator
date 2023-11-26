---
title: "配置 Doris BE 冷热分离存储"
weight: 660
---

本文描述了如何为 Doris BE on Kubernetes 配置多磁盘存储和冷热分离存储。

Doris BE 支持多个独立数据存储目录，比如通过同时挂载 SSD 和 HDD 存储介质来平衡热数据，冷数据的读写性能和成本。

在 [Doris 部署文档](https://doris.apache.org/docs/1.2/install/standard-deployment/#deploy-be)中描述了这部分的内容，Doris
Operator 提供了一种简单的方式来实现这一过程的配置，通过DorisCluster CRD 的 `spec.be.storage`  配置项。

{{< readfile file="/examples/be-multiple-storage/doris-cluster.yaml" code="true" lang="yaml" >}}

关于 StorageClass 和 PV 的制备，请参考：[配置 Storage Class](../../deploy/%E9%85%8D%E7%BD%AE-storage-class/)。
