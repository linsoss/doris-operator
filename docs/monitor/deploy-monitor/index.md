---
title: "Deploy Monitor for Doris Cluster"
weight: 410
---

This document describes how to monitor a Doris clusters deployed through Doris Operator.

## Configure the DorisMonitor

Doris Operator collects Doris cluster metrics through Prometheus and gathers Doris cluster logs using Loki, providing a
unified visualization interface through Grafana.

When creating a new Doris cluster through Doris Operator, it's possible to create and configure an independent
monitoring system for each Doris cluster. This monitoring system operates within the same namespace as the Doris cluster
and consists of four components: Prometheus, Grafana, Loki, and Promtail.

The `DorisInitializer` CR defines the configuration for Doris visualization components:

{{< details "A basic DorisMonitor CR sample" >}}

[doris-minitor.yaml](https://github.com/linsoss/doris-operator/blob/dev/examples/basic/doris-monitor.yaml)

{{< readfile file="/examples/basic/doris-monitor.yaml" code="true" lang="yaml" >}}

{{< /details >}}

{{< details "A advanced DorisMonitor CR sample" >}}

[doris-monitor.yaml](https://github.com/linsoss/doris-operator/blob/dev/examples/advanced/doris-monitor.yaml)

{{< readfile file="/examples/advanced/doris-monitor.yaml" code="true" lang="yaml" >}}

{{< /details >}}

{{< callout context="caution" title="Note" icon="rocket"  >}}

It is recommended to organize the Doris cluster's configuration under the `${cluster_name}` directory and save it
as `${cluster_name}/doris-monitor.yaml`.

{{< /callout >}}

### Storage

`spec.storageClassName` defines the storage type of the monitoring components. Refer to
the [storage configuration document](../configure-storage-class/).

```yaml
spec:
  # ...
  storageClassName: ${storageClassName}
```

`spec.<prometheus/grafana/loki>.requests.storage` defines the persistent storage size for Prometheus, Loki, and Grafana.
Please choose an appropriate size based on your data retention time. Below are recommended sizes for production
environments:

- prometheus: 50Gi or more
- loki: 50Gi or more
- grafana: 5Gi

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

### Data Retention Time

You can configure the data retention time for Prometheus and Loki components
using `spec.<prometheus/loki>.retentionTime`. When this value is not set, data in Prometheus and Loki will be retained
permanently on the respective bound PVCs.

The following example sets the data retention time for Prometheus and Loki to 15 days:

```yaml
spec:
  # ...
  prometheus:
    retentionTime: 15d
  loki:
    retentionTime: 15d
```

## Deploy the DorisMonitor

```yaml
kubectl apply -f ${cluster_name}/doris-monitor.yaml --namespace=${namespace}
```

View the status of the monitor components:

```yaml
kubectl get dorismonitor ${dorismonitor.metadata.name} -n ${namespace} -o yaml
```

## Access the DorisMonitor

### Access Grafana Dashboard

You can access the Grafana monitoring dashboard using `kubectl port-forward`:

```other
kubectl port-forward -n ${namespace} svc/${cluster_name}-grafana 3000:3000
```

Then open [http://localhost:3000](http://localhost:3000/) in your browser. The default username and password are
both `admin`.

You can also set `spec.grafana.service.type` to `NodePort` to access the monitoring dashboard through `NodePort`.

### Access Prometheus Monitoring Data

For cases where direct access to monitoring data is needed, you can access Prometheus using `kubectl port-forward`:

```other
kubectl port-forward -n ${namespace} svc/${cluster_name}-prometheus 9090:9090 
```

Then open [http://localhost:9090](http://localhost:9090/) in your browser or access this address through a client tool.

You can also set `spec.prometheus.service.type` to `NodePort` to access the monitoring data through `NodePort`.
