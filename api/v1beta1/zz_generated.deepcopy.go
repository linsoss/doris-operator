//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/autoscaling/v2"
	"k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AutoScalerRef) DeepCopyInto(out *AutoScalerRef) {
	*out = *in
	out.NamespacedName = in.NamespacedName
	out.TypeMeta = in.TypeMeta
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutoScalerRef.
func (in *AutoScalerRef) DeepCopy() *AutoScalerRef {
	if in == nil {
		return nil
	}
	out := new(AutoScalerRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AutoscalerRecStatus) DeepCopyInto(out *AutoscalerRecStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutoscalerRecStatus.
func (in *AutoscalerRecStatus) DeepCopy() *AutoscalerRecStatus {
	if in == nil {
		return nil
	}
	out := new(AutoscalerRecStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BESpec) DeepCopyInto(out *BESpec) {
	*out = *in
	in.DorisComponentSpec.DeepCopyInto(&out.DorisComponentSpec)
	if in.StorageClassName != nil {
		in, out := &in.StorageClassName, &out.StorageClassName
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BESpec.
func (in *BESpec) DeepCopy() *BESpec {
	if in == nil {
		return nil
	}
	out := new(BESpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BEStatus) DeepCopyInto(out *BEStatus) {
	*out = *in
	in.DorisComponentStatus.DeepCopyInto(&out.DorisComponentStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BEStatus.
func (in *BEStatus) DeepCopy() *BEStatus {
	if in == nil {
		return nil
	}
	out := new(BEStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BrokerSpec) DeepCopyInto(out *BrokerSpec) {
	*out = *in
	in.DorisComponentSpec.DeepCopyInto(&out.DorisComponentSpec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BrokerSpec.
func (in *BrokerSpec) DeepCopy() *BrokerSpec {
	if in == nil {
		return nil
	}
	out := new(BrokerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BrokerStatus) DeepCopyInto(out *BrokerStatus) {
	*out = *in
	in.DorisComponentStatus.DeepCopyInto(&out.DorisComponentStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BrokerStatus.
func (in *BrokerStatus) DeepCopy() *BrokerStatus {
	if in == nil {
		return nil
	}
	out := new(BrokerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CNAutoscalerRules) DeepCopyInto(out *CNAutoscalerRules) {
	*out = *in
	if in.Cpu != nil {
		in, out := &in.Cpu, &out.Cpu
		*out = new(UtilizationThresholdRange)
		(*in).DeepCopyInto(*out)
	}
	if in.Memory != nil {
		in, out := &in.Memory, &out.Memory
		*out = new(UtilizationThresholdRange)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CNAutoscalerRules.
func (in *CNAutoscalerRules) DeepCopy() *CNAutoscalerRules {
	if in == nil {
		return nil
	}
	out := new(CNAutoscalerRules)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CNAutoscalerSpec) DeepCopyInto(out *CNAutoscalerSpec) {
	*out = *in
	in.Replicas.DeepCopyInto(&out.Replicas)
	in.Rules.DeepCopyInto(&out.Rules)
	if in.ScalePeriodSeconds != nil {
		in, out := &in.ScalePeriodSeconds, &out.ScalePeriodSeconds
		*out = new(ScalePeriodSeconds)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CNAutoscalerSpec.
func (in *CNAutoscalerSpec) DeepCopy() *CNAutoscalerSpec {
	if in == nil {
		return nil
	}
	out := new(CNAutoscalerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CNAutoscalerStatus) DeepCopyInto(out *CNAutoscalerStatus) {
	*out = *in
	out.AutoscalerRecStatus = in.AutoscalerRecStatus
	in.CNAutoscalerSyncStatus.DeepCopyInto(&out.CNAutoscalerSyncStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CNAutoscalerStatus.
func (in *CNAutoscalerStatus) DeepCopy() *CNAutoscalerStatus {
	if in == nil {
		return nil
	}
	out := new(CNAutoscalerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CNAutoscalerSyncStatus) DeepCopyInto(out *CNAutoscalerSyncStatus) {
	*out = *in
	if in.ScaleUpHpaRef != nil {
		in, out := &in.ScaleUpHpaRef, &out.ScaleUpHpaRef
		*out = new(AutoScalerRef)
		**out = **in
	}
	if in.ScaleUpStatus != nil {
		in, out := &in.ScaleUpStatus, &out.ScaleUpStatus
		*out = new(v2.HorizontalPodAutoscalerStatus)
		(*in).DeepCopyInto(*out)
	}
	if in.ScaleDownHpaRef != nil {
		in, out := &in.ScaleDownHpaRef, &out.ScaleDownHpaRef
		*out = new(AutoScalerRef)
		**out = **in
	}
	if in.ScaleDownStatus != nil {
		in, out := &in.ScaleDownStatus, &out.ScaleDownStatus
		*out = new(v2.HorizontalPodAutoscalerStatus)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CNAutoscalerSyncStatus.
func (in *CNAutoscalerSyncStatus) DeepCopy() *CNAutoscalerSyncStatus {
	if in == nil {
		return nil
	}
	out := new(CNAutoscalerSyncStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CNSpec) DeepCopyInto(out *CNSpec) {
	*out = *in
	in.DorisComponentSpec.DeepCopyInto(&out.DorisComponentSpec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CNSpec.
func (in *CNSpec) DeepCopy() *CNSpec {
	if in == nil {
		return nil
	}
	out := new(CNSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CNStatus) DeepCopyInto(out *CNStatus) {
	*out = *in
	in.DorisComponentStatus.DeepCopyInto(&out.DorisComponentStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CNStatus.
func (in *CNStatus) DeepCopy() *CNStatus {
	if in == nil {
		return nil
	}
	out := new(CNStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DorisAutoscaler) DeepCopyInto(out *DorisAutoscaler) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	if in.objKey != nil {
		in, out := &in.objKey, &out.objKey
		*out = new(types.NamespacedName)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DorisAutoscaler.
func (in *DorisAutoscaler) DeepCopy() *DorisAutoscaler {
	if in == nil {
		return nil
	}
	out := new(DorisAutoscaler)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DorisAutoscaler) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DorisAutoscalerList) DeepCopyInto(out *DorisAutoscalerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DorisAutoscaler, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DorisAutoscalerList.
func (in *DorisAutoscalerList) DeepCopy() *DorisAutoscalerList {
	if in == nil {
		return nil
	}
	out := new(DorisAutoscalerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DorisAutoscalerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DorisAutoscalerSpec) DeepCopyInto(out *DorisAutoscalerSpec) {
	*out = *in
	if in.CN != nil {
		in, out := &in.CN, &out.CN
		*out = new(CNAutoscalerSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DorisAutoscalerSpec.
func (in *DorisAutoscalerSpec) DeepCopy() *DorisAutoscalerSpec {
	if in == nil {
		return nil
	}
	out := new(DorisAutoscalerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DorisAutoscalerStatus) DeepCopyInto(out *DorisAutoscalerStatus) {
	*out = *in
	if in.LastApplySpecHash != nil {
		in, out := &in.LastApplySpecHash, &out.LastApplySpecHash
		*out = new(string)
		**out = **in
	}
	out.ClusterRef = in.ClusterRef
	in.CN.DeepCopyInto(&out.CN)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DorisAutoscalerStatus.
func (in *DorisAutoscalerStatus) DeepCopy() *DorisAutoscalerStatus {
	if in == nil {
		return nil
	}
	out := new(DorisAutoscalerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DorisCluster) DeepCopyInto(out *DorisCluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	if in.objKey != nil {
		in, out := &in.objKey, &out.objKey
		*out = new(types.NamespacedName)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DorisCluster.
func (in *DorisCluster) DeepCopy() *DorisCluster {
	if in == nil {
		return nil
	}
	out := new(DorisCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DorisCluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DorisClusterList) DeepCopyInto(out *DorisClusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DorisCluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DorisClusterList.
func (in *DorisClusterList) DeepCopy() *DorisClusterList {
	if in == nil {
		return nil
	}
	out := new(DorisClusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DorisClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DorisClusterRecStatus) DeepCopyInto(out *DorisClusterRecStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DorisClusterRecStatus.
func (in *DorisClusterRecStatus) DeepCopy() *DorisClusterRecStatus {
	if in == nil {
		return nil
	}
	out := new(DorisClusterRecStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DorisClusterSpec) DeepCopyInto(out *DorisClusterSpec) {
	*out = *in
	if in.FE != nil {
		in, out := &in.FE, &out.FE
		*out = new(FESpec)
		(*in).DeepCopyInto(*out)
	}
	if in.BE != nil {
		in, out := &in.BE, &out.BE
		*out = new(BESpec)
		(*in).DeepCopyInto(*out)
	}
	if in.CN != nil {
		in, out := &in.CN, &out.CN
		*out = new(CNSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Broker != nil {
		in, out := &in.Broker, &out.Broker
		*out = new(BrokerSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.HadoopConf != nil {
		in, out := &in.HadoopConf, &out.HadoopConf
		*out = new(HadoopConfSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.ImagePullSecrets != nil {
		in, out := &in.ImagePullSecrets, &out.ImagePullSecrets
		*out = make([]v1.LocalObjectReference, len(*in))
		copy(*out, *in)
	}
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		*out = new(v1.Affinity)
		(*in).DeepCopyInto(*out)
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.StatefulSetUpdateStrategy != nil {
		in, out := &in.StatefulSetUpdateStrategy, &out.StatefulSetUpdateStrategy
		*out = new(appsv1.StatefulSetUpdateStrategyType)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DorisClusterSpec.
func (in *DorisClusterSpec) DeepCopy() *DorisClusterSpec {
	if in == nil {
		return nil
	}
	out := new(DorisClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DorisClusterStatus) DeepCopyInto(out *DorisClusterStatus) {
	*out = *in
	if in.LastApplySpecHash != nil {
		in, out := &in.LastApplySpecHash, &out.LastApplySpecHash
		*out = new(string)
		**out = **in
	}
	out.DorisClusterRecStatus = in.DorisClusterRecStatus
	in.DorisClusterSyncStatus.DeepCopyInto(&out.DorisClusterSyncStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DorisClusterStatus.
func (in *DorisClusterStatus) DeepCopy() *DorisClusterStatus {
	if in == nil {
		return nil
	}
	out := new(DorisClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DorisClusterSyncStatus) DeepCopyInto(out *DorisClusterSyncStatus) {
	*out = *in
	in.FE.DeepCopyInto(&out.FE)
	in.BE.DeepCopyInto(&out.BE)
	in.CN.DeepCopyInto(&out.CN)
	in.Broker.DeepCopyInto(&out.Broker)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DorisClusterSyncStatus.
func (in *DorisClusterSyncStatus) DeepCopy() *DorisClusterSyncStatus {
	if in == nil {
		return nil
	}
	out := new(DorisClusterSyncStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DorisComponentSpec) DeepCopyInto(out *DorisComponentSpec) {
	*out = *in
	in.ResourceRequirements.DeepCopyInto(&out.ResourceRequirements)
	if in.Configs != nil {
		in, out := &in.Configs, &out.Configs
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.HostAliases != nil {
		in, out := &in.HostAliases, &out.HostAliases
		*out = make([]v1.HostAlias, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		*out = new(v1.Affinity)
		(*in).DeepCopyInto(*out)
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.StatefulSetUpdateStrategy != nil {
		in, out := &in.StatefulSetUpdateStrategy, &out.StatefulSetUpdateStrategy
		*out = new(appsv1.StatefulSetUpdateStrategyType)
		**out = **in
	}
	if in.AdditionalEnvs != nil {
		in, out := &in.AdditionalEnvs, &out.AdditionalEnvs
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.AdditionalContainers != nil {
		in, out := &in.AdditionalContainers, &out.AdditionalContainers
		*out = make([]v1.Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.AdditionalVolumes != nil {
		in, out := &in.AdditionalVolumes, &out.AdditionalVolumes
		*out = make([]v1.Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.AdditionalVolumeMounts != nil {
		in, out := &in.AdditionalVolumeMounts, &out.AdditionalVolumeMounts
		*out = make([]v1.VolumeMount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DorisComponentSpec.
func (in *DorisComponentSpec) DeepCopy() *DorisComponentSpec {
	if in == nil {
		return nil
	}
	out := new(DorisComponentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DorisComponentStatus) DeepCopyInto(out *DorisComponentStatus) {
	*out = *in
	out.StatefulSetRef = in.StatefulSetRef
	if in.Members != nil {
		in, out := &in.Members, &out.Members
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ReadyMembers != nil {
		in, out := &in.ReadyMembers, &out.ReadyMembers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]appsv1.StatefulSetCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DorisComponentStatus.
func (in *DorisComponentStatus) DeepCopy() *DorisComponentStatus {
	if in == nil {
		return nil
	}
	out := new(DorisComponentStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DorisInitializer) DeepCopyInto(out *DorisInitializer) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	if in.objKey != nil {
		in, out := &in.objKey, &out.objKey
		*out = new(types.NamespacedName)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DorisInitializer.
func (in *DorisInitializer) DeepCopy() *DorisInitializer {
	if in == nil {
		return nil
	}
	out := new(DorisInitializer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DorisInitializer) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DorisInitializerList) DeepCopyInto(out *DorisInitializerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DorisInitializer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DorisInitializerList.
func (in *DorisInitializerList) DeepCopy() *DorisInitializerList {
	if in == nil {
		return nil
	}
	out := new(DorisInitializerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DorisInitializerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DorisInitializerRecStatus) DeepCopyInto(out *DorisInitializerRecStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DorisInitializerRecStatus.
func (in *DorisInitializerRecStatus) DeepCopy() *DorisInitializerRecStatus {
	if in == nil {
		return nil
	}
	out := new(DorisInitializerRecStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DorisInitializerSpec) DeepCopyInto(out *DorisInitializerSpec) {
	*out = *in
	if in.ImagePullSecrets != nil {
		in, out := &in.ImagePullSecrets, &out.ImagePullSecrets
		*out = make([]v1.LocalObjectReference, len(*in))
		copy(*out, *in)
	}
	if in.MaxRetry != nil {
		in, out := &in.MaxRetry, &out.MaxRetry
		*out = new(int32)
		**out = **in
	}
	in.ResourceRequirements.DeepCopyInto(&out.ResourceRequirements)
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DorisInitializerSpec.
func (in *DorisInitializerSpec) DeepCopy() *DorisInitializerSpec {
	if in == nil {
		return nil
	}
	out := new(DorisInitializerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DorisInitializerStatus) DeepCopyInto(out *DorisInitializerStatus) {
	*out = *in
	if in.LastApplySpecHash != nil {
		in, out := &in.LastApplySpecHash, &out.LastApplySpecHash
		*out = new(string)
		**out = **in
	}
	out.DorisInitializerRecStatus = in.DorisInitializerRecStatus
	in.DorisInitializerSyncStatus.DeepCopyInto(&out.DorisInitializerSyncStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DorisInitializerStatus.
func (in *DorisInitializerStatus) DeepCopy() *DorisInitializerStatus {
	if in == nil {
		return nil
	}
	out := new(DorisInitializerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DorisInitializerSyncStatus) DeepCopyInto(out *DorisInitializerSyncStatus) {
	*out = *in
	out.JobRef = in.JobRef
	in.JobStatus.DeepCopyInto(&out.JobStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DorisInitializerSyncStatus.
func (in *DorisInitializerSyncStatus) DeepCopy() *DorisInitializerSyncStatus {
	if in == nil {
		return nil
	}
	out := new(DorisInitializerSyncStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DorisMonitor) DeepCopyInto(out *DorisMonitor) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DorisMonitor.
func (in *DorisMonitor) DeepCopy() *DorisMonitor {
	if in == nil {
		return nil
	}
	out := new(DorisMonitor)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DorisMonitor) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DorisMonitorComponentStatus) DeepCopyInto(out *DorisMonitorComponentStatus) {
	*out = *in
	out.DeploymentRef = in.DeploymentRef
	out.PVCRef = in.PVCRef
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]appsv1.DeploymentCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DorisMonitorComponentStatus.
func (in *DorisMonitorComponentStatus) DeepCopy() *DorisMonitorComponentStatus {
	if in == nil {
		return nil
	}
	out := new(DorisMonitorComponentStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DorisMonitorList) DeepCopyInto(out *DorisMonitorList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DorisMonitor, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DorisMonitorList.
func (in *DorisMonitorList) DeepCopy() *DorisMonitorList {
	if in == nil {
		return nil
	}
	out := new(DorisMonitorList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DorisMonitorList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DorisMonitorSpec) DeepCopyInto(out *DorisMonitorSpec) {
	*out = *in
	if in.Prometheus != nil {
		in, out := &in.Prometheus, &out.Prometheus
		*out = new(PrometheusSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Grafana != nil {
		in, out := &in.Grafana, &out.Grafana
		*out = new(GrafanaSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.ImagePullSecrets != nil {
		in, out := &in.ImagePullSecrets, &out.ImagePullSecrets
		*out = make([]v1.LocalObjectReference, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DorisMonitorSpec.
func (in *DorisMonitorSpec) DeepCopy() *DorisMonitorSpec {
	if in == nil {
		return nil
	}
	out := new(DorisMonitorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DorisMonitorStatus) DeepCopyInto(out *DorisMonitorStatus) {
	*out = *in
	out.ClusterRef = in.ClusterRef
	if in.PrevSpec != nil {
		in, out := &in.PrevSpec, &out.PrevSpec
		*out = new(DorisMonitorSpec)
		(*in).DeepCopyInto(*out)
	}
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
	in.Prometheus.DeepCopyInto(&out.Prometheus)
	in.Grafana.DeepCopyInto(&out.Grafana)
	in.Loki.DeepCopyInto(&out.Loki)
	in.Promtail.DeepCopyInto(&out.Promtail)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DorisMonitorStatus.
func (in *DorisMonitorStatus) DeepCopy() *DorisMonitorStatus {
	if in == nil {
		return nil
	}
	out := new(DorisMonitorStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FESpec) DeepCopyInto(out *FESpec) {
	*out = *in
	in.DorisComponentSpec.DeepCopyInto(&out.DorisComponentSpec)
	if in.StorageClassName != nil {
		in, out := &in.StorageClassName, &out.StorageClassName
		*out = new(string)
		**out = **in
	}
	if in.Service != nil {
		in, out := &in.Service, &out.Service
		*out = new(FeServiceSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FESpec.
func (in *FESpec) DeepCopy() *FESpec {
	if in == nil {
		return nil
	}
	out := new(FESpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FEStatus) DeepCopyInto(out *FEStatus) {
	*out = *in
	out.ServiceRef = in.ServiceRef
	in.DorisComponentStatus.DeepCopyInto(&out.DorisComponentStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FEStatus.
func (in *FEStatus) DeepCopy() *FEStatus {
	if in == nil {
		return nil
	}
	out := new(FEStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FeServiceSpec) DeepCopyInto(out *FeServiceSpec) {
	*out = *in
	if in.QueryPort != nil {
		in, out := &in.QueryPort, &out.QueryPort
		*out = new(int32)
		**out = **in
	}
	if in.HttpPort != nil {
		in, out := &in.HttpPort, &out.HttpPort
		*out = new(int32)
		**out = **in
	}
	if in.ExternalTrafficPolicy != nil {
		in, out := &in.ExternalTrafficPolicy, &out.ExternalTrafficPolicy
		*out = new(v1.ServiceExternalTrafficPolicy)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FeServiceSpec.
func (in *FeServiceSpec) DeepCopy() *FeServiceSpec {
	if in == nil {
		return nil
	}
	out := new(FeServiceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GrafanaSpec) DeepCopyInto(out *GrafanaSpec) {
	*out = *in
	if in.Service != nil {
		in, out := &in.Service, &out.Service
		*out = new(MonitorServiceSpec)
		(*in).DeepCopyInto(*out)
	}
	in.ResourceRequirements.DeepCopyInto(&out.ResourceRequirements)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GrafanaSpec.
func (in *GrafanaSpec) DeepCopy() *GrafanaSpec {
	if in == nil {
		return nil
	}
	out := new(GrafanaSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GrafanaStatus) DeepCopyInto(out *GrafanaStatus) {
	*out = *in
	out.ServiceRef = in.ServiceRef
	in.DorisMonitorComponentStatus.DeepCopyInto(&out.DorisMonitorComponentStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GrafanaStatus.
func (in *GrafanaStatus) DeepCopy() *GrafanaStatus {
	if in == nil {
		return nil
	}
	out := new(GrafanaStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HadoopConfSpec) DeepCopyInto(out *HadoopConfSpec) {
	*out = *in
	if in.Hosts != nil {
		in, out := &in.Hosts, &out.Hosts
		*out = make([]HostnameIpItem, len(*in))
		copy(*out, *in)
	}
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HadoopConfSpec.
func (in *HadoopConfSpec) DeepCopy() *HadoopConfSpec {
	if in == nil {
		return nil
	}
	out := new(HadoopConfSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostnameIpItem) DeepCopyInto(out *HostnameIpItem) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostnameIpItem.
func (in *HostnameIpItem) DeepCopy() *HostnameIpItem {
	if in == nil {
		return nil
	}
	out := new(HostnameIpItem)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LokiSpec) DeepCopyInto(out *LokiSpec) {
	*out = *in
	in.ResourceRequirements.DeepCopyInto(&out.ResourceRequirements)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LokiSpec.
func (in *LokiSpec) DeepCopy() *LokiSpec {
	if in == nil {
		return nil
	}
	out := new(LokiSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LokiStatus) DeepCopyInto(out *LokiStatus) {
	*out = *in
	in.DorisMonitorComponentStatus.DeepCopyInto(&out.DorisMonitorComponentStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LokiStatus.
func (in *LokiStatus) DeepCopy() *LokiStatus {
	if in == nil {
		return nil
	}
	out := new(LokiStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MonitorServiceSpec) DeepCopyInto(out *MonitorServiceSpec) {
	*out = *in
	if in.ExternalTrafficPolicy != nil {
		in, out := &in.ExternalTrafficPolicy, &out.ExternalTrafficPolicy
		*out = new(v1.ServiceExternalTrafficPolicy)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MonitorServiceSpec.
func (in *MonitorServiceSpec) DeepCopy() *MonitorServiceSpec {
	if in == nil {
		return nil
	}
	out := new(MonitorServiceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespacedName) DeepCopyInto(out *NamespacedName) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespacedName.
func (in *NamespacedName) DeepCopy() *NamespacedName {
	if in == nil {
		return nil
	}
	out := new(NamespacedName)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrometheusSpec) DeepCopyInto(out *PrometheusSpec) {
	*out = *in
	if in.Service != nil {
		in, out := &in.Service, &out.Service
		*out = new(MonitorServiceSpec)
		(*in).DeepCopyInto(*out)
	}
	in.ResourceRequirements.DeepCopyInto(&out.ResourceRequirements)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrometheusSpec.
func (in *PrometheusSpec) DeepCopy() *PrometheusSpec {
	if in == nil {
		return nil
	}
	out := new(PrometheusSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrometheusStatus) DeepCopyInto(out *PrometheusStatus) {
	*out = *in
	out.ServiceRef = in.ServiceRef
	in.DorisMonitorComponentStatus.DeepCopyInto(&out.DorisMonitorComponentStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrometheusStatus.
func (in *PrometheusStatus) DeepCopy() *PrometheusStatus {
	if in == nil {
		return nil
	}
	out := new(PrometheusStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Promtail) DeepCopyInto(out *Promtail) {
	*out = *in
	out.DaemonSetRef = in.DaemonSetRef
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]appsv1.DaemonSetCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Promtail.
func (in *Promtail) DeepCopy() *Promtail {
	if in == nil {
		return nil
	}
	out := new(Promtail)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PromtailSpec) DeepCopyInto(out *PromtailSpec) {
	*out = *in
	in.ResourceRequirements.DeepCopyInto(&out.ResourceRequirements)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PromtailSpec.
func (in *PromtailSpec) DeepCopy() *PromtailSpec {
	if in == nil {
		return nil
	}
	out := new(PromtailSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReplicasRange) DeepCopyInto(out *ReplicasRange) {
	*out = *in
	if in.Min != nil {
		in, out := &in.Min, &out.Min
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReplicasRange.
func (in *ReplicasRange) DeepCopy() *ReplicasRange {
	if in == nil {
		return nil
	}
	out := new(ReplicasRange)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScalePeriodSeconds) DeepCopyInto(out *ScalePeriodSeconds) {
	*out = *in
	if in.ScaleUp != nil {
		in, out := &in.ScaleUp, &out.ScaleUp
		*out = new(int32)
		**out = **in
	}
	if in.ScaleDown != nil {
		in, out := &in.ScaleDown, &out.ScaleDown
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScalePeriodSeconds.
func (in *ScalePeriodSeconds) DeepCopy() *ScalePeriodSeconds {
	if in == nil {
		return nil
	}
	out := new(ScalePeriodSeconds)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UtilizationThresholdRange) DeepCopyInto(out *UtilizationThresholdRange) {
	*out = *in
	if in.Max != nil {
		in, out := &in.Max, &out.Max
		*out = new(int32)
		**out = **in
	}
	if in.Min != nil {
		in, out := &in.Min, &out.Min
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UtilizationThresholdRange.
func (in *UtilizationThresholdRange) DeepCopy() *UtilizationThresholdRange {
	if in == nil {
		return nil
	}
	out := new(UtilizationThresholdRange)
	in.DeepCopyInto(out)
	return out
}
