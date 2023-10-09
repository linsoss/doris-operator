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
	DefaultLokiImage = "grafana/loki:2.9.1"
)

var (
	LokiConfTmpl = template.NewTemplateOrPanic("loki-conf", "monitor/loki.yml")
)

type LokiTmplData struct {
	RetentionDeletesEnables bool
	RetentionPeriod         string
}

func GetLokiLabels(dorisClusterKey types.NamespacedName) map[string]string {
	return MakeResourceLabels(dorisClusterKey.Name, "loki")
}

func GetLokiConfigMapKey(monitorKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: monitorKey.Namespace,
		Name:      fmt.Sprintf("%s-loki-config", monitorKey.Name),
	}
}

func GetLokiServiceKey(monitorKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: monitorKey.Namespace,
		Name:      fmt.Sprintf("%s-loki", monitorKey.Name),
	}
}

func GetLokiPVCKey(monitorKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: monitorKey.Namespace,
		Name:      fmt.Sprintf("%s-loki-pvc", monitorKey.Name),
	}
}

func GetLokiDeploymentKey(monitorKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: monitorKey.Namespace,
		Name:      fmt.Sprintf("%s-loki", monitorKey.Name),
	}
}

func MakeLokiConfigMap(cr *dapi.DorisMonitor, scheme *runtime.Scheme) (*corev1.ConfigMap, error) {
	if cr.Spec.Cluster == "" || cr.Spec.DisableLoki {
		return nil, nil
	}
	clusterRef := types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Spec.Cluster,
	}
	configMapRef := GetLokiConfigMapKey(cr.ObjKey())
	labels := GetLokiLabels(clusterRef)

	// generate loki data
	lokiTmplData := LokiTmplData{
		RetentionDeletesEnables: cr.Spec.Loki.RetentionTime != nil,
		RetentionPeriod:         util.PointerDeRefer(cr.Spec.Loki.RetentionTime, "120h"),
	}
	lokiConfContent, err := template.ExecTemplate(LokiConfTmpl, lokiTmplData)
	if err != nil {
		return nil, util.MergeErrors(fmt.Errorf("fail to parse loki.conf template"), err)
	}

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapRef.Name,
			Namespace: configMapRef.Namespace,
			Labels:    labels,
		},
		Data: map[string]string{
			"loki.yml": lokiConfContent,
		},
	}
	_ = controllerutil.SetOwnerReference(cr, configMap, scheme)
	return configMap, nil
}

func MakeLokiService(cr *dapi.DorisMonitor, scheme *runtime.Scheme) *corev1.Service {
	if cr.Spec.Cluster == "" || cr.Spec.DisableLoki {
		return nil
	}
	clusterRef := types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Spec.Cluster,
	}
	serviceRef := GetLokiServiceKey(cr.ObjKey())
	labels := GetLokiLabels(clusterRef)

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceRef.Name,
			Namespace: serviceRef.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Name: "http-metrics",
				Port: 3100,
			}},
		},
	}
	_ = controllerutil.SetOwnerReference(cr, service, scheme)
	return service
}

func MakeLokiPVC(cr *dapi.DorisMonitor, scheme *runtime.Scheme) *corev1.PersistentVolumeClaim {
	if cr.Spec.Cluster == "" || cr.Spec.DisableLoki {
		return nil
	}
	clusterRef := types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Spec.Cluster,
	}
	pvcRef := GetLokiPVCKey(cr.ObjKey())
	labels := GetLokiLabels(clusterRef)

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pvcRef.Name,
			Namespace: pvcRef.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources:        cr.Spec.Loki.ResourceRequirements,
			StorageClassName: cr.Spec.StorageClassName,
		},
	}
	_ = controllerutil.SetOwnerReference(cr, pvc, scheme)
	return pvc
}

func MakeLokiDeployment(cr *dapi.DorisMonitor, scheme *runtime.Scheme) *appv1.Deployment {
	if cr.Spec.Cluster == "" || cr.Spec.DisableLoki {
		return nil
	}
	clusterRef := types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Spec.Cluster,
	}
	deploymentRef := GetLokiDeploymentKey(cr.ObjKey())
	configMapRef := GetLokiConfigMapKey(cr.ObjKey())
	pvcRef := GetLokiPVCKey(cr.ObjKey())
	labels := GetLokiLabels(clusterRef)
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
					Name: "config",
					VolumeSource: util.NewConfigMapItemsVolumeSource(
						configMapRef.Name, map[string]string{"loki.yml": "loki.yml"}),
				}, {
					Name: "data",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: pvcRef.Name}}},
			},
			Containers: []corev1.Container{{
				Name:            "loki",
				Image:           util.StringFallback(cr.Spec.Loki.Image, DefaultLokiImage),
				ImagePullPolicy: cr.Spec.ImagePullPolicy,
				Resources:       cr.Spec.Loki.ResourceRequirements,
				Args:            []string{"-config.file=/etc/loki/loki.yaml"},
				Ports: []corev1.ContainerPort{
					{
						Name:          "http-metrics",
						ContainerPort: 3100,
						Protocol:      corev1.ProtocolTCP,
					}, {
						Name:          "grpc",
						ContainerPort: 9095,
						Protocol:      corev1.ProtocolTCP,
					}},
				VolumeMounts: []corev1.VolumeMount{
					{Name: "config", MountPath: "/etc/loki"},
					{Name: "data", MountPath: "/data"},
				},
				ReadinessProbe: &corev1.Probe{
					ProbeHandler:        util.NewHttpGetProbeHandler("/ready", 3100),
					InitialDelaySeconds: 45,
					TimeoutSeconds:      1,
					PeriodSeconds:       10,
					SuccessThreshold:    1,
					FailureThreshold:    3,
				},
				LivenessProbe: &corev1.Probe{
					ProbeHandler:     util.NewHttpGetProbeHandler("/ready", 3100),
					TimeoutSeconds:   1,
					PeriodSeconds:    10,
					SuccessThreshold: 1,
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
