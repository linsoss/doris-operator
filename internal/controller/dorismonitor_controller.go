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
	appv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DorisMonitorReconciler reconciles a DorisMonitor object
type DorisMonitorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=al-assad.github.io,resources=dorismonitors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=al-assad.github.io,resources=dorismonitors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=al-assad.github.io,resources=dorismonitors/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=daemonsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles,verbs=get;create
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterrolebindings,verbs=get;create

func (r *DorisMonitorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	recCtx := reconciler.NewReconcileContext(r.Client, r.Scheme, ctx)

	// obtain CR
	cr := &dapi.DorisMonitor{}
	exist, err := recCtx.Exist(req.NamespacedName, cr)
	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}
	// skip reconciling process when it has been deleted
	if !exist {
		recCtx.Log.Info(fmt.Sprintf("DorisMonitor(%s) has been deleted", util.K8sObjKeyStr(req.NamespacedName)))
		return ctrl.Result{}, nil
	}
	rec := reconciler.DorisMonitorReconciler{ReconcileContext: recCtx, CR: cr}

	curSpecHash := util.Md5HashOr(cr.Spec, "")
	isFirstCreated := cr.Status.LastApplySpecHash == nil
	specHasChanged := isFirstCreated || *cr.Status.LastApplySpecHash != curSpecHash
	preRecCompleted := cr.Status.Stage == dapi.MnrOprStageCompleted

	if isFirstCreated && cr.Status.Stage == "" {
		recCtx.Log.Info(fmt.Sprintf("DorisMonitor(%s) is created for the first time", util.K8sObjKeyStr(req.NamespacedName)))
	}
	if specHasChanged {
		recCtx.Log.Info(fmt.Sprintf("DorisMonitor(%s) spec has been updated", util.K8sObjKeyStr(req.NamespacedName)))
	}

	// reconcile the sub resource of CR
	var recErr error
	if specHasChanged || !preRecCompleted {
		recRs := rec.Reconcile()
		recErr = recRs.Err
		cr.Status.DorisMonitorRecStatus = recRs.AsDorisClusterRecStatus()
		if recRs.Stage == dapi.MnrOprStageCompleted {
			cr.Status.LastApplySpecHash = &curSpecHash
		}
	}
	// sync the status of CR
	syncRs, syncErr := rec.Sync()
	cr.Status.DorisMonitorSyncStatus = syncRs
	// update status
	updateErr := r.Status().Update(ctx, cr)

	// merge error at different reconcile phases
	errSet := StCtrlErrSet{
		Rec:    recErr,
		Sync:   syncErr,
		Update: updateErr,
	}
	return errSet.AsResult()
}

// SetupWithManager sets up the controller with the Manager.
func (r *DorisMonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dapi.DorisMonitor{}).
		Owns(&appv1.DaemonSet{}).
		Owns(&appv1.Deployment{}).
		Complete(r)
}
