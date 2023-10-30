---
title: "部署 Doris 集群监控"
weight: 410
---

本文介绍如何对通过Doris Operator 部署的Doris 集群进行监控及日志集中查询。

## 配置 DorisMonitor

Doris Operator 通过 Prometheus 收集 Doris 集群指标，通过 Loki 收集 Doris 集群日志，并在 Grafana 提供统一的可视化界面。

在通过Doris Operator 创建新的 Doris 集群时，可以对于每个Doris 集群，创建、配置一套独立的监控系统，与Doris 集群运行在同一
Namespace，包括 Prometheus、Grafana、Loki、Promtail 四个组件。

`DorisInitializer` CR 定义了 Doris 可视化组件的配置：

{{< details "A basic DorisMonitor CR sample" >}}

[doris-monitor.yaml](https://github.com/linsoss/doris-operator/blob/dev/examples/basic/doris-monitor.yaml)

{{< readfile file="/examples/basic/doris-monitor.yaml" code="true" lang="yaml" >}}

{{< /details >}}

{{< details "A advanced DorisMonitor CR sample" >}}

[doris-monitor.yaml](https://github.com/linsoss/doris-operator/blob/dev/examples/advanced/doris-monitor.yaml)

{{< readfile file="/examples/advanced/doris-monitor.yaml" code="true" lang="yaml" >}}

{{< /details >}}

{{< callout context="caution" title="Note" icon="rocket"  >}}

建议在 `${cluster_name}` 目录下组织 Doris 集群的配置，并将其另存为 `${cluster_name}/doris-monitor.yaml`。

{{< /callout >}}

### 存储

`spec.storageClassName` 定义了监控组件的存储类型，参考[存储配置文档](../configure-storage-class/)。

```yaml
spec:
  # ...
  storageClassName: ${storageClassName}
```

`spec.<prometheus/grafana/loki>.requests.storage`  定义了 Prometheus、Loki、Grafana 的持久存储大小。请根据您的数据保留时间选择合适的大小，以下是生产环境的建议：

- prometheus： 50Gi 以上；
- loki：50Gi 以上；
- grafana：5Gi

```yaml
spec:
  # ...
  prometheus:
    requests:
      storage: 50Gi
  grafana:
    requests:
      storage: 5Gi
  loki:
    requests:
      storage: 50Gi
```

### 数据保留时间

可以通过 `spec.<prometheus/loki>.retentionTime`  配置 Prometheus，Loki 组件的数据保留时间，当不设置该值时，Prometheus 和
Loki 的数据会永久保留在对应绑定的 PVC 上。

以下示例设置了 prometheus、loki 的数据保留时间为 15 天：

```yaml
spec:
  # ...
  prometheus:
    retentionTime: 15d
  loki:
    retentionTime: 15d
```

## 部署 DorisMonitor

```yaml
kubectl apply -f ${cluster_name}/doris-monitor.yaml --namespace=${namespace}
```

查看 monitor 组件的运行情况：

```yaml
kubectl get dorismonitor ${dorismonitor_name} -n ${namespace} -o yaml
```

## 访问 DorisMonitor

### 访问 Grafana 面板

可以通过 `kubectl port-forward` 访问 Grafana 监控面板：

```bash
kubectl port-forward -n ${namespace} svc/${dorismonitor_name}-grafana 3000:3000
```

然后在浏览器中打开 [http://localhost:3000](http://localhost:3000/)，默认用户名和密码都为 `admin`。

也可以设置 `spec.grafana.service.type` 为 `NodePort`，通过 `NodePort`查看监控面板。

### 访问 Prometheus 监控数据

对于需要直接访问监控数据的情况，可以通过 `kubectl port-forward` 来访问 Prometheus：

```bash
kubectl port-forward -n ${namespace} svc/${dorismonitor_name}-prometheus 9090:9090 
```

然后在浏览器中打开 [http://localhost:9090](http://localhost:9090/)，或通过客户端工具访问此地址即可。

也可以设置 `spec.prometheus.service.type` 为 `NodePort`，通过 `NodePort` 访问监控数据。
