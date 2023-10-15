---
title: "销毁 Doris 集群"
weight: 620
---

本文描述了如何销毁 Kubernetes 集群上的 Doris 集群。

## 销毁 Doris 集群

要销毁使用 `DorisCluster` 管理的Doris 集群，执行以下命令：

```other
kubectl delete dc ${cluster_name} -n ${namespace}
```

如果集群中通过 `DorisMonitor` 部署了监控，要删除监控组件，可以执行以下命令：

```other
kubectl delete dorismonitor ${doris_monitor_name} -n ${namespace}
```

## 清除数据

上述销毁集群的命令后，Doris 集群中 FE/BE 持久化的数据仍然会保留。如果你不再需要那些数据，可以通过下面命令清除数据：

```shell
kubectl delete pvc --selector=app.kubernetes.io/name=${cluster_name},app.kubernetes.io/managed-by=doris-operator -n ${namespace}
```

该 namespace 上还会剩余一些 Doris CR 共用的 Kubernetes 对象，您可以通过以下命令删除：

```shell
kubectl delete secret,serviceaccount,rolebinding,role --selector=app.kubernetes.io/name=doris-cluster -n ${namespace}
```
