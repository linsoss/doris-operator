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

package transformer

import (
	"fmt"
	dapi "github.com/al-assad/doris-operator/api/v1beta1"
	"github.com/al-assad/doris-operator/internal/template"
	"github.com/al-assad/doris-operator/internal/util"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	DefaultPrometheusImage = "prom/prometheus:v2.37.8"
)

var (
	PrometheusConfTmpl = template.NewTemplateOrPanic("prometheus-conf", "monitor/prometheus.yml")
)

func GetMonitorPrometheusLabels(dorisClusterKey types.NamespacedName) map[string]string {
	return MakeResourceLabels(dorisClusterKey.Name, "prometheus")
}

func GetPrometheusConfigMapKey(monitorKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: monitorKey.Namespace,
		Name:      fmt.Sprintf("%s-promtheus-config", monitorKey.Name),
	}
}

func GetPrometheusServiceKey(monitorKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: monitorKey.Namespace,
		Name:      fmt.Sprintf("%s-prometheus", monitorKey.Name),
	}
}

func GetPrometheusPVCKey(monitorKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: monitorKey.Namespace,
		Name:      fmt.Sprintf("%s-prometheus-pvc", monitorKey.Name),
	}
}

func GetPrometheusDeploymentKey(monitorKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: monitorKey.Namespace,
		Name:      fmt.Sprintf("%s-prometheus", monitorKey.Name),
	}
}

func MakePrometheusConfigMap(cr *dapi.DorisMonitor, scheme *runtime.Scheme) (*corev1.ConfigMap, error) {
	if cr.Spec.Cluster == "" {
		return nil, nil
	}
	clusterRef := types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Spec.Cluster,
	}
	configMapRef := GetPrometheusConfigMapKey(cr.ObjKey())
	labels := GetMonitorPrometheusLabels(clusterRef)
	promConfContent, err := template.ExecTemplate(PrometheusConfTmpl, clusterRef)
	if err != nil {
		return nil, util.MergeErrors(fmt.Errorf("fail to parse prometheus.conf template"), err)
	}

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapRef.Name,
			Namespace: configMapRef.Namespace,
			Labels:    labels,
		},
		Data: map[string]string{
			"prometheus.yml": promConfContent,
		},
	}
	_ = controllerutil.SetOwnerReference(cr, configMap, scheme)
	return configMap, nil
}

func MakePrometheusService(cr *dapi.DorisMonitor, scheme *runtime.Scheme) *corev1.Service {
	if cr.Spec.Cluster == "" {
		return nil
	}
	clusterRef := types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Spec.Cluster,
	}
	serviceRef := GetPrometheusServiceKey(cr.ObjKey())
	labels := GetMonitorPrometheusLabels(clusterRef)

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceRef.Name,
			Namespace: serviceRef.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: labels,
		},
	}
	httpPort := corev1.ServicePort{
		Name: "http-port",
		Port: 9090,
	}
	crSvc := cr.Spec.Prometheus.Service
	if crSvc != nil {
		if crSvc.Type != "" {
			service.Spec.Type = crSvc.Type
		}
		if crSvc.ExternalTrafficPolicy != nil {
			service.Spec.ExternalTrafficPolicy = *crSvc.ExternalTrafficPolicy
		}
		if crSvc.HttpPort != nil {
			httpPort.NodePort = *crSvc.HttpPort
		}
	}
	service.Spec.Ports = []corev1.ServicePort{httpPort}
	_ = controllerutil.SetOwnerReference(cr, service, scheme)
	return service
}

func MakePrometheusPVC(cr *dapi.DorisMonitor, scheme *runtime.Scheme) *corev1.PersistentVolumeClaim {
	if cr.Spec.Cluster == "" {
		return nil
	}
	clusterRef := types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Spec.Cluster,
	}
	pvcRef := GetPrometheusPVCKey(cr.ObjKey())
	labels := GetMonitorPrometheusLabels(clusterRef)

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pvcRef.Name,
			Namespace: pvcRef.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			StorageClassName: cr.Spec.StorageClassName,
		},
	}
	storageRequest := cr.Spec.Prometheus.Requests.Storage()
	if storageRequest != nil {
		pvc.Spec.Resources.Requests = corev1.ResourceList{corev1.ResourceStorage: *storageRequest}
	}
	_ = controllerutil.SetOwnerReference(cr, pvc, scheme)
	return pvc
}

func MakePrometheusDeployment(cr *dapi.DorisMonitor, scheme *runtime.Scheme) *appv1.Deployment {
	if cr.Spec.Cluster == "" {
		return nil
	}
	clusterRef := types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Spec.Cluster,
	}
	deploymentRef := GetPrometheusDeploymentKey(cr.ObjKey())
	configMapRef := GetPrometheusConfigMapKey(cr.ObjKey())
	pvcRef := GetPrometheusPVCKey(cr.ObjKey())
	labels := GetMonitorPrometheusLabels(clusterRef)

	replicas := int32(1)
	// prometheus args
	promArgs := []string{
		"--config.file=/etc/prometheus/prometheus.yml",
		"--storage.tsdb.path=/data/prometheus",
		"--web.enable-lifecycle",
	}
	if cr.Spec.Prometheus.RetentionTime != nil {
		promArgs = append(promArgs,
			fmt.Sprintf("--storage.tsdb.retention.time=%s", *cr.Spec.Prometheus.RetentionTime),
		)
	}

	// pod template
	podTemplate := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      labels,
			Annotations: make(map[string]string),
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: MonitorNamespacedAccountName,
			ImagePullSecrets:   cr.Spec.ImagePullSecrets,
			Volumes: []corev1.Volume{
				{
					Name: "prometheus-config",
					VolumeSource: util.NewConfigMapItemsVolumeSource(
						configMapRef.Name, map[string]string{"prometheus.yml": "prometheus.yml"}),
				}, {
					Name: "prometheus-data",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: pvcRef.Name}}},
			},
			Containers: []corev1.Container{{
				Name:            "prometheus",
				Image:           util.StringFallback(cr.Spec.Prometheus.Image, DefaultPrometheusImage),
				ImagePullPolicy: cr.Spec.ImagePullPolicy,
				Resources:       formatContainerResourcesRequirement(cr.Spec.Prometheus.ResourceRequirements),
				Args:            promArgs,
				Ports: []corev1.ContainerPort{{
					Name:          "http-port",
					ContainerPort: 9090,
				}},
				VolumeMounts: []corev1.VolumeMount{
					{Name: "prometheus-config", MountPath: "/etc/prometheus"},
					{Name: "prometheus-data", MountPath: "/data/prometheus"},
				},
				ReadinessProbe: &corev1.Probe{
					ProbeHandler:     util.NewHttpGetProbeHandler("/-/ready", 9090),
					TimeoutSeconds:   1,
					PeriodSeconds:    5,
					SuccessThreshold: 1,
					FailureThreshold: 120,
				},
				LivenessProbe: &corev1.Probe{
					ProbeHandler:     util.NewHttpGetProbeHandler("/-/ready", 9090),
					TimeoutSeconds:   1,
					PeriodSeconds:    5,
					SuccessThreshold: 1,
					FailureThreshold: 30,
				},
			}},
		},
	}

	// deployment
	deployment := &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentRef.Name,
			Namespace: deploymentRef.Namespace,
			Labels:    labels,
		},
		Spec: appv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: podTemplate,
		},
	}
	_ = controllerutil.SetOwnerReference(cr, deployment, scheme)
	_ = controllerutil.SetControllerReference(cr, deployment, scheme)
	return deployment
}
