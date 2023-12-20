---
title: "Deploy Doris Cluster"
weight: 330
---

This document describes how to deploy a Doris cluster on general Kubernetes.

## Deploy the Doris cluster

1. Create `Namespace`:

    ```shell
    kubectl create namespace ${namespace}
    ```

2. Deploy the Doris cluster:

    ```shell
    kubectl apply -f ${cluster_name} -n ${namespace}
    ```

   {{< callout context="caution" title="Note" icon="rocket" >}}
   It is recommended to organize Doris cluster configurations under the `${cluster_name}` directory and save them
   as `${cluster_name}/doris-cluster.yaml`. After modifying the configuration and applying it to Kubernetes, The new
   configuration will be automatically applied to the Doris cluster.
   {{< /callout >}}

   {{< details "If your Kubernetes cluster cannot connect to the Internet" open >}}

   If the Kubernetes servers cannot connect to the internet, you need to download the Docker images used by the Doris
   cluster on a machine that has internet access and then upload them to the server. Afterwards, use `docker load` to
   install the Docker images on the server.

   Deploying a Doris cluster will require the following Docker images (assuming Doris cluster version is {{< param last_doris_image_version >}}):

    ```shell
    ghcr.io/linsoss/doris-fe:{{< param last_doris_image_version >}}
    ghcr.io/linsoss/doris-be:{{< param last_doris_image_version >}}
    ghcr.io/linsoss/doris-cn:{{< param last_doris_image_version >}}
    ghcr.io/linsoss/doris-broker:{{< param last_doris_image_version >}}
    prom/prometheus:v2.37.8
    grafana/grafana:9.5.2
    grafana/loki:2.9.1
    grafana/promtail:2.9.1
    tnir/mysqlclient:1.4.6
    busybox:1.36
    ```

   Then download all these images using the following command:

    ```shell
    docker pull ghcr.io/linsoss/doris-fe:{{< param last_doris_image_version >}}
    docker pull ghcr.io/linsoss/doris-be:{{< param last_doris_image_version >}}
    docker pull ghcr.io/linsoss/doris-cn:{{< param last_doris_image_version >}}
    docker pull ghcr.io/linsoss/doris-broker:{{< param last_doris_image_version >}}
    docker pull prom/prometheus:v2.37.8
    docker pull grafana/grafana:9.5.2
    docker pull grafana/loki:2.9.1
    docker pull grafana/promtail:2.9.1
    docker pull tnir/mysqlclient:1.4.6
    docker pull busybox:1.36
    
    docker save -o doris-fe-{{< param last_doris_image_version >}}.tar ghcr.io/linsoss/doris-fe:{{< param last_doris_image_version >}}
    docker save -o doris-be-{{< param last_doris_image_version >}}.tar ghcr.io/linsoss/doris-be:{{< param last_doris_image_version >}}
    docker save -o doris-cn-{{< param last_doris_image_version >}}.tar ghcr.io/linsoss/doris-cn:{{< param last_doris_image_version >}}
    docker save -o doris-broker-{{< param last_doris_image_version >}}.tar ghcr.io/linsoss/doris-broker:{{< param last_doris_image_version >}}
    docker save -o prometheus-v2.37.8.tar prom/prometheus:v2.37.8
    docker save -o grafana-9.5.2.tar grafana/grafana:9.5.2
    docker save -o loki-2.9.1.tar grafana/loki:2.9.1
    docker save -o promtail-2.9.1.tar grafana/promtail:2.9.1
    docker save -o mysqlclient-1.4.6.tar tnir/mysqlclient:1.4.6
    docker save -o busybox-1.36.tar busybox:1.36
    ```

   Next, upload these Docker images to the Kubernetes server and execute `docker load` to install these Docker images on
   the server:

    ```shell
    docker load -i doris-fe-{{< param last_doris_image_version >}}.tar
    docker load -i doris-be-{{< param last_doris_image_version >}}.tar
    docker load -i doris-cn-{{< param last_doris_image_version >}}.tar
    docker load -i doris-broker-{{< param last_doris_image_version >}}.tar
    docker load -i prometheus-v2.37.8.tar
    docker load -i grafana-9.5.2.tar
    docker load -i loki-2.9.1.tar
    docker load -i promtail-2.9.1.tar
    docker load -i mysqlclient-1.4.6.tar
    docker load -i busybox-1.36.tar
    ```

{{< /details >}}

3. View the DorisCluster CR status:

    ```shell
    kubectl get dc ${cluster_name} -n ${namespace} -o yaml
    ```

4. View the Pod status:

    ```shell
    kubectl get po -n ${namespace} -l app.kubernetes.io/instance=${cluster_name}
    ```

You can utilize the Doris Operator to deploy and manage multiple Doris clusters within a single Kubernetes cluster.
Repeat the steps outlined above and replace `cluster_name` with a different name for each cluster. These clusters
can exist within the same namespace or different namespaces based on your specific requirements and preferences.

## Initialize the Doris cluster

If you need to perform some initialization tasks after deploying the Doris cluster, such as modifying the initial
passwords for root and admin, or executing initialization SQL scripts, please refer to
the [Initialize Doris Cluster](../initialize-doris-cluster) documentation.

## Deploy the Doris monitoring components

If you require Doris Operator to automatically deploy observability components for your Doris cluster, including
Prometheus and Grafana for monitoring the Doris cluster, as well as Loki for centralized log collection and querying,
please refer to the [Deploy Doris Monitoring Components](../deploy-monitor-for-doris-cluster/) documentation.
