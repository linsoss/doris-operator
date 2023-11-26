---
title: "Configure Doris Cluster"
weight: 320
---

This document introduces how to configure a Doris cluster for production deployment.

## Configure resources

Before deploying a Doris cluster, it is necessary to configure the resources for each component of the cluster depending
on your needs.
FE, BE, CN and Broker are the core service components of a Doris cluster.
In a production environment,
you need to configure resources of these components according to their needs.
For details, refer
to [Hardware Recommendations](https://doris.apache.org/docs/dev/install/standard-deployment/#software-and-hardware-requirements).

To ensure the proper scheduling and stable operation of the components of the Doris cluster on Kubernetes, it is
recommended to set Guaranteed-level quality of service (QoS) by making `limits` equal to `requests` when configuring
resources.
For details, refer
to [Configure Quality of Service for Pods](https://kubernetes.io/docs/tasks/configure-pod-container/quality-service-pod/).

If you are using a NUMA-based CPU, you need to enable `Static`'s CPU management policy on the node for better
performance.
To allow the Doris cluster component to monopolize the corresponding CPU resources, the CPU quota
must be an integer greater than or equal to `1`, apart from setting Guaranteed-level QoS as mentioned above.
For details, refer
to [Control CPU Management Policies on the Node](https://kubernetes.io/docs/tasks/administer-cluster/cpu-management-policies).

## Configure Doris deployment

Configuring Doris Cluster via `DorisCluster` Custom Resource (CR):

{{< details "A basic DorisCluster CR sample" >}}
[doris-cluster.yaml](https://github.com/linsoss/doris-operator/blob/dev/examples/basic/doris-cluster.yaml)
{{< readfile file="/examples/basic/doris-cluster.yaml" code="true" lang="yaml" >}}
{{< /details >}}

{{< details "A advanced DorisCluster CR sample" >}}
[doris-cluster.yaml](https://github.com/linsoss/doris-operator/blob/dev/examples/advanced/doris-cluster.yaml)
{{< readfile file="/examples/advanced/doris-cluster.yaml" code="true" lang="yaml" >}}
{{< /details >}}

{{< callout context="caution" title="Note" icon="rocket" >}}
It is recommended to organize Doris cluster configurations under the `${cluster_name}` directory and save them
as `${cluster_name}/doris-cluster.yaml`. After modifying the configuration and applying it to Kubernetes, The new
configuration will be automatically applied to the Doris cluster.
{{< /callout >}}

### Cluster name

The cluster name can be configured by changing `metadata.name` in the `DorisCuster` CR.

### Version

Usually, components in a cluster are in the same version.
It is recommended to configure `spec.<fe/be/cn/broker>.baseImage` and `spec.version`.

Here are the formats of the parameters:

- `spec.version`: the format is `imageTag`, such as `2.0.2`
- `spec.<fe/be/cn/broker>.baseImage`: the format is `imageName`, such as `ghcr.io/linsoss/doris-fe`

Please note that you must use the Doris images built
using [doris-operator/images](https://github.com/linsoss/doris-operator/tree/dev/images).
Of course, you can also directly use the Doris images released by [linsoss](https://github.com/orgs/linsoss/packages) 😃:

| Component | Image                                                                                                 |
|-----------|-------------------------------------------------------------------------------------------------------|
| FE        | [ghcr.io/linsoss/doris-fe](https://github.com/linsoss/doris-operator/pkgs/container/doris-fe)         |
| BE        | [ghcr.io/linsoss/doris-be](https://github.com/linsoss/doris-operator/pkgs/container/doris-be)         |
| CN        | [ghcr.io/linsoss/doris-cn](https://github.com/linsoss/doris-operator/pkgs/container/doris-cn)         |
| Broker    | [ghcr.io/linsoss/doris-broker](https://github.com/linsoss/doris-operator/pkgs/container/doris-broker) |

### Storage

You can set the storage class by modifying `storageClassName` of each component in `${cluster_name}/doris-cluster.yaml`
and `${cluster_name}/doris-monitor.yaml`.

Different components of a Doris cluster have different disk requirements.
Before deploying a Doris cluster, refer to
the [Storage Configuration document](../configure-storage-class) to select
an appropriate storage class for each component according to the storage classes supported by the current Kubernetes
cluster and usage scenario.

If you need to configure cold-hot separation storage for Doris BE, you can refer
to [Cold-Hot Separation Storage for Doris BE](../../maintian/cold-hot-separation-storage-for-doris-be/).

### Doris configuration

You can configure parameters for various Doris components using `spec.<fe/be/cn/broker>.config`.

For example, if you want to modify the following configuration parameters for FE:

```yaml
prefer_compute_node_for_external_table=true
enable_spark_load=true
```

You would modify the `DorisCluster` configuration as follows:

```yaml
spec:
  fe:
    config:
      prefer_compute_node_for_external_table: 'true'
      enable_spark_load: 'true'
```

{{< callout context="caution" title="Note" icon="rocket" >}}
It's not necessary to set `enable_fqdn_mode` for FE.
Doris Operator will automatically set this parameter to true and inject it into the container.
{{< /callout >}}

### Service

By configuring `spec.fe.service`, you can define different Service types such as `ClusterIP` and `NodePort`. By default,
Doris Operator creates an additional `ClusterIP` type Service for FE.

- **ClusterIP**

  `ClusterIP` exposes the service via an internal IP in the cluster.
  When choosing this service type, the service can only be accessed within the cluster using ClusterIP or the service
  domain (`${cluster_name}-fe.${namespace}`).

    ```yaml
    spec:
      doris:
        service:
          type: ClusterIP
    ```

- **NodePort**

  During local testing, you can choose to expose the service via NodePort.
  Doris Operator will bind the SQL query port and Web UI port of FE to the NodePort.

  NodePort exposes the service using the node's IP and a static port.
  By accessing `NodeIP + NodePort`, you can access a NodePort service from outside the cluster.

    ```yaml
    spec:
      doris:
        service:
          type: NodePort
    ```

### Hadoop connection configuration

When the Doris cluster needs to connect to Hadoop, the relevant Hadoop configuration files are essential.
The `spec.hadoopConf` configuration item provides a convenient way to inject Hadoop configurations into FE, BE, CN, and
Broker components.

```yaml
spec:
  hadoopConf:
  # Host name and IP address of Hadoop cluster
  hosts:
    - ip: 10.233.123.189
      name: hadoop-01
    - ip: 10.233.123.179
      name: hadoop-02
    - ip: 10.233.123.179
      name: hadoop-03
  # Hadoop conf file contents
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

### High Availability of Physical Topology

Doris is a distributed database.
Here are three methods to maintain physical topology high availability for Doris on Kubernetes.

#### Use nodeSelector to schedule Pods

By configuring the `nodeSelector` field of each component, you can specify the specific nodes that the component Pods
are scheduled onto. For details on `nodeSelector`, refer
to [nodeSelector](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#nodeselector).

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

#### Use tolerations to schedule Pods

By configuring the `tolerations` field of each component, you can allow the component Pods to schedule onto nodes with
matching [taints](https://kubernetes.io/docs/reference/glossary/?all=true#term-taint). For details on taints and
tolerations, refer
to [Taints and Tolerations](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/).

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

#### Use affinity to schedule Pods

By configuring `PodAntiAffinity`, you can avoid the situation in which different instances of the same component are
deployed on the same physical topology node.
In this way, disaster recovery (high availability) is achieved.
For the user guide of Affinity,
see [Affinity & AntiAffinity](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#affinity-and-anti-affinity).

Here is an example of avoiding the scheduling of FE pod on the same physical node:

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

