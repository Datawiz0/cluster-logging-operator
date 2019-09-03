// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1

import (
	loggingv1 "github.com/openshift/elasticsearch-operator/pkg/apis/logging/v1"
	corev1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterCondition) DeepCopyInto(out *ClusterCondition) {
	*out = *in
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterCondition.
func (in *ClusterCondition) DeepCopy() *ClusterCondition {
	if in == nil {
		return nil
	}
	out := new(ClusterCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterLogging) DeepCopyInto(out *ClusterLogging) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterLogging.
func (in *ClusterLogging) DeepCopy() *ClusterLogging {
	if in == nil {
		return nil
	}
	out := new(ClusterLogging)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterLogging) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterLoggingList) DeepCopyInto(out *ClusterLoggingList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClusterLogging, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterLoggingList.
func (in *ClusterLoggingList) DeepCopy() *ClusterLoggingList {
	if in == nil {
		return nil
	}
	out := new(ClusterLoggingList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterLoggingList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterLoggingSpec) DeepCopyInto(out *ClusterLoggingSpec) {
	*out = *in
	in.Visualization.DeepCopyInto(&out.Visualization)
	in.LogStore.DeepCopyInto(&out.LogStore)
	in.Collection.DeepCopyInto(&out.Collection)
	in.Curation.DeepCopyInto(&out.Curation)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterLoggingSpec.
func (in *ClusterLoggingSpec) DeepCopy() *ClusterLoggingSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterLoggingSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterLoggingStatus) DeepCopyInto(out *ClusterLoggingStatus) {
	*out = *in
	in.Visualization.DeepCopyInto(&out.Visualization)
	in.LogStore.DeepCopyInto(&out.LogStore)
	in.Collection.DeepCopyInto(&out.Collection)
	in.Curation.DeepCopyInto(&out.Curation)
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]ClusterCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterLoggingStatus.
func (in *ClusterLoggingStatus) DeepCopy() *ClusterLoggingStatus {
	if in == nil {
		return nil
	}
	out := new(ClusterLoggingStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CollectionSpec) DeepCopyInto(out *CollectionSpec) {
	*out = *in
	in.Logs.DeepCopyInto(&out.Logs)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CollectionSpec.
func (in *CollectionSpec) DeepCopy() *CollectionSpec {
	if in == nil {
		return nil
	}
	out := new(CollectionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CollectionStatus) DeepCopyInto(out *CollectionStatus) {
	*out = *in
	in.Logs.DeepCopyInto(&out.Logs)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CollectionStatus.
func (in *CollectionStatus) DeepCopy() *CollectionStatus {
	if in == nil {
		return nil
	}
	out := new(CollectionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CurationSpec) DeepCopyInto(out *CurationSpec) {
	*out = *in
	in.CuratorSpec.DeepCopyInto(&out.CuratorSpec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CurationSpec.
func (in *CurationSpec) DeepCopy() *CurationSpec {
	if in == nil {
		return nil
	}
	out := new(CurationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CurationStatus) DeepCopyInto(out *CurationStatus) {
	*out = *in
	if in.CuratorStatus != nil {
		in, out := &in.CuratorStatus, &out.CuratorStatus
		*out = make([]CuratorStatus, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CurationStatus.
func (in *CurationStatus) DeepCopy() *CurationStatus {
	if in == nil {
		return nil
	}
	out := new(CurationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CuratorSpec) DeepCopyInto(out *CuratorSpec) {
	*out = *in
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(corev1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]corev1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CuratorSpec.
func (in *CuratorSpec) DeepCopy() *CuratorSpec {
	if in == nil {
		return nil
	}
	out := new(CuratorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CuratorStatus) DeepCopyInto(out *CuratorStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(map[string][]ClusterCondition, len(*in))
		for key, val := range *in {
			var outVal []ClusterCondition
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]ClusterCondition, len(*in))
				for i := range *in {
					(*in)[i].DeepCopyInto(&(*out)[i])
				}
			}
			(*out)[key] = outVal
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CuratorStatus.
func (in *CuratorStatus) DeepCopy() *CuratorStatus {
	if in == nil {
		return nil
	}
	out := new(CuratorStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ElasticsearchSpec) DeepCopyInto(out *ElasticsearchSpec) {
	*out = *in
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(corev1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]corev1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Storage.DeepCopyInto(&out.Storage)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ElasticsearchSpec.
func (in *ElasticsearchSpec) DeepCopy() *ElasticsearchSpec {
	if in == nil {
		return nil
	}
	out := new(ElasticsearchSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ElasticsearchStatus) DeepCopyInto(out *ElasticsearchStatus) {
	*out = *in
	if in.ReplicaSets != nil {
		in, out := &in.ReplicaSets, &out.ReplicaSets
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Deployments != nil {
		in, out := &in.Deployments, &out.Deployments
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.StatefulSets != nil {
		in, out := &in.StatefulSets, &out.StatefulSets
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.Cluster = in.Cluster
	if in.Pods != nil {
		in, out := &in.Pods, &out.Pods
		*out = make(map[ElasticsearchRoleType]PodStateMap, len(*in))
		for key, val := range *in {
			var outVal map[PodStateType][]string
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make(PodStateMap, len(*in))
				for key, val := range *in {
					var outVal []string
					if val == nil {
						(*out)[key] = nil
					} else {
						in, out := &val, &outVal
						*out = make([]string, len(*in))
						copy(*out, *in)
					}
					(*out)[key] = outVal
				}
			}
			(*out)[key] = outVal
		}
	}
	if in.ClusterConditions != nil {
		in, out := &in.ClusterConditions, &out.ClusterConditions
		*out = make([]loggingv1.ClusterCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.NodeConditions != nil {
		in, out := &in.NodeConditions, &out.NodeConditions
		*out = make(map[string][]loggingv1.ClusterCondition, len(*in))
		for key, val := range *in {
			var outVal []loggingv1.ClusterCondition
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]loggingv1.ClusterCondition, len(*in))
				for i := range *in {
					(*in)[i].DeepCopyInto(&(*out)[i])
				}
			}
			(*out)[key] = outVal
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ElasticsearchStatus.
func (in *ElasticsearchStatus) DeepCopy() *ElasticsearchStatus {
	if in == nil {
		return nil
	}
	out := new(ElasticsearchStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EventCollectionSpec) DeepCopyInto(out *EventCollectionSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EventCollectionSpec.
func (in *EventCollectionSpec) DeepCopy() *EventCollectionSpec {
	if in == nil {
		return nil
	}
	out := new(EventCollectionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EventCollectionStatus) DeepCopyInto(out *EventCollectionStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EventCollectionStatus.
func (in *EventCollectionStatus) DeepCopy() *EventCollectionStatus {
	if in == nil {
		return nil
	}
	out := new(EventCollectionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FluentdCollectorStatus) DeepCopyInto(out *FluentdCollectorStatus) {
	*out = *in
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Pods != nil {
		in, out := &in.Pods, &out.Pods
		*out = make(PodStateMap, len(*in))
		for key, val := range *in {
			var outVal []string
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]string, len(*in))
				copy(*out, *in)
			}
			(*out)[key] = outVal
		}
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(map[string][]ClusterCondition, len(*in))
		for key, val := range *in {
			var outVal []ClusterCondition
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]ClusterCondition, len(*in))
				for i := range *in {
					(*in)[i].DeepCopyInto(&(*out)[i])
				}
			}
			(*out)[key] = outVal
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FluentdCollectorStatus.
func (in *FluentdCollectorStatus) DeepCopy() *FluentdCollectorStatus {
	if in == nil {
		return nil
	}
	out := new(FluentdCollectorStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FluentdNormalizerStatus) DeepCopyInto(out *FluentdNormalizerStatus) {
	*out = *in
	if in.ReplicaSets != nil {
		in, out := &in.ReplicaSets, &out.ReplicaSets
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Pods != nil {
		in, out := &in.Pods, &out.Pods
		*out = make(PodStateMap, len(*in))
		for key, val := range *in {
			var outVal []string
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]string, len(*in))
				copy(*out, *in)
			}
			(*out)[key] = outVal
		}
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(map[string][]ClusterCondition, len(*in))
		for key, val := range *in {
			var outVal []ClusterCondition
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]ClusterCondition, len(*in))
				for i := range *in {
					(*in)[i].DeepCopyInto(&(*out)[i])
				}
			}
			(*out)[key] = outVal
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FluentdNormalizerStatus.
func (in *FluentdNormalizerStatus) DeepCopy() *FluentdNormalizerStatus {
	if in == nil {
		return nil
	}
	out := new(FluentdNormalizerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FluentdSpec) DeepCopyInto(out *FluentdSpec) {
	*out = *in
	in.OperandSpec.DeepCopyInto(&out.OperandSpec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FluentdSpec.
func (in *FluentdSpec) DeepCopy() *FluentdSpec {
	if in == nil {
		return nil
	}
	out := new(FluentdSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KibanaSpec) DeepCopyInto(out *KibanaSpec) {
	*out = *in
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(corev1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]corev1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.ProxySpec.DeepCopyInto(&out.ProxySpec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KibanaSpec.
func (in *KibanaSpec) DeepCopy() *KibanaSpec {
	if in == nil {
		return nil
	}
	out := new(KibanaSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KibanaStatus) DeepCopyInto(out *KibanaStatus) {
	*out = *in
	if in.ReplicaSets != nil {
		in, out := &in.ReplicaSets, &out.ReplicaSets
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Pods != nil {
		in, out := &in.Pods, &out.Pods
		*out = make(PodStateMap, len(*in))
		for key, val := range *in {
			var outVal []string
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]string, len(*in))
				copy(*out, *in)
			}
			(*out)[key] = outVal
		}
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(map[string][]ClusterCondition, len(*in))
		for key, val := range *in {
			var outVal []ClusterCondition
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]ClusterCondition, len(*in))
				for i := range *in {
					(*in)[i].DeepCopyInto(&(*out)[i])
				}
			}
			(*out)[key] = outVal
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KibanaStatus.
func (in *KibanaStatus) DeepCopy() *KibanaStatus {
	if in == nil {
		return nil
	}
	out := new(KibanaStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LogCollectionSpec) DeepCopyInto(out *LogCollectionSpec) {
	*out = *in
	in.FluentdSpec.DeepCopyInto(&out.FluentdSpec)
	in.PromTailSpec.DeepCopyInto(&out.PromTailSpec)
	in.RsyslogSpec.DeepCopyInto(&out.RsyslogSpec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LogCollectionSpec.
func (in *LogCollectionSpec) DeepCopy() *LogCollectionSpec {
	if in == nil {
		return nil
	}
	out := new(LogCollectionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LogCollectionStatus) DeepCopyInto(out *LogCollectionStatus) {
	*out = *in
	in.FluentdStatus.DeepCopyInto(&out.FluentdStatus)
	in.RsyslogStatus.DeepCopyInto(&out.RsyslogStatus)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LogCollectionStatus.
func (in *LogCollectionStatus) DeepCopy() *LogCollectionStatus {
	if in == nil {
		return nil
	}
	out := new(LogCollectionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LogStoreSpec) DeepCopyInto(out *LogStoreSpec) {
	*out = *in
	in.ElasticsearchSpec.DeepCopyInto(&out.ElasticsearchSpec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LogStoreSpec.
func (in *LogStoreSpec) DeepCopy() *LogStoreSpec {
	if in == nil {
		return nil
	}
	out := new(LogStoreSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LogStoreStatus) DeepCopyInto(out *LogStoreStatus) {
	*out = *in
	if in.ElasticsearchStatus != nil {
		in, out := &in.ElasticsearchStatus, &out.ElasticsearchStatus
		*out = make([]ElasticsearchStatus, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LogStoreStatus.
func (in *LogStoreStatus) DeepCopy() *LogStoreStatus {
	if in == nil {
		return nil
	}
	out := new(LogStoreStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NormalizerStatus) DeepCopyInto(out *NormalizerStatus) {
	*out = *in
	if in.FluentdStatus != nil {
		in, out := &in.FluentdStatus, &out.FluentdStatus
		*out = make([]FluentdNormalizerStatus, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NormalizerStatus.
func (in *NormalizerStatus) DeepCopy() *NormalizerStatus {
	if in == nil {
		return nil
	}
	out := new(NormalizerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OperandSpec) DeepCopyInto(out *OperandSpec) {
	*out = *in
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(corev1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]corev1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OperandSpec.
func (in *OperandSpec) DeepCopy() *OperandSpec {
	if in == nil {
		return nil
	}
	out := new(OperandSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in PodStateMap) DeepCopyInto(out *PodStateMap) {
	{
		in := &in
		*out = make(PodStateMap, len(*in))
		for key, val := range *in {
			var outVal []string
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]string, len(*in))
				copy(*out, *in)
			}
			(*out)[key] = outVal
		}
		return
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodStateMap.
func (in PodStateMap) DeepCopy() PodStateMap {
	if in == nil {
		return nil
	}
	out := new(PodStateMap)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PromTailSpec) DeepCopyInto(out *PromTailSpec) {
	*out = *in
	in.OperandSpec.DeepCopyInto(&out.OperandSpec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PromTailSpec.
func (in *PromTailSpec) DeepCopy() *PromTailSpec {
	if in == nil {
		return nil
	}
	out := new(PromTailSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProxySpec) DeepCopyInto(out *ProxySpec) {
	*out = *in
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(corev1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProxySpec.
func (in *ProxySpec) DeepCopy() *ProxySpec {
	if in == nil {
		return nil
	}
	out := new(ProxySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RsyslogCollectorStatus) DeepCopyInto(out *RsyslogCollectorStatus) {
	*out = *in
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Pods != nil {
		in, out := &in.Pods, &out.Pods
		*out = make(PodStateMap, len(*in))
		for key, val := range *in {
			var outVal []string
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]string, len(*in))
				copy(*out, *in)
			}
			(*out)[key] = outVal
		}
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(map[string][]ClusterCondition, len(*in))
		for key, val := range *in {
			var outVal []ClusterCondition
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]ClusterCondition, len(*in))
				for i := range *in {
					(*in)[i].DeepCopyInto(&(*out)[i])
				}
			}
			(*out)[key] = outVal
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RsyslogCollectorStatus.
func (in *RsyslogCollectorStatus) DeepCopy() *RsyslogCollectorStatus {
	if in == nil {
		return nil
	}
	out := new(RsyslogCollectorStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RsyslogSpec) DeepCopyInto(out *RsyslogSpec) {
	*out = *in
	in.OperandSpec.DeepCopyInto(&out.OperandSpec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RsyslogSpec.
func (in *RsyslogSpec) DeepCopy() *RsyslogSpec {
	if in == nil {
		return nil
	}
	out := new(RsyslogSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VisualizationSpec) DeepCopyInto(out *VisualizationSpec) {
	*out = *in
	in.KibanaSpec.DeepCopyInto(&out.KibanaSpec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VisualizationSpec.
func (in *VisualizationSpec) DeepCopy() *VisualizationSpec {
	if in == nil {
		return nil
	}
	out := new(VisualizationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VisualizationStatus) DeepCopyInto(out *VisualizationStatus) {
	*out = *in
	if in.KibanaStatus != nil {
		in, out := &in.KibanaStatus, &out.KibanaStatus
		*out = make([]KibanaStatus, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VisualizationStatus.
func (in *VisualizationStatus) DeepCopy() *VisualizationStatus {
	if in == nil {
		return nil
	}
	out := new(VisualizationStatus)
	in.DeepCopyInto(out)
	return out
}
