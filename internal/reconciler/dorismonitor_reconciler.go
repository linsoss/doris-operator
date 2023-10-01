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
	"fmt"
	dapi "github.com/al-assad/doris-operator/api/v1beta1"
	"github.com/al-assad/doris-operator/internal/transformer"
	"github.com/al-assad/doris-operator/internal/util"
	corev1 "k8s.io/api/core/v1"
)

var (
	PrometheusConfHashAnnotationKey = fmt.Sprintf("%s/prometheus-conf", dapi.GroupVersion.Group)
	GrafanaConfHashAnnotationKey    = fmt.Sprintf("%s/grafana-conf", dapi.GroupVersion.Group)
	LokiConfHashAnnotationKey       = fmt.Sprintf("%s/loki-conf", dapi.GroupVersion.Group)
	PromtailConfHashAnnotationKey   = fmt.Sprintf("%s/promtail-conf", dapi.GroupVersion.Group)
)

// DorisMonitorReconciler reconciles a DorisMonitor object
type DorisMonitorReconciler struct {
	ReconcileContext
	CR *dapi.DorisMonitor
}

type MonitorStageRecResult struct {
	Stage  dapi.DorisMonitorOprStage
	Status dapi.OprStageStatus
	Action dapi.OprStageAction
	Err    error
}

// Reconcile all sub components of DorisMonitor
func (r *DorisMonitorReconciler) Reconcile() MonitorStageRecResult {
	stages := []func() MonitorStageRecResult{
		r.recRbacResources,
		r.recPrometheusResources,
		r.recGrafanaResources,
		r.recLokiResources,
		r.recPromtailResources,
	}
	for _, stageFn := range stages {
		result := stageFn()
		if result.Err != nil {
			return result
		}
	}
	return MonitorStageRecResult{Stage: dapi.MnrOprStageCompleted, Status: dapi.StageResultSucceeded}
}

func (r *MonitorStageRecResult) AsDorisClusterRecStatus() dapi.DorisMonitorRecStatus {
	res := dapi.DorisMonitorRecStatus{
		Stage:       r.Stage,
		StageStatus: r.Status,
		StageAction: r.Action,
	}
	if r.Err != nil {
		res.LastMessage = r.Err.Error()
	}
	return res
}

func mnrStageSucc(stage dapi.DorisMonitorOprStage, action dapi.OprStageAction) MonitorStageRecResult {
	return MonitorStageRecResult{Stage: stage, Status: dapi.StageResultSucceeded, Action: action}
}

func mnrStageFail(stage dapi.DorisMonitorOprStage, action dapi.OprStageAction, err error) MonitorStageRecResult {
	return MonitorStageRecResult{Stage: stage, Status: dapi.StageResultSucceeded, Action: action, Err: err}
}

// reconcile rbac resources used by the DorisMonitor
func (r *DorisMonitorReconciler) recRbacResources() MonitorStageRecResult {
	action := dapi.StageActionApply
	// global cluster role
	if err := r.CreateWhenNotExist(transformer.MakeMonitorGlobalClusterRole()); err != nil {
		return mnrStageFail(dapi.MnrOprStageGlobalClusterRole, action, err)
	}
	// namespaced service account
	if err := r.CreateWhenNotExist(transformer.MakeMonitorNamespacedServiceAccount(r.CR.Namespace)); err != nil {
		return mnrStageFail(dapi.MnrOprStageNamespacedServiceAccount, action, err)
	}
	// namespaced role binding
	if err := r.CreateWhenNotExist(transformer.MakeMonitorNamespacedRoleBinding(r.CR.Namespace)); err != nil {
		return mnrStageFail(dapi.MnrOprStageNamespacedRoleBinding, action, err)
	}
	return mnrStageSucc(dapi.MnrOprStageRbac, action)
}

// reconcile rbac resources used by the DorisMonitor
func (r *DorisMonitorReconciler) recPrometheusResources() MonitorStageRecResult {
	action := dapi.StageActionApply
	// config map
	configMap, genConfErr := transformer.MakePrometheusConfigMap(r.CR, r.Schema)
	if genConfErr != nil {
		return mnrStageFail(dapi.MnrOprStagePrometheusConfigMap, action, genConfErr)
	}
	if err := r.CreateOrUpdate(configMap); err != nil {
		return mnrStageFail(dapi.MnrOprStagePrometheusConfigMap, action, err)
	}
	// service
	service := transformer.MakePrometheusService(r.CR, r.Schema)
	if err := r.CreateOrUpdate(service); err != nil {
		return mnrStageFail(dapi.MnrOprStagePrometheusService, action, err)
	}
	// pvc
	pvc := transformer.MakePrometheusPVC(r.CR, r.Schema)
	if err := r.CreateOrUpdate(pvc); err != nil {
		return mnrStageFail(dapi.MnrOprStagePrometheusPVC, action, err)
	}
	// deployment
	deployment := transformer.MakePrometheusDeployment(r.CR, r.Schema)
	deployment.Annotations[PrometheusConfHashAnnotationKey] = util.Md5HashOr(configMap.Data, "")
	if err := r.CreateOrUpdate(deployment); err != nil {
		return mnrStageFail(dapi.MnrOprStagePrometheusDeployment, action, err)
	}
	return mnrStageSucc(dapi.MnrOprStagePrometheus, action)
}

// reconcile grafana resources used by the DorisMonitor
func (r *DorisMonitorReconciler) recGrafanaResources() MonitorStageRecResult {
	action := dapi.StageActionApply
	// config map
	configMap, genConfErr := transformer.MakeGrafanaConfigMap(r.CR, r.Schema)
	if genConfErr != nil {
		return mnrStageFail(dapi.MnrOprStageGrafanaConfigMap, action, genConfErr)
	}
	if err := r.CreateOrUpdate(configMap); err != nil {
		return mnrStageFail(dapi.MnrOprStageGrafanaConfigMap, action, err)
	}
	// secret
	secret := transformer.MakeGrafanaSecret(r.CR, r.Schema)
	if err := r.CreateOrUpdate(secret); err != nil {
		return mnrStageFail(dapi.MnrOprStageGrafanaSecret, action, err)
	}
	// service
	service := transformer.MakeGrafanaService(r.CR, r.Schema)
	if err := r.CreateOrUpdate(service); err != nil {
		return mnrStageFail(dapi.MnrOprStageGrafanaService, action, err)
	}
	// pvc
	pvc := transformer.MakeGrafanaPVC(r.CR, r.Schema)
	if err := r.CreateOrUpdate(pvc); err != nil {
		return mnrStageFail(dapi.MnrOprStageGrafanaPVC, action, err)
	}
	// deployment
	deployment := transformer.MakeGrafanaDeployment(r.CR, r.Schema)
	deployment.Annotations[GrafanaConfHashAnnotationKey] = util.Md5HashOr(configMap.Data, "")
	if err := r.CreateOrUpdate(deployment); err != nil {
		return mnrStageFail(dapi.MnrOprStageGrafanaDeployment, action, err)
	}
	return mnrStageSucc(dapi.MnrOprStageGrafana, action)
}

// reconcile loki resources used by the DorisMonitor
func (r *DorisMonitorReconciler) recLokiResources() MonitorStageRecResult {
	// apply resources
	applyRes := func() MonitorStageRecResult {
		action := dapi.StageActionApply
		// configmap
		configMap, genErr := transformer.MakeLokiConfigMap(r.CR, r.Schema)
		if genErr != nil {
			return mnrStageFail(dapi.MnrOprStageLokiConfigMap, action, genErr)
		}
		if err := r.CreateOrUpdate(configMap); err != nil {
			return mnrStageFail(dapi.MnrOprStageLokiConfigMap, action, err)
		}
		// service
		service := transformer.MakeLokiService(r.CR, r.Schema)
		if err := r.CreateOrUpdate(service); err != nil {
			return mnrStageFail(dapi.MnrOprStageLokiService, action, err)
		}
		// pvc
		pvc := transformer.MakeLokiPVC(r.CR, r.Schema)
		if err := r.CreateOrUpdate(pvc); err != nil {
			return mnrStageFail(dapi.MnrOprStageLokiPVC, action, err)
		}
		// deployment
		deployment := transformer.MakeLokiDeployment(r.CR, r.Schema)
		deployment.Annotations[LokiConfHashAnnotationKey] = util.Md5HashOr(configMap.Data, "")
		if err := r.CreateOrUpdate(deployment); err != nil {
			return mnrStageFail(dapi.MnrOprStageLokiDeployment, action, err)
		}
		return mnrStageSucc(dapi.MnrOprStageLoki, action)
	}

	// delete resources
	deleteRes := func() MonitorStageRecResult {
		action := dapi.StageActionDelete
		// configmap
		configMapRef := transformer.GetLokiServiceKey(r.CR.ObjKey())
		if err := r.DeleteWhenExist(configMapRef, &corev1.ConfigMap{}); err != nil {
			return mnrStageFail(dapi.MnrOprStageLokiConfigMap, action, err)
		}
		// service
		serviceRef := transformer.GetLokiServiceKey(r.CR.ObjKey())
		if err := r.DeleteWhenExist(serviceRef, &corev1.Service{}); err != nil {
			return mnrStageFail(dapi.MnrOprStageLokiService, action, err)
		}
		// pvc
		pvcRef := transformer.GetLokiPVCKey(r.CR.ObjKey())
		if err := r.DeleteWhenExist(pvcRef, &corev1.PersistentVolumeClaim{}); err != nil {
			return mnrStageFail(dapi.MnrOprStageLokiPVC, action, err)
		}
		// deployment
		deploymentRef := transformer.GetLokiDeploymentKey(r.CR.ObjKey())
		if err := r.DeleteWhenExist(deploymentRef, &corev1.PersistentVolumeClaim{}); err != nil {
			return mnrStageFail(dapi.MnrOprStageLokiDeployment, action, err)
		}
		return mnrStageSucc(dapi.MnrOprStageLoki, action)
	}

	return util.Elvis(r.CR.Spec.DisableLoki, deleteRes, applyRes)()
}

// reconcile promtail resources used by the DorisMonitor
func (r *DorisMonitorReconciler) recPromtailResources() MonitorStageRecResult {
	// apply resources
	applyRes := func() MonitorStageRecResult {
		action := dapi.StageActionApply
		// configmap
		configMap, genErr := transformer.MakePromtailConfigMap(r.CR, r.Schema)
		if genErr != nil {
			return mnrStageFail(dapi.MnrOprStagePromtailConfigMap, action, genErr)
		}
		if err := r.CreateOrUpdate(configMap); err != nil {
			return mnrStageFail(dapi.MnrOprStagePromtailConfigMap, action, err)
		}
		// daemonset
		daemonSet := transformer.MakePromtailDaemonSet(r.CR, r.Schema)
		daemonSet.Annotations[PromtailConfHashAnnotationKey] = util.Md5HashOr(configMap.Data, "")
		if err := r.CreateOrUpdate(daemonSet); err != nil {
			return mnrStageFail(dapi.MnrOprStagePromtailDaemonSet, action, err)
		}
		return mnrStageSucc(dapi.MnrOprStagePromtail, action)
	}

	// delete resources
	deleteRes := func() MonitorStageRecResult {
		action := dapi.StageActionDelete
		// configmap
		configMapRef := transformer.GetPromtailConfigMapKey(r.CR.ObjKey())
		if err := r.DeleteWhenExist(configMapRef, &corev1.ConfigMap{}); err != nil {
			return mnrStageFail(dapi.MnrOprStagePromtailConfigMap, action, err)
		}
		// daemonset
		daemonSetRef := transformer.GetPromtailDaemonSetKey(r.CR.ObjKey())
		if err := r.DeleteWhenExist(daemonSetRef, &corev1.ConfigMap{}); err != nil {
			return mnrStageFail(dapi.MnrOprStagePromtailDaemonSet, action, err)
		}
		return mnrStageSucc(dapi.MnrOprStagePromtail, action)
	}

	return util.Elvis(r.CR.Spec.DisableLoki, deleteRes, applyRes)()
}
