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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// K8sNativeComponentSpec is the native K8s pod attributes specification.
// +k8s:openapi-gen=true
type K8sNativeComponentSpec struct {

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
	StatefulSetUpdateStrategy appsv1.StatefulSetUpdateStrategyType `json:"statefulSetUpdateStrategy,omitempty"`

	// Additional environment variables to set in the container
	// +optional
	AdditionalEnvs []corev1.EnvVar `json:"env,omitempty"`

	// Additional containers of the component.
	// +optional
	AdditionalContainers []corev1.Container `json:"additionalContainers,omitempty"`

	// Additional volumes of component pod.
	// +optional
	AdditionalVolumes []corev1.Volume `json:"additionalVolumes,omitempty"`

	// Additional volume mounts of component pod.
	// +optional
	AdditionalVolumeMounts []corev1.VolumeMount `json:"additionalVolumeMounts,omitempty"`
}

// ResourceRef contains the k8s resource ref info.
// +k8s:openapi-gen=true
type ResourceRef struct {
	Name string `json:"name,omitempty"`
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// DorisClusterRef is the reference to DorisCluster.
// +k8s:openapi-gen=true
type DorisClusterRef struct {
	ResourceRef `json:",inline"`
}

// ComponentPhase is the current state of member
type ComponentPhase string

const (
	StandbyPhase ComponentPhase = "Standby"
	NormalPhase  ComponentPhase = "Normal"
	UpgradePhase ComponentPhase = "Upgrade"
	ScalePhase   ComponentPhase = "Scale"
)
