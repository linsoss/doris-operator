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
	cr *dapi.DorisCluster
}

// ClusterStageRecResult represents the result of a stage reconciliation for DorisCluster
type ClusterStageRecResult struct {
	Stage  dapi.DorisClusterOprStage
	Status dapi.DorisClusterOprStageStatus
	Err    error
}

func recSuccess(stage dapi.DorisClusterOprStage) ClusterStageRecResult {
	return ClusterStageRecResult{Stage: stage, Status: dapi.OprStageStatusSucceeded, Err: nil}
}

func recFail(stage dapi.DorisClusterOprStage, err error) ClusterStageRecResult {
	return ClusterStageRecResult{Stage: stage, Status: dapi.OprStageStatusSucceeded, Err: err}
}

// Reconcile secret object that using to store the sql query account info
// that used by doris-operator.
func (r *DorisClusterReconciler) recOprAccountSecret() ClusterStageRecResult {
	secretRef := transformer.GetOprSqlAccountSecretName(r.cr)
	// check if secret exists
	exists, err := r.Exist(secretRef, &corev1.Secret{})
	if err != nil {
		return recFail(dapi.OprStageSqlAccountSecret, err)
	}
	// create secret if not exists
	if !exists {
		newSecret := transformer.MakeOprSqlAccountSecret(r.cr)
		if err := r.Create(r.ctx, newSecret); err != nil {
			return recFail(dapi.OprStageSqlAccountSecret, err)
		}
	}
	return recSuccess(dapi.OprStageSqlAccountSecret)
}

// Reconcile Doris FE component resources.
func (r *DorisClusterReconciler) recFeResources() ClusterStageRecResult {
	if r.cr.Spec.FE != nil {
		// apply resources
		// fe configmap
		configMap := transformer.MakeFeConfigMap(r.cr, r.schema)
		if err := r.CreateOrUpdate(configMap); err != nil {
			return recFail(dapi.OprStageFeConfigmap, err)
		}
		// fe service
		service := transformer.MakeFeService(r.cr, r.schema)
		if err := r.CreateOrUpdate(service); err != nil {
			return recFail(dapi.OprStageFeService, err)
		}
		peerService := transformer.MakeFePeerService(r.cr, r.schema)
		if err := r.CreateOrUpdate(peerService); err != nil {
			return recFail(dapi.OprStageFeService, err)
		}
		// fe statefulset
		statefulSet := transformer.MakeFeStatefulSet(r.cr, r.schema)
		statefulSet.Annotations[FeConfHashAnnotationKey] = util.MapMd5(configMap.Data)
		if err := r.CreateOrUpdate(statefulSet); err != nil {
			return recFail(dapi.OprStageFeStatefulSet, err)
		}
	} else {
		// delete resources
		// fe configmap
		configMapRef := transformer.GetFeConfigMapName(r.cr)
		if err := r.DeleteWhenExist(configMapRef, &corev1.ConfigMap{}); err != nil {
			return recFail(dapi.OprStageFeConfigmap, err)
		}
		// fe service
		serviceRef := transformer.GetFeServiceName(r.cr)
		if err := r.DeleteWhenExist(serviceRef, &corev1.Service{}); err != nil {
			return recFail(dapi.OprStageFeConfigmap, err)
		}
		peerServiceRef := transformer.GetFePeerServiceName(r.cr)
		if err := r.DeleteWhenExist(peerServiceRef, &corev1.Service{}); err != nil {
			return recFail(dapi.OprStageFeConfigmap, err)
		}
		// fe statefulset
		statefulsetRef := transformer.GetFeStatefulSetName(r.cr)
		if err := r.DeleteWhenExist(statefulsetRef, &appv1.StatefulSet{}); err != nil {
			return recFail(dapi.OprStageFeConfigmap, err)
		}
	}
	return recSuccess(dapi.OprStageFe)
}
