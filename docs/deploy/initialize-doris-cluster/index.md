---
title: "Initialize Doris Cluster"
weight: 340
---

This document describes how to initialize the passwords of `root` and `admin` accounts for Doris clusters on Kubernetes,
and how to initialize the database by executing SQL statements automatically in batch.

{{< callout context="caution" title="Note" icon="rocket" >}}

- The following steps apply only when you have created a cluster for the first time. Further configuration or
  modification after the initial cluster creation is not valid.

{{< /callout >}}

## Configure DorisInitializer

You can configure the account passwords and SQL initialization behavior for the Doris cluster by setting up
the `DorisInitializer` Custom Resource.

{{< details "A basic DorisInitializer CR sample" >}}

[doris-initializer.yaml](https://github.com/linsoss/doris-operator/blob/dev/examples/basic-initialize/doris-initializer.yaml)

{{< readfile file="/examples/basic-initialize/doris-initializer.yaml" code="true" lang="yaml" >}}

{{< /details >}}

{{< details "A advanced DorisInitializer CR sample" >}}

[doris-initializer.yaml](https://github.com/linsoss/doris-operator/blob/dev/examples/advanced/doris-initializer.yaml)

{{< readfile file="/examples/advanced/doris-initializer.yaml" code="true" lang="yaml" >}}

{{< /details >}}

{{< callout context="caution" title="Note" icon="rocket"  >}}

It is recommended to organize Doris cluster configurations under the `${cluster_name}` directory and save them
as `${cluster_name}/doris-initializer.yaml`.

{{< /callout >}}

### Set the Name of the Doris Cluster

In the `${cluster_name}/doris-initializer.yaml` file, modify the `spec.cluster` field to correspond to
the `metadata.name` of the DorisCluster CR:

```yaml
spec:
  cluster: ${cluster_name}
```

### Setting the Initial Account Passwords

When the Doris cluster is created, the `root` and `admin` accounts are created by default, but their passwords are
empty, which can pose some security risks. You can set the initial passwords as follows:

```yaml
spec:
  # ...
  rootPassword: "your password"
  adminPassword: "your password"
```

### Set the Initialization SQL Script

The cluster can also automatically execute the SQL statements in batch in `initSqlScript` during the initialization.
This function can be used to create some databases or tables for the cluster and perform user privilege management
operations.

```yaml
spec:
  # ...
  initSqlScript: |
    CREATE DATABASE app;
    GRANT ALL PRIVILEGES ON app.* TO 'developer'@'%';
```

## Initialize the cluster

```shell
kubectl apply -f ${cluster_name}/doris-initializer.yaml --namespace=${namespace}
```

The above command will automatically create an initialization Job. Once initialization is complete, the Pod's status
will change to Completed.

To check the status of the initialization task, run the following command:

```shell
kubectl get dorisinitializer ${dorisinitializer.metadata.name} -n ${namespace} -o yaml
```

