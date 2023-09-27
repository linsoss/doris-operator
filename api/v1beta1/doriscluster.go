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

package v1beta1

import (
	"fmt"
	"github.com/al-assad/doris-operator/api/v1beta1/defaulting"
	"github.com/al-assad/doris-operator/internal/translator"
	"k8s.io/apimachinery/pkg/types"
)

func (r *DorisCluster) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.Namespace,
		Name:      r.Name,
	}
}

func (r *DorisCluster) GetOprSqlAccountSecretName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.Namespace,
		Name:      fmt.Sprintf("%s-opr-account", r.Name),
	}
}

// ------------ FE resources ------------

func (r *DorisCluster) GetFeComponentLabels() map[string]string {
	return translator.MakeResourceLabels(r.Name, "fe")
}

func (r *DorisCluster) GetFeConfigMapName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.Namespace,
		Name:      fmt.Sprintf("%s-fe-config", r.Name),
	}
}

func (r *DorisCluster) GetFeServiceName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.Namespace,
		Name:      fmt.Sprintf("%s-fe", r.Name),
	}
}

func (r *DorisCluster) GetFePeerServiceName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.Namespace,
		Name:      fmt.Sprintf("%s-fe-peer", r.Name),
	}
}

func (r *DorisCluster) GetFeStatefulSetName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.Namespace,
		Name:      fmt.Sprintf("%s-fe", r.Name),
	}
}

func (r *DorisCluster) GetFeHttpPort() int32 {
	if r.Spec.FE == nil {
		return defaulting.DefaultFeHttpPort
	}
	return translator.GetPortValueFromRawConf(r.Spec.FE.Configs, "http_port", defaulting.DefaultFeHttpPort)
}

func (r *DorisCluster) GetFeQueryPort() int32 {
	if r.Spec.FE == nil {
		return defaulting.DefaultFeQueryPort
	}
	return translator.GetPortValueFromRawConf(r.Spec.FE.Configs, "query_port", defaulting.DefaultFeQueryPort)
}

func (r *DorisCluster) GetFeRpcPort() int32 {
	if r.Spec.FE == nil {
		return defaulting.DefaultFeRpcPort
	}
	return translator.GetPortValueFromRawConf(r.Spec.FE.Configs, "query_port", defaulting.DefaultFeRpcPort)
}

func (r *DorisCluster) GetFeEditLogPort() int32 {
	if r.Spec.FE == nil {
		return defaulting.DefaultFeEditLogPort
	}
	return translator.GetPortValueFromRawConf(r.Spec.FE.Configs, "edit_log_port", defaulting.DefaultFeEditLogPort)
}

// ------------ BE resources ------------
