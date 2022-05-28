/*
Copyright 2022.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LTargetSpec defines the desired state of LTarget
type LTargetSpec struct {
	// InternalDelvePortOrName points to the target's ContainerPort (or name).
	InternalDelvePortOrName intstr.IntOrString `json:"internalDelvePortOrName"`
	// ExternalDelvePort points to the IDE's listening port.
	ExternalDelvePort int32 `json:"externalDelvePort"`
	// LTargetLabel is the selector used to filter out the parent pod.
	LTargetLabel map[string]string `json:"lTargetLabel"`
}

// LTargetStatus defines the observed state of LTarget
type LTargetStatus struct {
	// ConnectionStatus exhibits the current status of LTarget's connection to the IDE.
	ConnectionStatus string `json:"connectionStatus"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LTarget is the Schema for the ltargets API
type LTarget struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LTargetSpec   `json:"spec,omitempty"`
	Status LTargetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LTargetList contains a list of LTarget
type LTargetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LTarget `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LTarget{}, &LTargetList{})
}
