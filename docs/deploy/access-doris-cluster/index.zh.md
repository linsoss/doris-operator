---
title: "访问 Doris 集群"
weight: 350
---

Doris Cluster 对外的服务是通过 FE 组件暴露的，包括 SQL 交互和 Web UI。

Service 可以根据场景配置不同的类型，比如 `ClusterIP`、`NodePort`、`LoadBalancer` 等，对于不同的类型可以有不同的访问方式。

可以通过如下命令获取 Doris FE Service 信息：

```shell
kubectl get svc ${serviceName} -n ${namespace}
```

示例：

```shell
kubectl get svc basic-fe -n default
NAME       TYPE       CLUSTER-IP    EXTERNAL-IP   PORT(S)                         AGE
basic-fe   NodePort   10.233.7.47   <none>        8030:31851/TCP,9030:31068/TCP   3d2h
```

上述示例描述了 `default` namespace 下 `basic-fe` 服务的信息，类型为 `NodePort`，ClusterIP 为 `10.233.7.47`，ServicePort
为 `8030` 和 `9030`，对应的 NodePort 分别为 `31851` 和 `31068`。

## ClusterIP

`ClusterIP` 是通过集群的内部 IP 暴露服务，选择该类型的服务时，只能在集群内部访问，可以通过如下方式访问：

- ClusterIP + ServicePort
- Service 域名 (`${serviceName}.${namespace}`) + ServicePort

## NodePort

NodePort 是通过节点的 IP 和静态端口暴露服务。通过请求 `NodeIP + NodePort`，可以从集群的外部访问一个 NodePort 服务。

查看 Service 分配的 Node Port，可通过获取 Doris 的 Service 对象来获知：

```shell
kubectl -n ${namespace} get svc ${cluster_name}-fe -ojsonpath="{.spec.ports[?(@.name=='query-port')].nodePort}{'\n'}"
```

