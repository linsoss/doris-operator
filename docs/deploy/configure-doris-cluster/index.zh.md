---
title: "配置 Doris 集群"
weight: 320
---

本文档介绍了如何配置生产可用的 Doris 集群，

## 资源配置

部署前需要根据实际情况和需求，为 Doris 集群各个组件配置资源，其中 FE、BE、CN、Broker 是 Doris
集群的核心服务组件，在生产环境下它们的资源配置还需要按组件要求指定，具体参考：[资源配置推荐](https://doris.apache.org/docs/dev/install/standard-deployment/#software-and-hardware-requirements)。

为了保证 Doris 集群的组件在 Kubernetes 中合理的调度和稳定的运行，建议为其设置 Guaranteed 级别的 QoS，通过在配置资源时让
limits 等于 requests 来实现,
具体参考：[配置 QoS](https://kubernetes.io/docs/tasks/configure-pod-container/quality-service-pod/)。

如果使用 NUMA 架构的 CPU，为了获得更好的性能，需要在节点上开启 `Static` 的 CPU 管理策略。为了 Doris 集群组件能独占相应的
CPU
资源，除了为其设置上述 Guaranteed 级别的 QoS 外，还需要保证 CPU 的配额必须是大于或等于 1
的整数。具体参考: [CPU 管理策略](https://kubernetes.io/docs/tasks/administer-cluster/cpu-management-policies)。

## 部署配置

通过配置 `DorisCluster` CR 来配置 Doris 集群：

{{< details "简要的 DorisCluster CR 示例" >}}
[doris-cluster.yaml](https://github.com/linsoss/doris-operator/blob/dev/examples/basic/doris-cluster.yaml)
{{< readfile file="/examples/basic/doris-cluster.yaml" code="true" lang="yaml" >}}
{{< /details >}}

{{< details "完整的 DorisCluster CR 示例" >}}
[doris-cluster.yaml](https://github.com/linsoss/doris-operator/blob/dev/examples/advanced/doris-cluster.yaml)
{{< readfile file="/examples/advanced/doris-cluster.yaml" code="true" lang="yaml" >}}
{{< /details >}}

{{< callout context="caution" title="Note" icon="rocket"  >}}
建议在 `${cluster_name}` 目录下组织 Doris 集群的配置，并将其另存为 `${cluster_name}/doris-cluster.yaml`。修改配置并提交后，会自动应用到
Doris 集群中。
{{< /callout >}}

### 集群名称

通过更改 `DorisCuster` CR 中的 `metadata.name` 来配置集群名称。

### 版本

正常情况下，集群内的各组件应该使用相同版本，所以一般建议配置 `spec.<fe/be/cn/broker>.baseImage` + `spec.version` 即可。

相关参数的格式如下：

- `spec.version`，格式为 `imageTag`，例如 `2.0.2`
- `spec.<fe/be/cn/broker>.baseImage`，格式为 `imageName`，例如 `ghcr.io/linsoss/doris-fe` ；

请注意必须使用 [doris-operator/images](https://github.com/linsoss/doris-operator/tree/dev/images)  进行构建的 Doris
组件镜像，当然您也可以直接使用 linsoss 发布的 doris 组件镜像 😃：

| Component | Image                                                                                                 |
|-----------|-------------------------------------------------------------------------------------------------------|
| FE        | [ghcr.io/linsoss/doris-fe](https://github.com/linsoss/doris-operator/pkgs/container/doris-fe)         |
| BE        | [ghcr.io/linsoss/doris-be](https://github.com/linsoss/doris-operator/pkgs/container/doris-be)         |
| CN        | [ghcr.io/linsoss/doris-cn](https://github.com/linsoss/doris-operator/pkgs/container/doris-cn)         |
| Broker    | [ghcr.io/linsoss/doris-broker](https://github.com/linsoss/doris-operator/pkgs/container/doris-broker) |

### 存储

如果需要设置存储类型，可以修改 `${cluster_name}/doris-cluster.yaml` 中各组件的 `storageClassName` 字段。

Doris 集群不同组件对磁盘的要求不一样，所以部署集群前，要根据当前 Kubernetes
集群支持的存储类型以及使用场景，参考[存储配置文档](../%E9%85%8D%E7%BD%AE-storage-class/)为 Doris 集群各组件选择合适的存储类型。

### Doris 组件配置参数

可以通过 `spec.<fe/be/cn/broker>.config`  来配置各个组件的参数。

比如想修改 FE 以下配置参数：

```yaml
prefer_compute_node_for_external_table=true
enable_spark_load=true
```

则修改 `DorisCluster` 的以下配置：

```yaml
spec:
  fe:
    config:
      prefer_compute_node_for_external_table: 'true'
      enable_spark_load: 'true'
```

{{< callout context="caution" title="Note" icon="rocket" >}}
并不需要为 FE 设置 enable_fqdn_mode，Doris Operator 会强制自动将该参数设置为 true 并注入容器。
{{< /callout >}}

### 配置 Doris 服务

通过配置 `spec.fe.service` 定义不同的 Service 类型，如 `ClusterIP` 、 `NodePort`。默认情况下 Doris Operator 会为 FE
创建一个额外的 `ClusterIP` 类型 Service。

- **ClusterIP**

  `ClusterIP` 是通过集群的内部 IP 暴露服务，选择该类型的服务时，只能在集群内部访问，使用 ClusterIP 或者 Service
  域名（`${cluster_name}-fe.${namespace}`）访问。

    ```yaml
    spec:
      doris:
        service:
          type: ClusterIP
    ```

- **NodePort**

  在本地测试时候，可选择通过 NodePort 暴露，Doris Operator 会绑定 FE 的 SQL 查询端口和 Web UI 端口到 NodePort。

  NodePort 是通过节点的 IP 和静态端口暴露服务。通过请求 `NodeIP + NodePort`，可以从集群的外部访问一个 NodePort 服务。

    ```yaml
    spec:
      doris:
        service:
          type: NodePort
    ```

### Hadoop 连接配置

当 Doris 集群需要连接 Hadoop，相关的 Hadoop 配置文件是必不可少的，`spec.hadoopConf` 配置项提供了方便的向 FE、BE、CN、Broke 注入
Hadoop 配置的方式。

```yaml
spec:
  hadoopConf:
  # Hadoop 集群的 hostname-ip 映射
  hosts:
    - ip: 10.233.123.189
      name: hadoop-01
    - ip: 10.233.123.179
      name: hadoop-02
    - ip: 10.233.123.179
      name: hadoop-03
  # Hadoop 配置文件内容
  configs:
    hdfs-site.xml: |
      <configuration>
      ...
      </configuration>
    hive-site.xml: |
      <configuration>
      ...
      </configuration>
```

### 物理拓扑高可用

Doris 是一个分布式数据库，以下介绍 3 种方式来为维持 Doris 在 Kubernetes 上的物理拓扑高可用。

#### 通过 nodeSelector 约束调度实例

通过各组件配置的 `nodeSelector` 字段，可以约束组件的实例只能调度到特定的节点上。关于 `nodeSelector`
的更多说明，请参阅 [nodeSelector](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#nodeselector)。

```yaml
apiVersion: al-assad.github.io/v1beta1
kind: DorisCluster
# ...
spec:
  fe:
    nodeSelector:
      node-role.kubernetes.io/fe: true
    # ...
  be:
    nodeSelector:
      node-role.kubernetes.io/be: true
    # ...
  cn:
    nodeSelector:
      node-role.kubernetes.io/cn: true
    # ...
  broker:
    nodeSelector:
      node-role.kubernetes.io/broker: true
```

#### 通过 tolerations 调度实例

通过各组件配置的 `tolerations`
字段，可以允许组件的实例能够调度到带有与之匹配的[污点](https://kubernetes.io/docs/reference/glossary/?all=true#term-taint) (
Taint)
的节点上。关于污点与容忍度的更多说明，请参阅 [Taints and Tolerations](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/)。

```yaml
apiVersion: al-assad.github.io/v1beta1
kind: DorisCluster
# ...
spec:
  fe:
    tolerations:
      - effect: NoSchedule
        key: dedicated
        operator: Equal
        value: fe
    # ...
  be:
    tolerations:
      - effect: NoSchedule
        key: dedicated
        operator: Equal
        value: be
    # ...
  cn:
    tolerations:
      - effect: NoSchedule
        key: dedicated
        operator: Equal
        value: cn
    # ...
  broker:
    tolerations:
      - effect: NoSchedule
        key: dedicated
        operator: Equal
        value: broker
    # ...
```

#### 通过 affinity 调度实例

配置 `PodAntiAffinity` 能尽量避免同一组件的不同实例部署到同一个物理拓扑节点上，从而达到高可用的目的。关于 Affinity
的使用说明，请参阅 [Affinity & AntiAffinity](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#affinity-and-anti-affinity)。

下面是一个避免 FE 实例调度到同一个物理节点的例子：

```yaml
affinity:
  podAntiAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      - labelSelector:
          matchExpressions:
            - key: app.kubernetes.io/component
              operator: In
              values:
                - fe
            - key: app.kubernetes.io/instance
              operator: In
              values:
                - ${name}
        topologyKey: kubernetes.io/hostname
```







