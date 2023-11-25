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
	DefaultBeHeartbeatServicePort = 9050
	DefaultBePort                 = 9060
	DefaultBeWebserverPort        = 8040
	DefaultBeBrpcPort             = 8060

	BeProbeTimeoutSec = 200
)

func GetBeComponentLabels(dorisClusterKey types.NamespacedName) map[string]string {
	return MakeResourceLabels(dorisClusterKey.Name, "be")
}

func GetBeConfigMapKey(dorisClusterKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: dorisClusterKey.Namespace,
		Name:      fmt.Sprintf("%s-be-config", dorisClusterKey.Name),
	}
}

func GetBeServiceKey(dorisClusterKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: dorisClusterKey.Namespace,
		Name:      fmt.Sprintf("%s-be", dorisClusterKey.Name),
	}
}

func GetBePeerServiceKey(dorisClusterKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: dorisClusterKey.Namespace,
		Name:      fmt.Sprintf("%s-be-peer", dorisClusterKey.Name),
	}
}

func GetBeStatefulSetKey(dorisClusterKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: dorisClusterKey.Namespace,
		Name:      fmt.Sprintf("%s-be", dorisClusterKey.Name),
	}
}

func GetBeImage(r *dapi.DorisCluster) string {
	version := util.StringFallback(r.Spec.BE.Version, r.Spec.Version)
	return fmt.Sprintf("%s:%s", r.Spec.BE.BaseImage, version)
}

func GetBeHeartbeatServicePort(cr *dapi.DorisCluster) int32 {
	if cr.Spec.BE == nil {
		return DefaultBeHeartbeatServicePort
	}
	return getPortValueFromRawConf(cr.Spec.BE.Configs, "be_port", DefaultBeHeartbeatServicePort)
}

func GetBePort(cr *dapi.DorisCluster) int32 {
	if cr.Spec.BE == nil {
		return DefaultBePort
	}
	return getPortValueFromRawConf(cr.Spec.BE.Configs, "be_port", DefaultBePort)
}

func GetBeWebserverPort(cr *dapi.DorisCluster) int32 {
	if cr.Spec.BE == nil {
		return DefaultBeWebserverPort
	}
	return getPortValueFromRawConf(cr.Spec.BE.Configs, "webserver_port", DefaultBeWebserverPort)
}

func GetBeBrpcPort(cr *dapi.DorisCluster) int32 {
	if cr.Spec.BE == nil {
		return DefaultBeBrpcPort
	}
	return getPortValueFromRawConf(cr.Spec.BE.Configs, "brpc_port", DefaultBeBrpcPort)
}

func GetBeExpectPodNames(dorisClusterKey types.NamespacedName, replicas int32) []string {
	stsName := GetBeStatefulSetKey(dorisClusterKey).Name
	var expectPods []string
	for i := 0; i < int(replicas); i++ {
		expectPods = append(expectPods, fmt.Sprintf("%s-%d", stsName, i))
	}
	return expectPods
}

func MakeBeConfigMap(cr *dapi.DorisCluster, scheme *runtime.Scheme) *corev1.ConfigMap {
	if cr.Spec.BE == nil {
		return nil
	}

	configs := util.MapFallback(cr.Spec.BE.Configs, make(map[string]string))
	configs["be_node_role"] = "mix"
	configMapRef := GetBeConfigMapKey(cr.ObjKey())
	data := map[string]string{
		"be.conf": dumpCppBasedComponentConf(configs),
	}
	// merge hadoop config data
	if cr.Spec.HadoopConf != nil {
		data = util.MergeMaps(cr.Spec.HadoopConf.Config, data)
	}
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapRef.Name,
			Namespace: configMapRef.Namespace,
			Labels:    GetBeComponentLabels(cr.ObjKey()),
		},
		Data: data,
	}
	_ = controllerutil.SetOwnerReference(cr, configMap, scheme)
	return configMap
}

func MakeBeService(cr *dapi.DorisCluster, scheme *runtime.Scheme) *corev1.Service {
	if cr.Spec.BE == nil {
		return nil
	}
	serviceRef := GetBeServiceKey(cr.ObjKey())
	beLabels := GetBeComponentLabels(cr.ObjKey())
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceRef.Name,
			Namespace: serviceRef.Namespace,
			Labels:    beLabels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: "webserver-port",
					Port: GetBeWebserverPort(cr),
				},
			},
			Selector: beLabels,
			Type:     corev1.ServiceTypeClusterIP,
		},
	}
	_ = controllerutil.SetOwnerReference(cr, service, scheme)
	return service
}

func MakeBePeerService(cr *dapi.DorisCluster, scheme *runtime.Scheme) *corev1.Service {
	if cr.Spec.BE == nil {
		return nil
	}
	serviceRef := GetBePeerServiceKey(cr.ObjKey())
	beLabels := GetBeComponentLabels(cr.ObjKey())
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceRef.Name,
			Namespace: serviceRef.Namespace,
			Labels:    beLabels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Name: "webserver-port", Port: GetBeWebserverPort(cr)},
				{Name: "heart-port", Port: GetBeHeartbeatServicePort(cr)},
				{Name: "be-port", Port: GetBePort(cr)},
				{Name: "brpc-port", Port: GetBeBrpcPort(cr)},
			},
			Selector:  beLabels,
			ClusterIP: "None",
		},
	}
	_ = controllerutil.SetOwnerReference(cr, service, scheme)
	return service
}

func MakeBeStatefulSet(cr *dapi.DorisCluster, scheme *runtime.Scheme) *appv1.StatefulSet {
	if cr.Spec.BE == nil {
		return nil
	}
	statefulSetRef := GetBeStatefulSetKey(cr.ObjKey())
	accountSecretRef := GetOprSqlAccountSecretKey(cr.ObjKey())
	beLabels := GetBeComponentLabels(cr.ObjKey())

	// volume claim template
	pvcTemplate := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: "be-storage",
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			StorageClassName: cr.Spec.BE.StorageClassName,
		},
	}
	storageRequest := cr.Spec.BE.Requests.Storage()
	if storageRequest != nil {
		pvcTemplate.Spec.Resources.Requests = corev1.ResourceList{
			corev1.ResourceStorage: *storageRequest,
		}
	}

	// pod template: volumes
	volumes := []corev1.Volume{
		{Name: "conf", VolumeSource: util.NewConfigMapVolumeSource(GetBeConfigMapKey(cr.ObjKey()).Name)},
		{Name: "be-log", VolumeSource: util.NewEmptyDirVolumeSource()},
	}
	// merge addition volumes defined by user
	volumes = append(volumes, cr.Spec.BE.AdditionalVolumes...)

	// pod template: main container
	mainContainer := corev1.Container{
		Name:            "be",
		Image:           GetBeImage(cr),
		ImagePullPolicy: cr.Spec.ImagePullPolicy,
		Resources:       formatContainerResourcesRequirement(cr.Spec.BE.ResourceRequirements),
		Ports: []corev1.ContainerPort{
			{Name: "webserver-port", ContainerPort: GetBeWebserverPort(cr)},
			{Name: "heart-port", ContainerPort: GetBeHeartbeatServicePort(cr)},
			{Name: "be-port", ContainerPort: GetBePort(cr)},
			{Name: "brpc-port", ContainerPort: GetBeBrpcPort(cr)},
		},
		Env: []corev1.EnvVar{
			{Name: "FE_SVC", Value: GetFeServiceKey(cr.ObjKey()).Name},
			{Name: "FE_QUERY_PORT", Value: strconv.Itoa(int(GetFeQueryPort(cr)))},
			{Name: "ACC_USER", ValueFrom: util.NewEnvVarSecretSource(accountSecretRef.Name, "user")},
			{Name: "ACC_PWD", ValueFrom: util.NewEnvVarSecretSource(accountSecretRef.Name, "password")},
			{Name: "BE_PROBE_TIMEOUT", Value: strconv.Itoa(BeProbeTimeoutSec)},
		},
		VolumeMounts: []corev1.VolumeMount{
			{Name: "conf", MountPath: "/etc/apache-doris/be/"},
			{Name: "be-storage", MountPath: "/opt/apache-doris/be/storage"},
			{Name: "be-log", MountPath: "/opt/apache-doris/be/log"},
		},
		ReadinessProbe: &corev1.Probe{
			ProbeHandler:     util.NewTcpSocketProbeHandler(GetBeHeartbeatServicePort(cr)),
			TimeoutSeconds:   1,
			PeriodSeconds:    5,
			SuccessThreshold: 1,
			FailureThreshold: 3,
		},
		LivenessProbe: &corev1.Probe{
			ProbeHandler:        util.NewTcpSocketProbeHandler(GetBeHeartbeatServicePort(cr)),
			InitialDelaySeconds: 20,
			TimeoutSeconds:      1,
			PeriodSeconds:       5,
			SuccessThreshold:    1,
			FailureThreshold:    5,
		},
	}
	// pod template: init container
	privileged := true
	initContainer := corev1.Container{
		Name:            "sysctl",
		Image:           GetBusyBoxImage(cr),
		Command:         []string{"sysctl", "-w", "vm.max_map_count=2000000"},
		SecurityContext: &corev1.SecurityContext{Privileged: &privileged},
	}
	// pod template: merge additional pod containers configs defined by user
	mainContainer.Env = append(mainContainer.Env, cr.Spec.BE.AdditionalEnvs...)
	mainContainer.VolumeMounts = append(mainContainer.VolumeMounts, cr.Spec.BE.AdditionalVolumeMounts...)
	containers := append([]corev1.Container{mainContainer}, cr.Spec.BE.AdditionalContainers...)

	// pod template: host alias
	var hostAlias []corev1.HostAlias
	if cr.Spec.HadoopConf != nil {
		hostAlias = mergeHostAlias(cr.Spec.HadoopConf.Hosts, cr.Spec.BE.HostAliases)
	} else {
		hostAlias = cr.Spec.BE.HostAliases
	}

	// pod template
	podTemplate := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      beLabels,
			Annotations: MakePrometheusAnnotations("/metrics", GetBeWebserverPort(cr)),
		},
		Spec: corev1.PodSpec{
			Volumes:            volumes,
			Containers:         containers,
			InitContainers:     []corev1.Container{initContainer},
			ImagePullSecrets:   cr.Spec.ImagePullSecrets,
			ServiceAccountName: util.StringFallback(cr.Spec.BE.ServiceAccount, cr.Spec.ServiceAccount),
			NodeSelector:       util.MapFallback(cr.Spec.BE.NodeSelector, cr.Spec.NodeSelector),
			Affinity:           util.PointerFallback(cr.Spec.BE.Affinity, cr.Spec.Affinity),
			Tolerations:        util.ArrayFallback(cr.Spec.BE.Tolerations, cr.Spec.Tolerations),
			PriorityClassName:  util.StringFallback(cr.Spec.BE.PriorityClassName, cr.Spec.PriorityClassName),
			HostAliases:        hostAlias,
		},
	}

	// update strategy
	updateStg := appv1.StatefulSetUpdateStrategy{
		Type: util.PointerFallbackAndDeRefer(
			cr.Spec.BE.StatefulSetUpdateStrategy,
			cr.Spec.StatefulSetUpdateStrategy,
			appv1.RollingUpdateStatefulSetStrategyType),
	}

	// statefulset
	statefulSet := &appv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      statefulSetRef.Name,
			Namespace: statefulSetRef.Namespace,
			Labels:    beLabels,
		},
		Spec: appv1.StatefulSetSpec{
			Replicas:             &cr.Spec.BE.Replicas,
			ServiceName:          GetBePeerServiceKey(cr.ObjKey()).Name,
			Selector:             &metav1.LabelSelector{MatchLabels: beLabels},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{pvcTemplate},
			Template:             podTemplate,
			UpdateStrategy:       updateStg,
		},
	}

	_ = controllerutil.SetOwnerReference(cr, statefulSet, scheme)
	_ = controllerutil.SetControllerReference(cr, statefulSet, scheme)
	return statefulSet
}
