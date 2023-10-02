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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

// Sync all subcomponents status.
func (r *DorisClusterReconciler) Sync() (dapi.DorisClusterSyncStatus, error) {
	syncRes := &dapi.DorisClusterSyncStatus{}
	errCtr := &util.MultiError{}

	type SyncStatus = *dapi.DorisClusterSyncStatus
	type ErrCollector = *util.MultiError
	type MuteFn = func(SyncStatus, ErrCollector)

	// collect fe, be, cn, broker status in parallel
	syncFns := []func() MuteFn{
		func() MuteFn {
			status, err := r.syncFeStatus()
			return func(s SyncStatus, c ErrCollector) {
				s.FE = status
				c.Collect(err)
			}
		},
		func() MuteFn {
			status, err := r.syncBeStatus()
			return func(s SyncStatus, c ErrCollector) {
				s.BE = status
				c.Collect(err)
			}
		},
		func() MuteFn {
			status, err := r.syncCnStatus()
			return func(s SyncStatus, c ErrCollector) {
				s.CN = status
				c.Collect(err)
			}
		},
		func() MuteFn {
			status, err := r.syncBrokerStatus()
			return func(s SyncStatus, c ErrCollector) {
				s.Broker = status
				c.Collect(err)
			}
		},
	}
	muteFnRes := util.ParallelRun(syncFns...)
	for _, muteFn := range muteFnRes {
		muteFn(syncRes, errCtr)
	}

	// eval allReady state
	allReady, err := r.inferIsDorisClusterAllReady()
	if !errCtr.Collect(err) {
		syncRes.AllReady = allReady
	}
	return *syncRes, errCtr.Dry()
}

// infer whether the DorisCluster components are all ready.
func (r *DorisClusterReconciler) inferIsDorisClusterAllReady() (bool, error) {
	if r.CR.Spec.FE != nil {
		if int(r.CR.Spec.FE.Replicas) != len(r.CR.Status.FE.ReadyMembers) {
			return false, nil
		}
	}
	if r.CR.Spec.BE != nil {
		if int(r.CR.Spec.BE.Replicas) != len(r.CR.Status.BE.ReadyMembers) {
			return false, nil
		}
	}
	if r.CR.Spec.Broker != nil {
		if int(r.CR.Spec.Broker.Replicas) != len(r.CR.Status.Broker.ReadyMembers) {
			return false, nil
		}
	}
	if r.CR.Spec.CN != nil {
		autoScale, err := r.FindRefDorisAutoScaler(r.CR.ObjKey())
		if err != nil {
			return false, err
		}
		if autoScale == nil {
			// when not exist DorisAutoScaler
			if int(r.CR.Spec.CN.Replicas) != len(r.CR.Status.CN.ReadyMembers) {
				return false, nil
			}
		} else {
			// when exist DorisAutoScaler
			if len(r.CR.Status.CN.ReadyMembers) < 1 {
				return false, nil
			}
		}
	}
	return true, nil
}

// sync FE status
func (r *DorisClusterReconciler) syncFeStatus() (dapi.FEStatus, error) {
	if r.CR.Spec.FE == nil {
		return dapi.FEStatus{}, nil
	}
	feStatus := util.PointerDeRefer(r.CR.Status.FE.DeepCopy(), dapi.FEStatus{})
	feStatus.ServiceRef = dapi.NewNamespacedName(tran.GetFeServiceKey(r.CR.ObjKey()))
	statefulSetRef := tran.GetFeStatefulSetKey(r.CR.ObjKey())
	image := tran.GetFeImage(r.CR)

	err := r.fillDorisComponentStatus(&feStatus.DorisComponentStatus, statefulSetRef, image)
	return feStatus, err
}

// sync BE status
func (r *DorisClusterReconciler) syncBeStatus() (dapi.BEStatus, error) {
	if r.CR.Spec.BE == nil {
		return dapi.BEStatus{}, nil
	}
	beStatus := util.PointerDeRefer(r.CR.Status.BE.DeepCopy(), dapi.BEStatus{})
	statefulSetRef := tran.GetBeStatefulSetKey(r.CR.ObjKey())
	image := tran.GetBeImage(r.CR)

	err := r.fillDorisComponentStatus(&beStatus.DorisComponentStatus, statefulSetRef, image)
	return beStatus, err
}

// sync CN status
func (r *DorisClusterReconciler) syncCnStatus() (dapi.CNStatus, error) {
	if r.CR.Spec.CN == nil {
		return dapi.CNStatus{}, nil
	}
	cnStatus := util.PointerDeRefer(r.CR.Status.CN.DeepCopy(), dapi.CNStatus{})
	statefulSetRef := tran.GetCnStatefulSetKey(r.CR.ObjKey())
	image := tran.GetCnImage(r.CR)

	err := r.fillDorisComponentStatus(&cnStatus.DorisComponentStatus, statefulSetRef, image)
	return cnStatus, err
}

// sync Broker status
func (r *DorisClusterReconciler) syncBrokerStatus() (dapi.BrokerStatus, error) {
	if r.CR.Spec.Broker == nil {
		return dapi.BrokerStatus{}, nil
	}
	status := util.PointerDeRefer(r.CR.Status.Broker.DeepCopy(), dapi.BrokerStatus{})
	image := tran.GetBrokerImage(r.CR)
	statefulSetRef := tran.GetBrokerStatefulSetKey(r.CR.ObjKey())

	err := r.fillDorisComponentStatus(&status.DorisComponentStatus, statefulSetRef, image)
	return status, err
}

func (r *DorisClusterReconciler) fillDorisComponentStatus(
	baseStatus *dapi.DorisComponentStatus,
	statefulSetKey types.NamespacedName,
	image string) error {

	baseStatus.Image = image
	baseStatus.StatefulSetRef = dapi.NewNamespacedName(statefulSetKey)

	// collect members status via ref statefulset
	sts := &appv1.StatefulSet{}
	if err := r.Find(statefulSetKey, sts); err != nil {
		return err
	}
	if sts != nil {
		baseStatus.Members = r.getComponentMembers(sts)
		baseStatus.Conditions = sts.Status.Conditions
		readyMembers, err := r.getComponentReadyMembers(r.CR.Namespace, tran.GetFeComponentLabels(r.CR.ObjKey()))
		if err != nil {
			return err
		}
		baseStatus.ReadyMembers = readyMembers
	}
	return nil
}

func (r *DorisClusterReconciler) getComponentMembers(sts *appv1.StatefulSet) []string {
	replicas := sts.Status.Replicas
	members := make([]string, replicas)
	for i := 0; i < int(replicas); i++ {
		members[i] = sts.Name + "-" + strconv.Itoa(i) + "." + sts.Namespace
	}
	return members
}

func (r *DorisClusterReconciler) getComponentReadyMembers(namespace string, componentLabels map[string]string) ([]string, error) {
	readyMembers := make([]string, 0)
	podList := &corev1.PodList{}
	listOptions := &client.ListOptions{
		Namespace:     namespace,
		LabelSelector: labels.Set(componentLabels).AsSelector(),
	}
	if err := r.List(r.Ctx, podList, listOptions); err != nil {
		return readyMembers, err
	}
	for _, pod := range podList.Items {
		if util.IsPodReady(pod) {
			readyMembers = append(readyMembers, pod.Name+"."+pod.Namespace)
		}
	}
	return readyMembers, nil
}
