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
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// Operator Doris SQL account resources

func GetOprSqlAccountSecretKey(dorisClusterKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: dorisClusterKey.Namespace,
		Name:      fmt.Sprintf("%s-opr-account", dorisClusterKey.Name),
	}
}

// MakeOprSqlAccountSecret generates a Secret for the operator SQL account.
func MakeOprSqlAccountSecret(cr *dapi.DorisCluster) *corev1.Secret {
	secretRef := GetOprSqlAccountSecretKey(cr.ObjKey())
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
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

// Doris Monitor RBAC resources

const (
	MonitorClusterRoleName           = "doris-monitor"
	MonitorNamespacedRoleBindingName = "doris-monitor-binding"
	MonitorNamespacedAccountName     = "doris-monitor"
)

func MakeMonitorGlobalRole() *rbacv1.ClusterRole {
	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:   MonitorClusterRoleName,
			Labels: MakeResourceLabels("", "monitor"),
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"nodes", "nodes/proxy", "services", "endpoints", "pods"},
				Verbs:     []string{"get", "list", "watch"},
			}, {
				APIGroups: []string{"extensions"},
				Resources: []string{"ingresses"},
				Verbs:     []string{"get", "list", "watch"},
			}, {
				NonResourceURLs: []string{"/metrics"},
				Verbs:           []string{"get"},
			},
		},
	}
	return clusterRole
}

func MakeMonitorNamespacedAccount(namespace string) *corev1.ServiceAccount {
	account := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      MonitorNamespacedAccountName,
			Namespace: namespace,
			Labels:    MakeResourceLabels("", "monitor"),
		},
	}
	return account
}

func MakeMonitorNamespacedRoleBinding(namespace string) *rbacv1.RoleBinding {
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      MonitorNamespacedRoleBindingName,
			Namespace: namespace,
			Labels:    MakeResourceLabels("", "monitor"),
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			APIGroup: "rbac.authorization.k8s.io",
			Name:     MonitorClusterRoleName,
		},
		Subjects: []rbacv1.Subject{{
			Kind:      "ServiceAccount",
			Name:      MonitorNamespacedAccountName,
			Namespace: namespace,
		}},
	}
	return roleBinding
}
