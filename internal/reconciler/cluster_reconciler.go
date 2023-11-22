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
	tran "github.com/al-assad/doris-operator/internal/transformer"
	"github.com/al-assad/doris-operator/internal/util"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	FeConfHashAnnotationKey     = fmt.Sprintf("%s/fe-config", dapi.GroupVersion.Group)
	BeConfHashAnnotationKey     = fmt.Sprintf("%s/be-config", dapi.GroupVersion.Group)
	CnConfHashAnnotationKey     = fmt.Sprintf("%s/cn-config", dapi.GroupVersion.Group)
	BrokerConfHashAnnotationKey = fmt.Sprintf("%s/broker-config", dapi.GroupVersion.Group)
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

// Reconcile all sub components
func (r *DorisClusterReconciler) Reconcile() ClusterStageRecResult {
	stages := []func() ClusterStageRecResult{
		r.recOprAccountSecret,
		r.recFeResources,
		r.recBeResources,
		r.recCnResources,
		r.recBrokerResources,
	}
	for _, fn := range stages {
		result := fn()
		if result.Err != nil {
			return result
		}
	}
	return ClusterStageRecResult{Stage: dapi.StageComplete, Status: dapi.StageResultSucceeded}
}

func (r *ClusterStageRecResult) AsDorisClusterRecStatus() dapi.DorisClusterRecStatus {
	res := dapi.DorisClusterRecStatus{
		Stage:       r.Stage,
		StageStatus: r.Status,
		StageAction: r.Action,
	}
	if r.Err != nil {
		res.LastMessage = r.Err.Error()
	}
	return res
}

func clusterStageSucc(stage dapi.DorisClusterOprStage, action dapi.OprStageAction) ClusterStageRecResult {
	return ClusterStageRecResult{Stage: stage, Status: dapi.StageResultSucceeded, Action: action}
}

func clusterStageFail(stage dapi.DorisClusterOprStage, action dapi.OprStageAction, err error) ClusterStageRecResult {
	return ClusterStageRecResult{Stage: stage, Status: dapi.StageResultFailed, Action: action, Err: err}
}

// reconcile secret object that using to store the sql query account info
// that used by doris-operator.
func (r *DorisClusterReconciler) recOprAccountSecret() ClusterStageRecResult {
	action := dapi.StageActionApply
	// create secret if not exists
	secret := tran.MakeOprSqlAccountSecret(r.CR)
	if err := r.CreateWhenNotExist(secret, &corev1.Secret{}); err != nil {
		return clusterStageFail(dapi.StageSqlAccountSecret, action, err)
	}
	return clusterStageSucc(dapi.StageSqlAccountSecret, action)
}

// reconcile Doris FE component resources.
func (r *DorisClusterReconciler) recFeResources() ClusterStageRecResult {

	// apply resources
	applyRes := func() ClusterStageRecResult {
		action := dapi.StageActionApply
		// fe configmap
		configMap := tran.MakeFeConfigMap(r.CR, r.Schema)
		if err := r.CreateOrUpdate(configMap, &corev1.ConfigMap{}); err != nil {
			return clusterStageFail(dapi.StageFeConfigmap, action, err)
		}
		// fe service
		service := tran.MakeFeService(r.CR, r.Schema)
		if err := r.CreateOrUpdate(service, &corev1.Service{}); err != nil {
			return clusterStageFail(dapi.StageFeService, action, err)
		}
		peerService := tran.MakeFePeerService(r.CR, r.Schema)
		if err := r.CreateOrUpdate(peerService, &corev1.Service{}); err != nil {
			return clusterStageFail(dapi.StageFeService, action, err)
		}
		// fe statefulset
		statefulSet := tran.MakeFeStatefulSet(r.CR, r.Schema)
		statefulSet.Spec.Template.Annotations[FeConfHashAnnotationKey] = util.Md5HashOr(configMap.Data, "")
		if err := r.CreateOrUpdate(statefulSet, &appv1.StatefulSet{}); err != nil {
			return clusterStageFail(dapi.StageFeStatefulSet, action, err)
		}
		return clusterStageSucc(dapi.StageFe, action)
	}

	// delete resources
	deleteRes := func() ClusterStageRecResult {
		action := dapi.StageActionDelete
		// fe statefulset
		statefulsetRef := tran.GetFeStatefulSetKey(r.CR.ObjKey())
		if err := r.DeleteWhenExist(statefulsetRef, &appv1.StatefulSet{}); err != nil {
			return clusterStageFail(dapi.StageFeStatefulSet, action, err)
		}
		// fe service
		serviceRef := tran.GetFeServiceKey(r.CR.ObjKey())
		if err := r.DeleteWhenExist(serviceRef, &corev1.Service{}); err != nil {
			return clusterStageFail(dapi.StageFeService, action, err)
		}
		peerServiceRef := tran.GetFePeerServiceKey(r.CR.ObjKey())
		if err := r.DeleteWhenExist(peerServiceRef, &corev1.Service{}); err != nil {
			return clusterStageFail(dapi.StageFeService, action, err)
		}
		// fe configmap
		configMapRef := tran.GetFeConfigMapKey(r.CR.ObjKey())
		if err := r.DeleteWhenExist(configMapRef, &corev1.ConfigMap{}); err != nil {
			return clusterStageFail(dapi.StageFeConfigmap, action, err)
		}
		return clusterStageSucc(dapi.StageFe, action)
	}

	return util.Elvis(r.CR.Spec.FE != nil, applyRes, deleteRes)()
}

// reconcile Doris BE component resources.
func (r *DorisClusterReconciler) recBeResources() ClusterStageRecResult {

	// apply resources
	applyRes := func() ClusterStageRecResult {
		action := dapi.StageActionApply
		// be configmap
		configMap := tran.MakeBeConfigMap(r.CR, r.Schema)
		if err := r.CreateOrUpdate(configMap, &corev1.ConfigMap{}); err != nil {
			return clusterStageFail(dapi.StageBeConfigmap, action, err)
		}
		// be service
		service := tran.MakeBeService(r.CR, r.Schema)
		if err := r.CreateOrUpdate(service, &corev1.Service{}); err != nil {
			return clusterStageFail(dapi.StageBeService, action, err)
		}
		peerService := tran.MakeBePeerService(r.CR, r.Schema)
		if err := r.CreateOrUpdate(peerService, &corev1.Service{}); err != nil {
			return clusterStageFail(dapi.StageBeService, action, err)
		}
		// be statefulset
		statefulSet := tran.MakeBeStatefulSet(r.CR, r.Schema)
		statefulSet.Spec.Template.Annotations[BeConfHashAnnotationKey] = util.Md5HashOr(configMap.Data, "")
		if err := r.CreateOrUpdate(statefulSet, &appv1.StatefulSet{}); err != nil {
			return clusterStageFail(dapi.StageBeStatefulSet, action, err)
		}
		return clusterStageSucc(dapi.StageBe, action)
	}

	// delete resources
	deleteRes := func() ClusterStageRecResult {
		action := dapi.StageActionDelete
		// be statefulset
		statefulsetRef := tran.GetBeStatefulSetKey(r.CR.ObjKey())
		if err := r.DeleteWhenExist(statefulsetRef, &appv1.StatefulSet{}); err != nil {
			return clusterStageFail(dapi.StageBeStatefulSet, action, err)
		}
		// be service
		serviceRef := tran.GetBeServiceKey(r.CR.ObjKey())
		if err := r.DeleteWhenExist(serviceRef, &corev1.Service{}); err != nil {
			return clusterStageFail(dapi.StageBeService, action, err)
		}
		peerServiceRef := tran.GetBePeerServiceKey(r.CR.ObjKey())
		if err := r.DeleteWhenExist(peerServiceRef, &corev1.Service{}); err != nil {
			return clusterStageFail(dapi.StageBeService, action, err)
		}
		// be configmap
		configMapRef := tran.GetBeConfigMapKey(r.CR.ObjKey())
		if err := r.DeleteWhenExist(configMapRef, &corev1.ConfigMap{}); err != nil {
			return clusterStageFail(dapi.StageBeConfigmap, action, err)
		}
		return clusterStageSucc(dapi.StageBe, action)
	}

	return util.Elvis(r.CR.Spec.BE != nil, applyRes, deleteRes)()
}

// reconcile Doris CN component resources.
func (r *DorisClusterReconciler) recCnResources() ClusterStageRecResult {

	// apply resources
	applyRes := func() ClusterStageRecResult {
		action := dapi.StageActionApply
		// cn configmap
		configMap := tran.MakeCnConfigMap(r.CR, r.Schema)
		if err := r.CreateOrUpdate(configMap, &corev1.ConfigMap{}); err != nil {
			return clusterStageFail(dapi.StageCnConfigmap, action, err)
		}
		// cn service
		service := tran.MakeCnService(r.CR, r.Schema)
		if err := r.CreateOrUpdate(service, &corev1.Service{}); err != nil {
			return clusterStageFail(dapi.StageCnService, action, err)
		}
		peerService := tran.MakeCnPeerService(r.CR, r.Schema)
		if err := r.CreateOrUpdate(peerService, &corev1.Service{}); err != nil {
			return clusterStageFail(dapi.StageCnService, action, err)
		}

		// cn statefulset
		statefulSet := tran.MakeCnStatefulSet(r.CR, r.Schema)
		statefulSet.Spec.Template.Annotations[CnConfHashAnnotationKey] = util.Md5HashOr(configMap.Data, "")
		// when the corresponding DorisAutoScaler resource exists,
		// the replica of statefulset would not be overridden
		autoScaler, err := r.FindRefDorisAutoScaler(client.ObjectKeyFromObject(r.CR))
		if err != nil {
			return clusterStageFail(dapi.StageCnStatefulSet, action, err)
		}
		if autoScaler != nil {
			statefulSet.Spec.Replicas = nil
		}
		if err := r.CreateOrUpdate(statefulSet, &appv1.StatefulSet{}); err != nil {
			return clusterStageFail(dapi.StageCnStatefulSet, action, err)
		}

		return clusterStageSucc(dapi.StageCn, action)
	}

	// delete resources
	deleteRes := func() ClusterStageRecResult {
		action := dapi.StageActionDelete
		// cn statefulset
		statefulsetRef := tran.GetCnStatefulSetKey(r.CR.ObjKey())
		if err := r.DeleteWhenExist(statefulsetRef, &appv1.StatefulSet{}); err != nil {
			return clusterStageFail(dapi.StageCnStatefulSet, action, err)
		}
		// cn service
		serviceRef := tran.GetCnServiceKey(r.CR.ObjKey())
		if err := r.DeleteWhenExist(serviceRef, &corev1.Service{}); err != nil {
			return clusterStageFail(dapi.StageCnService, action, err)
		}
		peerServiceRef := tran.GetCnPeerServiceKey(r.CR.ObjKey())
		if err := r.DeleteWhenExist(peerServiceRef, &corev1.Service{}); err != nil {
			return clusterStageFail(dapi.StageCnService, action, err)
		}
		// cn configmap
		configMapRef := tran.GetCnConfigMapKey(r.CR.ObjKey())
		if err := r.DeleteWhenExist(configMapRef, &corev1.ConfigMap{}); err != nil {
			return clusterStageFail(dapi.StageCnConfigmap, action, err)
		}
		return clusterStageSucc(dapi.StageCn, action)
	}

	return util.Elvis(r.CR.Spec.CN != nil, applyRes, deleteRes)()
}

// Reconcile Doris Broker component resources.
func (r *DorisClusterReconciler) recBrokerResources() ClusterStageRecResult {

	// apply resources
	applyRes := func() ClusterStageRecResult {
		action := dapi.StageActionApply
		// broker configmap
		configMap := tran.MakeBrokerConfigMap(r.CR, r.Schema)
		if err := r.CreateOrUpdate(configMap, &corev1.ConfigMap{}); err != nil {
			return clusterStageFail(dapi.StageBrokerConfigmap, action, err)
		}
		// broker service
		peerService := tran.MakeBrokerPeerService(r.CR, r.Schema)
		if err := r.CreateOrUpdate(peerService, &corev1.Service{}); err != nil {
			return clusterStageFail(dapi.StageBrokerService, action, err)
		}
		// broker statefulset
		statefulSet := tran.MakeBrokerStatefulSet(r.CR, r.Schema)
		statefulSet.Spec.Template.Annotations[BrokerConfHashAnnotationKey] = util.Md5HashOr(configMap.Data, "")
		if err := r.CreateOrUpdate(statefulSet, &appv1.StatefulSet{}); err != nil {
			return clusterStageFail(dapi.StageBrokerStatefulSet, action, err)
		}
		return clusterStageSucc(dapi.StageBroker, action)
	}

	// delete resources
	deleteRes := func() ClusterStageRecResult {
		action := dapi.StageActionDelete
		// broker statefulset
		statefulsetRef := tran.GetBrokerStatefulSetKey(r.CR.ObjKey())
		if err := r.DeleteWhenExist(statefulsetRef, &appv1.StatefulSet{}); err != nil {
			return clusterStageFail(dapi.StageBrokerStatefulSet, action, err)
		}
		// broker service
		peerServiceRef := tran.GetBrokerPeerServiceKey(r.CR.ObjKey())
		if err := r.DeleteWhenExist(peerServiceRef, &corev1.Service{}); err != nil {
			return clusterStageFail(dapi.StageBrokerService, action, err)
		}
		// broker configmap
		configMapRef := tran.GetBrokerConfigMapKey(r.CR.ObjKey())
		if err := r.DeleteWhenExist(configMapRef, &corev1.ConfigMap{}); err != nil {
			return clusterStageFail(dapi.StageBrokerConfigmap, action, err)
		}
		return clusterStageSucc(dapi.StageBroker, action)
	}

	return util.Elvis(r.CR.Spec.Broker != nil, applyRes, deleteRes)()
}
