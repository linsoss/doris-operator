---
title: "Get Started"
weight: 130
---

This document introduces how to create a simple Kubernetes cluster and use it to deploy a basic test Doris cluster
usingDoris Operator.

{{< callout context="caution" title="Caution" icon="alert-triangle" >}}

This document is for demonstration purposes only. **Do not** follow it in production environments. For deployment in
production environments, refer to the instructions in [See also](https://#see-also).

{{< /callout >}}

## Step 1: Create a test Kubernetes

[kind](https://kind.sigs.k8s.io/) is suitable for running a local Kubernetes cluster using Docker containers.

{{< tabs "Install kind & kubectl" >}}

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

Create a Kubernetes cluster：

```shell
kind create cluster
```

Check whether the Kubernetes cluster was created successfully:

```shell
kubectl cluster-info
```

## Step 2: Deploy Doris Operator

Doris Operator supports two installation methods: [Kustomized](../../installation/kustomized-installation/)
and [Helm](../../installation/helm-installation/). It is recommended to use Kustomized for installation.

{{< tabs "Install flux-cli " >}}

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

Install Doris Operator：

```shell
mkdir doris-operator
flux pull artifact oci://ghcr.io/linsoss/kustomize/doris-operator:1.0.0 --output ./doris-operator/
kubectl apply -k doris-operator
```

To confirm that the Doris Operator components are running, run the following command:

```shell
kubectl get pods -n doris-operator-system
```

Once all the Pods are in the "Running" state, you can proceed to the next step.

## Step 3: Deploy a Doris cluster and monitor

- **Deploy Doris Cluster**

    ```shell
    kubectl create ns doris
    kubectl apply -n doris -f https://raw.githubusercontent.com/linsoss/doris-operator/dev/examples/basic/doris-cluster.yaml 
    ```

- **Deploy Doris Cluster Monitor**

   ```shell
   kubectl apply -n doris -f https://raw.githubusercontent.com/linsoss/doris-operator/dev/examples/basic/doris-monitor.yaml
   ```

- **View the Pod Status**

    ```shell
    watch kubectl get po -n doris
    ```

  Once all components' pods are started, each component type (FE, BE, CN, Broker) should be in the Running state.

## Step 4: Connect to the Doris Cluster

- **Connect to Doris SQL Service**

  Forward the FE Service in Kubernetes for local access:

  ```shell
  kubectl port-forward -n doris svc/doris-fe 8030:8030 > /dev/null 2>&1 &
  ```

  You can use the `mysql` command-line tool to connect to Doris.

  ```shell
  mysql -h 127.0.0.1 -P 9030 -u root
  ```

- **Access Doris FE UI**

  Forward the HTTP port of the FE Service in Kubernetes:

  ```shell
  kubectl port-forward -n doris svc/doris-fe 8030:8030 > /dev/null 2>&1 &
  ```

  Then open [http://localhost:8030](http://localhost:3000/) in your browser.

- Access Grafana Dashboard

  Forward the Grafana Service in Kubernetes for local access:

  ```shell
  kubectl port-forward -n doris svc/basic-monitor-grafana 3000:3000 > /dev/null 2>&1 &
  ```

  Then open [http://localhost:3000](http://localhost:3000/) in your browser. The default username and password are
  both `admin`.

  {{< callout context="caution" title="Explore Doris Grafana" icon="rocket" >}}

    - [View Doris Grafana Dashboard](../../monitor/view-doris-grafana-dashboard/)
    - [Query Doris Logs in Grafana](../../monitor/query-doris-logs-in-grafana/)

  {{< /callout >}}

## Step 5: Destroy the Doris and Kubernetes cluster

After testing, you may want to destroy the Doris cluster and Kubernetes cluster.

- **Stop kubectl Port Forwarding**

  If you still have `kubectl` processes running that forward ports, terminate them:

  ```shell
  pgrep -lfa kubectl
  ```

- **Destroy the Doris Cluster**

  ```shell
  # Delete Doris Cluster
  kubectl delete dc basic -n doris
  # Delete Doris Monitor
  kubectl delete dm basic-monitor -n doris
  # Delete PV data and other resources
  kubectl delete pvc,secret,serviceaccount,rolebinding,role --selector=app.kubernetes.io/name=doris-cluster -n doris
  ```

- **Destroy Kubernetes Cluster**

  ```shell
  kind delete cluster
  ```

## See also

If you are interested in deploying a Doris cluster in production environments, refer to the following documents:

- [Configure Storage Class](../../deploy/configure-storage-class/)
- [Configure Doris Cluster](../../deploy/configure-doris-cluster/)
- [Deploy Doris Cluster](../../deploy/deploy-doris-cluster/)
- [Initialize Doris Cluster](../../deploy/initialize-doris-cluster/)
- [Deploy Monitor for Doris Cluster](../../monitor/deploy-monitor-for-doris-cluster/)
- [Automatic Scaling Doris Cluster](../../scale/automatic-scaling-doris-cluster/)
