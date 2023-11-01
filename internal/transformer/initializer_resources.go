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
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strconv"
)

const (
	DefaultMysqlclientImage = "tnir/mysqlclient:1.4.6"
)

var (
	InitializerCheckConnScriptContent = template.ReadOrPanic("initializer/check-conn.sh")
	InitializerStartScriptContent     = template.ReadOrPanic("initializer/start-script.py")
)

func GetInitializerLabels(dorisClusterName string) map[string]string {
	return MakeResourceLabels(dorisClusterName, "initializer")
}

func GetInitializerConfigMapKey(initKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: initKey.Namespace,
		Name:      fmt.Sprintf("%s-initr-conf", initKey.Name),
	}
}

func GetInitializerSecretKey(initKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: initKey.Namespace,
		Name:      fmt.Sprintf("%s-initr-secret", initKey.Name),
	}
}

func GetInitializerJobKey(initKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: initKey.Namespace,
		Name:      fmt.Sprintf("%s-initializer", initKey.Name),
	}
}

func GetInitializerImage(r *dapi.DorisInitializer) string {
	return util.StringFallback(r.Spec.Image, DefaultMysqlclientImage)
}

func MakeInitializerSecret(cr *dapi.DorisInitializer, scheme *runtime.Scheme) *corev1.Secret {
	if cr.Spec.Cluster == "" {
		return nil
	}
	ref := GetInitializerSecretKey(cr.ObjKey())
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ref.Name,
			Namespace: ref.Namespace,
			Labels:    GetInitializerLabels(cr.Spec.Cluster),
		},
		Type: corev1.SecretTypeOpaque,
		StringData: map[string]string{
			"root":  cr.Spec.RootPassword,
			"admin": cr.Spec.AdminPassword,
		},
	}
	_ = controllerutil.SetOwnerReference(cr, secret, scheme)
	return secret
}

func MakeInitializerConfigMap(cr *dapi.DorisInitializer, scheme *runtime.Scheme) *corev1.ConfigMap {
	if cr.Spec.Cluster == "" {
		return nil
	}
	configMapRef := GetInitializerConfigMapKey(cr.ObjKey())
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapRef.Name,
			Namespace: configMapRef.Namespace,
			Labels:    GetInitializerLabels(cr.Spec.Cluster),
		},
		Data: map[string]string{
			"init.sql":             cr.Spec.SqlScript,
			"check-conn-script.sh": InitializerCheckConnScriptContent,
			"start-script.py":      InitializerStartScriptContent,
		},
	}
	_ = controllerutil.SetOwnerReference(cr, configMap, scheme)
	return configMap

}

func MakeInitializerJob(cr *dapi.DorisInitializer, feSvcQueryPort int32, scheme *runtime.Scheme) *batchv1.Job {
	if cr.Spec.Cluster == "" {
		return nil
	}
	clusterRef := types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Spec.Cluster,
	}
	jobRef := GetInitializerJobKey(cr.ObjKey())
	secretRef := GetInitializerSecretKey(cr.ObjKey())
	configMapRef := GetInitializerConfigMapKey(cr.ObjKey())
	feSvcRef := GetFeServiceKey(clusterRef)
	accountSecretRef := GetOprSqlAccountSecretKey(clusterRef)

	initLabels := GetInitializerLabels(cr.Spec.Cluster)
	image := GetInitializerImage(cr)

	// pod template: volumes
	volumes := []corev1.Volume{
		{
			Name: "password",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{SecretName: secretRef.Name}},
		}, {
			Name: "init-sql",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{Name: configMapRef.Name},
					Items:                []corev1.KeyToPath{{Key: "init.sql", Path: "init.sql"}},
				}},
		}, {
			Name: "check-conn-script",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{Name: configMapRef.Name},
					Items:                []corev1.KeyToPath{{Key: "check-conn-script.sh", Path: "check_conn.sh"}},
				}},
		}, {
			Name: "start-script",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{Name: configMapRef.Name},
					Items:                []corev1.KeyToPath{{Key: "start-script.py", Path: "start.py"}},
				}},
		},
	}
	// pod template: init container
	initContainer := corev1.Container{
		Name:            "wait",
		Image:           image,
		ImagePullPolicy: cr.Spec.ImagePullPolicy,
		Command:         []string{"sh", "/usr/local/bin/check_conn.sh"},
		VolumeMounts: []corev1.VolumeMount{{
			Name:      "check-conn-script",
			MountPath: "/usr/local/bin/check_conn.sh",
			SubPath:   "check_conn.sh",
			ReadOnly:  true,
		}},
		Env: []corev1.EnvVar{
			{
				Name:  "FE_SVC",
				Value: feSvcRef.Name,
			}, {
				Name:  "FE_QUERY_PORT",
				Value: strconv.Itoa(int(feSvcQueryPort)),
			},
		},
	}
	// pod template: main container
	mainContainer := corev1.Container{
		Name:            "mysql-client",
		Image:           image,
		ImagePullPolicy: cr.Spec.ImagePullPolicy,
		Command:         []string{"python", "/usr/local/bin/start.py"},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "start-script",
				MountPath: "/usr/local/bin/start.py",
				SubPath:   "start.py",
				ReadOnly:  true,
			}, {
				Name:      "password",
				MountPath: "/etc/doris/password",
				ReadOnly:  true,
			}, {
				Name:      "init-sql",
				MountPath: "/etc/doris/init.sql",
				SubPath:   "init.sql",
				ReadOnly:  true,
			},
		},
		Env: []corev1.EnvVar{
			{
				Name:  "FE_SVC",
				Value: feSvcRef.Name,
			}, {
				Name:  "FE_QUERY_PORT",
				Value: strconv.Itoa(int(feSvcQueryPort)),
			}, {
				Name:      "ACC_USER",
				ValueFrom: util.NewEnvVarSecretSource(accountSecretRef.Name, "user"),
			}, {
				Name:      "ACC_PWD",
				ValueFrom: util.NewEnvVarSecretSource(accountSecretRef.Name, "password"),
			},
		},
	}

	// pod template
	podTemplate := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      initLabels,
			Annotations: make(map[string]string),
		},
		Spec: corev1.PodSpec{
			RestartPolicy:    corev1.RestartPolicyNever,
			Volumes:          volumes,
			InitContainers:   []corev1.Container{initContainer},
			Containers:       []corev1.Container{mainContainer},
			ImagePullSecrets: cr.Spec.ImagePullSecrets,
			Tolerations:      cr.Spec.Tolerations,
			NodeSelector:     cr.Spec.NodeSelector,
		},
	}

	// job
	parallelism := int32(1)
	backoffLimit := cr.Spec.MaxRetry
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobRef.Name,
			Namespace: jobRef.Namespace,
			Labels:    initLabels,
		},
		Spec: batchv1.JobSpec{
			Parallelism:  &parallelism,
			BackoffLimit: backoffLimit,
			Template:     podTemplate,
		},
	}
	_ = controllerutil.SetOwnerReference(cr, job, scheme)
	_ = controllerutil.SetControllerReference(cr, job, scheme)
	return job
}
