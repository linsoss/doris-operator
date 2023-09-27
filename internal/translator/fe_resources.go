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
	"fmt"
	dapi "github.com/al-assad/doris-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	FEComponentName = "fe"
)

func GetFEConfigMapName(dorisRef types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: dorisRef.Namespace,
		Name:      fmt.Sprintf("%s-fe-config", dorisRef.Name),
	}
}

func GetFEServiceName(dorisRef types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: dorisRef.Namespace,
		Name:      fmt.Sprintf("%s-fe", dorisRef.Name),
	}
}

func GetFEPeerServiceName(dorisRef types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: dorisRef.Namespace,
		Name:      fmt.Sprintf("%s-fe-peer", dorisRef.Name),
	}
}

func GetFEStatefulSetName(dorisRef types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: dorisRef.Namespace,
		Name:      fmt.Sprintf("%s-fe", dorisRef.Name),
	}
}

func genFELabels(dorisRef types.NamespacedName) map[string]string {
	return MakeResourceLabels(dorisRef.Name, FEComponentName)
}

// MakeFEConfigMap generates a ConfigMap for the FE component
func MakeFEConfigMap(dorisRef types.NamespacedName, cr *dapi.DorisCluster, scheme *runtime.Scheme) *corev1.ConfigMap {
	configMapRef := GetFEConfigMapName(dorisRef)
	configMap := &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name:      configMapRef.Name,
			Namespace: configMapRef.Namespace,
			Labels:    genFELabels(dorisRef),
		},
		Data: map[string]string{
			//"fe.conf": dumpJavaBasedComponentConf(cr.Spec.FE.Configs),
		},
	}
	_ = controllerutil.SetOwnerReference(cr, configMap, scheme)
	return configMap
}
