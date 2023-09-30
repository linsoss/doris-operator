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

package util

import (
	acv2 "k8s.io/api/autoscaling/v2"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewEmptyDirVolumeSource() corev1.VolumeSource {
	return corev1.VolumeSource{
		EmptyDir: &corev1.EmptyDirVolumeSource{
			Medium: corev1.StorageMediumDefault},
	}
}

func NewConfigMapVolumeSource(configMapName string) corev1.VolumeSource {
	return corev1.VolumeSource{
		ConfigMap: &corev1.ConfigMapVolumeSource{
			LocalObjectReference: corev1.LocalObjectReference{Name: configMapName},
		},
	}
}

func NewEnvVarSecretSource(secretName string, key string) *corev1.EnvVarSource {
	return &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: secretName},
			Key:                  key,
		},
	}
}

func NewTcpSocketProbeHandler(tcpPort int32) corev1.ProbeHandler {
	return corev1.ProbeHandler{
		TCPSocket: &corev1.TCPSocketAction{
			Port: intstr.FromInt(int(tcpPort)),
		},
	}
}

func NewResourceAvgUtilizationMetricSpec(name corev1.ResourceName, avgUnit *int32) acv2.MetricSpec {
	return acv2.MetricSpec{
		Type: acv2.ResourceMetricSourceType,
		Resource: &acv2.ResourceMetricSource{
			Name: name,
			Target: acv2.MetricTarget{
				Type:               acv2.UtilizationMetricType,
				AverageUtilization: avgUnit,
			},
		},
	}
}

func IsPodReady(pod corev1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == "Ready" && condition.Status == "True" {
			return true
		}
	}
	return false
}

func IsJobComplete(job batchv1.Job) bool {
	for _, condition := range job.Status.Conditions {
		if condition.Type == "Complete" && condition.Status == "True" {
			return true
		}
	}
	return false
}

func IsJobFailed(job batchv1.Job) bool {
	for _, condition := range job.Status.Conditions {
		if condition.Type == "Failed" && condition.Status == "True" {
			return true
		}
	}
	return false
}
