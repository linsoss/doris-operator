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
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

var (
	FeConfHashAnnotationKey = fmt.Sprintf("%s/fe-config", dapi.GroupVersion.Group)
)

// DorisClusterReconciler reconciles a DorisCluster object
type DorisClusterReconciler struct {
	ReconcileContext
	CR *dapi.DorisCluster
}

// ClusterStageRecResult represents the result of a stage reconciliation for DorisCluster
type ClusterStageRecResult struct {
	Stage  dapi.DorisClusterOprStage
	Status dapi.OprStageStatus
	Action dapi.OprStageAction
	Err    error
}

func clusterStageSucc(stage dapi.DorisClusterOprStage, action dapi.OprStageAction) ClusterStageRecResult {
	return ClusterStageRecResult{Stage: stage, Status: dapi.StageResultSucceeded, Action: action}
}

func clusterStageFail(stage dapi.DorisClusterOprStage, action dapi.OprStageAction, err error) ClusterStageRecResult {
	return ClusterStageRecResult{Stage: stage, Status: dapi.StageResultSucceeded, Action: action, Err: err}
}

// Reconcile secret object that using to store the sql query account info
// that used by doris-operator.
func (r *DorisClusterReconciler) recOprAccountSecret() ClusterStageRecResult {
	secretRef := transformer.GetOprSqlAccountSecretName(r.CR)
	action := dapi.StageActionApply
	// check if secret exists
	exists, err := r.Exist(secretRef, &corev1.Secret{})
	if err != nil {
		return clusterStageFail(dapi.StageSqlAccountSecret, action, err)
	}
	// create secret if not exists
	if !exists {
		newSecret := transformer.MakeOprSqlAccountSecret(r.CR)
		if err := r.Create(r.Ctx, newSecret); err != nil {
			return clusterStageFail(dapi.StageSqlAccountSecret, action, err)
		}
	}
	return clusterStageSucc(dapi.StageSqlAccountSecret, action)
}

// Reconcile Doris FE component resources.
func (r *DorisClusterReconciler) recFeResources() ClusterStageRecResult {
	if r.CR.Spec.FE != nil {
		// apply resources
		action := dapi.StageActionApply
		// fe configmap
		configMap := transformer.MakeFeConfigMap(r.CR, r.Schema)
		if err := r.CreateOrUpdate(configMap); err != nil {
			return clusterStageFail(dapi.StageFeConfigmap, action, err)
		}
		// fe service
		service := transformer.MakeFeService(r.CR, r.Schema)
		if err := r.CreateOrUpdate(service); err != nil {
			return clusterStageFail(dapi.StageFeService, action, err)
		}
		peerService := transformer.MakeFePeerService(r.CR, r.Schema)
		if err := r.CreateOrUpdate(peerService); err != nil {
			return clusterStageFail(dapi.StageFeService, action, err)
		}
		// fe statefulset
		statefulSet := transformer.MakeFeStatefulSet(r.CR, r.Schema)
		statefulSet.Annotations[FeConfHashAnnotationKey] = util.MapMd5(configMap.Data)
		if err := r.CreateOrUpdate(statefulSet); err != nil {
			return clusterStageFail(dapi.StageFeStatefulSet, action, err)
		}
		return clusterStageSucc(dapi.StageFe, action)
	} else {
		// delete resources
		action := dapi.StageActionDelete
		// fe configmap
		configMapRef := transformer.GetFeConfigMapName(r.CR)
		if err := r.DeleteWhenExist(configMapRef, &corev1.ConfigMap{}); err != nil {
			return clusterStageFail(dapi.StageFeConfigmap, action, err)
		}
		// fe service
		serviceRef := transformer.GetFeServiceName(r.CR)
		if err := r.DeleteWhenExist(serviceRef, &corev1.Service{}); err != nil {
			return clusterStageFail(dapi.StageFeConfigmap, action, err)
		}
		peerServiceRef := transformer.GetFePeerServiceName(r.CR)
		if err := r.DeleteWhenExist(peerServiceRef, &corev1.Service{}); err != nil {
			return clusterStageFail(dapi.StageFeConfigmap, action, err)
		}
		// fe statefulset
		statefulsetRef := transformer.GetFeStatefulSetName(r.CR)
		if err := r.DeleteWhenExist(statefulsetRef, &appv1.StatefulSet{}); err != nil {
			return clusterStageFail(dapi.StageFeConfigmap, action, err)
		}
		return clusterStageSucc(dapi.StageFe, action)
	}
}
