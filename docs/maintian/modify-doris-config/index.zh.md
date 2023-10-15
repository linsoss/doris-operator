---
title: "修改 Doris 集群配置"
weight: 610
---

Doris 集群自身支持通过 SQL 对 FE、BE 等组件进行在线配置变更，无需重启集群组件。

但是，对于部署在 Kubernetes 中的 Doris 集群，部分组件在升级或者重启后，配置项会被 DorisCluster CR 中的配置项覆盖，导致在线变更的配置失效。

因此，如果需要持久化修改配置，你需要在 DorisCluster CR 中直接修改配置项。

1. 参考[配置 Doris 组件](../../deploy/configure-doris-cluster/#doris-configuration)中的参数，修改集群的 DorisCluster
   CR 中各组件配置：

    ```shell
    kubectl edit dc ${cluster_name} -n ${namespace}
    ```

2. 查看配置修改后的更新进度：

   ```shell
   watch kubectl -n ${namespace} get pod -o wide
   ```

