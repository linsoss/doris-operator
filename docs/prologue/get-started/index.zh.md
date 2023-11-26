---
title: "快速上手"
weight: 130
---

本文档介绍了如何创建一个简单的 Kubernetes 集群，部署 Doris Operator，并使用Doris Operator 部署Doris 集群。

## 第 1 步：创建 Kubernetes 测试集群

[kind](https://kind.sigs.k8s.io/) 十分适合用于使用 Docker 容器作为集群节点运行本地 Kubernetes 集群。

以下命令快速安装 kind 和 kubectl:
{{< tabs "安装 kind 和 kubectl" >}}
{{< tab "Mac" >}}

```shell
brew install kind
brew install kubectl
```

{{< /tab >}}
{{< tab "Linux" >}}

```shell
curl -Lo ./kind "https://kind.sigs.k8s.io/dl/v0.20.0/kind-darwin-amd64"
chmod +x ./kind
mv ./kind /bin/kind

curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/linux/amd64/kubectl
chmod +x ./kubectl
mv ./kubectl /bin/kubectl
```

{{< /tab >}}
{{< /tabs >}}

创建 Kubernetes 集群：

```shell
kind create cluster
```

检查 Kubernetes 集群是否创建成功：

```shell
kubectl cluster-info
```

## 第 2 步：部署 Doris Operator

Doris Operator 支持 [Kustomized](../../installation/kustomized-installation/)
和 [Helm](../../installation/helm-installation/) 两种安装方式，推荐使用 Kustomized 安装。

{{< tabs "安装 flux-cli " >}}

{{< tab "Mac" >}}

```shell
brew install fluxcd/tap/flux
```

{{< /tab >}}

{{< tab "Linux" >}}

```shell
curl -s https://fluxcd.io/install.sh | sudo bash
```

{{< /tab >}}

{{< /tabs >}}

安装 Doris Operator：

```shell
mkdir doris-operator
flux pull artifact oci://ghcr.io/linsoss/kustomize/doris-operator:1.0.3 --output ./doris-operator/
kubectl apply -k doris-operator
```

检查Doris Operator 组件是否正常运行：

```shell
kubectl get pods -n doris-operator-system
```

所有的 pods 都处于 `Running` 状态时，继续下一步。

## 第 3 步：部署 Doris 集群和监控

- **部署Doris 集群**

    ```shell
    kubectl create ns doris
    kubectl apply -n doris -f https://raw.githubusercontent.com/linsoss/doris-operator/dev/examples/basic/doris-cluster.yaml 
    ```

  如果访问 ghcr 网速较慢，可以使用 dockerproxy 代理的镜像：

    ```shell
    kubectl apply -n doris -f https://raw.githubusercontent.com/linsoss/doris-operator/dev/examples/basic_cn_special/doris-cluster.yaml
    ```

- **部署 Doris 集群监控**

   ```shell
   kubectl apply -n doris -f https://raw.githubusercontent.com/linsoss/doris-operator/dev/examples/basic/doris-monitor.yaml
   ```

  如果访问 ghcr 网速较慢，可以使用 dockerproxy 代理的镜像：

   ```shell
   kubectl apply -n doris -f https://raw.githubusercontent.com/linsoss/doris-operator/dev/examples/basic_cn_special/doris-monitor.yaml
   ```

- **查看 Pod 状态**

   ```shell
   watch kubectl get po -n doris
   ```

  所有组件的 Pod 都启动后，每种类型组件（FE，BE，CN，Broker）都会处于 Running 状态。

## 第 4 步：连接 Doris 集群

- **连接 Doris SQL 服务**

  转发 Kuebrnetes 中的 FE Service，以便本地访问：

    ```shell
    kubectl port-forward -n doris svc/doris-fe 9030:9030 > /dev/null 2>&1 &
    ```

  您可以直接使用 `mysql` 命令行工具连接Doris 进行操作。

    ```shell
    mysql -h 127.0.0.1 -P 9030 -u root
    ```

- **访问 Doris FE UI**

  转发 Kuebrnetes 中 FE Service 的 HTTP 端口：

    ```shell
    kubectl port-forward -n doris svc/doris-fe 8030:8030 > /dev/null 2>&1 &
    ```

  在浏览器中打开 [http://localhost:8030](http://localhost:3000/)

- **访问 Grafana 面板**

  转发 Kubernetes 中的 Grafana Service，以便本地访问：

    ```shell
    kubectl port-forward -n doris svc/basic-monitor-grafana 3000:3000 > /dev/null 2>&1 &
    ```

  然后在浏览器中打开 [http://localhost:3000](http://localhost:3000/)，默认用户名和密码都为 `admin`。

- **探索 Doris Grafana**

    - [查看 Doris Grafana 仪表盘](../../monitor/%E6%9F%A5%E7%9C%8B-doris-grafana-%E4%BB%AA%E8%A1%A8%E7%9B%98/)
    - [在 Grafana 查询 Doris 日志](../../monitor/%E5%9C%A8-grafana-%E6%9F%A5%E8%AF%A2-doris-%E6%97%A5%E5%BF%97/)

## 第 5 步：销毁 Doris 和 Kubernetes

完成测试后，您可能希望销毁 Doris 集群和 Kubernetes 集群。

- **停止 kubectl 的端口转发**

  如果你仍在运行正在转发端口的 `kubectl` 进程，请终止它们：

    ```shell
    pgrep -lfa kubectl
    ```

- 销毁 Doris 集群

    ```shell
    # 删除 Doris Cluster
    kubectl delete dc basic -n doris
    # 删除 Doris Monitor
    kubectl delete dm basic-monitor -n doris
    # 删除 PV 数据和其他资源
    kubectl delete pvc,secret,serviceaccount,rolebinding,role --selector=app.kubernetes.io/name=doris-cluster -n doris
    ```

- 销毁 Kubernetes 集群

    ```shell
    kind delete cluster
    ```

## 探索更多

如果你想在生产环境部署，请参考以下文档：

- [配置 Storage Class](../../deploy/%E9%85%8D%E7%BD%AE-storage-class/)
- [配置 Doris 集群](../../deploy/%E9%85%8D%E7%BD%AE-doris-%E9%9B%86%E7%BE%A4/)
- [部署 Doris 集群](../../deploy/%E9%83%A8%E7%BD%B2-doris-%E9%9B%86%E7%BE%A4/)
- [初始化 Doris 集群](../../deploy/%E5%88%9D%E5%A7%8B%E5%8C%96-doris-%E9%9B%86%E7%BE%A4/)
- [部署 Doris 集群监控](../../monitor/%E9%83%A8%E7%BD%B2-doris-%E9%9B%86%E7%BE%A4%E7%9B%91%E6%8E%A7/)
- [自动扩缩容 Doris 集群](../../scale/%E8%87%AA%E5%8A%A8%E6%89%A9%E7%BC%A9%E5%AE%B9-doris-%E9%9B%86%E7%BE%A4/)
