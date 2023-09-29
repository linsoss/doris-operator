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
	"github.com/al-assad/doris-operator/internal/reconciler"
	"github.com/al-assad/doris-operator/internal/util"
	acv2 "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"

	dapi "github.com/al-assad/doris-operator/api/v1beta1"
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

	// obtain CR and skip reconciling process when it has been deleted
	cr := &dapi.DorisAutoscaler{}
	if err := recCtx.Find(req.NamespacedName, cr); err != nil {
		return ctrl.Result{Requeue: true}, err
	}
	if cr == nil {
		recCtx.Log.Info("DorisAutoscaler has been deleted")
		return ctrl.Result{}, nil
	}
	rec := reconciler.DorisAutoScalerReconciler{ReconcileContext: recCtx, CR: cr}

	// reconcile the sub resource of CR when the spec of it has been changed,
	// or the previous reconcile phase has not been completed.
	var recErr error
	if cr.Status.PrevSpec == nil || !reflect.DeepEqual(cr.Spec, *cr.Status.PrevSpec) || cr.Status.CN.Phase != dapi.AutoScalePhaseCompleted {
		recErr = rec.Reconcile()
		if recErr == nil {
			cr.Status.CN.Phase = dapi.AutoScalePhaseCompleted
		} else {
			cr.Status.CN.LastMessage = recErr.Error()
			// recovery the phase of CR
			switch recErr.(type) {
			case *reconciler.DorisClusterNotExist:
				cr.Status.CN.Phase = dapi.AutoScalePhasePending
			case *reconciler.ClusterAlreadyBoundAutoscaler:
				cr.Status.CN.Phase = dapi.AutoScalePhasePending
			default:
				cr.Status.CN.Phase = dapi.AutoScalePhaseUpdating
			}
		}
		cr.Status.CN.LastTransitionTime = metav1.Now()
	}
	// sync the status of CR
	syncErr := rec.Sync()
	// update status
	updateErr := r.Status().Update(ctx, cr)

	// merge error at different reconcile phases
	mergedErr := util.MergeErrorsWithTag(map[string]error{
		"reconcile":     recErr,
		"sync":          syncErr,
		"update-status": updateErr,
	})
	if mergedErr != nil {
		return ctrl.Result{Requeue: true}, mergedErr
	} else {
		return ctrl.Result{}, nil
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *DorisAutoscalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dapi.DorisAutoscaler{}).
		Owns(&acv2.HorizontalPodAutoscaler{}).
		Complete(r)
}
