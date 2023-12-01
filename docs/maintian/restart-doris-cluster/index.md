---
title: "Restart Doris Cluster"
weight: 605
---

{{< callout context="caution" title="Caution" icon="alert-triangle" >}}
The following feature is available in Doris-Operator v1.0.4 and later.
{{< /callout >}}

When using the Doris cluster, if you need to restart the cluster, you can modify the cluster configuration by using
`kubectl edit dc ${name} -n ${namespace}`.
To indicate the desired rolling restart for Doris cluster components, add the
annotation `al-assad.github.io/restartedAt` to the Spec of the respective Doris cluster components, with the Value set
to the current time.
In the following example, annotations are set for the FE, BE, CN, and Broker components, indicating
a rolling restart for all Pods of these three Doris cluster components.
Depending on the actual situation, you can set annotations only for specific components.

```yaml
apiVersion: al-assad.github.io/v1beta1
kind: DorisCluster
metadata:
  name: basic
spec:
  fe:
    ...
    annotations:
      al-assad.github.io/restartedAt: "2023-12-01T12:00"
  be:
    ...
    annotations:
      al-assad.github.io/restartedAt: "2023-12-01T12:00"
  cn:
    ...
    annotations:
      al-assad.github.io/restartedAt: "2023-12-01T12:00"
  broker:
    ...
    annotations:
      al-assad.github.io/restartedAt: "2023-12-01T12:00"
```
