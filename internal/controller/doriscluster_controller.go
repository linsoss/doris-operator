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
	dapi "github.com/al-assad/doris-operator/api/v1beta1"
	"github.com/al-assad/doris-operator/internal/reconciler"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DorisClusterReconciler reconciles a DorisCluster object
type DorisClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=al-assad.github.io,resources=dorisclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=al-assad.github.io,resources=dorisclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=al-assad.github.io,resources=dorisclusters/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

func (r *DorisClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	recCtx := reconciler.NewReconcileContext(r.Client, r.Scheme, ctx)

	// obtain CR and skip reconciling process when it has been deleted
	cr := &dapi.DorisCluster{}
	crExist, err := recCtx.Exist(req.NamespacedName, cr)
	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}
	if !crExist {
		recCtx.Log.Info("DorisCluster has been deleted")
		return ctrl.Result{}, nil
	}
	// reconcile CR when the spec of it has been changed
	crSpecHasChanged := cr.Status.PrevSpec == nil || !reflect.DeepEqual(cr.Spec, *cr.Status.PrevSpec)
	// TODO reconcile when middle stage fail
	if crSpecHasChanged {

		//rec := reconciler.DorisClusterReconciler{
		//	ReconcileContext: recCtx,
		//	CR:               cr}

		cr.Status.PrevSpec = cr.Spec.DeepCopy()

	}
	// sync the status of CR
	// TODO sync logic

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DorisClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dapi.DorisCluster{}).
		Complete(r)
}
