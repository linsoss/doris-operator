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
	"context"
	dapi "github.com/al-assad/doris-operator/api/v1beta1"
	"github.com/al-assad/doris-operator/internal/transformer"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DorisClusterReconciler reconciles a DorisCluster object
type DorisClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Req    ctrl.Request
	Ctx    context.Context
	CR     *dapi.DorisCluster
}

// ClusterStageRecResult represents the result of a stage reconciliation for DorisCluster
type ClusterStageRecResult struct {
	Stage   dapi.DorisClusterOprStage
	Status  dapi.DorisClusterOprStageStatus
	Message string
	Err     error
}

func recSuccess(stage dapi.DorisClusterOprStage) ClusterStageRecResult {
	return ClusterStageRecResult{Stage: stage, Status: dapi.OprStageStatusSucceeded, Err: nil}
}

func recFail(stage dapi.DorisClusterOprStage, err error) ClusterStageRecResult {
	return ClusterStageRecResult{Stage: stage, Status: dapi.OprStageStatusSucceeded, Err: err}
}

// Reconcile secret object that using to store the sql query account info
// that used by doris-operator.
func (r *DorisClusterReconciler) recOprAccountSecret() ClusterStageRecResult {
	secretRef := transformer.GetOprSqlAccountSecretName(r.CR)
	secret := &corev1.Secret{}
	secretExists := false

	if err := r.Get(r.Ctx, secretRef, secret); err != nil {
		if errors.IsNotFound(err) {
			secretExists = false
		} else {
			return recSuccess(dapi.OprStageSqlAccountSecret)
		}
	}
	// create secret if not exists
	if !secretExists {
		newSecret := transformer.MakeOprSqlAccountSecret(r.CR)
		if err := r.Create(r.Ctx, newSecret); err != nil {
			return recFail(dapi.OprStageSqlAccountSecret, err)
		}
	}
	return recSuccess(dapi.OprStageSqlAccountSecret)
}

// Reconcile Doris FE resources.
func recFeResources() ClusterStageRecResult {
	// FE configmap
	// FE service
	// FE statefulset

	return recSuccess("")
}
