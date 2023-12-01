---
title: "重启 Doris 集群"
weight: 605
---

{{< callout context="caution" title="注意" icon="alert-triangle" >}}
以下功能在 Doris-Operator v1.0.4 及以上版本中可用。
{{< /callout >}}

在使用 Doris 集群的过程中，如果需要对集群进行重启，可以通过 `kubectl edit dc ${name} -n ${namespace}` 修改集群配置，为期望优雅滚动重启的
Doris 集群组件 Spec 添加 annotation `al-assad.github.io/restartedAt`，Value 设置为当前时间。

以下示例中为组件 FE、BE、CN、Broker 都设置了 annotation，表示将优雅滚动重启以上三个 Doris 集群组件的所有
Pod。可以根据实际情况，只为某个组件设置 annotation。

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

