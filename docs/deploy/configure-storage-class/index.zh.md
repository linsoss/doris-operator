---
title: "配置 Storage Class"
weight: 310
---

Doris 集群中 FE、BE 以及监控等组件需要使用将数据持久化的存储。Kubernetes
上的数据持久化需要使用 [PersistentVolume (PV)](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)
。Kubernetes 提供多种[存储类型](https://kubernetes.io/docs/concepts/storage/volumes/)，主要分为两大类：

- 网络存储

  存储介质不在当前节点，而是通过网络方式挂载到当前节点。一般有多副本冗余提供高可用保证，在节点出现故障时，对应网络存储可以再挂载到其它节点继续使用。


- 本地存储

  存储介质在当前节点，通常能提供比网络存储更低的延迟，但没有多副本冗余，一旦节点出故障，数据就有可能丢失。如果是 IDC
  服务器，节点故障可以一定程度上对数据进行恢复，但公有云上使用本地盘的虚拟机在节点故障后，数据是**无法找回**的。

PV 一般由系统管理员或 volume provisioner 自动创建，PV 与 Pod
是通过 [PersistentVolumeClaim (PVC)](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims)
进行关联的。普通用户在使用 PV 时并不需要直接创建 PV，而是通过 PVC 来申请使用 PV，对应的 volume provisioner 根据 PVC 创建符合要求的
PV，并将 PVC 与该 PV 进行绑定。

{{< callout context="caution" title="注意" icon="alert-triangle" >}}
为了数据安全，任何情况下都不要直接删除 PV，除非对 volume provisioner 原理非常清楚。手动删除 PV 可能导致非预期的行为。
{{< /callout >}}

## Doris 集群推荐存储类型

BE 要求存储有较低的读写延迟，所以生产环境强烈推荐使用本地 SSD 存储。

FE 作为存储集群元信息的数据库，并不是 IO 密集型应用，所以一般本地普通 SAS 盘或网络 SSD 存储（例如 AWS 上 gp2 类型的 EBS
存储卷，Google Cloud 上的持久化 SSD 盘）就可以满足要求。

监控组件工具，由于自身没有做多副本冗余，所以为保证可用性，推荐用网络存储。

## 网络 PV 配置

Kubernetes 1.11
及以上的版本支持[网络 PV 的动态扩容](https://kubernetes.io/blog/2018/07/12/resizing-persistent-volumes-using-kubernetes/)
，但用户需要为相应的 `StorageClass` 开启动态扩容支持。

开启动态扩容后，通过下面方式对 PV 进行扩容：

1. 修改 PVC 大小

   假设之前 PVC 大小是 10 Gi，现在需要扩容到 100 Gi

    ```other
    kubectl patch pvc -n ${namespace} ${pvc_name} -p '{"spec": {"resources": {"requests": {"storage": "100Gi"}}}}'
    ```

2. 查看 PV 扩容成功

   扩容成功后，通过 `kubectl get pvc -n ${namespace} ${pvc_name}` 显示的大小仍然是初始大小，但查看 PV 大小会显示已经扩容到预期的大小。

    ```other
    kubectl get pv | grep ${pvc_name}
    ```

## 本地 PV 配置

Kubernetes
当前支持静态分配的本地存储。可使用 [local-static-provisioner](https://github.com/kubernetes-sigs/sig-storage-local-static-provisioner)
项目中的 `local-volume-provisioner` 程序创建本地存储对象。

### 第 1 步：准备本地存储

- BE
  数据使用的盘，可通过[普通挂载](https://github.com/kubernetes-sigs/sig-storage-local-static-provisioner/blob/master/docs/operations.md#use-a-whole-disk-as-a-filesystem-pv)
  方式将盘挂载到 `/mnt/ssd` 目录。

  出于性能考虑，推荐 BE 独占一个磁盘，并且推荐磁盘类型为 SSD。


- FE
  数据使用的盘，可以参考[步骤](https://github.com/kubernetes-sigs/sig-storage-local-static-provisioner/blob/master/docs/operations.md#sharing-a-disk-filesystem-by-multiple-filesystem-pvs)
  挂载盘，创建目录，并将新建的目录以 bind mount 方式挂载到 `/mnt/sharedssd` 目录下。
-

给监控数据使用的盘，可以参考[步骤](https://github.com/kubernetes-sigs/sig-storage-local-static-provisioner/blob/master/docs/operations.md#sharing-a-disk-filesystem-by-multiple-filesystem-pvs)
挂载盘，创建目录，并将新建的目录以 bind mount 方式挂载到 `/mnt/monitoring` 目录下。

上述的 `/mnt/ssd`、`/mnt/sharedssd`、`/mnt/monitoring` 是 local-volume-provisioner 使用的发现目录（discovery
directory），local-volume-provisioner 会为发现目录下的每一个子目录创建对应的 PV。

### 第 2 步：部署 local-volume-provisioner

1. 下载 local-volume-provisioner 部署文件。

    ```shell
   wget https://raw.githubusercontent.com/linsoss/doris-operator/dev/examples/local-pv/local-volume-provisioner.yaml
   ```

2. 如果你使用的发现路径与 **第 1 步：准备本地存储** 中的示例一致，可跳过这一步。如果你使用与上一步中不同路径的发现目录，需要修改
   ConfigMap 和 DaemonSet 定义。

- 修改 ConfigMap 定义中的 `data.storageClassMap` 字段：

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

  关于 local-volume-provisioner
  更多的配置项，参考文档 [Configuration](https://github.com/kubernetes-sigs/sig-storage-local-static-provisioner/blob/master/docs/provisioner.md#configuration) 。

  - 修改 DaemonSet 定义中的 `volumes` 与 `volumeMounts` 字段，以确保发现目录能够挂载到 Pod 中的对应目录：

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

3. 部署 local-volume-provisioner 程序。

    ```other
    kubectl apply -f local-volume-provisioner.yaml
    ```

4. 检查 Pod 和 PV 状态。

    ```other
    kubectl get po -n kube-system -l app=local-volume-provisioner && \
    kubectl get pv | grep -e ssd-storage -e shared-ssd-storage -e monitoring-storage -e backup-storage
    ```

   `local-volume-provisioner` 会为发现目录下的每一个挂载点创建一个 PV。

更多信息，可参阅 [Kubernetes 本地存储](https://kubernetes.io/docs/concepts/storage/volumes/#local)
和 [local-static-provisioner 文档](https://github.com/kubernetes-sigs/sig-storage-local-static-provisioner#overview)。

### 最佳实践

- 本地 PV 的路径是本地存储卷的唯一标示符。为了保证唯一性并避免冲突，推荐使用设备的 UUID 来生成唯一的路径。
- 如果想要 IO 隔离，建议每个存储卷使用一块物理盘，在硬件层隔离。
- 如果想要容量隔离，建议每个存储卷一个分区使用一块物理盘，或者每个存储卷使用一块物理盘。

更多信息，可参阅 local-static-provisioner
的[最佳实践文档](https://github.com/kubernetes-sigs/sig-storage-local-static-provisioner/blob/master/docs/best-practices.md)。

## 数据安全

一般情况下 PVC 在使用完删除后，与其绑定的 PV 会被 provisioner
清理回收再放入资源池中被调度使用。为避免数据意外丢失，可在全局配置 `StorageClass` 的回收策略 (reclaim policy) 为 `Retain`
或者只将某个 PV 的回收策略修改为 `Retain`。`Retain` 模式下，PV 不会自动被回收。

- **全局配置**

  `StorageClass` 的回收策略一旦创建就不能再修改，所以只能在创建时进行设置。如果创建时没有设置，可以再创建相同 provisioner
  的 `StorageClass`，例如 GKE 上默认的 fe 类型的 `StorageClass` 默认保留策略是 `Delete`，可以再创建一个名为 `fe-standard`
  的保留策略是 `Retain` 的存储类型，并在创建 Doris 集群时将相应组件的 `storageClassName` 修改为 `fe-standard`。

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

- **配置单个 PV**

    ```other
    kubectl patch pv ${pv_name} -p '{"spec":{"persistentVolumeReclaimPolicy":"Retain"}}'
    ```

### 删除 PV 以及对应的数据

PV 保留策略是 `Retain` 时，如果确认某个 PV 的数据可以被删除，需要严格按照下面的操作顺序来删除 PV 以及对应的数据：

1. 删除 PV 对应的 PVC 对象：

    ```other
    kubectl delete pvc ${pvc_name} --namespace=${namespace}
    ```

2. 设置 PV 的保留策略为 `Delete`，PV 会被自动删除并回收：

    ```other
    kubectl patch pv ${pv_name} -p '{"spec":{"persistentVolumeReclaimPolicy":"Delete"}}'
    ```

要了解更多关于 PV
的保留策略可参考[修改 PV 保留策略](https://kubernetes.io/docs/tasks/administer-cluster/change-pv-reclaim-policy/)。
