---
title: "恢复误删除的 Doris 集群"
weight: 630
---

如果你使用 `DorisCluster` 意外删除了 Doris 集群，可参考本文介绍的方法恢复集群。

DorisOperator 使用 PV (Persistent Volume)、PVC (Persistent Volume Claim)
来存储持久化的数据，如果不小心使用 `kubectl delete doriscluster` 意外删除了 Doris 集群，PV/PVC 对象以及数据都会保留下来，以最大程度保证数据安全。

此时你可以使用 `kubectl create` 命令来创建一个同名同配置的集群，之前保留下来未被删除的 PV/PVC 以及数据会被复用：

```shell
kubectl -n ${namespace} create -f doris-cluster.yaml
```

