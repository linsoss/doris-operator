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
	"k8s.io/apimachinery/pkg/types"
)

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

//func MakeGrafanaConfigMap(cr *dapi.DorisMonitor, scheme *runtime.Scheme) (*corev1.ConfigMap, error) {
//	if cr.Spec.Cluster == "" {
//		return nil, nil
//	}
//	clusterRef := types.NamespacedName{
//		Namespace: cr.Namespace,
//		Name:      cr.Spec.Cluster,
//	}
//	configMapRef := GetPrometheusConfigMapKey(cr.ObjKey())
//	labels := GetMonitorPrometheusLabels(clusterRef)
//	promConfContent, err := template.ExecTemplate(PrometheusConfTmpl, clusterRef)
//	if err != nil {
//		return nil, util.MergeErrors(fmt.Errorf("fail to parse prometheus.conf template"), err)
//	}
//	// merge hadoop config data
//	configMap := &corev1.ConfigMap{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      configMapRef.Name,
//			Namespace: configMapRef.Namespace,
//			Labels:    labels,
//		},
//		Data: map[string]string{
//			"prometheus.yml": promConfContent,
//		},
//	}
//	_ = controllerutil.SetOwnerReference(cr, configMap, scheme)
//	return configMap, nil
//}
