/*
Copyright 2023 @ Linying Assad <linying@apache.org>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	dapi "github.com/al-assad/doris-operator/api/v1beta1"
	"github.com/al-assad/doris-operator/internal/reconciler"
	"github.com/al-assad/doris-operator/internal/util"
	acv2 "k8s.io/api/autoscaling/v2"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DorisAutoscalerReconciler reconciles a DorisAutoscaler object
type DorisAutoscalerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=al-assad.github.io,resources=dorisautoscalers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=al-assad.github.io,resources=dorisautoscalers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=al-assad.github.io,resources=dorisautoscalers/finalizers,verbs=update
//+kubebuilder:rbac:groups=autoscaling/v2,resources=horizontalpodautoscalers,verbs=get;list;watch;create;update;patch;delete

func (r *DorisAutoscalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	recCtx := reconciler.NewReconcileContext(r.Client, r.Scheme, ctx)

	// obtain CR
	cr := &dapi.DorisAutoscaler{}
	exist, err := recCtx.Exist(req.NamespacedName, cr)
	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}
	// skip reconciling process when it has been deleted
	if !exist {
		recCtx.Log.Info(fmt.Sprintf("DorisAutoscaler(%s) has been deleted", util.K8sObjKeyStr(req.NamespacedName)))
		return ctrl.Result{}, nil
	}
	rec := reconciler.DorisAutoScalerReconciler{ReconcileContext: recCtx, CR: cr}

	curSpecHash := util.Md5HashOr(cr.Spec, "")
	isFirstCreated := cr.Status.LastApplySpecHash == nil
	specHasChanged := isFirstCreated || *cr.Status.LastApplySpecHash != curSpecHash
	preRecCompleted := cr.Status.CN.Phase == dapi.AutoScalePhaseCompleted

	if isFirstCreated && cr.Status.CN.Phase == "" {
		recCtx.Log.Info(fmt.Sprintf("DorisAutoscaler(%s) is created for the first time", util.K8sObjKeyStr(req.NamespacedName)))
	}
	if specHasChanged {
		recCtx.Log.Info(fmt.Sprintf("DorisAutoscaler(%s) spec has been updated", util.K8sObjKeyStr(req.NamespacedName)))
	}

	// reconcile the sub resources
	var recErr error
	if isFirstCreated || specHasChanged || !preRecCompleted {
		recRs, err := rec.Reconcile()
		recErr = err
		cr.Status.CN.AutoscalerRecStatus = recRs
		// when reconcile process competed success, update the last apply spec hash
		if err == nil {
			cr.Status.LastApplySpecHash = &curSpecHash
		}
	}
	// sync the status of CR
	syncRs, syncErr := rec.Sync()
	cr.Status.CN.CNAutoscalerSyncStatus = syncRs
	// update the status of CR
	updateErr := r.Status().Update(ctx, cr)

	// merged error as result
	errSet := StCtrlErrSet{
		Rec:    recErr,
		Sync:   syncErr,
		Update: updateErr,
	}
	return errSet.AsResult()
}

// SetupWithManager sets up the controller with the Manager.
func (r *DorisAutoscalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dapi.DorisAutoscaler{}).
		Owns(&acv2.HorizontalPodAutoscaler{}).
		Complete(r)
}
