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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

// Sync all sub components status.
func (r *DorisClusterReconciler) Sync() (dapi.DorisClusterSyncStatus, error) {
	syncRes := dapi.DorisClusterSyncStatus{}
	errCtr := util.MultiError{}
	// todo modified to execute code in parallel
	sync(r.syncFeStatus, errCtr, func(s dapi.FEStatus) { syncRes.FE = s })
	sync(r.syncBeStatus, errCtr, func(s dapi.BEStatus) { syncRes.BE = s })
	sync(r.syncCnStatus, errCtr, func(s dapi.CNStatus) { syncRes.CN = s })
	sync(r.syncBrokerStatus, errCtr, func(s dapi.BrokerStatus) { syncRes.Broker = s })
	return syncRes, errCtr.Dry()
}

func sync[T any](syncFn func() (T, error), errCtr util.MultiError, rightFn func(t T)) {
	status, err := syncFn()
	if err != nil {
		errCtr.Collect(err)
	}
	rightFn(status)
}

// sync FE status
func (r *DorisClusterReconciler) syncFeStatus() (dapi.FEStatus, error) {
	if r.CR.Spec.FE == nil {
		return dapi.FEStatus{}, nil
	}
	feStatus := util.PointerDeRefer(r.CR.Status.FE.DeepCopy(), dapi.FEStatus{})
	feStatus.Image = transformer.GetFeImage(r.CR)
	feStatus.ServiceRef = dapi.NewNamespacedName(transformer.GetFeServiceKey(r.CR.ObjKey()))
	statefulSetRef := transformer.GetFeStatefulSetKey(r.CR.ObjKey())
	feStatus.StatefulSetRef = dapi.NewNamespacedName(statefulSetRef)

	// collect members status via ref statefulset
	sts := &appv1.StatefulSet{}
	if err := r.Find(statefulSetRef, sts); err != nil {
		return feStatus, err
	}
	if sts != nil {
		feStatus.Members = r.getComponentMembers(sts)
		feStatus.Conditions = sts.Status.Conditions
		readyMembers, err := r.getComponentReadyMembers(r.CR.Namespace, transformer.GetFeComponentLabels(r.CR.ObjKey()))
		if err != nil {
			return feStatus, err
		}
		feStatus.ReadyMembers = readyMembers
	}
	return feStatus, nil
}

// sync BE status
func (r *DorisClusterReconciler) syncBeStatus() (dapi.BEStatus, error) {
	if r.CR.Spec.BE == nil {
		return dapi.BEStatus{}, nil
	}
	beStatus := util.PointerDeRefer(r.CR.Status.BE.DeepCopy(), dapi.BEStatus{})
	beStatus.Image = transformer.GetBeImage(r.CR)
	statefulSetRef := transformer.GetBeStatefulSetKey(r.CR)
	beStatus.StatefulSetRef = dapi.NewNamespacedName(statefulSetRef)
	// collect members status via ref statefulset
	sts := &appv1.StatefulSet{}
	if err := r.Find(statefulSetRef, sts); err != nil {
		return beStatus, err
	}
	if sts != nil {
		beStatus.Members = r.getComponentMembers(sts)
		beStatus.Conditions = sts.Status.Conditions
		readyMembers, err := r.getComponentReadyMembers(r.CR.Namespace, transformer.GetBeComponentLabels(r.CR.ObjKey()))
		if err != nil {
			return beStatus, err
		}
		beStatus.ReadyMembers = readyMembers
	}
	return beStatus, nil
}

// sync CN status
func (r *DorisClusterReconciler) syncCnStatus() (dapi.CNStatus, error) {
	if r.CR.Spec.CN == nil {
		return dapi.CNStatus{}, nil
	}
	cnStatus := util.PointerDeRefer(r.CR.Status.CN.DeepCopy(), dapi.CNStatus{})
	cnStatus.Image = transformer.GetCnImage(r.CR)
	statefulSetRef := transformer.GetCnStatefulSetKey(r.CR.ObjKey())
	cnStatus.StatefulSetRef = dapi.NewNamespacedName(statefulSetRef)
	// collect members status via ref statefulset
	sts := &appv1.StatefulSet{}
	if err := r.Find(statefulSetRef, sts); err != nil {
		return cnStatus, err
	}
	if sts != nil {
		cnStatus.Members = r.getComponentMembers(sts)
		cnStatus.Conditions = sts.Status.Conditions
		readyMembers, err := r.getComponentReadyMembers(r.CR.Namespace, transformer.GetCnComponentLabels(r.CR.ObjKey()))
		if err != nil {
			return cnStatus, err
		}
		cnStatus.ReadyMembers = readyMembers
	}
	return cnStatus, nil
}

// sync Broker status
func (r *DorisClusterReconciler) syncBrokerStatus() (dapi.BrokerStatus, error) {
	if r.CR.Spec.Broker == nil {
		return dapi.BrokerStatus{}, nil
	}
	status := util.PointerDeRefer(r.CR.Status.Broker.DeepCopy(), dapi.BrokerStatus{})
	status.Image = transformer.GetBrokerImage(r.CR)
	statefulSetRef := transformer.GetBrokerStatefulSetKey(r.CR.ObjKey())
	status.StatefulSetRef = dapi.NewNamespacedName(statefulSetRef)
	// collect members status via ref statefulset
	sts := &appv1.StatefulSet{}
	if err := r.Find(statefulSetRef, sts); err != nil {
		return status, err
	}
	if sts != nil {
		status.Members = r.getComponentMembers(sts)
		status.Conditions = sts.Status.Conditions
		readyMembers, err := r.getComponentReadyMembers(r.CR.Namespace, transformer.GetBrokerComponentLabels(r.CR.ObjKey()))
		if err != nil {
			return status, err
		}
		status.ReadyMembers = readyMembers
	}
	return status, nil
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
