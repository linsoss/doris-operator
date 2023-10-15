---
title: "Access Doris Cluster"
weight: 350
---

Doris Cluster's external service is exposed through the FE component, including SQL interaction and the Web UI.

The Service can be configured with different types based on the scenario, such as `ClusterIP`, `NodePort`, etc., each
having different access methods.

You can obtain Doris FE Service information using the following command:

```shell
kubectl get svc ${serviceName} -n ${namespace}
```

For example：

```shell
kubectl get svc basic-fe -n default
NAME       TYPE       CLUSTER-IP    EXTERNAL-IP   PORT(S)                         AGE
basic-fe   NodePort   10.233.7.47   <none>        8030:31851/TCP,9030:31068/TCP   3d2h
```

The above example describes the information of the `basic-fe` service under the `default` namespace. The type
is `NodePort`, ClusterIP is `10.233.7.47`, and the ServicePort is `8030` and `9030`. The corresponding NodePorts
are `31851` and `31068`.

## ClusterIP

`ClusterIP` exposes the service via an internal IP within the cluster. When choosing this service type, you can only
access it within the cluster using the following methods:

- ClusterIP + ServicePort
- Service domain name (`${serviceName}.${namespace}`) + ServicePort

{{< callout context="caution" title="Note" icon="rocket" >}}
You can directly access the Doris cluster on your local machine through `kubectl port-forward`:

```shell
kubectl port-forward -n ${namespace} svc/doris-fe 9030:9030 
```

Access through mysql-client:

```shell
mysql -h localhost -P9030 -u root -p
```

{{< /callout >}}

## NodePort

NodePort exposes the service using the node's IP and a static port. By accessing `NodeIP + NodePort`, you can access a
NodePort service from outside the cluster.

To view the Node Ports allocated by the Service, you can use the following command to get information about Doris
Service object:

```shell
kubectl -n ${namespace} get svc ${cluster_name}-fe -ojsonpath="{.spec.ports[?(@.name=='query-port')].nodePort}{'\n'}"
```

