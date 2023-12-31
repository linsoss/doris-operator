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
	"errors"
	"fmt"
	dapi "github.com/al-assad/doris-operator/api/v1beta1"
	tran "github.com/al-assad/doris-operator/internal/transformer"
	"github.com/al-assad/doris-operator/internal/util"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"time"
)

var (
	InitializerConfHashAnnotationKey = fmt.Sprintf("%s/initr-config", dapi.GroupVersion.Group)
)

type DorisInitializerReconciler struct {
	ReconcileContext
	CR *dapi.DorisInitializer
}

type PendingError struct {
	Reason string
}

func NewPendingError(format string, a ...any) *PendingError {
	return &PendingError{Reason: fmt.Sprintf(format, a...)}
}
func (e PendingError) Error() string {
	return e.Reason
}

// Reconcile initializer resources
func (r *DorisInitializerReconciler) Reconcile() (dapi.DorisInitializerRecStatus, error) {
	if r.CR.Spec.Cluster == "" {
		return dapi.DorisInitializerRecStatus{Phase: dapi.InitializeRecCompleted}, nil
	}
	clusterRef := types.NamespacedName{
		Namespace: r.CR.Namespace,
		Name:      r.CR.Spec.Cluster,
	}

	apply := func() error {
		// check if target DorisCluster exist
		clusterCr := &dapi.DorisCluster{}
		exist, err := r.Exist(clusterRef, clusterCr)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("target DorisCluster[name=%s][namespace=%s] not exist",
				clusterRef.Name, clusterRef.Name)
		}
		// check if target DorisCluster already bound another DorisAutoscaler
		bound, bErr := r.FindRefDorisInitializer(clusterRef)
		if bErr != nil {
			return bErr
		}
		if bound != nil && bound.Name != r.CR.Name && bound.Namespace != r.CR.Namespace {
			return NewPendingError("target DorisCluster already bound another DorisInitializer[name=%s][namespace=%s]",
				bound.Name, bound.Name)
		}
		// check if target DorisCluster already to write data
		isDorisReadyToWrite := len(clusterCr.Status.FE.ReadyMembers) > 0 && len(clusterCr.Status.BE.ReadyMembers) > 0
		if !isDorisReadyToWrite {
			return fmt.Errorf("target DorisCluster[name=%s][namespace=%s] is not ready to write data",
				clusterRef.Name, clusterRef.Name)
		}

		// secret
		if secret := tran.MakeInitializerSecret(r.CR, r.Schema); secret != nil {
			if err := r.CreateOrUpdate(secret, &corev1.Secret{}); err != nil {
				return err
			}
		}
		// configmap
		configMap := tran.MakeInitializerConfigMap(r.CR, r.Schema)
		if configMap != nil {
			if err := r.CreateOrUpdate(configMap, &corev1.ConfigMap{}); err != nil {
				return err
			}
		}
		// replace job
		feQueryPort := tran.GetFeQueryPort(clusterCr)
		if job := tran.MakeInitializerJob(r.CR, feQueryPort, r.Schema); job != nil {
			job.Spec.Template.Annotations[InitializerConfHashAnnotationKey] = util.Md5HashOr(configMap.Data, "")
			if err := r.Replace(job, &batchv1.Job{}, 30*time.Second); err != nil {
				return err
			}
		}
		return nil
	}

	err := apply()
	if err == nil {
		return dapi.DorisInitializerRecStatus{Phase: dapi.InitializeRecCompleted}, nil
	} else if errors.As(err, &PendingError{}) {
		return dapi.DorisInitializerRecStatus{Phase: dapi.InitializeRecWaiting}, err
	} else {
		return dapi.DorisInitializerRecStatus{Phase: dapi.InitializeRecFailed}, err
	}
}

// Sync initializer resources status
func (r *DorisInitializerReconciler) Sync() (dapi.DorisInitializerSyncStatus, error) {
	status := util.PointerDeRefer(r.CR.Status.DorisInitializerSyncStatus.DeepCopy(), dapi.DorisInitializerSyncStatus{})
	if r.CR.Spec.Cluster == "" {
		return status, nil
	}

	jobRef := tran.GetInitializerJobKey(r.CR.ObjKey())
	job := &batchv1.Job{}
	exist, err := r.Exist(jobRef, job)
	if err != nil {
		return status, err
	}
	if exist {
		status.JobRef = dapi.NewNamespacedName(jobRef)
		status.JobStatus = job.Status
		inferJobState := func() dapi.InitializeJobStatus {
			if util.IsJobComplete(*job) {
				return dapi.InitializeJobCompleted
			}
			if util.IsJobFailed(*job) {
				return dapi.InitializeJobFailed
			}
			if job.Status.StartTime == nil {
				return dapi.InitializeJobPending
			}
			return dapi.InitializeJobRunning
		}
		status.Status = inferJobState()
	}
	return status, nil
}
