---
title: "Automatic Scaling Doris Cluster"
weight: 520
---

This document describes how to enable automatic horizontal scaling for CN (Compute Node) in a Doris cluster on
Kubernetes, which based on the actual workload.

## Prerequisites

Automatic scaling with the Doris Operator requires Kubernetes version 1.22+ and the installation
of [Metric Server](https://kubernetes.io/docs/tasks/debug/debug-cluster/resource-metrics-pipeline/) on the Kubernetes
cluster.

## Configuring DorisAutoscaler

You can configure the automatic scaling behavior of the Doris cluster by configuring the `DorisInitializer` custom
resource (CR).

{{< details "A DorisAutoscaler CR sample" >}}

[doris-autoscaler.yaml](https://github.com/linsoss/doris-operator/blob/dev/examples/basic-autoscale/doris-autoscaler.yaml)

{{< readfile file="/examples/basic-autoscale/doris-autoscaler.yaml" code="true" lang="yaml" >}}

{{< /details >}}

{{< callout context="caution" title="Note" icon="rocket" >}}

It's recommended to organize the Doris cluster configuration under the `${cluster_name}` directory and save it
as `${cluster_name}/doris-autoscaler.yaml`.

{{< /callout >}}

### Replica Limit

`spec.cn.replicas` defines the maximum and minimum replica limits for CN's automatic scaling. In the following example,
the maximum number of replicas for CN scaling is limited to 5, and the minimum for scaling ins is 1.

```yaml
spec:
  cn:
    # ...
    replicas:
      min: 1
      max: 5
```

### Scaling Rules

`spec.cn.rules` defines the scaling rules for CN based on CPU and memory metrics.

```yaml
spec:
  cn:
    # ...
    rules:
      cpu:
        max: 90
        min: 20
        memory:
          max: 80
          min: 20
```

In the above example, `max` and `min` values are in percentages, e.g., 90 represents 90%. The DorisAutoscaler
dynamically scales in or out based on the overall CPU and memory utilization of CN nodes. Taking CPU as an example:

- When the overall average CPU usage of the CN cluster continuously exceeds `cpu.max` for a period, it automatically
  adds a replica until the next assessment does not exceed `cpu.max`.
- When the overall average CPU usage of the CN cluster falls below `cpu.min` for a period, it automatically removes a
  replica until the next assessment shows a CPU usage above `cpu.min`.

## Apply DorisAutoscaler

```shell
kubectl apply -f ${cluster_name}/doris-autoscale.yaml --namespace=${namespace}
```

View the status of the DorisAutoscaler:

```shell
kubectl get dorisautoscaler ${dorisautoscaler_name} -n ${namespace} -o yaml
```

## Delete DorisAutoScaler

```shell
kubectl delete -f ${cluster_name}/doris-autoscale.yaml --namespace=${namespace}
```

