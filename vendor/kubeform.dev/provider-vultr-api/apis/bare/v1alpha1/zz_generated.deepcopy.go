//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright AppsCode Inc. and Contributors

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	apiv1alpha1 "kubeform.dev/apimachinery/api/v1alpha1"

	v1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	apiv1 "kmodules.xyz/client-go/api/v1"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetalServer) DeepCopyInto(out *MetalServer) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetalServer.
func (in *MetalServer) DeepCopy() *MetalServer {
	if in == nil {
		return nil
	}
	out := new(MetalServer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MetalServer) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetalServerList) DeepCopyInto(out *MetalServerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MetalServer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetalServerList.
func (in *MetalServerList) DeepCopy() *MetalServerList {
	if in == nil {
		return nil
	}
	out := new(MetalServerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MetalServerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetalServerSpec) DeepCopyInto(out *MetalServerSpec) {
	*out = *in
	if in.State != nil {
		in, out := &in.State, &out.State
		*out = new(MetalServerSpecResource)
		(*in).DeepCopyInto(*out)
	}
	in.Resource.DeepCopyInto(&out.Resource)
	out.ProviderRef = in.ProviderRef
	if in.SecretRef != nil {
		in, out := &in.SecretRef, &out.SecretRef
		*out = new(v1.LocalObjectReference)
		**out = **in
	}
	if in.BackendRef != nil {
		in, out := &in.BackendRef, &out.BackendRef
		*out = new(v1.LocalObjectReference)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetalServerSpec.
func (in *MetalServerSpec) DeepCopy() *MetalServerSpec {
	if in == nil {
		return nil
	}
	out := new(MetalServerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetalServerSpecResource) DeepCopyInto(out *MetalServerSpecResource) {
	*out = *in
	if in.Timeouts != nil {
		in, out := &in.Timeouts, &out.Timeouts
		*out = new(apiv1alpha1.ResourceTimeout)
		(*in).DeepCopyInto(*out)
	}
	if in.ActivationEmail != nil {
		in, out := &in.ActivationEmail, &out.ActivationEmail
		*out = new(bool)
		**out = **in
	}
	if in.AppID != nil {
		in, out := &in.AppID, &out.AppID
		*out = new(int64)
		**out = **in
	}
	if in.CpuCount != nil {
		in, out := &in.CpuCount, &out.CpuCount
		*out = new(int64)
		**out = **in
	}
	if in.DateCreated != nil {
		in, out := &in.DateCreated, &out.DateCreated
		*out = new(string)
		**out = **in
	}
	if in.DefaultPassword != nil {
		in, out := &in.DefaultPassword, &out.DefaultPassword
		*out = new(string)
		**out = **in
	}
	if in.Disk != nil {
		in, out := &in.Disk, &out.Disk
		*out = new(string)
		**out = **in
	}
	if in.EnableIpv6 != nil {
		in, out := &in.EnableIpv6, &out.EnableIpv6
		*out = new(bool)
		**out = **in
	}
	if in.GatewayV4 != nil {
		in, out := &in.GatewayV4, &out.GatewayV4
		*out = new(string)
		**out = **in
	}
	if in.Hostname != nil {
		in, out := &in.Hostname, &out.Hostname
		*out = new(string)
		**out = **in
	}
	if in.ImageID != nil {
		in, out := &in.ImageID, &out.ImageID
		*out = new(string)
		**out = **in
	}
	if in.Label != nil {
		in, out := &in.Label, &out.Label
		*out = new(string)
		**out = **in
	}
	if in.MacAddress != nil {
		in, out := &in.MacAddress, &out.MacAddress
		*out = new(int64)
		**out = **in
	}
	if in.MainIP != nil {
		in, out := &in.MainIP, &out.MainIP
		*out = new(string)
		**out = **in
	}
	if in.NetmaskV4 != nil {
		in, out := &in.NetmaskV4, &out.NetmaskV4
		*out = new(string)
		**out = **in
	}
	if in.Os != nil {
		in, out := &in.Os, &out.Os
		*out = new(string)
		**out = **in
	}
	if in.OsID != nil {
		in, out := &in.OsID, &out.OsID
		*out = new(int64)
		**out = **in
	}
	if in.Plan != nil {
		in, out := &in.Plan, &out.Plan
		*out = new(string)
		**out = **in
	}
	if in.Ram != nil {
		in, out := &in.Ram, &out.Ram
		*out = new(string)
		**out = **in
	}
	if in.Region != nil {
		in, out := &in.Region, &out.Region
		*out = new(string)
		**out = **in
	}
	if in.ReservedIpv4 != nil {
		in, out := &in.ReservedIpv4, &out.ReservedIpv4
		*out = new(string)
		**out = **in
	}
	if in.ScriptID != nil {
		in, out := &in.ScriptID, &out.ScriptID
		*out = new(string)
		**out = **in
	}
	if in.SnapshotID != nil {
		in, out := &in.SnapshotID, &out.SnapshotID
		*out = new(string)
		**out = **in
	}
	if in.SshKeyIDS != nil {
		in, out := &in.SshKeyIDS, &out.SshKeyIDS
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(string)
		**out = **in
	}
	if in.Tag != nil {
		in, out := &in.Tag, &out.Tag
		*out = new(string)
		**out = **in
	}
	if in.UserData != nil {
		in, out := &in.UserData, &out.UserData
		*out = new(string)
		**out = **in
	}
	if in.V6MainIP != nil {
		in, out := &in.V6MainIP, &out.V6MainIP
		*out = new(string)
		**out = **in
	}
	if in.V6Network != nil {
		in, out := &in.V6Network, &out.V6Network
		*out = new(string)
		**out = **in
	}
	if in.V6NetworkSize != nil {
		in, out := &in.V6NetworkSize, &out.V6NetworkSize
		*out = new(int64)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetalServerSpecResource.
func (in *MetalServerSpecResource) DeepCopy() *MetalServerSpecResource {
	if in == nil {
		return nil
	}
	out := new(MetalServerSpecResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetalServerStatus) DeepCopyInto(out *MetalServerStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]apiv1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetalServerStatus.
func (in *MetalServerStatus) DeepCopy() *MetalServerStatus {
	if in == nil {
		return nil
	}
	out := new(MetalServerStatus)
	in.DeepCopyInto(out)
	return out
}
