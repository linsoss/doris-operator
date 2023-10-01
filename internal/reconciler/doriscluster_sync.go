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
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

// Sync all subcomponents status.
func (r *DorisClusterReconciler) Sync() (dapi.DorisClusterSyncStatus, error) {
	syncRes := dapi.DorisClusterSyncStatus{}
	errCtr := &util.MultiError{}

	// todo modified to execute code in parallel
	util.CollectFnErr(errCtr, r.syncFeStatus, func(s dapi.FEStatus) { syncRes.FE = s })
	util.CollectFnErr(errCtr, r.syncBeStatus, func(s dapi.BEStatus) { syncRes.BE = s })
	util.CollectFnErr(errCtr, r.syncCnStatus, func(s dapi.CNStatus) { syncRes.CN = s })
	util.CollectFnErr(errCtr, r.syncBrokerStatus, func(s dapi.BrokerStatus) { syncRes.Broker = s })

	return syncRes, errCtr.Dry()
}

// sync FE status
func (r *DorisClusterReconciler) syncFeStatus() (dapi.FEStatus, error) {
	if r.CR.Spec.FE == nil {
		return dapi.FEStatus{}, nil
	}
	feStatus := util.PointerDeRefer(r.CR.Status.FE.DeepCopy(), dapi.FEStatus{})
	feStatus.ServiceRef = dapi.NewNamespacedName(transformer.GetFeServiceKey(r.CR.ObjKey()))
	statefulSetRef := transformer.GetFeStatefulSetKey(r.CR.ObjKey())
	image := transformer.GetFeImage(r.CR)

	err := r.fillDorisComponentStatus(&feStatus.DorisComponentStatus, statefulSetRef, image)
	return feStatus, err
}

// sync BE status
func (r *DorisClusterReconciler) syncBeStatus() (dapi.BEStatus, error) {
	if r.CR.Spec.BE == nil {
		return dapi.BEStatus{}, nil
	}
	beStatus := util.PointerDeRefer(r.CR.Status.BE.DeepCopy(), dapi.BEStatus{})
	statefulSetRef := transformer.GetBeStatefulSetKey(r.CR)
	image := transformer.GetBeImage(r.CR)

	err := r.fillDorisComponentStatus(&beStatus.DorisComponentStatus, statefulSetRef, image)
	return beStatus, err
}

// sync CN status
func (r *DorisClusterReconciler) syncCnStatus() (dapi.CNStatus, error) {
	if r.CR.Spec.CN == nil {
		return dapi.CNStatus{}, nil
	}
	cnStatus := util.PointerDeRefer(r.CR.Status.CN.DeepCopy(), dapi.CNStatus{})
	statefulSetRef := transformer.GetCnStatefulSetKey(r.CR.ObjKey())
	image := transformer.GetCnImage(r.CR)

	err := r.fillDorisComponentStatus(&cnStatus.DorisComponentStatus, statefulSetRef, image)
	return cnStatus, err
}

// sync Broker status
func (r *DorisClusterReconciler) syncBrokerStatus() (dapi.BrokerStatus, error) {
	if r.CR.Spec.Broker == nil {
		return dapi.BrokerStatus{}, nil
	}
	status := util.PointerDeRefer(r.CR.Status.Broker.DeepCopy(), dapi.BrokerStatus{})
	image := transformer.GetBrokerImage(r.CR)
	statefulSetRef := transformer.GetBrokerStatefulSetKey(r.CR.ObjKey())

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
		readyMembers, err := r.getComponentReadyMembers(r.CR.Namespace, transformer.GetFeComponentLabels(r.CR.ObjKey()))
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
