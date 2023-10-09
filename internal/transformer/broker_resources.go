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
	"strconv"
)

const (
	DefaultBrokerIpcPort  = 8000
	BrokerProbeTimeoutSec = 100
)

var DefaultBrokerLog4jContent = template.ReadOrPanic("broker/log4j.properties")

func GetBrokerComponentLabels(dorisClusterKey types.NamespacedName) map[string]string {
	return MakeResourceLabels(dorisClusterKey.Name, "broker")
}

func GetBrokerConfigMapKey(dorisClusterKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: dorisClusterKey.Namespace,
		Name:      fmt.Sprintf("%s-broker-config", dorisClusterKey.Name),
	}
}

func GetBrokerPeerServiceKey(dorisClusterKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: dorisClusterKey.Namespace,
		Name:      fmt.Sprintf("%s-broker-peer", dorisClusterKey.Name),
	}
}

func GetBrokerStatefulSetKey(dorisClusterKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: dorisClusterKey.Namespace,
		Name:      fmt.Sprintf("%s-broker", dorisClusterKey.Name),
	}
}

func GetBrokerImage(r *dapi.DorisCluster) string {
	version := util.StringFallback(r.Spec.Broker.Version, r.Spec.Version)
	return fmt.Sprintf("%s:%s", r.Spec.Broker.BaseImage, version)
}

func GetBrokerIpcPort(cr *dapi.DorisCluster) int32 {
	if cr.Spec.Broker == nil {
		return DefaultBrokerIpcPort
	}
	return getPortValueFromRawConf(cr.Spec.Broker.Configs, "broker_ipc_port", DefaultBrokerIpcPort)
}

func GetBrokerExpectPodNames(dorisClusterKey types.NamespacedName, replicas int32) []string {
	stsName := GetBrokerStatefulSetKey(dorisClusterKey).Name
	var expectPods []string
	for i := 0; i < int(replicas); i++ {
		expectPods = append(expectPods, fmt.Sprintf("%s-%d", stsName, i))
	}
	return expectPods
}

func MakeBrokerConfigMap(cr *dapi.DorisCluster, scheme *runtime.Scheme) *corev1.ConfigMap {
	if cr.Spec.Broker == nil {
		return nil
	}
	configMapRef := GetBrokerConfigMapKey(cr.ObjKey())
	configs := util.MapFallback(cr.Spec.Broker.Configs, make(map[string]string))
	data := map[string]string{
		"apache_hdfs_broker.conf": dumpJavaBasedComponentConf(configs),
		"log4j.properties":        DefaultBrokerLog4jContent,
	}
	// merge hadoop config data
	if cr.Spec.HadoopConf != nil {
		data = util.MergeMaps(cr.Spec.HadoopConf.Config, data)
	}
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapRef.Name,
			Namespace: configMapRef.Namespace,
			Labels:    GetBrokerComponentLabels(cr.ObjKey()),
		},
		Data: data,
	}
	_ = controllerutil.SetOwnerReference(cr, configMap, scheme)
	return configMap
}

func MakeBrokerPeerService(cr *dapi.DorisCluster, scheme *runtime.Scheme) *corev1.Service {
	if cr.Spec.Broker == nil {
		return nil
	}
	serviceRef := GetBrokerPeerServiceKey(cr.ObjKey())
	brokerLabels := GetBrokerComponentLabels(cr.ObjKey())
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceRef.Name,
			Namespace: serviceRef.Namespace,
			Labels:    brokerLabels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Name: "ipc-port", Port: GetBrokerIpcPort(cr)},
			},
			Selector:  brokerLabels,
			ClusterIP: "None",
		},
	}
	_ = controllerutil.SetOwnerReference(cr, service, scheme)
	return service
}

func MakeBrokerStatefulSet(cr *dapi.DorisCluster, scheme *runtime.Scheme) *appv1.StatefulSet {
	if cr.Spec.Broker == nil {
		return nil
	}
	statefulSetRef := GetBrokerStatefulSetKey(cr.ObjKey())
	accountSecretRef := GetOprSqlAccountSecretKey(cr.ObjKey())
	brokerLabels := GetBrokerComponentLabels(cr.ObjKey())

	// pod template: volumes
	volumes := []corev1.Volume{
		{Name: "conf", VolumeSource: util.NewConfigMapVolumeSource(GetBrokerConfigMapKey(cr.ObjKey()).Name)},
	}
	// merge addition volumes defined by user
	volumes = append(volumes, cr.Spec.Broker.AdditionalVolumes...)

	// pod template: main container
	mainContainer := corev1.Container{
		Name:            "broker",
		Image:           GetBrokerImage(cr),
		ImagePullPolicy: cr.Spec.ImagePullPolicy,
		Resources:       formatContainerResourcesRequirement(cr.Spec.Broker.ResourceRequirements),
		Ports: []corev1.ContainerPort{
			{Name: "ipc-port", ContainerPort: GetBrokerIpcPort(cr)},
		},
		Env: []corev1.EnvVar{
			{Name: "FE_SVC", Value: GetFeServiceKey(cr.ObjKey()).Name},
			{Name: "FE_QUERY_PORT", Value: strconv.Itoa(int(GetFeQueryPort(cr)))},
			{Name: "ACC_USER", ValueFrom: util.NewEnvVarSecretSource(accountSecretRef.Name, "user")},
			{Name: "ACC_PWD", ValueFrom: util.NewEnvVarSecretSource(accountSecretRef.Name, "password")},
			{Name: "PROBE_TIMEOUT", Value: strconv.Itoa(BrokerProbeTimeoutSec)},
		},
		VolumeMounts: []corev1.VolumeMount{
			{Name: "conf", MountPath: "/opt/apache-doris/broker/conf"},
		},
		ReadinessProbe: &corev1.Probe{
			ProbeHandler:     util.NewTcpSocketProbeHandler(GetBrokerIpcPort(cr)),
			TimeoutSeconds:   1,
			PeriodSeconds:    5,
			SuccessThreshold: 1,
			FailureThreshold: 3,
		},
	}
	// pod template: merge additional pod containers configs defined by user
	mainContainer.Env = append(mainContainer.Env, cr.Spec.Broker.AdditionalEnvs...)
	mainContainer.VolumeMounts = append(mainContainer.VolumeMounts, cr.Spec.Broker.AdditionalVolumeMounts...)
	containers := append([]corev1.Container{mainContainer}, cr.Spec.Broker.AdditionalContainers...)

	// pod template: host alias
	var hostAlias []corev1.HostAlias
	if cr.Spec.HadoopConf != nil {
		hostAlias = mergeHostAlias(cr.Spec.HadoopConf.Hosts, cr.Spec.Broker.HostAliases)
	} else {
		hostAlias = cr.Spec.Broker.HostAliases
	}

	// pod template
	podTemplate := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      brokerLabels,
			Annotations: make(map[string]string),
		},
		Spec: corev1.PodSpec{
			Volumes:            volumes,
			Containers:         containers,
			ImagePullSecrets:   cr.Spec.ImagePullSecrets,
			ServiceAccountName: util.StringFallback(cr.Spec.Broker.ServiceAccount, cr.Spec.ServiceAccount),
			Affinity:           util.PointerFallback(cr.Spec.Broker.Affinity, cr.Spec.Affinity),
			Tolerations:        util.ArrayFallback(cr.Spec.Broker.Tolerations, cr.Spec.Tolerations),
			PriorityClassName:  util.StringFallback(cr.Spec.Broker.PriorityClassName, cr.Spec.PriorityClassName),
			HostAliases:        hostAlias,
		},
	}

	// update strategy
	updateStg := appv1.StatefulSetUpdateStrategy{
		Type: util.PointerFallbackAndDeRefer(
			cr.Spec.Broker.StatefulSetUpdateStrategy,
			cr.Spec.StatefulSetUpdateStrategy,
			appv1.RollingUpdateStatefulSetStrategyType),
	}

	// statefulset
	statefulSet := &appv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      statefulSetRef.Name,
			Namespace: statefulSetRef.Namespace,
			Labels:    brokerLabels,
		},
		Spec: appv1.StatefulSetSpec{
			Replicas:       &cr.Spec.Broker.Replicas,
			ServiceName:    GetBrokerPeerServiceKey(cr.ObjKey()).Name,
			Selector:       &metav1.LabelSelector{MatchLabels: brokerLabels},
			Template:       podTemplate,
			UpdateStrategy: updateStg,
		},
	}

	_ = controllerutil.SetOwnerReference(cr, statefulSet, scheme)
	_ = controllerutil.SetControllerReference(cr, statefulSet, scheme)
	return statefulSet
}
