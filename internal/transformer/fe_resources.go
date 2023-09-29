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
	"github.com/al-assad/doris-operator/internal/util"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strconv"
)

const (
	DefaultFeHttpPort    = 8030
	DefaultFeEditLogPort = 9010
	DefaultFeRpcPort     = 9020
	DefaultFeQueryPort   = 9030
)

func GetFeComponentLabels(dorisClusterKey types.NamespacedName) map[string]string {
	return MakeResourceLabels(dorisClusterKey.Name, "fe")
}

func GetFeConfigMapKey(dorisClusterKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: dorisClusterKey.Namespace,
		Name:      fmt.Sprintf("%s-fe-config", dorisClusterKey.Name),
	}
}

func GetFeServiceKey(dorisClusterKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: dorisClusterKey.Namespace,
		Name:      fmt.Sprintf("%s-fe", dorisClusterKey.Name),
	}
}

func GetFePeerServiceKey(dorisClusterKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: dorisClusterKey.Namespace,
		Name:      fmt.Sprintf("%s-fe-peer", dorisClusterKey.Name),
	}
}

func GetFeStatefulSetKey(dorisClusterKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: dorisClusterKey.Namespace,
		Name:      fmt.Sprintf("%s-fe", dorisClusterKey.Name),
	}
}

func GetFeImage(r *dapi.DorisCluster) string {
	version := util.StringFallback(r.Spec.FE.Version, r.Spec.Version)
	return fmt.Sprintf("%s:%s", r.Spec.FE.BaseImage, version)
}

func GetFeHttpPort(cr *dapi.DorisCluster) int32 {
	if cr.Spec.FE == nil {
		return DefaultFeHttpPort
	}
	return getPortValueFromRawConf(cr.Spec.FE.Configs, "http_port", DefaultFeHttpPort)
}

func GetFeQueryPort(cr *dapi.DorisCluster) int32 {
	if cr.Spec.FE == nil {
		return DefaultFeQueryPort
	}
	return getPortValueFromRawConf(cr.Spec.FE.Configs, "query_port", DefaultFeQueryPort)
}

func GetFeRpcPort(cr *dapi.DorisCluster) int32 {
	if cr.Spec.FE == nil {
		return DefaultFeRpcPort
	}
	return getPortValueFromRawConf(cr.Spec.FE.Configs, "query_port", DefaultFeRpcPort)
}

func GetFeEditLogPort(cr *dapi.DorisCluster) int32 {
	if cr.Spec.FE == nil {
		return DefaultFeEditLogPort
	}
	return getPortValueFromRawConf(cr.Spec.FE.Configs, "edit_log_port", DefaultFeEditLogPort)
}

func MakeFeConfigMap(cr *dapi.DorisCluster, scheme *runtime.Scheme) *corev1.ConfigMap {
	if cr.Spec.FE == nil {
		return nil
	}
	configMapRef := GetFeConfigMapKey(cr.ObjKey())
	data := map[string]string{
		"fe.conf": dumpJavaBasedComponentConf(cr.Spec.FE.Configs),
	}
	// merge hadoop config data
	if cr.Spec.HadoopConf != nil {
		data = util.MergeMaps(cr.Spec.HadoopConf.Config, data)
	}
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapRef.Name,
			Namespace: configMapRef.Namespace,
			Labels:    GetFeComponentLabels(cr.ObjKey()),
		},
		Data: data,
	}
	_ = controllerutil.SetOwnerReference(cr, configMap, scheme)
	return configMap
}

func MakeFeService(cr *dapi.DorisCluster, scheme *runtime.Scheme) *corev1.Service {
	if cr.Spec.FE == nil {
		return nil
	}
	serviceRef := GetFeServiceKey(cr.ObjKey())
	feLabels := GetFeComponentLabels(cr.ObjKey())
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceRef.Name,
			Namespace: serviceRef.Namespace,
			Labels:    feLabels,
		},
		Spec: corev1.ServiceSpec{
			Selector: feLabels,
			Type:     corev1.ServiceTypeClusterIP,
		},
	}
	httpPort := corev1.ServicePort{
		Name: "http-port",
		Port: GetFeHttpPort(cr),
	}
	queryPort := corev1.ServicePort{
		Name: "query-port",
		Port: GetFeQueryPort(cr),
	}
	// When the user specifies a service type
	crSvc := cr.Spec.FE.Service
	if crSvc != nil {
		if crSvc.Type != "" {
			service.Spec.Type = crSvc.Type
		}
		if crSvc.ExternalTrafficPolicy != nil {
			service.Spec.ExternalTrafficPolicy = *crSvc.ExternalTrafficPolicy
		}
		if cr.Spec.FE.Service.QueryPort != nil {
			httpPort.NodePort = *crSvc.QueryPort
		}
		if cr.Spec.FE.Service.HttpPort != nil {
			queryPort.NodePort = *crSvc.HttpPort
		}
	}
	service.Spec.Ports = []corev1.ServicePort{httpPort, queryPort}
	_ = controllerutil.SetOwnerReference(cr, service, scheme)
	return service
}

func MakeFePeerService(cr *dapi.DorisCluster, scheme *runtime.Scheme) *corev1.Service {
	if cr.Spec.FE == nil {
		return nil
	}
	serviceRef := GetFePeerServiceKey(cr.ObjKey())
	feLabels := GetFeComponentLabels(cr.ObjKey())
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceRef.Name,
			Namespace: serviceRef.Namespace,
			Labels:    feLabels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Name: "http-port", Port: GetFeHttpPort(cr)},
				{Name: "edit-log-port", Port: GetFeEditLogPort(cr)},
				{Name: "rpc-port", Port: GetFeRpcPort(cr)},
			},
			Selector:  feLabels,
			ClusterIP: "None",
		},
	}
	_ = controllerutil.SetOwnerReference(cr, service, scheme)
	return service
}

func MakeFeStatefulSet(cr *dapi.DorisCluster, scheme *runtime.Scheme) *appv1.StatefulSet {
	if cr.Spec.FE == nil {
		return nil
	}
	statefulSetRef := GetFeStatefulSetKey(cr.ObjKey())
	configMapRef := GetFeConfigMapKey(cr.ObjKey())
	accountSecretRef := GetOprSqlAccountSecretKey(cr.ObjKey())
	feLabels := GetFeComponentLabels(cr.ObjKey())

	// volume claim template
	pvcTemplate := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: "fe-meta",
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			StorageClassName: cr.Spec.FE.StorageClassName,
		},
	}
	storageRequest := cr.Spec.FE.Requests.Storage()
	if storageRequest != nil {
		pvcTemplate.Spec.Resources.Requests = corev1.ResourceList{
			corev1.ResourceStorage: *storageRequest,
		}
	}

	// pod template: volumes
	volumes := []corev1.Volume{
		{Name: "conf", VolumeSource: util.NewConfigMapVolumeSource(configMapRef.Name)},
		{Name: "fe-log", VolumeSource: util.NewEmptyDirVolumeSource()},
	}
	// merge addition volumes defined by user
	volumes = append(volumes, cr.Spec.FE.AdditionalVolumes...)

	// pod template: main container
	mainContainer := corev1.Container{
		Name:            "fe",
		Image:           GetFeImage(cr),
		ImagePullPolicy: cr.Spec.ImagePullPolicy,
		Resources: corev1.ResourceRequirements{
			Requests: cr.Spec.FE.Requests,
		},
		Ports: []corev1.ContainerPort{
			{Name: "http-port", ContainerPort: GetFeHttpPort(cr)},
			{Name: "edit-log-port", ContainerPort: GetFeEditLogPort(cr)},
			{Name: "rpc-port", ContainerPort: GetFeRpcPort(cr)},
			{Name: "query-port", ContainerPort: GetFeQueryPort(cr)},
		},
		Env: []corev1.EnvVar{
			{Name: "FE_SVC", Value: GetFeServiceKey(cr.ObjKey()).Name},
			{Name: "ACC_USER", ValueFrom: util.NewEnvVarSecretSource(accountSecretRef.Name, "user")},
			{Name: "ACC_PWD", ValueFrom: util.NewEnvVarSecretSource(accountSecretRef.Name, "password")},
		},
		VolumeMounts: []corev1.VolumeMount{
			{Name: "conf", MountPath: "/etc/apache-doris/fe/"},
			{Name: "fe-meta", MountPath: "/opt/apache-doris/fe/doris-meta"},
			{Name: "fe-log", MountPath: "/opt/apache-doris/fe/log"},
		},
		ReadinessProbe: &corev1.Probe{
			ProbeHandler:        util.NewTcpSocketProbeHandler(GetFeQueryPort(cr)),
			InitialDelaySeconds: 3,
			TimeoutSeconds:      1,
			PeriodSeconds:       5,
			SuccessThreshold:    1,
			FailureThreshold:    3,
		},
	}
	// pod template: merge additional pod containers configs defined by user
	mainContainer.Env = append(mainContainer.Env, cr.Spec.FE.AdditionalEnvs...)
	mainContainer.VolumeMounts = append(mainContainer.VolumeMounts, cr.Spec.FE.AdditionalVolumeMounts...)
	containers := append([]corev1.Container{mainContainer}, cr.Spec.FE.AdditionalContainers...)

	// pod template: host alias
	var hostAlias []corev1.HostAlias
	if cr.Spec.HadoopConf != nil {
		hostAlias = mergeHostAlias(cr.Spec.HadoopConf.Hosts, cr.Spec.FE.HostAliases)
	} else {
		hostAlias = cr.Spec.FE.HostAliases
	}

	// pod template
	podTemplate := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: feLabels,
			Annotations: map[string]string{
				PrometheusPathAnnoKey:   "/metrics",
				PrometheusPortAnnoKey:   strconv.Itoa(int(GetFeHttpPort(cr))),
				PrometheusScrapeAnnoKey: "true",
			},
		},
		Spec: corev1.PodSpec{
			Volumes:            volumes,
			Containers:         containers,
			ImagePullSecrets:   cr.Spec.ImagePullSecrets,
			ServiceAccountName: util.StringFallback(cr.Spec.FE.ServiceAccount, cr.Spec.ServiceAccount),
			Affinity:           util.PointerFallback(cr.Spec.FE.Affinity, cr.Spec.Affinity),
			Tolerations:        util.ArrayFallback(cr.Spec.FE.Tolerations, cr.Spec.Tolerations),
			PriorityClassName:  util.StringFallback(cr.Spec.FE.PriorityClassName, cr.Spec.PriorityClassName),
			HostAliases:        hostAlias,
		},
	}

	// update strategy
	updateStg := appv1.StatefulSetUpdateStrategy{
		Type: util.PointerFallbackAndDeRefer(
			cr.Spec.FE.StatefulSetUpdateStrategy, cr.Spec.StatefulSetUpdateStrategy,
			appv1.RollingUpdateStatefulSetStrategyType),
	}

	// statefulset
	statefulSet := &appv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      statefulSetRef.Name,
			Namespace: statefulSetRef.Namespace,
			Labels:    feLabels,
		},
		Spec: appv1.StatefulSetSpec{
			Replicas:             &cr.Spec.FE.Replicas,
			ServiceName:          GetFePeerServiceKey(cr.ObjKey()).Name,
			Selector:             &metav1.LabelSelector{MatchLabels: feLabels},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{pvcTemplate},
			Template:             podTemplate,
			UpdateStrategy:       updateStg,
		},
	}

	_ = controllerutil.SetOwnerReference(cr, statefulSet, scheme)
	return statefulSet
}
