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
	DefaultGrafanaImage = "grafana/grafana:9.5.2"
)

var (
	GrafanaDataSourceConfTmpl    = template.NewTemplateOrPanic("grafana-datasource-conf", "monitor/grafana-datasource.yml")
	GrafanaDashboardConfContent  = template.ReadOrPanic("monitor/grafana-dashboard.yml")
	GrafanaDashboardsConfContent = template.ReadOrPanic("monitor/grafana-dashboards.json")
)

type GrafanaDataSourceTmplData struct {
	PrometheusName      string
	PrometheusNamespace string
	LokiName            string
	LokiNamespace       string
}

func GetGrafanaLabels(dorisClusterKey types.NamespacedName) map[string]string {
	return MakeResourceLabels(dorisClusterKey.Name, "grafana")
}

func GetGrafanaSecretKey(monitorKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: monitorKey.Namespace,
		Name:      fmt.Sprintf("%s-grafana-config", monitorKey.Name),
	}
}
func GetGrafanaConfigMapKey(monitorKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: monitorKey.Namespace,
		Name:      fmt.Sprintf("%s-grafana-secret", monitorKey.Name),
	}
}

func GetGrafanaServiceKey(monitorKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: monitorKey.Namespace,
		Name:      fmt.Sprintf("%s-grafana", monitorKey.Name),
	}
}

func GetGrafanaPVCKey(monitorKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: monitorKey.Namespace,
		Name:      fmt.Sprintf("%s-grafana-pvc", monitorKey.Name),
	}
}

func GetGrafanaDeploymentKey(monitorKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: monitorKey.Namespace,
		Name:      fmt.Sprintf("%s-grafana", monitorKey.Name),
	}
}

func MakeGrafanaConfigMap(cr *dapi.DorisMonitor, scheme *runtime.Scheme) (*corev1.ConfigMap, error) {
	if cr.Spec.Cluster == "" {
		return nil, nil
	}
	clusterRef := types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Spec.Cluster,
	}
	configMapRef := GetGrafanaConfigMapKey(cr.ObjKey())
	prometheusSvcRef := GetPrometheusServiceKey(cr.ObjKey())
	lokiSvcRef := GetLokiServiceKey(cr.ObjKey())
	labels := GetGrafanaLabels(clusterRef)

	// generate grafana datasource template
	tmplData := GrafanaDataSourceTmplData{
		PrometheusName:      prometheusSvcRef.Name,
		PrometheusNamespace: prometheusSvcRef.Namespace,
		LokiName:            lokiSvcRef.Name,
		LokiNamespace:       lokiSvcRef.Namespace,
	}
	datasourceData, err := template.ExecTemplate(GrafanaDataSourceConfTmpl, tmplData)
	if err != nil {
		return nil, util.MergeErrors(fmt.Errorf("fail to parse grafana datasource.yaml template"), err)
	}

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapRef.Name,
			Namespace: configMapRef.Namespace,
			Labels:    labels,
		},
		Data: map[string]string{
			"datasource.yml":  datasourceData,
			"dashboard.yml":   GrafanaDashboardConfContent,
			"dashboards.json": GrafanaDashboardsConfContent,
		},
	}
	_ = controllerutil.SetOwnerReference(cr, configMap, scheme)
	return configMap, nil
}

func MakeGrafanaSecret(cr *dapi.DorisMonitor, scheme *runtime.Scheme) *corev1.Secret {
	if cr.Spec.Cluster == "" {
		return nil
	}
	clusterRef := types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Spec.Cluster,
	}
	secretRef := GetGrafanaSecretKey(cr.ObjKey())
	labels := GetGrafanaLabels(clusterRef)

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretRef.Name,
			Namespace: secretRef.Namespace,
			Labels:    labels,
		},
		Type: corev1.SecretTypeOpaque,
		StringData: map[string]string{
			"user":     util.StringFallback(cr.Spec.Grafana.AdminUser, "admin"),
			"password": util.StringFallback(cr.Spec.Grafana.AdminPassword, "admin"),
		},
	}
	_ = controllerutil.SetOwnerReference(cr, secret, scheme)
	return secret
}

func MakeGrafanaService(cr *dapi.DorisMonitor, scheme *runtime.Scheme) *corev1.Service {
	if cr.Spec.Cluster == "" {
		return nil
	}
	clusterRef := types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Spec.Cluster,
	}
	serviceRef := GetGrafanaServiceKey(cr.ObjKey())
	labels := GetGrafanaLabels(clusterRef)

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
	crSvc := cr.Spec.Grafana.Service
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

func MakeGrafanaPVC(cr *dapi.DorisMonitor, scheme *runtime.Scheme) *corev1.PersistentVolumeClaim {
	if cr.Spec.Cluster == "" {
		return nil
	}
	clusterRef := types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Spec.Cluster,
	}
	pvcRef := GetGrafanaPVCKey(cr.ObjKey())
	labels := GetGrafanaLabels(clusterRef)

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pvcRef.Name,
			Namespace: pvcRef.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources:        cr.Spec.Grafana.ResourceRequirements,
			StorageClassName: cr.Spec.StorageClassName,
		},
	}
	_ = controllerutil.SetOwnerReference(cr, pvc, scheme)
	return pvc
}

func MakeGrafanaDeployment(cr *dapi.DorisMonitor, scheme *runtime.Scheme) *appv1.Deployment {
	if cr.Spec.Cluster == "" {
		return nil
	}
	clusterRef := types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Spec.Cluster,
	}
	deploymentRef := GetGrafanaDeploymentKey(cr.ObjKey())
	configMapRef := GetGrafanaConfigMapKey(cr.ObjKey())
	secretRef := GetGrafanaSecretKey(cr.ObjKey())
	pvcRef := GetGrafanaPVCKey(cr.ObjKey())
	labels := GetGrafanaLabels(clusterRef)

	replicas := int32(1)

	// pod template
	podTemplate := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      labels,
			Annotations: make(map[string]string),
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: cr.Spec.ServiceAccount,
			ImagePullSecrets:   cr.Spec.ImagePullSecrets,
			Volumes: []corev1.Volume{
				{
					Name: "grafana-datasource",
					VolumeSource: util.NewConfigMapItemsVolumeSource(configMapRef.Name, map[string]string{
						"datasource.yml": "doris-cluster-datasource.yml",
					}),
				}, {
					Name: "grafana-dashboard",
					VolumeSource: util.NewConfigMapItemsVolumeSource(configMapRef.Name, map[string]string{
						"dashboards.json": "doris-cluster-dashboards.json",
						"dashboard.yml":   "dashboard.yml",
					}),
				}, {
					Name: "grafana-data",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: pvcRef.Name}}},
			},
			Containers: []corev1.Container{{
				Name:            "grafana",
				Image:           util.StringFallback(cr.Spec.Grafana.Image, DefaultGrafanaImage),
				ImagePullPolicy: cr.Spec.ImagePullPolicy,
				Resources:       cr.Spec.Grafana.ResourceRequirements,
				Ports: []corev1.ContainerPort{{
					Name:          "http-port",
					ContainerPort: 3000,
				}},
				Env: []corev1.EnvVar{
					{Name: "GF_PATHS_DATA", Value: "/data/grafana"},
					{Name: "GF_SECURITY_ADMIN_USER", ValueFrom: util.NewEnvVarSecretSource(secretRef.Name, "user")},
					{Name: "GF_SECURITY_ADMIN_PASSWORD", ValueFrom: util.NewEnvVarSecretSource(secretRef.Name, "password")},
				},
				VolumeMounts: []corev1.VolumeMount{
					{Name: "grafana-datasource", MountPath: "/etc/grafana/provisioning/datasources"},
					{Name: "grafana-dashboard", MountPath: "/etc/grafana/provisioning/dashboards"},
					{Name: "grafana-data", MountPath: "/data/grafana"},
				},
				ReadinessProbe: &corev1.Probe{
					ProbeHandler:     util.NewHttpGetProbeHandler("/api/health", 3000),
					TimeoutSeconds:   5,
					PeriodSeconds:    10,
					SuccessThreshold: 5,
					FailureThreshold: 3,
				},
				LivenessProbe: &corev1.Probe{
					ProbeHandler:     util.NewHttpGetProbeHandler("/api/health", 3000),
					TimeoutSeconds:   5,
					PeriodSeconds:    10,
					SuccessThreshold: 5,
					FailureThreshold: 3,
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
