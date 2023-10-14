---
title: "简介"
weight: 110
---

Doris Operator 致力于简化在 [Kubernetes](https://kubernetes.io/) 上管理 [Apache Doris](https://github.com/apache/doris)
集群的流程，自动化 Doris 集群的运维任务和附属监控设施的管理，让 Doris 成为真正的云原生数据库。

![image](arch.png)

## 版本要求

- Kubernetes: 1.16+
- Apache Doris: 2.0.0+

## 主要特性

Doris Operator 包含了以下关键特性：

- **Kubernetes 包管理支持**

  Doris Operator 支持通过 [Helm](https://helm.sh/) 或 [Kustomize](https://kustomize.io/) 安装 ，您只需要一个命令即可轻松部署。


- **滚动更新 Doris 集群**

  有序平滑地滚动更新 Doris 集群，实现 Doris 集群不停机地更新配置/升级版本。


- **安全扩缩容 Doris 集群**

  Doris Operator 为 Doris 提供云上的安全水平可拓展性。


- **计算节点自动扩缩容**

  根据 Doris 计算负载压力，自动水平扩缩容 Doris 集群计算节点。


- **自动故障恢复**

  当集群发生故障时，Doris Operator 会自动为 Doris 集群实施故障故障恢复。


- **自动创建 Doris 集群可观测性设施**

  自动部署对 Doris 集群的 Prometheus、Grafana 监控组件和 Loki 日志组件，以维护 Doris 集群的可观测性。


- **自动数据备份 (开发中)**

  提供用户友好的自定义、自动重试的周期性数据备份。


- **多租户支持**

  用户可以在一个 Kubernetes 集群上轻松部署和管理多个 Doris 集群。
