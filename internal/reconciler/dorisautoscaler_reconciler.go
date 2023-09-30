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
	acv2 "k8s.io/api/autoscaling/v2"
	"k8s.io/apimachinery/pkg/types"
)

// DorisAutoScalerReconciler reconciles a DorisCluster object
type DorisAutoScalerReconciler struct {
	ReconcileContext
	CR *dapi.DorisAutoscaler
}

type hpaType = acv2.HorizontalPodAutoscaler

// Reconcile hpa resources
func (r *DorisAutoScalerReconciler) Reconcile() (dapi.AutoscalerRecStatus, error) {
	if r.CR.Spec.Cluster == "" {
		return dapi.AutoscalerRecStatus{}, nil
	}

	// delete cn auto scaler
	deleteHpa := func() error {
		if err := r.DeleteWhenExist(transformer.GetCnScaleUpHpaKey(r.CR.ObjKey()), &hpaType{}); err != nil {
			return err
		}
		if err := r.DeleteWhenExist(transformer.GetCnScaleDownHpaKey(r.CR.ObjKey()), &hpaType{}); err != nil {
			return err
		}
		return nil
	}

	// apply cn auto scaler
	applyHpa := func() error {
		clusterRef := types.NamespacedName{
			Namespace: r.CR.Namespace,
			Name:      r.CR.Spec.Cluster,
		}
		// check if target DorisCluster exist
		exist, err := r.Exist(clusterRef, &dapi.DorisCluster{})
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("target DorisCluster[name=%s][namespace=%s] not exist",
				clusterRef.Name, clusterRef.Name)
		}
		// check if target DorisCluster already bound another DorisAutoscaler
		bound, bErr := r.FindRefDorisAutoScaler(clusterRef)
		if bErr != nil {
			return err
		}
		if bound != nil {
			return fmt.Errorf("target DorisCluster already bound another DorisAutoscaler[name=%s][namespace=%s]",
				bound.Name, bound.Name)
		}
		// apply hpa resources
		if cnUpHpa := transformer.MakeCnScaleUpHpa(r.CR, r.Schema); cnUpHpa != nil {
			if err := r.CreateOrUpdate(cnUpHpa); err != nil {
				return err
			}
		}
		if cnDownHpa := transformer.MakeCnScaleDownHpa(r.CR, r.Schema); cnDownHpa != nil {
			if err := r.CreateOrUpdate(cnDownHpa); err != nil {
				return err
			}
		}
		return nil
	}

	if err := util.Elvis(r.CR.Spec.CN != nil, applyHpa, deleteHpa)(); err != nil {
		return dapi.AutoscalerRecStatus{
			Phase:   dapi.AutoScalePhaseFailed,
			Message: err.Error(),
		}, err
	}
	return dapi.AutoscalerRecStatus{Phase: dapi.AutoScalePhaseCompleted}, nil
}

// Sync status of hpa resources
func (r *DorisAutoScalerReconciler) Sync() (dapi.CNAutoscalerSyncStatus, error) {
	status := util.PointerDeRefer(r.CR.Status.CN.CNAutoscalerSyncStatus.DeepCopy(), dapi.CNAutoscalerSyncStatus{})

	syncCnUpHpa := func() error {
		hpaRef := transformer.GetCnScaleUpHpaKey(r.CR.ObjKey())
		hpa := &hpaType{}
		if err := r.Find(hpaRef, hpa); err != nil {
			return nil
		}
		if hpa != nil {
			status.ScaleUpHpaRef = &dapi.AutoScalerRef{
				NamespacedName: dapi.NewNamespacedName(hpaRef),
				TypeMeta:       hpa.TypeMeta,
			}
			status.ScaleUpStatus = &hpa.Status
		}
		return nil
	}

	syncCnDownHpa := func() error {
		hpaRef := transformer.GetCnScaleDownHpaKey(r.CR.ObjKey())
		hpa := &hpaType{}
		if err := r.Find(hpaRef, hpa); err != nil {
			return nil
		}
		if hpa != nil {
			status.ScaleDownHpaRef = &dapi.AutoScalerRef{
				NamespacedName: dapi.NewNamespacedName(hpaRef),
				TypeMeta:       hpa.TypeMeta,
			}
			status.ScaleDownStatus = &hpa.Status
		}
		return nil
	}

	cnUpErr := syncCnUpHpa()
	cnDownErr := syncCnDownHpa()
	return status, util.MergeErrors(cnUpErr, cnDownErr)
}
