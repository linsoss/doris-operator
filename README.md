![image-banner](docs/img/banner.png)

# Doris Operator

The Doris Operator is designed to streamlines the management of [Apache Doris](https://github.com/apache/doris) clusters
on [Kubernetes](https://kubernetes.io/), automating operational tasks and monitoring of the Doris cluster, with the
primary goal of transforming Doris into a truly **cloud-native** database.

😆 Find out more on [our official website](https://linsoss.github.io/doris-operator).

![image-arch](docs/img/arch.png)

## Some Convincing Benefits

The Doris Operator encompasses the following key features:

- **Kubernetes Package Management Support**

  Doris Operator supports installation via [Helm](https://helm.sh/) or [Kustomize](https://kustomize.io/), enabling easy
  deployment with a single command.

- **Rolling Updates for Doris Cluster**

  Orchestrates a smooth rolling update for the Doris cluster, allowing for non-disruptive configuration updates and
  version upgrades.

- **Secure Scalability of Doris Cluster**

  The Doris Operator facilitates horizontal scalability for Doris in the cloud, ensuring a secure and efficient scaling
  process.

- **Automated Compute Node Scaling Based on Load**

  Automatically adjusts the cluster's compute nodes based on Doris load, optimizing performance through horizontal
  scaling.

- **Automatic Failover**

  In case of cluster failures, the Doris Operator seamlessly initiates automatic failover procedures to ensure
  uninterrupted service.

- **Automatic Monitoring Setup Upon Creation of Doris Cluster**

  Automatically deploys monitoring components such as [Prometheus](https://prometheus.io/)
  and [Grafana](https://grafana.com/) for monitoring, as well as [Loki](https://grafana.com/oss/loki/) for logging, to
  maintain the observability of the Doris cluster.

- **Automatic Data Backup (WIP)**

  Provides a user-friendly and customizable periodic data backup mechanism with automatic retry capabilities.

- **Multi-tenancy Support**

  Allows users to effortlessly deploy and manage multiple Doris clusters on a single Kubernetes cluster, promoting
  efficient multi-tenant utilization.

## Deploying a Doris Cluster in 3 minutes!

You can follow our [Get Started](https://linsoss.github.io/doris-operator/docs/prologue/get-started/) guide to quickly
start a testing Kubernetes cluster and play with Doris Operator on your own machine.

## Documentation

- [English](https://linsoss.github.io/doris-operator/docs/prologue/introduction/)
- [Chinese](https://linsoss.github.io/doris-operator/zh/docs/prologue/%E7%AE%80%E4%BB%8B/)

## License

Doris Operator is under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.

