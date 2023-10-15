---
title: "Recovery a deleted Doris Cluster"
weight: 630
---

If you have mistakenly deleted a Doris cluster using `DorisCluster`, you can use the method introduced in this document
to recover the cluster.

Doris Operator uses PV (Persistent Volume) and PVC (Persistent Volume Claim) to store persistent data. If you
accidentally delete a cluster using `kubectl delete tc`, the PV/PVC objects and data are still retained to ensure data
safety.

To recover the deleted cluster, use the `kubectl create` command to create a cluster that has the same name and
configuration as the deleted one. In the new cluster, the retained PV/PVC and data are reused.

```shell
kubectl -n ${namespace} create -f doris-cluster.yaml
```

