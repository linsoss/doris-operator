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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DorisAutoscaler is the Schema for the Doris cluster autoscaling API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type DorisAutoscaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DorisAutoscalerSpec   `json:"spec,omitempty"`
	Status DorisAutoscalerStatus `json:"status,omitempty"`
}

// DorisAutoscalerList contains a list of DorisAutoscaler
// +kubebuilder:object:root=true
type DorisAutoscalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DorisAutoscaler `json:"items"`
}

// DorisAutoscalerSpec defines the desired state of DorisAutoscaler
type DorisAutoscalerSpec struct {
	// Reference of the target DorisCluster
	Cluster *DorisClusterRef  `json:"cluster"`
	CN      *CNAutoscalerSpec `json:"cn,omitempty"`
}

// CNAutoscalerSpec contains autoscaling details of CN components.
type CNAutoscalerSpec struct {
	// Reference of the target DorisCluster
	Cluster *DorisClusterRef `json:"cluster"`

	// The range of replicas for automatic scaling
	Replicas *ReplicasRange `json:"replicas,omitempty"`

	// The metric rules for automatic scaling
	Rules *CNAutoscalerRules `json:"rules,omitempty"`

	// ScalePeriodSeconds indicates the length of time in the past for which the k8s HPA scale policy must hold true
	// Ref: https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#scaling-policies
	// +optional
	ScalePeriodSeconds *ScalePeriodSeconds `json:"scalePeriodSeconds,omitempty"`

	// Whether to disable scale down
	// Default to false
	// +optional
	DisableScaleDown bool `json:"disableScaleDown,omitempty"`
}

// CNAutoscalerRules contains metric rules for automatic scaling.
type CNAutoscalerRules struct {
	// Rules for scaling based on cpu usage percentage of CN pods
	// +optional
	Cpu *UtilizationThresholdRange `json:"cpu,omitempty"`

	// Rules for scaling based on memory usage percentage of CN pods
	// +optional
	Memory *UtilizationThresholdRange `json:"memory,omitempty"`
}

type ReplicasRange struct {
	Max int32 `json:"max,omitempty"`
	Min int32 `json:"min,omitempty"`
}

type UtilizationThresholdRange struct {
	Max int32 `json:"max,omitempty"`
	Min int32 `json:"min,omitempty"`
}

type ScalePeriodSeconds struct {
	ScaleUp   int `json:"scaleUp,omitempty"`
	ScaleDown int `json:"scaleDown,omitempty"`
}

// DorisAutoscalerStatus defines the observed state of DorisAutoscaler
type DorisAutoscalerStatus struct {
	CN *CNAutoscalerStatus `json:"cn,omitempty"`
}

// CNAutoscalerStatus defines the observed state of CN autoscaler
type CNAutoscalerStatus struct {
	Phase AutoScalePhase `json:"phase,omitempty"`
	// Last time the condition transit from one Phase to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// +optional
	Reason string `json:"reason,omitempty"`
	// +optional
	Message string `json:"message,omitempty"`
	// Scale up HPA
	// +nullable
	ScaleUp *AutoScalerRef `json:"scaleUp,omitempty"`
	// +nullable
	// Scale down HPA
	ScaleDown *AutoScalerRef `json:"scaleDown,omitempty"`
}

// AutoScalePhase is the current state of autoscaler
type AutoScalePhase string

const (
	AutoScalePhaseStandby ComponentPhase = "Standby"
	AutoScalPhaseUpgrade  ComponentPhase = "Upgrade"
	AutoScalePhaseNormal  ComponentPhase = "Normal"
)

// AutoScalerRef identifies a K8s HPA resource.
type AutoScalerRef struct {
	*ResourceRef     `json:",inline"`
	*metav1.TypeMeta `json:",inline"`
}

func init() {
	SchemeBuilder.Register(&DorisAutoscaler{}, &DorisAutoscalerList{})
}
