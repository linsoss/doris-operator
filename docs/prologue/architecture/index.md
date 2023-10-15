---
title: "Operator Architecture"
weight: 120
---

This document describes the architecture of Doris Operator and how it works.

## Architecture

![image](arch.png)

`DorisCluster`, `DorisInitializer`, `DorisMonitor`, `DorisClusterAutoscaler`are custom resources defined by CRD.

- `DorisCluster`describes the desired state of the Doris cluster;
- `DorisMonitor`describes the monitoring components of the Doris cluster;
- `DorisInitializer` describes the desired initialization Job of the Doris cluster;
- `DorisClusterAutoscaler` describes the automatic scaling of the Doris cluster;

The following components are responsible for the orchestration and scheduling logic in a Doris cluster:

- `doris-controller-manager` is a set of custom controllers in Kubernetes. These controllers constantly compare the
  desired state recorded in the ``DorisCluster` object with the actual state of the Doris cluster. They adjust the
  resources in Kubernetes to drive the Doris cluster to meet the desired state and complete the corresponding control
  logic according to other CRs;
- `doris-admission-webhook` is a dynamic admission controller in Kubernetes, which completes the modification,
  verification, operation, and maintenance of Pod, StatefulSet, and other related resources.

## Control flow

The diagram below illustrates the control flow analysis of the Doris Operator:

![image](control-flow.png)

The overall control flow is described as follows:

1. The user creates a `DorisCluster` object and other CR objects through kubectl, such as `DorisMonitor`;
2. Doris Operator watches `DorisCluster` and other related objects, and constantly adjust
   the `StatefulSet`, `Deployment`, `Service`, and other objects of FE, BE, Broker, Monitor or other components based on
   the actual state of the cluster;
3. Kubernetes' native controllers create, update, or delete the corresponding `Pod` based on objects such
   as `StatefulSet`, `Deployment`, and `Job`;

Based on the above declarative control flow, Doris Operator automatically performs health check and fault recovery for
the cluster nodes. You can easily modify the `DorisCluster` object declaration to perform operations such as deployment,
upgrade, and scaling.
