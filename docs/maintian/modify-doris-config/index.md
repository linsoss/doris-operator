---
title: "Modify Doris Cluster Configuration"
weight: 610
---

Doris cluster supports online configuration changes for components such as FE and BE using SQL, without the need to
restart the components.

However, in a Doris cluster deployed on Kubernetes, certain components have their configuration overridden by the
configuration in the DorisCluster custom resource (CR) after an upgrade or restart. This can result in the loss of the
online changed configuration.

Therefore, to persistently modify the configuration, you need to directly modify the configuration items in the
DorisCluster CR.

1. Refer to the parameters
   in [Configure Doris Components](../../deploy/configure-doris-cluster/#doris-configuration) and modify the
   configuration for various components in the DorisCluster CR:

    ```shell
    kubectl edit dc ${cluster_name} -n ${namespace}
    ```

2. View the progress of the updates after modifying the configuration:

    ```shell
    watch kubectl -n ${namespace} get pod -o wide
    ```