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

package v1beta1

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// DorisInitializer is the Schema for the doris initializers API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type DorisInitializer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DorisInitializerSpec   `json:"spec,omitempty"`
	Status            DorisInitializerStatus `json:"status,omitempty"`
	objKey            *types.NamespacedName  `json:"-"`
}

// DorisInitializerList contains a list of DorisInitializer
// +kubebuilder:object:root=true
type DorisInitializerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DorisInitializer `json:"items"`
}

// DorisInitializerSpec defines the desired state of doris initializer Job
// +k8s:openapi-gen=true
type DorisInitializerSpec struct {
	// name of the target DorisCluster at the same namespace
	Cluster string `json:"cluster"`

	// Image tags of the python-sql executor container
	Image string `json:"image"`

	// ImagePullPolicy of Doris cluster Pods
	// +optional
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling
	// any of the images.
	// +optional
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// MaxRetry of the Doris initializer Job
	MaxRetry *int32 `json:"maxRetry"`

	// Specify a Service Account
	// +optional
	ServiceAccount string `json:"serviceAccount,omitempty"`

	// Default Doris root user password
	// +optional
	RootPassword string `json:"rootPassword,omitempty"`

	// Default Doris root user password
	// +optional
	AdminPassword string `json:"adminPassword,omitempty"`

	// Doris initialize sqls
	// +optional
	SqlScript string `json:"initSqlScript,omitempty"`

	// +optional
	corev1.ResourceRequirements `json:",inline"`

	// Tolerations of the Doris initializer Pod
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
}

// DorisInitializerStatus defines the observed state of DorisInitializer
// +k8s:openapi-gen=true
type DorisInitializerStatus struct {
	LastApplySpecHash          *string `json:"LastApplySpecHash,omitempty"`
	DorisInitializerRecStatus  `json:",inline"`
	DorisInitializerSyncStatus `json:",inline"`
}

type DorisInitializerRecStatus struct {
	Phase   InitializeRecPhase `json:"phase,omitempty"`
	Message string             `json:"Message,omitempty"`
}

type DorisInitializerSyncStatus struct {
	JobRef            NamespacedName      `json:"jobRef,omitempty"`
	Status            InitializeJobStatus `json:"status,omitempty"`
	batchv1.JobStatus `json:",inline"`
}

type InitializeRecPhase string

const (
	InitializeRecFailed    InitializeRecPhase = "failed"
	InitializeRecCompleted InitializeRecPhase = "completed"
)

type InitializeJobStatus string

const (
	InitializeJobPending   InitializeJobStatus = "pending"
	InitializeJobRunning   InitializeJobStatus = "running"
	InitializeJobCompleted InitializeJobStatus = "completed"
	InitializeJobFailed    InitializeJobStatus = "failed"
)

func init() {
	SchemeBuilder.Register(&DorisInitializer{}, &DorisInitializerList{})
}
