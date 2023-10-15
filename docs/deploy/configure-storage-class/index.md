---
title: "Configure Storage Class"
weight: 310
---

Doris cluster components such as FE, BE and monitoring components require persistent storage for data. To achieve this
on Kubernetes, you need to use [PersistentVolume (PV)](https://kubernetes.io/docs/concepts/storage/persistent-volumes/).
Kubernetes supports different types of [storage classes](https://kubernetes.io/docs/concepts/storage/volumes/), which
can be categorized into two main types:

- Network storage

  Network storage is not located on the current node but is mounted to the node through the network. It usually has
  redundant replicas to ensure high availability. In the event of a node failure, the corresponding network storage can
  be remounted to another node for continued use.


- Local storage

  Local storage is located on the current node and typically provides lower latency compared to network storage.
  However, it does not have redundant replicas, so data might be lost if the node fails. If the node is an IDC server,
  data can be partially restored, but if it is a virtual machine using local disk on a public cloud, data cannot be
  retrieved after a node failure.

PVs are automatically created by the system administrator or volume provisioner. PVs and Pods are bound
by [PersistentVolumeClaim (PVC)](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims).
Instead of creating a PV directly, users request to use a PV through a PVC. The corresponding volume provisioner creates
a PV that meets the requirements of the PVC and then binds the PV to the PVC.

{{< callout context="caution" title="Caution" icon="alert-triangle" >}}
Do not delete a PV under any circumstances unless you are familiar with the underlying volume provisioner. Manually
deleting a PV can result in orphaned volumes and unexpected behavior.
{{< /callout >}}

## Recommended storage classes for Doris clusters

In order to achieve low read and write latency, it is strongly recommended to use local SSD storage for the BE in the
production environment.

The FE, which serves as the database storing cluster metadata, is not an IO-intensive application. Hence, typically,
regular local SAS disks or network SSD storage (such as AWS gp2 EBS volumes or Google Cloud persistent SSD disks) should
suffice for its requirements.

For monitoring components, it's recommended to use network storage to ensure availability since they do not inherently
have built-in redundancy through multiple replicas.

## Network PV configuration

Starting from Kubernetes 1.11, volume expansion of network PV is supported. However, you need to run the following
command to enable volume expansion for the corresponding `StorageClass`:

```other
kubectl patch storageclass ${storage_class} -p '{"allowVolumeExpansion": true}'
```

After enabling volume expansion, you can expand the PV using the following method:

1. Edit the PersistentVolumeClaim (PVC) object:

   Suppose the PVC is currently 10 Gi and you need to expand it to 100 Gi.

    ```other
    kubectl patch pvc -n ${namespace} ${pvc_name} -p '{"spec": {"resources": {"requests": {"storage": "100Gi"}}}}'
    ```

2. View the size of the PV:

   After the expansion, the size displayed by running `kubectl get pvc -n ${namespace} ${pvc_name}` still shows the
   original size. However, if you run the following command to view the size of the PV, it shows that the size has been
   expanded to the expected value.

    ```other
    kubectl get pv | grep ${pvc_name}
    ```

## Local PV configuration

Kubernetes currently supports statically allocated local storage. You can use the `local-volume-provisioner` program
from the [local-static-provisioner](https://github.com/kubernetes-sigs/sig-storage-local-static-provisioner) project to
create local storage objects.

### Step 1: Pre-allocate local storage

- For the disks used by BE data, you can mount the disk to the `/mnt/ssd` directory using
  the [regular mounting](https://github.com/kubernetes-sigs/sig-storage-local-static-provisioner/blob/master/docs/operations.md#use-a-whole-disk-as-a-filesystem-pv)
  method.

  For performance reasons, it is recommended to dedicate a disk for BE and recommend using SSD disk types.


- For the disks used by FE data, you can follow
  the [steps](https://github.com/kubernetes-sigs/sig-storage-local-static-provisioner/blob/master/docs/operations.md#sharing-a-disk-filesystem-by-multiple-filesystem-pvs)
  to mount the disk, create a directory, and mount the newly created directory to `/mnt/sharedssd` using bind mount.
- For the disks used by monitoring data, you can follow
  the [steps](https://github.com/kubernetes-sigs/sig-storage-local-static-provisioner/blob/master/docs/operations.md#sharing-a-disk-filesystem-by-multiple-filesystem-pvs)
  to mount the disk, create a directory, and mount the newly created directory to `/mnt/monitoring` using bind mount.

The `/mnt/ssd`, `/mnt/sharedssd`, and `/mnt/monitoring` mentioned above are the discovery directories used
by `local-volume-provisioner`. The `local-volume-provisioner` will create corresponding PVs for each subdirectory under
the discovery directory.

### Step 2: Deploy the Local-Volume-Provisioner

1. Download the local-volume-provisioner deployment file.

    ```shell
   wget https://raw.githubusercontent.com/linsoss/doris-operator/dev/examples/local-pv/local-volume-provisioner.yaml
   ```

2. If your discovery path matches the example in *Step 1: Pre-allocate local storage*, you can skip this step. If you
   are using a different discovery directory
   path than the previous step, you need to modify the ConfigMap and DaemonSet definitions.

- Modify the `data.storageClassMap` field in the ConfigMap definition:

  ```yaml
  apiVersion: v1
  kind: ConfigMap
  metadata:
    name: local-provisioner-config
    namespace: kube-system
  data:
    # ...
    storageClassMap: |
      ssd-storage: # for BE
        hostDir: /mnt/ssd
        mountDir: /mnt/ssd
      shared-ssd-storage: # for FE
        hostDir: /mnt/sharedssd
        mountDir: /mnt/sharedssd
      monitoring-storage: # for moniting data
        hostDir: /mnt/monitoring
        mountDir: /mnt/monitoring
  ```

  For more configuration options for `local-volume-provisioner`, refer to
  the [Configuration](https://github.com/kubernetes-sigs/sig-storage-local-static-provisioner/blob/master/docs/provisioner.md#configuration)
  documentation.

  - Modify the `volumes` and `volumeMounts` fields in the DaemonSet definition to ensure that the discovery directory
    can be mounted to the corresponding directory in the Pod:

    ```yaml
    ......
          volumeMounts:
            - mountPath: /mnt/ssd
              name: local-ssd
              mountPropagation: "HostToContainer"
            - mountPath: /mnt/sharedssd
              name: local-sharedssd
              mountPropagation: "HostToContainer"
            - mountPath: /mnt/monitoring
              name: local-monitoring
              mountPropagation: "HostToContainer"
      volumes:
        - name: local-ssd
          hostPath:
            path: /mnt/ssd
        - name: local-sharedssd
          hostPath:
            path: /mnt/sharedssd
        - name: local-backup
          hostPath:
            path: /mnt/backup
        - name: local-monitoring
          hostPath:
            path: /mnt/monitoring
    ......
    ```

3. Deploy the local-volume-provisioner program.

    ```other
    kubectl apply -f local-volume-provisioner.yaml
    ```

4. Check the Pod and PV status.

    ```shell
    kubectl get po -n kube-system -l app=local-volume-provisioner && \
    kubectl get pv | grep -e ssd-storage -e shared-ssd-storage -e monitoring-storage -e backup-storage
    ```

   The `local-volume-provisioner` will create a PV for each mount point under the discovery directory.

   For more information, refer to [Kubernetes Local Storage](https://kubernetes.io/docs/concepts/storage/volumes/#local)
   and
   the [local-static-provisioner documentation](https://github.com/kubernetes-sigs/sig-storage-local-static-provisioner#overview).

### Best practices

- The unique identifier for a local PV is its path. To avoid conflicts, it is recommended to generate a unique path
  using the UUID of the device.
- To ensure I/O isolation, it is recommended to use a dedicated physical disk per PV for hardware-based isolation.
- For capacity isolation, it is recommended to use either a partition per PV or a physical disk per PV.

For more information on local PV on Kubernetes, refer to
the [Best Practices](https://github.com/kubernetes-sigs/sig-storage-local-static-provisioner/blob/master/docs/best-practices.md)
document.

## Data safety

In general, when a PVC is deleted and no longer in use, the PV bound to it is reclaimed and placed in the resource pool
for scheduling by the provisioner. To prevent accidental data loss, you can configure the reclaim policy of
the `StorageClass` to `Retain` globally or change the reclaim policy of a single PV to `Retain`. With the `Retain`
policy, a PV is not automatically reclaimed.

- **To configure globally:**

  The reclaim policy of a `StorageClass` is set at creation time and cannot be updated once created. If it is not set
  during creation, you can create another `StorageClass` with the same provisioner. For example, the default reclaim
  policy of the `StorageClass` for persistent disks on Google Kubernetes Engine (GKE) is `Delete`. You can create
  another `StorageClass` named `fe-standard` with a reclaim policy of `Retain` and change the `storageClassName` of the
  corresponding component to `fe-standard` when creating a Doris cluster.

    ```yaml
    apiVersion: storage.k8s.io/v1
    kind: StorageClass
    metadata:
      name: fe-standard
    parameters:
       type: fe-standard
    provisioner: kubernetes.io/gce-fe
    reclaimPolicy: Retain
    volumeBindingMode: Immediate
    ```

- **To configure a single PV:**

    ```other
    kubectl patch pv ${pv_name} -p '{"spec":{"persistentVolumeReclaimPolicy":"Retain"}}'
    ```

### Delete PV and data

When the reclaim policy of PVs is set to `Retain`, if you have confirmed that the data of a PV can be deleted, you can
delete the PV and its corresponding data by following these steps:

1. Delete the PVC object corresponding to the PV:

    ```other
    kubectl delete pvc ${pvc_name} --namespace=${namespace}
    ```

2. Set the reclaim policy of the PV to `Delete`. This automatically deletes and reclaims the PV.

    ```other
    kubectl patch pv ${pv_name} -p '{"spec":{"persistentVolumeReclaimPolicy":"Delete"}}'
    ```

For more details, refer to
the [Change the Reclaim Policy of a PersistentVolume](https://kubernetes.io/docs/tasks/administer-cluster/change-pv-reclaim-policy/)
document.
