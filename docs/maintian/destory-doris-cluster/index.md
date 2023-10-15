---
title: "Destroy a Doris Cluster"
weight: 620
---

This document describes how to destroy Doris clusters on Kubernetes.

## Destroy a Doris cluster managed by DorisCluster

To destroy a Doris cluster managed by `DorisCluster`, run the following command:

```other
kubectl delete dc ${cluster_name} -n ${namespace}
```

If you deploy the monitor in the cluster using `DorisMonitor`, run the following command to delete the monitor
component:

```other
kubectl delete dorismonitor ${doris_monitor_name} -n ${namespace}
```

## Delete data

After executing the command to destroy the cluster mentioned above, the persistent data in FE/BE of the Doris cluster
will still be retained. If you no longer need this data, you can clear it using the following command:

```shell
kubectl delete pvc --selector=app.kubernetes.io/name=${cluster_name},app.kubernetes.io/managed-by=doris-operator -n ${namespace}
```

There will still be some Kubernetes objects shared by Doris CRs on the namespace. You can delete them using the
following command:

```shell
kubectl delete secret,serviceaccount,rolebinding,role --selector=app.kubernetes.io/name=doris-cluster -n ${namespace}
```

