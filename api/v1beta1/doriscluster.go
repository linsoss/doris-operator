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
	"k8s.io/apimachinery/pkg/types"
)

func (r *DorisCluster) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.Namespace,
		Name:      r.Name,
	}
}

func (r *DorisCluster) GetFeImage() string {
	var version string
	if r.Spec.FE.Version != "" {
		version = r.Spec.FE.Version
	} else {
		version = r.Spec.Version
	}
	return fmt.Sprintf("%s:%s", r.Spec.FE.BaseImage, version)
}

func (r *DorisCluster) GetBeImage() string {
	var version string
	if r.Spec.BE.Version != "" {
		version = r.Spec.BE.Version
	} else {
		version = r.Spec.Version
	}
	return fmt.Sprintf("%s:%s", r.Spec.BE.BaseImage, version)
}

func (r *DorisCluster) GetCnImage() string {
	var version string
	if r.Spec.CN.Version != "" {
		version = r.Spec.CN.Version
	} else {
		version = r.Spec.Version
	}
	return fmt.Sprintf("%s:%s", r.Spec.CN.BaseImage, version)
}

func (r *DorisCluster) GetBrokerImage() string {
	var version string
	if r.Spec.Broker.Version != "" {
		version = r.Spec.Broker.Version
	} else {
		version = r.Spec.Version
	}
	return fmt.Sprintf("%s:%s", r.Spec.Broker.BaseImage, version)
}
