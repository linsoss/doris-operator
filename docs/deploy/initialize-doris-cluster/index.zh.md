---
title: "初始化 Doris 集群"
weight: 340
---

本文介绍如何对 Kubernetes 上的 Doris 集群进行初始化 root, admin 账号和密码设置，以及批量自动执行 SQL 语句对数据库进行初始化。

{{< callout context="caution" title="Note" icon="rocket" >}}

- 以下功能只在 Doris 集群创建后第一次执行起作用，执行完以后再修改不会生效。
- 初始化过程不依赖 `root` 帐户，创建 Doris 集群后可以随意手动修改 `root` 帐户密码。

{{< /callout >}}

## 配置 DorisInitializer

通过配置 `DorisInitializer` CR 来配置 Doris 集群的帐号密码、SQL 初始化行为。

{{< details "A basic DorisInitializer CR sample" >}}

[doris-initializer.yaml](https://github.com/linsoss/doris-operator/blob/dev/examples/basic-initialize/doris-initializer.yaml)

{{< readfile file="/examples/basic-initialize/doris-initializer.yaml" code="true" lang="yaml" >}}

{{< /details >}}

{{< details "A advanced DorisInitializer CR sample" >}}

[doris-initializer.yaml](https://github.com/linsoss/doris-operator/blob/dev/examples/advanced/doris-initializer.yaml)

{{< readfile file="/examples/advanced/doris-initializer.yaml" code="true" lang="yaml" >}}

{{< /details >}}

{{< callout context="caution" title="Note" icon="rocket"  >}}

建议在 `${cluster_name}` 目录下组织 Doris 集群的配置，并将其另存为 `${cluster_name}/doris-initializer.yaml`。

{{< /callout >}}

### 设置 Doris 集群的名称

在 `${cluster_name}/doris-initializer.yaml` 文件中，修改`spec.cluster` 字段为对应 DorisCluster CR 的 `metadata.name`:

```yaml
spec:
  cluster: ${cluster_name}
```

### 设置初始化账号密码

集群创建时默认会创建 `root` 、`admin` 账号，但是密码为空，这会带来一些安全性问题。可以通过如下步骤设置初始密码。

```yaml
spec:
  # ...
  rootPassword: "your password"
  adminPassword: "your password"
```

### 设置初始化 SQL 脚本

集群在初始化过程还可以自动执行 `initSqlScript` 中的 SQL 语句用于初始化，该功能可以用于默认给集群创建一些 database 或者
table，并且执行一些用户权限管理类的操作。

```yaml
spec:
  # ...
  initSqlScript: |
    CREATE DATABASE app;
    GRANT ALL PRIVILEGES ON app.* TO 'developer'@'%';
```

## 执行初始化

```shell
kubectl apply -f ${cluster_name}/doris-initializer.yaml --namespace=${namespace}
```

以上命令会自动创建一个初始化的 Job，初始化完成后 Pod 状态会变成 Completed。

查看初始化任务运行情况：

```shell
kubectl get dorisinitializer ${dorisinitializer.metadata.name} -n ${namespace} -o yaml
```

