---
title: "Manual Scaling Doris Cluster"
weight: 510
---

This document introduces how to horizontally and vertically scale a Doris cluster on Kubernetes.

## Horizontal Scaling

Horizontal scaling for Doris involves increasing or decreasing the number of Pods to achieve the scaling of the cluster.
When scaling the Doris cluster, the Doris components are scaled based on the provided `replicas` value.

To scale FE, BE, CN, Broker horizontally, use kubectl to modify `spec.<fe/be/cn/broker>.replicas` in the `DorisCluster`
object of the cluster to a desired value.

1. Modify the `replicas` value for the Doris cluster components as needed. For example, the following command sets
   the `replicas` value for BE to 3:

    ```shell
    kubectl patch -n ${namespace} tc ${cluster_name} --type merge --patch '{"spec":{"be":{"replicas":3}}}'
    ```

2. Check whether your configuration has been updated in the corresponding Doris cluster on Kubernetes.

    ```shell
    kubectl get doriscluster ${cluster_name} -n ${namespace} -o yaml
    ```

3. Observe whether the `DorisCluster` Pods have been increased or decreased.

    ```shell
    watch kubectl -n ${namespace} get pod -o wide
    ```

For FE, BE, and Broker, it might take 30 seconds to 1 minute to complete scaling in or out.

{{< callout context="caution" title="Note" icon="rocket" >}}

- The number of FE replicas should be odd, such as 1, 3, 5.
- When the FE and BE components scale in, the PVC of the deleted node will be retained.

{{< /callout >}}

### Decommiss Doris components before scaling in

In the current version of Doris Operator, when scaling in, the operator will not automatically delete the component's
metadata in FE. You need to manually decommission the component before scaling in. You can refer to the following SQL
operations.

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

  BE nodes should be safely decommissioned
  using [DECOMMISSION](https://doris.apache.org/docs/dev/admin-manual/cluster-management/elastic-expansion/#delete-be-nodes).

    ```sql
    SHOW BACKENDS\G;
    ALTER SYSTEM DECOMMISSION BACKEND "<be-host>:9050"
    ```

  Check `SystemDecommissioned` status for the BE node using `SHOW PROC '/backends';`. When `TabletNum` is 0, it
  indicates a successful decommission.

  If the remaining BE node's disk does not have enough space to accommodate the data migration, you can cancel
  DECOMMISSION using the following command, and Doris will redistribute the load later.

    ```sql
    CANCEL DECOMMISSION BACKEND "<be-host>:9050"
    ```

- **Broker**

    ```sql
    SHOW BROKER\G;
    ALTER SYSTEM DROP ALL BROKER "<broker-pod-name>:8000"
    ```

In future versions of Doris Operator, this process will be automated to simplify the scaling-in operation.

## Vertical Capacity Scaling

Vertical scaling involves adjusting the resource limits of Pods, either by increasing or decreasing them, to achieve
cluster scaling. Vertical scaling is essentially a rolling upgrade process for Pods.

To vertically scale FE, BE, CN, or Broker, modify the `spec.<fe/be/cn/broker>.resources` in the `DorisCluster` object
corresponding to the cluster using `kubectl`.

You can monitor the progress of vertical scaling using the following command. Vertical scaling is complete when all Pods
have been rebuilt and entered the `Running` state.

```shell
watch kubectl -n ${namespace} get pod -o wide
```

