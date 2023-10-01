/*
 *
 * Copyright 2023 @ Linying Assad <linying@apache.org>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * /
 */

package reconciler

import (
	dapi "github.com/al-assad/doris-operator/api/v1beta1"
	"github.com/al-assad/doris-operator/internal/transformer"
	"github.com/al-assad/doris-operator/internal/util"
	appv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
)

// Sync all subcomponents status of DorisMonitor
func (r *DorisMonitorReconciler) Sync() (dapi.DorisMonitorSyncStatus, error) {
	syncRes := dapi.DorisMonitorSyncStatus{}
	errCtr := &util.MultiError{}

	// todo modified to execute code in parallel
	util.CollectFnErr(errCtr, r.syncPrometheusStatus, func(s dapi.PrometheusStatus) { syncRes.Prometheus = s })
	util.CollectFnErr(errCtr, r.syncGrafanaStatus, func(s dapi.GrafanaStatus) { syncRes.Grafana = s })
	util.CollectFnErr(errCtr, r.syncLokiStatus, func(s dapi.LokiStatus) { syncRes.Loki = s })
	util.CollectFnErr(errCtr, r.syncPromtailStatus, func(s dapi.PromtailStatus) { syncRes.Promtail = s })

	return syncRes, errCtr.Dry()
}

func (r *DorisMonitorReconciler) syncPrometheusStatus() (dapi.PrometheusStatus, error) {
	status := util.PointerDeRefer(r.CR.Status.Prometheus.DeepCopy(), dapi.PrometheusStatus{})
	serviceKey := transformer.GetPrometheusServiceKey(r.CR.ObjKey())
	pvcKey := transformer.GetPrometheusPVCKey(r.CR.ObjKey())
	deployKey := transformer.GetPrometheusDeploymentKey(r.CR.ObjKey())

	err := r.fillMonitorComponentStatus(&status.DorisMonitorComponentStatus, serviceKey, pvcKey, deployKey)
	return status, err
}

func (r *DorisMonitorReconciler) syncGrafanaStatus() (dapi.GrafanaStatus, error) {
	status := util.PointerDeRefer(r.CR.Status.Grafana.DeepCopy(), dapi.GrafanaStatus{})
	serviceKey := transformer.GetGrafanaServiceKey(r.CR.ObjKey())
	pvcKey := transformer.GetGrafanaPVCKey(r.CR.ObjKey())
	deployKey := transformer.GetGrafanaDeploymentKey(r.CR.ObjKey())

	err := r.fillMonitorComponentStatus(&status.DorisMonitorComponentStatus, serviceKey, pvcKey, deployKey)
	return status, err
}

func (r *DorisMonitorReconciler) syncLokiStatus() (dapi.LokiStatus, error) {
	status := util.PointerDeRefer(r.CR.Status.Loki.DeepCopy(), dapi.LokiStatus{})
	serviceKey := transformer.GetLokiServiceKey(r.CR.ObjKey())
	pvcKey := transformer.GetLokiPVCKey(r.CR.ObjKey())
	deployKey := transformer.GetLokiDeploymentKey(r.CR.ObjKey())

	err := r.fillMonitorComponentStatus(&status.DorisMonitorComponentStatus, serviceKey, pvcKey, deployKey)
	return status, err
}

func (r *DorisMonitorReconciler) syncPromtailStatus() (dapi.PromtailStatus, error) {
	status := util.PointerDeRefer(r.CR.Status.Promtail.DeepCopy(), dapi.PromtailStatus{})
	daemonsetKey := transformer.GetPromtailDaemonSetKey(r.CR.ObjKey())
	status.DaemonSetRef = dapi.NewNamespacedName(daemonsetKey)

	// Get daemonset status
	daemonset := &appv1.DaemonSet{}
	if err := r.Find(daemonsetKey, daemonset); err != nil {
		return status, err
	}
	if daemonset != nil {
		status.Ready = daemonset.Status.NumberReady > 0
		status.Conditions = daemonset.Status.Conditions
	}
	return status, nil
}

func (r *DorisMonitorReconciler) fillMonitorComponentStatus(
	baseStatus *dapi.DorisMonitorComponentStatus,
	serviceKey types.NamespacedName,
	pvcKey types.NamespacedName,
	deploymentKey types.NamespacedName) error {

	baseStatus.ServiceRef = dapi.NewNamespacedName(serviceKey)
	baseStatus.PVCRef = dapi.NewNamespacedName(pvcKey)
	baseStatus.DeploymentRef = dapi.NewNamespacedName(deploymentKey)

	// Get deployment status
	deploy := &appv1.Deployment{}
	if err := r.Find(deploymentKey, deploy); err != nil {
		return err
	}
	if deploy != nil {
		baseStatus.Ready = deploy.Status.ReadyReplicas > 0
		baseStatus.Conditions = deploy.Status.Conditions
	}
	return nil
}