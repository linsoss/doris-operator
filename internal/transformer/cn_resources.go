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
	CnProbeTimeoutSec = 200
)

func GetCnComponentLabels(dorisClusterKey types.NamespacedName) map[string]string {
	return MakeResourceLabels(dorisClusterKey.Name, "cn")
}

func GetCnConfigMapKey(dorisClusterKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: dorisClusterKey.Namespace,
		Name:      fmt.Sprintf("%s-cn-config", dorisClusterKey.Name),
	}
}

func GetCnServiceKey(dorisClusterKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: dorisClusterKey.Namespace,
		Name:      fmt.Sprintf("%s-cn", dorisClusterKey.Name),
	}
}

func GetCnPeerServiceKey(dorisClusterKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: dorisClusterKey.Namespace,
		Name:      fmt.Sprintf("%s-cn-peer", dorisClusterKey.Name),
	}
}

func GetCnStatefulSetKey(dorisClusterKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: dorisClusterKey.Namespace,
		Name:      fmt.Sprintf("%s-cn", dorisClusterKey.Name),
	}
}

func GetCnImage(r *dapi.DorisCluster) string {
	version := util.StringFallback(r.Spec.CN.Version, r.Spec.Version)
	return fmt.Sprintf("%s:%s", r.Spec.CN.BaseImage, version)
}

func GetCnHeartbeatServicePort(cr *dapi.DorisCluster) int32 {
	if cr.Spec.CN == nil {
		return DefaultBeHeartbeatServicePort
	}
	return getPortValueFromRawConf(cr.Spec.CN.Configs, "be_port", DefaultBeHeartbeatServicePort)
}

func GetCnPort(cr *dapi.DorisCluster) int32 {
	if cr.Spec.CN == nil {
		return DefaultBePort
	}
	return getPortValueFromRawConf(cr.Spec.CN.Configs, "be_port", DefaultBePort)
}

func GetCnWebserverPort(cr *dapi.DorisCluster) int32 {
	if cr.Spec.CN == nil {
		return DefaultBeWebserverPort
	}
	return getPortValueFromRawConf(cr.Spec.CN.Configs, "webserver_port", DefaultBeWebserverPort)
}

func GetCnBrpcPort(cr *dapi.DorisCluster) int32 {
	if cr.Spec.CN == nil {
		return DefaultBeBrpcPort
	}
	return getPortValueFromRawConf(cr.Spec.CN.Configs, "brpc_port", DefaultBeBrpcPort)
}

func GetCnExpectPodNames(dorisClusterKey types.NamespacedName, replicas int32) []string {
	stsName := GetCnStatefulSetKey(dorisClusterKey).Name
	var expectPods []string
	for i := 0; i < int(replicas); i++ {
		expectPods = append(expectPods, fmt.Sprintf("%s-%d", stsName, i))
	}
	return expectPods
}

func MakeCnConfigMap(cr *dapi.DorisCluster, scheme *runtime.Scheme) *corev1.ConfigMap {
	if cr.Spec.CN == nil {
		return nil
	}
	configMapRef := GetCnConfigMapKey(cr.ObjKey())
	data := map[string]string{
		"be.conf": dumpCppBasedComponentConf(cr.Spec.CN.Configs),
	}
	// merge hadoop config data
	if cr.Spec.HadoopConf != nil {
		data = util.MergeMaps(cr.Spec.HadoopConf.Config, data)
	}
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapRef.Name,
			Namespace: configMapRef.Namespace,
			Labels:    GetCnComponentLabels(cr.ObjKey()),
		},
		Data: data,
	}
	_ = controllerutil.SetOwnerReference(cr, configMap, scheme)
	return configMap
}

func MakeCnService(cr *dapi.DorisCluster, scheme *runtime.Scheme) *corev1.Service {
	if cr.Spec.CN == nil {
		return nil
	}
	serviceRef := GetCnServiceKey(cr.ObjKey())
	cnLabels := GetCnComponentLabels(cr.ObjKey())
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceRef.Name,
			Namespace: serviceRef.Namespace,
			Labels:    cnLabels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: "webserver-port",
					Port: GetCnWebserverPort(cr),
				},
			},
			Selector: cnLabels,
			Type:     corev1.ServiceTypeClusterIP,
		},
	}
	_ = controllerutil.SetOwnerReference(cr, service, scheme)
	return service
}

func MakeCnPeerService(cr *dapi.DorisCluster, scheme *runtime.Scheme) *corev1.Service {
	if cr.Spec.CN == nil {
		return nil
	}
	serviceRef := GetCnPeerServiceKey(cr.ObjKey())
	cnLabels := GetCnComponentLabels(cr.ObjKey())
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceRef.Name,
			Namespace: serviceRef.Namespace,
			Labels:    cnLabels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Name: "webserver-port", Port: GetCnWebserverPort(cr)},
				{Name: "heart-port", Port: GetCnHeartbeatServicePort(cr)},
				{Name: "be-port", Port: GetCnPort(cr)},
				{Name: "brpc-port", Port: GetCnBrpcPort(cr)},
			},
			Selector:  cnLabels,
			ClusterIP: "None",
		},
	}
	_ = controllerutil.SetOwnerReference(cr, service, scheme)
	return service
}

func MakeCnStatefulSet(cr *dapi.DorisCluster, scheme *runtime.Scheme) *appv1.StatefulSet {
	if cr.Spec.CN == nil {
		return nil
	}
	statefulSetRef := GetCnStatefulSetKey(cr.ObjKey())
	accountSecretRef := GetOprSqlAccountSecretKey(cr.ObjKey())
	configMapRef := GetCnConfigMapKey(cr.ObjKey())
	cnLabels := GetCnComponentLabels(cr.ObjKey())

	// pod template: volumes
	volumes := []corev1.Volume{
		{Name: "conf", VolumeSource: util.NewConfigMapVolumeSource(configMapRef.Name)},
		{Name: "cn-log", VolumeSource: util.NewEmptyDirVolumeSource()},
	}
	// merge addition volumes defined by user
	volumes = append(volumes, cr.Spec.CN.AdditionalVolumes...)

	// pod template: main container
	mainContainer := corev1.Container{
		Name:            "cn",
		Image:           GetCnImage(cr),
		ImagePullPolicy: cr.Spec.ImagePullPolicy,
		Resources: corev1.ResourceRequirements{
			Requests: cr.Spec.CN.Requests,
		},
		Ports: []corev1.ContainerPort{
			{Name: "webserver-port", ContainerPort: GetCnWebserverPort(cr)},
			{Name: "heart-port", ContainerPort: GetCnHeartbeatServicePort(cr)},
			{Name: "be-port", ContainerPort: GetCnPort(cr)},
			{Name: "brpc-port", ContainerPort: GetCnBrpcPort(cr)},
		},
		Env: []corev1.EnvVar{
			{Name: "FE_SVC", Value: GetFeServiceKey(cr.ObjKey()).Name},
			{Name: "FE_QUERY_PORT", Value: strconv.Itoa(int(GetFeQueryPort(cr)))},
			{Name: "ACC_USER", ValueFrom: util.NewEnvVarSecretSource(accountSecretRef.Name, "user")},
			{Name: "ACC_PWD", ValueFrom: util.NewEnvVarSecretSource(accountSecretRef.Name, "password")},
			{Name: "CN_PROBE_TIMEOUT", Value: strconv.Itoa(CnProbeTimeoutSec)},
		},
		VolumeMounts: []corev1.VolumeMount{
			{Name: "conf", MountPath: "/etc/apache-doris/be/"},
			{Name: "cn-log", MountPath: "/opt/apache-doris/be/log"},
		},
		ReadinessProbe: &corev1.Probe{
			ProbeHandler:     util.NewTcpSocketProbeHandler(GetCnHeartbeatServicePort(cr)),
			TimeoutSeconds:   1,
			PeriodSeconds:    5,
			SuccessThreshold: 1,
			FailureThreshold: 3,
		},
	}
	// pod template: merge additional pod containers configs defined by user
	mainContainer.Env = append(mainContainer.Env, cr.Spec.CN.AdditionalEnvs...)
	mainContainer.VolumeMounts = append(mainContainer.VolumeMounts, cr.Spec.CN.AdditionalVolumeMounts...)
	containers := append([]corev1.Container{mainContainer}, cr.Spec.CN.AdditionalContainers...)

	// pod template: host alias
	var hostAlias []corev1.HostAlias
	if cr.Spec.HadoopConf != nil {
		hostAlias = mergeHostAlias(cr.Spec.HadoopConf.Hosts, cr.Spec.CN.HostAliases)
	} else {
		hostAlias = cr.Spec.CN.HostAliases
	}

	// pod template
	podTemplate := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      cnLabels,
			Annotations: MakePrometheusAnnotations("/metrics", GetCnHeartbeatServicePort(cr)),
		},
		Spec: corev1.PodSpec{
			Volumes:            volumes,
			Containers:         containers,
			ImagePullSecrets:   cr.Spec.ImagePullSecrets,
			ServiceAccountName: util.StringFallback(cr.Spec.CN.ServiceAccount, cr.Spec.ServiceAccount),
			Affinity:           util.PointerFallback(cr.Spec.CN.Affinity, cr.Spec.Affinity),
			Tolerations:        util.ArrayFallback(cr.Spec.CN.Tolerations, cr.Spec.Tolerations),
			PriorityClassName:  util.StringFallback(cr.Spec.CN.PriorityClassName, cr.Spec.PriorityClassName),
			HostAliases:        hostAlias,
		},
	}

	// update strategy
	updateStg := appv1.StatefulSetUpdateStrategy{
		Type: util.PointerFallbackAndDeRefer(
			cr.Spec.CN.StatefulSetUpdateStrategy,
			cr.Spec.StatefulSetUpdateStrategy,
			appv1.RollingUpdateStatefulSetStrategyType),
	}

	// statefulset
	statefulSet := &appv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      statefulSetRef.Name,
			Namespace: statefulSetRef.Namespace,
			Labels:    cnLabels,
		},
		Spec: appv1.StatefulSetSpec{
			Replicas:       &cr.Spec.CN.Replicas,
			ServiceName:    GetCnPeerServiceKey(cr.ObjKey()).Name,
			Selector:       &metav1.LabelSelector{MatchLabels: cnLabels},
			Template:       podTemplate,
			UpdateStrategy: updateStg,
		},
	}

	_ = controllerutil.SetOwnerReference(cr, statefulSet, scheme)
	return statefulSet
}
