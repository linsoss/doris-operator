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
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"

	dapi "github.com/al-assad/doris-operator/api/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DorisInitializerReconciler reconciles a DorisInitializer object
type DorisInitializerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=al-assad.github.io,resources=dorisinitializers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=al-assad.github.io,resources=dorisinitializers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=al-assad.github.io,resources=dorisinitializers/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete

func (r *DorisInitializerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	recCtx := reconciler.NewReconcileContext(r.Client, r.Scheme, ctx)

	// obtain DorisInitializerReconciler CR and skip reconciling process when it has been deleted
	cr := &dapi.DorisInitializer{}
	if err := recCtx.Find(req.NamespacedName, cr); err != nil {
		return ctrl.Result{Requeue: true}, err
	}
	if cr == nil {
		recCtx.Log.Info("DorisInitializer has been deleted")
		return ctrl.Result{}, nil
	}
	rec := reconciler.DorisInitializerReconciler{ReconcileContext: recCtx, CR: cr}

	// reconcile the sub resources
	curSpecHash := util.Md5HashOr(cr.Spec, "")
	preRecNotCompleted := cr.Status.Phase != dapi.InitializeRecCompleted
	specHasChanged := cr.Status.LastApplySpecHash == nil || *cr.Status.LastApplySpecHash != curSpecHash

	var recErr error
	if preRecNotCompleted || specHasChanged {
		recRs, err := rec.Reconcile()
		recErr = err
		cr.Status.DorisInitializerRecStatus = recRs
		// when reconcile succeed, update the last apply sepc hash
		if err == nil {
			cr.Status.LastApplySpecHash = &curSpecHash
		}
	}
	// sync the status of CR
	syncRs, syncErr := rec.Sync()
	cr.Status.DorisInitializerSyncStatus = syncRs
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
func (r *DorisInitializerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dapi.DorisInitializer{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Secret{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
