---
title: "Cold-Hot Separation Storage for Doris BE"
weight: 660
---

This document describes how to configure multi-disk storage and cold-hot separation storage for Doris BE on Kubernetes.

Doris BE supports multiple independent data storage directories, balancing the read and write performance as well as the
cost of hot and cold data by simultaneously mounting SSD and HDD storage medium.

The [Doris deployment documentation](https://doris.apache.org/docs/1.2/install/standard-deployment/#deploy-be) provides
details on this aspect. The Doris Operator offers a straightforward way to configure this process through
the `spec.be.storage` item in the DorisCluster CRD.

{{< readfile file="/examples/be-multiple-storage/doris-cluster.yaml" code="true" lang="yaml" >}}

For the preparation of StorageClass and PV, please refer
to: [Configuring Storage Class](../../deploy/configure-storage-class/).
