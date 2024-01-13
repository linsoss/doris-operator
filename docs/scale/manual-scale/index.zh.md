---
title: "手动扩缩容 Doris 集群"
weight: 510
---

本文介绍如何对部署在 Kubernetes 上的 Doris 集群进行手动水平扩缩容和垂直扩缩容。

## 水平扩缩容

Doris 水平扩缩容操作指的是通过增加或减少 Pod 的数量，来达到集群扩缩容的目的。扩缩容 Doris 集群时，会按照填入的 `replicas`
值，对 Doris 组件进行扩缩容操作。

如果要对 FE、BE、CN、Broker 进行水平扩缩容，可以使用 kubectl 修改集群所对应的 `DorisCluster`
对象中的 `spec.<fe/be/cn/broker>.replicas`至期望值。

1. 按需修改 Doris 集群组件的 `replicas` 值。例如，执行以下命令可将 BE 的 `replicas` 值设置为 3：

    ```shell
    kubectl patch -n ${namespace} tc ${cluster_name} --type merge --patch '{"spec":{"be":{"replicas":3}}}'
    ```

2. 查看 Kubernetes 集群中对应的 Doris 集群是否更新到了你期望的配置。

    ```shell
    kubectl get doriscluster ${cluster_name} -n ${namespace} -o yaml
    ```

3. 观察 `DorisCluster` Pod 是否新增或者减少。

    ```shell
    watch kubectl -n ${namespace} get pod -o wide
    ```

FE、BE、Broker 通常需要 30 秒 - 1分钟完成扩容或缩容。

{{< callout context="caution" title="Note" icon="rocket"  >}}

- FE 副本应该是奇数，如 1、3、5；
- FE、BE 组件在缩容过程中被删除的节点的 PVC 会保留；

{{< /callout >}}

### 缩容前下线组件

当前版本的 Doris Operator 缩容时不会主动删除该组件的在 FE 元数据，需要您在缩容前手动下线该组件，可以参考以下 SQL 操作。

- **FE**

    ```sql
    SHOW FRONTENDS\G;
    ALTER SYSTEM DROP FOLLOWER "<fe-host>:9010"
    ```

- **CN**

    ```sql
    SHOW BACKENDS\G;
    ALTER SYSTEM DROP BACKEND "<cn-host>:9050"
    ```

- **BE**

  BE
  节点应该通过 [DECOMMISSION](https://doris.apache.org/docs/dev/admin-manual/cluster-management/elastic-expansion/#delete-be-nodes)
  安全下线。

    ```sql
    SHOW BACKENDS\G;
    ALTER SYSTEM DECOMMISSION BACKEND "<be-host>:9050"
    ```

  通过 `SHOW PROC '/backends';` 看到该 BE 节点的 `SystemDecommissioned` 状态为
  true。表示该节点正在进行下线，当其中的 `TabletNum` 为 0 时，表示下线成功。

  如果剩余的 BE 节点磁盘不足以容纳迁移的数据，可以通过以下命令取消 DECOMMISSION，后续 Doris 会重新进行负载均衡。

    ```sql
    CANCEL DECOMMISSION BACKEND "<be-host>:9050"
    ```

- **Broker**

    ```sql
    SHOW BROKER\G;
    ALTER SYSTEM DROP ALL BROKER "<broker-pod-name>:8000"
    ```

在 Doris Operator 后续版本将会自动实现这一过程，以简化整个缩容操作。

## 垂直扩容容量

垂直扩缩容操作指的是通过增加或减少 Pod 的资源限制，来达到集群扩缩容的目的。垂直扩缩容本质上是 Pod 滚动升级的过程。

要对 FE、BE、CN、Broker 进行垂直扩缩容，通过 kubectl 修改集群所对应的 `DorisCluster`
对象的 `spec.<fe/be/cn/broker>.resources` 至期望值。

可以通过以下命令查看垂直扩缩容进度，当所有 Pod 都重建完毕进入 `Running` 状态后，垂直扩缩容完成。

```shell
watch kubectl -n ${namespace} get pod -o wide
```

