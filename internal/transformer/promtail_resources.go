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
	DefaultPromtailImage = "grafana/promtail:2.9.1"
)

var (
	PromtailConfTmpl = template.NewTemplateOrPanic("promtail-conf", "monitor/promtail.yml")
)

type PromtailTmplData struct {
	LokiName         string
	LokiNamespace    string
	ClusterName      string
	ClusterNamespace string
}

func GetPromtailLabels(dorisClusterKey types.NamespacedName) map[string]string {
	return MakeResourceLabels(dorisClusterKey.Name, "promtail")
}

func GetPromtailConfigMapKey(monitorKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: monitorKey.Namespace,
		Name:      fmt.Sprintf("%s-promtail-config", monitorKey.Name),
	}
}

func GetPromtailDaemonSetKey(monitorKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: monitorKey.Namespace,
		Name:      fmt.Sprintf("%s-promtail", monitorKey.Name),
	}
}

func MakePromtailConfigMap(cr *dapi.DorisMonitor, scheme *runtime.Scheme) (*corev1.ConfigMap, error) {
	if cr.Spec.Cluster == "" || cr.Spec.DisableLoki {
		return nil, nil
	}
	clusterRef := types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Spec.Cluster,
	}
	configMapRef := GetPromtailConfigMapKey(cr.ObjKey())
	lokiRef := GetLokiServiceKey(cr.ObjKey())
	labels := GetPromtailLabels(clusterRef)

	// gen config template
	tmplData := PromtailTmplData{
		LokiName:         lokiRef.Name,
		LokiNamespace:    lokiRef.Namespace,
		ClusterName:      clusterRef.Name,
		ClusterNamespace: clusterRef.Namespace,
	}
	confContent, err := template.ExecTemplate(PromtailConfTmpl, tmplData)
	if err != nil {
		return nil, util.MergeErrors(fmt.Errorf("fail to parse promtail.conf template"), err)
	}

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapRef.Name,
			Namespace: configMapRef.Namespace,
			Labels:    labels,
		},
		Data: map[string]string{
			"promtail.yaml": confContent,
		},
	}
	_ = controllerutil.SetOwnerReference(cr, configMap, scheme)
	return configMap, nil
}

func MakePromtailDaemonSet(cr *dapi.DorisMonitor, scheme *runtime.Scheme) *appv1.DaemonSet {
	if cr.Spec.Cluster == "" || cr.Spec.DisableLoki {
		return nil
	}
	clusterRef := types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Spec.Cluster,
	}
	daemonSetRef := GetPromtailDaemonSetKey(cr.ObjKey())
	configMapRef := GetPromtailConfigMapKey(cr.ObjKey())
	labels := GetPromtailLabels(clusterRef)

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
					Name: "config",
					VolumeSource: util.NewConfigMapItemsVolumeSource(
						configMapRef.Name, map[string]string{"promtail.yaml": "promtail.yaml"}),
				}, {
					Name:         "run",
					VolumeSource: util.NewHostPathVolumeSource("/run/promtail"),
				}, {
					Name:         "containers",
					VolumeSource: util.NewHostPathVolumeSource("/var/lib/docker/containers"),
				}, {
					Name:         "pods",
					VolumeSource: util.NewHostPathVolumeSource("/var/log/pods"),
				},
			},
			Containers: []corev1.Container{{
				Name:            "promtail",
				Image:           util.StringFallback(cr.Spec.Promtail.Image, DefaultPromtailImage),
				ImagePullPolicy: cr.Spec.ImagePullPolicy,
				Resources:       cr.Spec.Promtail.ResourceRequirements,
				Args:            []string{"-config.file=/etc/promtail/promtail.yaml"},
				Ports: []corev1.ContainerPort{{
					Name:          "http-metrics",
					ContainerPort: 3101,
					Protocol:      corev1.ProtocolTCP,
				}},
				Env: []corev1.EnvVar{{
					Name: "HOSTNAME",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{FieldPath: "spec.nodeName"}},
				}},
				VolumeMounts: []corev1.VolumeMount{
					{Name: "config", MountPath: "/etc/promtail"},
					{Name: "run", MountPath: "/run/promtail"},
					{Name: "containers", MountPath: "/var/lib/docker/containers"},
					{Name: "pods", MountPath: "/var/log/pods", ReadOnly: true},
				},
			}},
		},
	}

	daemonSet := &appv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      daemonSetRef.Name,
			Namespace: daemonSetRef.Namespace,
			Labels:    labels,
		},
		Spec: appv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: podTemplate,
		},
	}
	_ = controllerutil.SetOwnerReference(cr, daemonSet, scheme)
	_ = controllerutil.SetControllerReference(cr, daemonSet, scheme)
	return daemonSet
}
