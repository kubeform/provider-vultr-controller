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

// Code generated by Kubeform. DO NOT EDIT.

package v1alpha1

import (
	base "kubeform.dev/apimachinery/api/v1alpha1"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	"sigs.k8s.io/cli-utils/pkg/kstatus/status"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`

type MetalServer struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MetalServerSpec   `json:"spec,omitempty"`
	Status            MetalServerStatus `json:"status,omitempty"`
}

type MetalServerSpec struct {
	State *MetalServerSpecResource `json:"state,omitempty" tf:"-"`

	Resource MetalServerSpecResource `json:"resource" tf:"resource"`

	UpdatePolicy base.UpdatePolicy `json:"updatePolicy,omitempty" tf:"-"`

	TerminationPolicy base.TerminationPolicy `json:"terminationPolicy,omitempty" tf:"-"`

	ProviderRef core.LocalObjectReference `json:"providerRef" tf:"-"`

	SecretRef *core.LocalObjectReference `json:"secretRef,omitempty" tf:"-"`

	BackendRef *core.LocalObjectReference `json:"backendRef,omitempty" tf:"-"`
}

type MetalServerSpecResource struct {
	ID string `json:"id,omitempty" tf:"id,omitempty"`

	// +optional
	ActivationEmail *bool `json:"activationEmail,omitempty" tf:"activation_email"`
	// +optional
	AppID *int64 `json:"appID,omitempty" tf:"app_id"`
	// +optional
	CpuCount *int64 `json:"cpuCount,omitempty" tf:"cpu_count"`
	// +optional
	DateCreated *string `json:"dateCreated,omitempty" tf:"date_created"`
	// +optional
	DefaultPassword *string `json:"-" sensitive:"true" tf:"default_password"`
	// +optional
	Disk *string `json:"disk,omitempty" tf:"disk"`
	// +optional
	EnableIpv6 *bool `json:"enableIpv6,omitempty" tf:"enable_ipv6"`
	// +optional
	GatewayV4 *string `json:"gatewayV4,omitempty" tf:"gateway_v4"`
	// +optional
	Hostname *string `json:"hostname,omitempty" tf:"hostname"`
	// +optional
	ImageID *string `json:"imageID,omitempty" tf:"image_id"`
	// +optional
	Label *string `json:"label,omitempty" tf:"label"`
	// +optional
	MacAddress *int64 `json:"macAddress,omitempty" tf:"mac_address"`
	// +optional
	MainIP *string `json:"mainIP,omitempty" tf:"main_ip"`
	// +optional
	NetmaskV4 *string `json:"netmaskV4,omitempty" tf:"netmask_v4"`
	// +optional
	Os *string `json:"os,omitempty" tf:"os"`
	// +optional
	OsID *int64  `json:"osID,omitempty" tf:"os_id"`
	Plan *string `json:"plan" tf:"plan"`
	// +optional
	Ram    *string `json:"ram,omitempty" tf:"ram"`
	Region *string `json:"region" tf:"region"`
	// +optional
	ReservedIpv4 *string `json:"reservedIpv4,omitempty" tf:"reserved_ipv4"`
	// +optional
	ScriptID *string `json:"scriptID,omitempty" tf:"script_id"`
	// +optional
	SnapshotID *string `json:"snapshotID,omitempty" tf:"snapshot_id"`
	// +optional
	SshKeyIDS []string `json:"sshKeyIDS,omitempty" tf:"ssh_key_ids"`
	// +optional
	Status *string `json:"status,omitempty" tf:"status"`
	// +optional
	Tag *string `json:"tag,omitempty" tf:"tag"`
	// +optional
	UserData *string `json:"userData,omitempty" tf:"user_data"`
	// +optional
	V6MainIP *string `json:"v6MainIP,omitempty" tf:"v6_main_ip"`
	// +optional
	V6Network *string `json:"v6Network,omitempty" tf:"v6_network"`
	// +optional
	V6NetworkSize *int64 `json:"v6NetworkSize,omitempty" tf:"v6_network_size"`
}

type MetalServerStatus struct {
	// Resource generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// +optional
	Phase status.Status `json:"phase,omitempty"`
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// MetalServerList is a list of MetalServers
type MetalServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of MetalServer CRD objects
	Items []MetalServer `json:"items,omitempty"`
}
