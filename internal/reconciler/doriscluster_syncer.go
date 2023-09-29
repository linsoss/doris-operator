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

// DorisClusterSyncer sync the status of a DorisCluster object
type DorisClusterSyncer struct {
	ReconcileContext
	CR *dapi.DorisCluster
}

// Sync all sub components status.
func (r *DorisClusterSyncer) Sync() error {
	syncs := []func(cr *dapi.DorisCluster) error{
		r.syncFeStatus,
		r.syncBeStatus,
		r.syncCnStatus,
		r.syncBrokerStatus,
	}
	errs := make([]error, 0)
	for _, sync := range syncs {
		err := sync(r.CR)
		errs = append(errs, err)
		return util.MergeErrors(errs...)
	}
	return nil
}

// sync FE status
func (r *DorisClusterSyncer) syncFeStatus(cr *dapi.DorisCluster) error {
	if cr.Spec.FE == nil && cr.Status.Stage == dapi.StageComplete {
		return nil
	}
	feStatus := dapi.FEStatus{}
	feStatus.Image = transformer.GetFeImage(cr)
	feStatus.ServiceRef = dapi.NewNamespacedName(transformer.GetFeServiceKey(cr.ObjKey()))
	statefulSetRef := transformer.GetFeStatefulSetKey(cr.ObjKey())
	feStatus.StatefulSetRef = dapi.NewNamespacedName(statefulSetRef)

	// collect members status via ref statefulset
	sts := &appv1.StatefulSet{}
	if err := r.Find(statefulSetRef, sts); err != nil {
		return err
	}
	if sts != nil {
		feStatus.Members = r.getComponentMembers(sts)
		feStatus.Conditions = sts.Status.Conditions
		readyMembers, err := r.getComponentReadyMembers(cr.Namespace, transformer.GetFeComponentLabels(cr.ObjKey()))
		if err != nil {
			return err
		}
		feStatus.ReadyMembers = readyMembers
	}
	cr.Status.FE = feStatus
	return nil
}

// sync BE status
func (r *DorisClusterSyncer) syncBeStatus(cr *dapi.DorisCluster) error {
	if cr.Spec.BE == nil && cr.Status.Stage == dapi.StageComplete {
		return nil
	}
	beStatus := dapi.BEStatus{}
	beStatus.Image = transformer.GetBeImage(cr)
	statefulSetRef := transformer.GetBeStatefulSetKey(cr)
	beStatus.StatefulSetRef = dapi.NewNamespacedName(statefulSetRef)
	// collect members status via ref statefulset
	sts := &appv1.StatefulSet{}
	if err := r.Find(statefulSetRef, sts); err != nil {
		return err
	}
	if sts != nil {
		beStatus.Members = r.getComponentMembers(sts)
		beStatus.Conditions = sts.Status.Conditions
		readyMembers, err := r.getComponentReadyMembers(cr.Namespace, transformer.GetBeComponentLabels(cr.ObjKey()))
		if err != nil {
			return err
		}
		beStatus.ReadyMembers = readyMembers
	}
	cr.Status.BE = beStatus
	return nil
}

// sync CN status
func (r *DorisClusterSyncer) syncCnStatus(cr *dapi.DorisCluster) error {
	if cr.Spec.CN == nil && cr.Status.Stage == dapi.StageComplete {
		return nil
	}
	cnStatus := dapi.CNStatus{}
	cnStatus.Image = transformer.GetCnImage(cr)
	statefulSetRef := transformer.GetCnStatefulSetKey(cr.ObjKey())
	cnStatus.StatefulSetRef = dapi.NewNamespacedName(statefulSetRef)
	// collect members status via ref statefulset
	sts := &appv1.StatefulSet{}
	if err := r.Find(statefulSetRef, sts); err != nil {
		return err
	}
	if sts != nil {
		cnStatus.Members = r.getComponentMembers(sts)
		cnStatus.Conditions = sts.Status.Conditions
		readyMembers, err := r.getComponentReadyMembers(cr.Namespace, transformer.GetCnComponentLabels(cr.ObjKey()))
		if err != nil {
			return err
		}
		cnStatus.ReadyMembers = readyMembers
	}
	cr.Status.CN = cnStatus
	return nil
}

// sync Broker status
func (r *DorisClusterSyncer) syncBrokerStatus(cr *dapi.DorisCluster) error {
	if cr.Spec.Broker == nil && cr.Status.Stage == dapi.StageComplete {
		return nil
	}
	status := dapi.BrokerStatus{}
	status.Image = transformer.GetBrokerImage(cr)
	statefulSetRef := transformer.GetBrokerStatefulSetKey(cr.ObjKey())
	status.StatefulSetRef = dapi.NewNamespacedName(statefulSetRef)
	// collect members status via ref statefulset
	sts := &appv1.StatefulSet{}
	if err := r.Find(statefulSetRef, sts); err != nil {
		return err
	}
	if sts != nil {
		status.Members = r.getComponentMembers(sts)
		status.Conditions = sts.Status.Conditions
		readyMembers, err := r.getComponentReadyMembers(cr.Namespace, transformer.GetBrokerComponentLabels(cr.ObjKey()))
		if err != nil {
			return err
		}
		status.ReadyMembers = readyMembers
	}
	cr.Status.Broker = status
	return nil
}

func (r *DorisClusterSyncer) getComponentMembers(sts *appv1.StatefulSet) []string {
	replicas := sts.Status.Replicas
	members := make([]string, replicas)
	for i := 0; i < int(replicas); i++ {
		members[i] = sts.Name + "-" + strconv.Itoa(i) + "." + sts.Namespace
	}
	return members
}

func (r *DorisClusterSyncer) getComponentReadyMembers(namespace string, componentLabels map[string]string) ([]string, error) {
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
