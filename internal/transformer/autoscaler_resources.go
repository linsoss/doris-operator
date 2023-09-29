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
	acv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	DefaultHpaPeriodSeconds int32 = 60
)

func GetCnAutoscalerLabels(dorisAutoScalerKey types.NamespacedName) map[string]string {
	return MakeResourceLabels(dorisAutoScalerKey.Name, "cn-autoscaler")
}

func GetCnScaleUpHpaKey(dorisAutoScalerKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: dorisAutoScalerKey.Namespace,
		Name:      fmt.Sprintf("%s-cn-scaleup", dorisAutoScalerKey.Name),
	}
}

func GetCnScaleDownHpaKey(dorisAutoScalerKey types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Namespace: dorisAutoScalerKey.Namespace,
		Name:      fmt.Sprintf("%s-cn-scaledown", dorisAutoScalerKey.Name),
	}
}

func MakeCnScaleUpHpa(cr *dapi.DorisAutoscaler, scheme *runtime.Scheme) *acv2.HorizontalPodAutoscaler {
	if cr.Spec.CN == nil {
		return nil
	}
	cpuRuleExist := cr.Spec.CN.Rules.Cpu != nil && cr.Spec.CN.Rules.Cpu.Max != nil
	memRuleExist := cr.Spec.CN.Rules.Memory != nil && cr.Spec.CN.Rules.Memory.Max != nil
	// empty rule means no autoscaling
	if !cpuRuleExist && !memRuleExist {
		return nil
	}

	hpaRef := GetCnScaleUpHpaKey(cr.ObjKey())
	clusterRef := types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Spec.Cluster,
	}
	// hpa metrics
	metricsList := make([]acv2.MetricSpec, 0, 2)
	if cpuRuleExist {
		metricsList = append(metricsList,
			util.NewResourceAvgUtilizationMetricSpec(corev1.ResourceCPU, cr.Spec.CN.Rules.Cpu.Max))
	}
	if memRuleExist {
		metricsList = append(metricsList,
			util.NewResourceAvgUtilizationMetricSpec(corev1.ResourceMemory, cr.Spec.CN.Rules.Memory.Max))
	}
	// hpa behavior
	selectPolicy := acv2.MaxChangePolicySelect
	periodSec := DefaultHpaPeriodSeconds
	if cr.Spec.CN.ScalePeriodSeconds != nil && cr.Spec.CN.ScalePeriodSeconds.ScaleUp != nil {
		periodSec = *cr.Spec.CN.ScalePeriodSeconds.ScaleUp
	}
	behavior := &acv2.HorizontalPodAutoscalerBehavior{
		ScaleUp: &acv2.HPAScalingRules{
			SelectPolicy: &selectPolicy,
			Policies: []acv2.HPAScalingPolicy{{
				Type:          acv2.PodsScalingPolicy,
				Value:         1,
				PeriodSeconds: periodSec,
			}},
		},
	}
	// hpa resource
	hpa := &acv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: hpaRef.Namespace,
			Name:      hpaRef.Name,
			Labels:    GetCnAutoscalerLabels(cr.ObjKey()),
		},
		Spec: acv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: acv2.CrossVersionObjectReference{
				Kind:       "StatefulSet",
				APIVersion: "apps/v1",
				Name:       GetCnStatefulSetKey(clusterRef).Name,
			},
			MaxReplicas: cr.Spec.CN.Replicas.Max,
			MinReplicas: cr.Spec.CN.Replicas.Min,
			Metrics:     metricsList,
			Behavior:    behavior,
		},
	}
	_ = controllerutil.SetOwnerReference(cr, hpa, scheme)
	return hpa
}

func MakeCnScaleDownHpa(cr *dapi.DorisAutoscaler, scheme *runtime.Scheme) *acv2.HorizontalPodAutoscaler {
	if cr.Spec.CN == nil || cr.Spec.CN.DisableScaleDown {
		return nil
	}
	cpuRuleExist := cr.Spec.CN.Rules.Cpu != nil && cr.Spec.CN.Rules.Cpu.Min != nil
	memRuleExist := cr.Spec.CN.Rules.Memory != nil && cr.Spec.CN.Rules.Memory.Min != nil
	// empty rule means no autoscaling
	if !cpuRuleExist && !memRuleExist {
		return nil
	}

	hpaRef := GetCnScaleDownHpaKey(cr.ObjKey())
	clusterRef := types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Spec.Cluster,
	}
	// hpa metrics
	metricsList := make([]acv2.MetricSpec, 0, 2)
	if cpuRuleExist {
		metricsList = append(metricsList,
			util.NewResourceAvgUtilizationMetricSpec(corev1.ResourceCPU, cr.Spec.CN.Rules.Cpu.Min))
	}
	if memRuleExist {
		metricsList = append(metricsList,
			util.NewResourceAvgUtilizationMetricSpec(corev1.ResourceMemory, cr.Spec.CN.Rules.Memory.Min))
	}
	// hpa behavior
	selectPolicy := acv2.MinChangePolicySelect
	periodSec := DefaultHpaPeriodSeconds
	if cr.Spec.CN.ScalePeriodSeconds != nil && cr.Spec.CN.ScalePeriodSeconds.ScaleDown != nil {
		periodSec = *cr.Spec.CN.ScalePeriodSeconds.ScaleDown
	}
	behavior := &acv2.HorizontalPodAutoscalerBehavior{
		ScaleDown: &acv2.HPAScalingRules{
			SelectPolicy: &selectPolicy,
			Policies: []acv2.HPAScalingPolicy{{
				Type:          acv2.PodsScalingPolicy,
				Value:         1,
				PeriodSeconds: periodSec,
			}},
		},
	}
	// hpa resource
	hpa := &acv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: hpaRef.Namespace,
			Name:      hpaRef.Name,
			Labels:    GetCnAutoscalerLabels(cr.ObjKey()),
		},
		Spec: acv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: acv2.CrossVersionObjectReference{
				Kind:       "StatefulSet",
				APIVersion: "apps/v1",
				Name:       GetCnStatefulSetKey(clusterRef).Name,
			},
			MaxReplicas: cr.Spec.CN.Replicas.Max,
			MinReplicas: cr.Spec.CN.Replicas.Min,
			Metrics:     metricsList,
			Behavior:    behavior,
		},
	}
	_ = controllerutil.SetOwnerReference(cr, hpa, scheme)
	return hpa
}
