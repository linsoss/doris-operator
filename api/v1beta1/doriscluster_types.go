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

// DorisCluster is the Schema for the doris clusters API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type DorisCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DorisClusterSpec   `json:"spec,omitempty"`
	Status DorisClusterStatus `json:"status,omitempty"`
}

// DorisClusterList contains a list of DorisCluster
// +kubebuilder:object:root=true
type DorisClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DorisCluster `json:"items"`
}

// DorisClusterSpec describes the attributes that a user creates on a Doris cluster.
// +k8s:openapi-gen=true
type DorisClusterSpec struct {
	FE         *FESpec         `json:"fe,omitempty"`
	BE         *BESpec         `json:"be,omitempty"`
	CN         *CNSpec         `json:"cn,omitempty"`
	Broker     *BrokerSpec     `json:"broker,omitempty"`
	HadoopConf *HadoopConfSpec `json:"hadoopConf,omitempty"`

	// Doris cluster image version
	Version string `json:"version"`

	// ImagePullPolicy of Doris cluster Pods
	// +optional
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling
	// any of the images.
	// +optional
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// Specify a Service Account
	// +optional
	ServiceAccount string `json:"serviceAccount,omitempty"`

	// Annotations of Doris cluster pods that will be merged with component annotation settings.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// Affinity for pod scheduling of Doris cluster.
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// Tolerations are applied to Doris cluster pods, allowing pods to be scheduled onto nodes with matching taints.
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// Specify pod priorities of pods in Doris cluster, default to empty.
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`

	// Update strategy of Doris cluster StatefulSet.
	// +optional
	StatefulSetUpdateStrategy apps.StatefulSetUpdateStrategyType `json:"statefulSetUpdateStrategy,omitempty"`
}

// FESpec contains details of FE members.
// +k8s:openapi-gen=true
type FESpec struct {
	DorisComponentSpec `json:",inline"`

	// The storageClassName of the persistent volume for TiDB data storage.
	// Defaults to Kubernetes default storage class.
	// +optional
	StorageClassName string `json:"storageClassName,omitempty"`

	// Service defines a Kubernetes service of FE
	Service *FeServiceSpec `json:"service,omitempty"`
}

// BESpec contains details of BE members.
// +k8s:openapi-gen=true
type BESpec struct {
	DorisComponentSpec `json:",inline"`

	// The storageClassName of the persistent volume for TiDB data storage.
	// Defaults to Kubernetes default storage class.
	// +optional
	StorageClassName string `json:"storageClassName,omitempty"`
}

// CNSpec contains details of CN members.
// +k8s:openapi-gen=true
type CNSpec struct {
	DorisComponentSpec `json:",inline"`
}

// BrokerSpec contains details of Broker members.
// +k8s:openapi-gen=true
type BrokerSpec struct {
	DorisComponentSpec `json:",inline"`
}

// HadoopConfSpec contains the configuration needed for doris to connect to the Hadoop cluster.
// +k8s:openapi-gen=true
type HadoopConfSpec struct {
	// Hostname-IP mapping of Hadoop cluster nodes.
	Hosts []HostnameIpItem `json:"hostAliases"`
	// Hadoop configuration files.
	Config map[string]string `json:"config,omitempty"`
}

// HostnameIpItem define Hostname-IP kv item
// +k8s:openapi-gen=true
type HostnameIpItem struct {
	IP   string `json:"ip"`
	Name string `json:"name"`
}

// FeServiceSpec defines `.fe.service` field of `DorisCluster.spec`.
// +k8s:openapi-gen=true
type FeServiceSpec struct {
	// Type of the real kubernetes service
	Type corev1.ServiceType `json:"type,omitempty"`

	// ExternalTrafficPolicy of the service
	// Optional: Defaults to omitted
	// +optional
	ExternalTrafficPolicy *corev1.ServiceExternalTrafficPolicyType `json:"externalTrafficPolicy,omitempty"`

	// Expose the FE query port
	// Optional: Defaults to 0
	// +optional
	QueryPort int `json:"queryPort,omitempty"`

	// Expose the FE http port
	// Optional: Defaults to 0
	// +optional
	HttpPort int `json:"httpPort,omitempty"`
}

// DorisComponentSpec is the base component spec.
// +k8s:openapi-gen=true
type DorisComponentSpec struct {
	//Base image of the component
	BaseImage string `json:"baseImage"`

	// Type of the real kubernetes service
	Version string `json:"version"`

	// The desired ready replicas
	// +kubebuilder:validation:Minimum=0
	Replicas int32 `json:"replicas"`

	// Defines the specification of resource cpu, mem, storage.
	corev1.ResourceRequirements `json:",inline"`

	// Additional Doris component configuration
	// Ref:
	// - https://doris.apache.org/docs/dev/admin-manual/config/fe-config/
	// - https://doris.apache.org/docs/dev/admin-manual/config/be-config/
	// +optional
	Configs map[string]string `json:"config,omitempty"`

	// HostAliases is an optional list of hosts and IPs that will be injected into the pod's hosts
	// file if specified.
	// +optional
	HostAliases []corev1.HostAlias `json:"hostAliases,omitempty"`

	// Additional component spec
	K8sNativeComponentSpec `json:",inline"`
}

// DorisClusterStatus defines the observed state of DorisCluster
// +k8s:openapi-gen=true
type DorisClusterStatus struct {
	FE         FEStatus                `json:"fe,omitempty"`
	BE         BEStatus                `json:"be,omitempty"`
	CN         CNStatus                `json:"cn,omitempty"`
	Broker     BrokerStatus            `json:"broker,omitempty"`
	Conditions []DorisClusterCondition `json:"conditions,omitempty"`

	// The reference of the DorisInitializer
	// +nullable
	InitializerRef ResourceRef `json:"initializerRef,omitempty"`

	// Whether the initialization process defined by DorisInitializer has been completed once.
	// +nullable
	Initialized bool `json:"initialized,omitempty"`

	// The reference of the DorisAutoscaler
	// +nullable
	AutoScalerRef AutoScalerRef `json:"autoScalerRef,omitempty"`

	// The reference of the DorisMonitor
	// +nullable
	MonitorRef ResourceRef `json:"monitorRef,omitempty"`
}

// FEStatus represents the current state of Doris FE
type FEStatus struct {
	ServiceRef      ResourceRef `json:"serviceName,omitempty"`
	ComponentStatus `json:",inline"`
}

// BEStatus represents the current state of Doris BE
type BEStatus struct {
	ComponentStatus `json:",inline"`
}

// CNStatus represents the current state of Doris CN
type CNStatus struct {
	ComponentStatus `json:",inline"`
}

// BrokerStatus represents the current state of Doris Broker
type BrokerStatus struct {
	ComponentStatus `json:",inline"`
}

type ComponentStatus struct {
	StatefulSetRef ResourceRef                `json:"statefulSetRef,omitempty"`
	StatefulSet    *apps.StatefulSetStatus    `json:"statefulSet,omitempty"`
	Image          string                     `json:"image,omitempty"`
	Phase          ComponentPhase             `json:"phase,omitempty"`
	Members        map[string]ComponentMember `json:"members,omitempty"`
	Conditions     []metav1.Condition         `json:"conditions,omitempty"`
}

// ComponentMember represents the current state of a component member
type ComponentMember struct {
	Name      string `json:"name"`
	Available bool   `json:"available"`
	// +nullable
	CreatedAt metav1.Time `json:"createdAt,omitempty"`
	// +nullable
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
}

// DorisClusterCondition describes the state of a tidb cluster at a certain point.
type DorisClusterCondition struct {
	Type               DorisClusterConditionType `json:"type"`
	Status             corev1.ConditionStatus    `json:"status"`
	LastTransitionTime metav1.Time               `json:"lastTransitionTime,omitempty"`
	Reason             string                    `json:"reason,omitempty"`
	Message            string                    `json:"message,omitempty"`
}

// DorisClusterConditionType represents a tidb cluster condition value.
type DorisClusterConditionType string

const (
	// DorisClusterReady indicates that the tidb cluster is ready or not.
	DorisClusterReady DorisClusterConditionType = "Ready"
)

func init() {
	SchemeBuilder.Register(&DorisCluster{}, &DorisClusterList{})
}
