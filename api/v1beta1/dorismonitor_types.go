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
)

// DorisMonitor is the Schema for the Doris cluster monitors API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type DorisMonitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DorisMonitorSpec   `json:"spec,omitempty"`
	Status DorisMonitorStatus `json:"status,omitempty"`
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
	// Reference of the target DorisCluster
	Cluster *DorisClusterRef `json:"cluster"`

	Prometheus *PrometheusSpec `json:"prometheus,omitempty"`
	Grafana    *GrafanaSpec    `json:"grafana,omitempty"`

	// ImagePullPolicy of DorisMonitor Pods
	// +optional
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling
	// any of the images.
	// +optional
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// Replicas is the number of desired replicas.
	// Defaults to 1
	Replicas int32 `json:"replicas"`

	// Persistent determines DorisMinitor persists monitor data or not
	// Defaults to true
	// +optional
	Persistent bool `json:"persistent,omitempty"`

	// StorageClassName of the persistent volume for monitor data storage.
	// Kubernetes default storage class is used if not setting this field.
	// +optional
	StorageClassName string `json:"storageClassName,omitempty"`

	// Tolerations of the Doris initializer Pod
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
}

// PrometheusSpec defines the desired state of Prometheus
// +k8s:openapi-gen=true
type PrometheusSpec struct {
	// Image tags of the prometheus
	Image string `json:"image"`

	// Prometheus data retention time
	// +optional
	RetentionTime string `json:"retentionTime,omitempty"`

	// LogLevel is Prometheus log level
	// +optional
	LogLevel string `json:"logLevel,omitempty"`

	// CommandOptions is the additional prometheus CLI command option
	// Ref: https://prometheus.io/docs/prometheus/latest/configuration/configuration/
	// +optional
	CommandOptions *[]string `json:"commandOptions,omitempty"`

	// Service defines a Kubernetes service of Prometheus
	// +optional
	Service *MonitorServiceSpec `json:"service,omitempty"`

	// +optional
	corev1.ResourceRequirements `json:",inline"`
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
}

// MonitorServiceSpec defines the service of prometheus and grafana
// +k8s:openapi-gen=true
type MonitorServiceSpec struct {
	// Type of the real kubernetes service
	Type corev1.ServiceType `json:"type,omitempty"`

	// Expose the http query port of prometheus or grafana
	// Optional: Defaults to 0
	HttpPort int `json:"httpPort,omitempty"`

	// ExternalTrafficPolicy of the service
	// Optional: Defaults to omitted
	// +optional
	ExternalTrafficPolicy *corev1.ServiceExternalTrafficPolicyType `json:"externalTrafficPolicy,omitempty"`
}

// DorisMonitorStatus defines the observed state of DorisMonitor
// +k8s:openapi-gen=true
type DorisMonitorStatus struct {
	StatefulSetRef       NamespacedName `json:"statefulSetRef,omitempty"`
	PrometheusServiceRef NamespacedName `json:"prometheusServiceRef,omitempty"`
	GrafanaServiceRef    NamespacedName `json:"grafanaServiceRef,omitempty"`
	// +nullable
	StatefulSet *apps.StatefulSetStatus `json:"statefulSet,omitempty"`
}

func init() {
	SchemeBuilder.Register(&DorisMonitor{}, &DorisMonitorList{})
}
