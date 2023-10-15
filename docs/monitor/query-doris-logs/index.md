---
title: "Query Doris Logs in Grafana"
weight: 430
---

Doris Operator automatically collects logs from all Doris cluster components into Loki when DorisMonitor is deployed.
You can directly search these logs in Grafana.

## Go to Logs Queries UI

![image](img/img.png)

![image](img/img_1.png)

![image](img/img_2.png)

## LogQL

LogQL is used in Loki for querying logs. For more information on using LogQL, please refer to Loki Log queries.

## Quick Filtering of Doris Component Logs

If you're not familiar with LogQL, you can conveniently filter Doris component logs through Grafana's interactive
interface.

Here's a demonstration of how to filter logs for a specific Doris FE component instance "basic-fe-0":

1. Select the label "component":
   ![image](img/img_3.png)
2. Choose the target Doris component, in this case, select "fe":
   ![image](img/img_4.png)
3. Click "+" to add a new filter.
   ![image](img/img_5.png)
4. Select the label "instance":
   ![image](img/img_6.png)
5. Choose the Pod instance of the target component:
   ![image](img/img_7.png)
6. Click "Run query" to get the log results.
   ![image](img/img_8.png)

Through Grafana Loki, you can easily perform actions like querying all component error logs and keyword searches without
having to manually run "kubectl logs" on each pod.