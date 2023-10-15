---
title: "自动扩缩容 Doris 集群"
weight: 520
---

本文介绍了如何让 Kubernetes 上的 Doris 集群的 CN 计算节点根据实际负载自动扩缩容。

## 前置要求

Doris Operator 的自动扩缩容要求 Kubernetes 版本 1.22 以及以上，且 Kubernetes
集群上已经安装 [Metric Server](https://kubernetes.io/docs/tasks/debug/debug-cluster/resource-metrics-pipeline/)。

## 配置 DorisAutoscaler

可以通过配置 `DorisInitializer` CR 来配置 Doris 集群的自动扩缩容行为。

{{< details "A DorisAutoscaler CR sample" >}}

[doris-autoscaler.yaml](https://github.com/linsoss/doris-operator/blob/dev/examples/basic-autoscale/doris-autoscaler.yaml)

{{< readfile file="/examples/basic-autoscale/doris-autoscaler.yaml" code="true" lang="yaml" >}}

{{< /details >}}

{{< callout context="caution" title="Note" icon="rocket"  >}}

建议在 `${cluster_name}` 目录下组织 Doris 集群的配置，并将其另存为 `${cluster_name}/doris-autoscaler.yaml`。

{{< /callout >}}

### 副本数量限制

`spec.cn.replicas`  定义了 CN 自动扩缩容的最大、最小副本数量限制。以下例子中限制了 CN 扩容最大的副本数量为 5，缩容的最小副本为
1。

```yaml
spec:
  cn:
    # ...
    replicas:
      min: 1
      max: 5
```

### 扩缩容规则

`spec.cn.rules` 定义了 CN 扩缩容的依据规则，支持根据 CPU 和内存的指标评估。

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

以上例子中，其中 `max`  和 `min`  的值为百分比，比如 90 代表 90%。DorisAutoscaler 会分别根据 CN 节点整体的 CPU
和内存的利用率进行动态扩缩容， 以 CPU 为例子：

- 当 CN 集群的整体平均CPU 使用率在一段时间内持续大于 `cpu.max` 时，将自动增加一个副本，直到下一轮计算评估不大于 `cpu.max`。
- 当 CN 集群的整体平均 CPU 占用率在一段时间内小于 `cpu.min`时，将自动移除一个副本，直到下一轮计算的 CPU
  占用率高于该 `cpu.min`。

## 执行DorisAutoscaler

```shell
kubectl apply -f ${cluster_name}/doris-autoscale.yaml --namespace=${namespace}
```

查看 DorisAutoscaler 的运行情况：

```shell
kubectl get dorisautoscaler ${dorisautoscaler_name} -n ${namespace} -o yaml
```

## 删除 DorisAutoScaler

```shell
kubectl delete -f ${cluster_name}/doris-autoscale.yaml --namespace=${namespace}
```

