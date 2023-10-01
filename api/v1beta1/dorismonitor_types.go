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
	apps "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// DorisMonitor is the Schema for the Doris cluster monitors API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type DorisMonitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DorisMonitorSpec      `json:"spec,omitempty"`
	Status            DorisMonitorStatus    `json:"status,omitempty"`
	objKey            *types.NamespacedName `json:"-"`
}

// DorisMonitorList contains a list of DorisMonitor
// +kubebuilder:object:root=true
type DorisMonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DorisMonitor `json:"items"`
}

// DorisMonitorSpec defines the desired state of DorisMonitor
// +k8s:openapi-gen=true
type DorisMonitorSpec struct {
	// Name of the target DorisCluster
	Cluster string `json:"cluster"`

	Prometheus *PrometheusSpec `json:"prometheus,omitempty"`
	Grafana    *GrafanaSpec    `json:"grafana,omitempty"`
	Loki       *LokiSpec       `json:"loki,omitempty"`
	Promtail   *PromtailSpec   `json:"promtail,omitempty"`

	// DisableLoki to disable Loki for log collection
	// Default to false
	DisableLoki bool `json:"enableLoki,omitempty"`

	// ImagePullPolicy of DorisMonitor Pods
	// +optional
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling
	// any of the images.
	// +optional
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// StorageClassName of the persistent volume for monitor data storage.
	// Kubernetes default storage class is used if not setting this field.
	// +optional
	StorageClassName *string `json:"storageClassName,omitempty"`

	// Specify a Service Account
	// +optional
	ServiceAccount string `json:"serviceAccount,omitempty"`
}

// PrometheusSpec defines the desired state of Prometheus
// +k8s:openapi-gen=true
type PrometheusSpec struct {
	// Image tags of the prometheus
	Image string `json:"image"`

	// Prometheus data retention time
	// +optional
	RetentionTime *string `json:"retentionTime,omitempty"`

	// Service defines a Kubernetes service of Prometheus
	// +optional
	Service *MonitorServiceSpec `json:"service,omitempty"`

	// +optional
	corev1.ResourceRequirements `json:",inline"`

	// +optional
	StorageClassName string `json:"storageClassName,omitempty"`
}

// GrafanaSpec defines the desired state of Grafana
// +k8s:openapi-gen=true
type GrafanaSpec struct {
	// Image tags of the grafana
	Image string `json:"image"`

	// Default grafana admin username
	// +optional
	AdminUser string `json:"adminUser,omitempty"`

	// Default grafana admin password
	// +optional
	AdminPassword string `json:"adminPassword,omitempty"`

	// +optional
	Service *MonitorServiceSpec `json:"service,omitempty"`

	// +optional
	corev1.ResourceRequirements `json:",inline"`

	// +optional
	StorageClassName string `json:"storageClassName,omitempty"`
}

// LokiSpec defines the desired state of Loki
type LokiSpec struct {
	// Image tags of the loki
	Image string `json:"image"`

	// Loki chunks retention time
	// When it is empty, do not enable retention deletion of Loki.
	// +optional
	RetentionTime *string `json:"retentionTime,omitempty"`

	// +optional
	corev1.ResourceRequirements `json:",inline"`

	// +optional
	StorageClassName string `json:"storageClassName,omitempty"`
}

// PromtailSpec defines the desired state of Promtail
type PromtailSpec struct {
	// Image tags of the promtail
	Image string `json:"image"`
	// +optional
	corev1.ResourceRequirements `json:",inline"`
}

// MonitorServiceSpec defines the service of prometheus and grafana
// +k8s:openapi-gen=true
type MonitorServiceSpec struct {
	// Type of the real kubernetes service
	// Only ClusterIP and NodePort support is available.
	Type corev1.ServiceType `json:"type,omitempty"`

	// Expose the http query port of prometheus or grafana
	// Optional: Defaults to 0
	HttpPort *int32 `json:"httpPort,omitempty"`

	// ExternalTrafficPolicy of the service
	// Optional: Defaults to omitted
	// +optional
	ExternalTrafficPolicy *corev1.ServiceExternalTrafficPolicyType `json:"externalTrafficPolicy,omitempty"`
}

// DorisMonitorStatus defines the observed state of DorisMonitor
// +k8s:openapi-gen=true
type DorisMonitorStatus struct {
	LastApplySpecHash      *string `json:"lastApplySpecHash,omitempty"`
	DorisMonitorRecStatus  `json:",inline"`
	DorisMonitorSyncStatus `json:",inline"`
}

type DorisMonitorRecStatus struct {
	Stage       DorisMonitorOprStage `json:"stage,omitempty"`
	StageStatus OprStageStatus       `json:"stageStatus,omitempty"`
	StageAction OprStageAction       `json:"stageAction,omitempty"`
	LastMessage string               `json:"lastMessage,omitempty"`
}

type DorisMonitorSyncStatus struct {
	Prometheus PrometheusStatus `json:"prometheus,omitempty"`
	Grafana    GrafanaStatus    `json:"grafana,omitempty"`
	Loki       LokiStatus       `json:"loki,omitempty"`
	Promtail   PromtailStatus   `json:"promtail,omitempty"`
}

type DorisMonitorOprStage string

const (
	MnrOprStageRbac                     DorisMonitorOprStage = "rbac"
	MnrOprStageGlobalClusterRole        DorisMonitorOprStage = "global-rbac/ClusterRole"
	MnrOprStageNamespacedServiceAccount DorisMonitorOprStage = "rbac/ServiceAccount"
	MnrOprStageNamespacedRoleBinding    DorisMonitorOprStage = "rbac/RoleBinding"

	MnrOprStagePrometheus           DorisMonitorOprStage = "prometheus"
	MnrOprStagePrometheusConfigMap  DorisMonitorOprStage = "prometheus/ConfigMap"
	MnrOprStagePrometheusService    DorisMonitorOprStage = "prometheus/Service"
	MnrOprStagePrometheusPVC        DorisMonitorOprStage = "prometheus/PersistentVolumeClaim"
	MnrOprStagePrometheusDeployment DorisMonitorOprStage = "prometheus/Deployment"

	MnrOprStageGrafana           DorisMonitorOprStage = "grafana"
	MnrOprStageGrafanaSecret     DorisMonitorOprStage = "grafana/Secret"
	MnrOprStageGrafanaConfigMap  DorisMonitorOprStage = "grafana/ConfigMap"
	MnrOprStageGrafanaService    DorisMonitorOprStage = "grafana/Service"
	MnrOprStageGrafanaPVC        DorisMonitorOprStage = "grafana/PersistentVolumeClaim"
	MnrOprStageGrafanaDeployment DorisMonitorOprStage = "grafana/Deployment"

	MnrOprStagePromtail          DorisMonitorOprStage = "promtail"
	MnrOprStagePromtailConfigMap DorisMonitorOprStage = "promtail/ConfigMap"
	MnrOprStagePromtailDaemonSet DorisMonitorOprStage = "promtail/DemonSet"

	MnrOprStageLoki           DorisMonitorOprStage = "loki"
	MnrOprStageLokiConfigMap  DorisMonitorOprStage = "loki/ConfigMap"
	MnrOprStageLokiService    DorisMonitorOprStage = "loki/Service"
	MnrOprStageLokiPVC        DorisMonitorOprStage = "loki/PersistentVolumeClaim"
	MnrOprStageLokiDeployment DorisMonitorOprStage = "loki/Deployment"

	MnrRecStageComplete DorisMonitorOprStage = "completed"
)

// PrometheusStatus represents the current state of the prometheus
type PrometheusStatus struct {
	DorisMonitorComponentStatus `json:",inline"`
}

// GrafanaStatus represents the current state of the grafana
type GrafanaStatus struct {
	DorisMonitorComponentStatus `json:",inline"`
}

// LokiStatus represents the current state of the loki
type LokiStatus struct {
	DorisMonitorComponentStatus `json:",inline"`
}

// PromtailStatus represents the current state of the promtail
type PromtailStatus struct {
	DaemonSetRef NamespacedName            `json:"daemonSetRef,omitempty"`
	Ready        bool                      `json:"ready,omitempty"`
	Conditions   []apps.DaemonSetCondition `json:"conditions,omitempty"`
}

// DorisMonitorComponentStatus defines the status of the doris monitor component
type DorisMonitorComponentStatus struct {
	ServiceRef    NamespacedName             `json:"serviceRef,omitempty"`
	DeploymentRef NamespacedName             `json:"deploymentRef,omitempty"`
	PVCRef        NamespacedName             `json:"pvcRef,omitempty"`
	Ready         bool                       `json:"ready,omitempty"`
	Conditions    []apps.DeploymentCondition `json:"conditions,omitempty"`
}

func init() {
	SchemeBuilder.Register(&DorisMonitor{}, &DorisMonitorList{})
}
