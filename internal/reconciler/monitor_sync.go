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
	tran "github.com/al-assad/doris-operator/internal/transformer"
	"github.com/al-assad/doris-operator/internal/util"
	appv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
)

// Sync all subcomponents status of DorisMonitor
func (r *DorisMonitorReconciler) Sync() (dapi.DorisMonitorSyncStatus, error) {
	syncRes := &dapi.DorisMonitorSyncStatus{}
	errCtr := &util.MultiError{}

	type SyncStatus = *dapi.DorisMonitorSyncStatus
	type ErrCollector = *util.MultiError
	type MuteFn = func(SyncStatus, ErrCollector)

	// collect prometheus, grafana, loki, promtail status in parallel
	syncFns := []func() MuteFn{
		func() MuteFn {
			status, err := r.syncPrometheusStatus()
			return func(s SyncStatus, c ErrCollector) {
				s.Prometheus = status
				c.Collect(err)
			}
		},
		func() MuteFn {
			status, err := r.syncGrafanaStatus()
			return func(s SyncStatus, c ErrCollector) {
				s.Grafana = status
				c.Collect(err)
			}
		},
		func() MuteFn {
			status, err := r.syncLokiStatus()
			return func(s SyncStatus, c ErrCollector) {
				s.Loki = status
				c.Collect(err)
			}
		},
		func() MuteFn {
			status, err := r.syncPromtailStatus()
			return func(s SyncStatus, c ErrCollector) {
				s.Promtail = status
				c.Collect(err)
			}
		},
	}
	muteFnRes := util.ParallelRun(syncFns...)
	for _, muteFn := range muteFnRes {
		muteFn(syncRes, errCtr)
	}

	return *syncRes, errCtr.Dry()
}

func (r *DorisMonitorReconciler) syncPrometheusStatus() (dapi.PrometheusStatus, error) {
	status := util.PointerDeRefer(r.CR.Status.Prometheus.DeepCopy(), dapi.PrometheusStatus{})
	serviceKey := tran.GetPrometheusServiceKey(r.CR.ObjKey())
	pvcKey := tran.GetPrometheusPVCKey(r.CR.ObjKey())
	deployKey := tran.GetPrometheusDeploymentKey(r.CR.ObjKey())

	err := r.fillMonitorComponentStatus(&status.DorisMonitorComponentStatus, serviceKey, pvcKey, deployKey)
	return status, err
}

func (r *DorisMonitorReconciler) syncGrafanaStatus() (dapi.GrafanaStatus, error) {
	status := util.PointerDeRefer(r.CR.Status.Grafana.DeepCopy(), dapi.GrafanaStatus{})
	serviceKey := tran.GetGrafanaServiceKey(r.CR.ObjKey())
	pvcKey := tran.GetGrafanaPVCKey(r.CR.ObjKey())
	deployKey := tran.GetGrafanaDeploymentKey(r.CR.ObjKey())

	err := r.fillMonitorComponentStatus(&status.DorisMonitorComponentStatus, serviceKey, pvcKey, deployKey)
	return status, err
}

func (r *DorisMonitorReconciler) syncLokiStatus() (dapi.LokiStatus, error) {
	status := util.PointerDeRefer(r.CR.Status.Loki.DeepCopy(), dapi.LokiStatus{})
	serviceKey := tran.GetLokiServiceKey(r.CR.ObjKey())
	pvcKey := tran.GetLokiPVCKey(r.CR.ObjKey())
	deployKey := tran.GetLokiDeploymentKey(r.CR.ObjKey())

	err := r.fillMonitorComponentStatus(&status.DorisMonitorComponentStatus, serviceKey, pvcKey, deployKey)
	return status, err
}

func (r *DorisMonitorReconciler) syncPromtailStatus() (dapi.PromtailStatus, error) {
	status := util.PointerDeRefer(r.CR.Status.Promtail.DeepCopy(), dapi.PromtailStatus{})
	daemonsetKey := tran.GetPromtailDaemonSetKey(r.CR.ObjKey())
	status.DaemonSetRef = dapi.NewNamespacedName(daemonsetKey)

	// Get daemonset status
	daemonset := &appv1.DaemonSet{}
	exist, err := r.Exist(daemonsetKey, daemonset)
	if err != nil {
		return status, err
	}
	if exist {
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
	exist, err := r.Exist(deploymentKey, deploy)
	if err != nil {
		return err
	}
	if exist {
		baseStatus.Ready = deploy.Status.ReadyReplicas > 0
		baseStatus.Conditions = deploy.Status.Conditions
	}
	return nil
}
