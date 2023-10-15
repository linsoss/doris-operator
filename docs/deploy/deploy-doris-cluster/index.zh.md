---
title: "部署 Doris 集群"
weight: 330
---

在部署 Doris 集群之前，需要先配置Doris 集群。请参阅[在 Kubernetes 中配置 Doris 集群](../configure-doris-cluster)。

## 部署 Doris 集群

配置Doris 集群后，请按照以下步骤部署Doris 集群。

1. 创建 `Namespace`：

    ```shell
    kubectl create namespace ${namespace}
    ```

2. 部署 Doris 集群：

    ```shell
    kubectl apply -f ${cluster_name} -n ${namespace}
    ```
   {{< callout context="caution" title="Note" icon="rocket"  >}}
   建议在 `${cluster_name}` 目录下组织 Doris 集群的配置，并将其另存为 `${cluster_name}/doris-cluster.yaml`
   。修改配置并提交后，会自动应用到
   Doris 集群中。
   {{< /callout >}}

   {{< details "如果您的 Kubernetes 无法连接外网 ？" open >}}

   如果 Kubernetes 服务器无法连接外网，需要在有外网的机器上将 Doris 集群用到的 Docker
   镜像下载下来并上传到服务器上，然后使用 `docker load` 将 Docker 镜像安装到服务器上。

   部署一套 Doris 集群会用到下面这些 Docker 镜像（假设 Doris 集群版本为 2.0.2）:

    ```shell
    ghcr.io/linsoss/doris-fe:2.0.2
    ghcr.io/linsoss/doris-be:2.0.2
    ghcr.io/linsoss/doris-cn:2.0.2
    ghcr.io/linsoss/doris-broker:2.0.2
    prom/prometheus:v2.37.8
    grafana/grafana:9.5.2
    grafana/loki:2.9.1
    grafana/promtail:2.9.1
    tnir/mysqlclient:1.4.6
    busybox:1.36
    ```

   接下来通过下面的命令将所有这些镜像下载下来：

    ```shell
    docker pull ghcr.io/linsoss/doris-fe:2.0.2
    docker pull ghcr.io/linsoss/doris-be:2.0.2
    docker pull ghcr.io/linsoss/doris-cn:2.0.2
    docker pull ghcr.io/linsoss/doris-broker:2.0.2
    docker pull prom/prometheus:v2.37.8
    docker pull grafana/grafana:9.5.2
    docker pull grafana/loki:2.9.1
    docker pull grafana/promtail:2.9.1
    docker pull tnir/mysqlclient:1.4.6
    docker pull busybox:1.36
    
    docker save -o doris-fe-2.0.2.tar ghcr.io/linsoss/doris-fe:2.0.2
    docker save -o doris-be-2.0.2.tar ghcr.io/linsoss/doris-be:2.0.2
    docker save -o doris-cn-2.0.2.tar ghcr.io/linsoss/doris-cn:2.0.2
    docker save -o doris-broker-2.0.2.tar ghcr.io/linsoss/doris-broker:2.0.2
    docker save -o prometheus-v2.37.8.tar prom/prometheus:v2.37.8
    docker save -o grafana-9.5.2.tar grafana/grafana:9.5.2
    docker save -o loki-2.9.1.tar grafana/loki:2.9.1
    docker save -o promtail-2.9.1.tar grafana/promtail:2.9.1
    docker save -o mysqlclient-1.4.6.tar tnir/mysqlclient:1.4.6
    docker save -o busybox-1.36.tar busybox:1.36
    ```

   接下来将这些 Docker 镜像上传到服务器上，并执行 `docker load` 将这些 Docker 镜像安装到服务器上：

    ```shell
    docker load -i doris-fe-2.0.2.tar
    docker load -i doris-be-2.0.2.tar
    docker load -i doris-cn-2.0.2.tar
    docker load -i doris-broker-2.0.2.tar
    docker load -i prometheus-v2.37.8.tar
    docker load -i grafana-9.5.2.tar
    docker load -i loki-2.9.1.tar
    docker load -i promtail-2.9.1.tar
    docker load -i mysqlclient-1.4.6.tar
    docker load -i busybox-1.36.tar
    ```

   {{< /details >}}

3. 通过以下命令查看 DorisCluster CR 状态：

    ```shell
    kubectl get dc ${cluster_name} -n ${namespace} -o yaml
    ```

4. 通过下面命令查看 Pod 状态：

    ```shell
    kubectl get po -n ${namespace} -l app.kubernetes.io/instance=${cluster_name}
    ```

单个 Kubernetes 集群中可以利用 Doris Operator 部署管理多套 Doris 集群，重复以上步骤并将 `cluster_name`
替换成不同名字即可。不同集群既可以在相同 `namespace` 中，也可以在不同 `namespace` 中，可根据实际需求进行选择。

## 初始化 Doris 集群

如果要在部署完 Doris 集群后做一些初始化工作，比如修改 root、admin 初始密码，执行初始化 SQL
脚本，请参考 [初始化 Doris 集群](../initialize-doris-cluster) 文档。

## 部署 Doris 监控组件

如果您需要 Doris Operator 为您的 Doris 集群自动部署可观测性组件，包含自动监控该 Doris 集群的 Prometheus、Grafana
，以及自动收集日志并集中查询的 Loki，请参考 [部署 Doris 集群监控](../deploy-monitor-for-doris-cluster) 文档。
