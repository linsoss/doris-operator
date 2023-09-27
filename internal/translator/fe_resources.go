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

package translator

import (
	dapi "github.com/al-assad/doris-operator/api/v1beta1"
	"github.com/al-assad/doris-operator/internal/util"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// MakeOprSqlAccountSecret generates a Secret for the operator SQL account.
func MakeOprSqlAccountSecret(cr *dapi.DorisCluster) *corev1.Secret {
	secretRef := cr.GetOprSqlAccountSecretName()
	secret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      secretRef.Name,
			Namespace: secretRef.Namespace,
			Labels:    MakeResourceLabels(cr.Name, ""),
		},
		Type: corev1.SecretTypeOpaque,
		StringData: map[string]string{
			"user":     "k8sopr",
			"password": GenerateRandomDorisPassword(15),
		},
	}
	return secret
}

func MakeFeConfigMap(cr *dapi.DorisCluster, scheme *runtime.Scheme) *corev1.ConfigMap {
	if cr.Spec.FE == nil {
		return nil
	}
	configMapRef := cr.GetFeConfigMapName()
	data := map[string]string{
		"fe.conf": dumpJavaBasedComponentConf(cr.Spec.FE.Configs),
	}
	// merge hadoop config data
	if cr.Spec.HadoopConf != nil {
		data = util.MergeMaps(cr.Spec.HadoopConf.Config, data)
	}
	configMap := &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name:      configMapRef.Name,
			Namespace: configMapRef.Namespace,
			Labels:    cr.GetFeComponentLabels(),
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
	serviceRef := cr.GetFeServiceName()
	feLabels := cr.GetFeComponentLabels()
	service := &corev1.Service{
		ObjectMeta: v1.ObjectMeta{
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
		Port: cr.GetFeHttpPort(),
	}
	queryPort := corev1.ServicePort{
		Name: "query-port",
		Port: cr.GetFeQueryPort(),
	}

	// When the user specifies a service type
	if cr.Spec.FE.Service != nil {
		if cr.Spec.FE.Service.Type != "" {
			service.Spec.Type = cr.Spec.FE.Service.Type
		}
		if cr.Spec.FE.Service.ExternalTrafficPolicy != nil {
			service.Spec.ExternalTrafficPolicy = *cr.Spec.FE.Service.ExternalTrafficPolicy
		}
		if cr.Spec.FE.Service.QueryPort != 0 {
			httpPort.NodePort = cr.Spec.FE.Service.QueryPort
		}
		if cr.Spec.FE.Service.HttpPort != 0 {
			queryPort.NodePort = cr.Spec.FE.Service.HttpPort
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
	serviceRef := cr.GetFePeerServiceName()
	feLabels := cr.GetFeComponentLabels()
	service := &corev1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:      serviceRef.Name,
			Namespace: serviceRef.Namespace,
			Labels:    feLabels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: "edit-log-port",
					Port: cr.GetFeEditLogPort(),
				}, {
					Name: "rpc-port",
					Port: cr.GetFeRpcPort(),
				},
			},
			Selector:  feLabels,
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: "None",
		},
	}
	_ = controllerutil.SetOwnerReference(cr, service, scheme)
	return service
}
