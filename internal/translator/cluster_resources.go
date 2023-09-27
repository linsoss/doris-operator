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
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func GetOprSqlAccountSecretName(dorisRef types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: dorisRef.Namespace,
		Name:      fmt.Sprintf("%s-opr-account", dorisRef.Name),
	}
}

func MakeOprSqlAccountSecret(dorisRef types.NamespacedName) *corev1.Secret {
	secretRef := GetOprSqlAccountSecretName(dorisRef)
	secret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      secretRef.Name,
			Namespace: secretRef.Namespace,
			Labels:    MakeResourceLabels(dorisRef.Name, ""),
		},
		Type: corev1.SecretTypeOpaque,
		StringData: map[string]string{
			"user":     "k8sopr",
			"password": GenerateRandomDorisPassword(15),
		},
	}
	return secret
}
