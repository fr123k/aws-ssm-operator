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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ParameterStoreSpec defines the desired state of ParameterStore
type ParameterStoreSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	ValueFrom ValueFrom `json:"valueFrom"`
}

type ValueFrom struct {
	// +kubebuilder:validation:Optional
	ParameterStoreRef ParameterStoreRef `json:"parameterStoreRef"`
	// +kubebuilder:validation:Optional
	ParametersStoreRef []ParametersStoreRef `json:"parametersStoreRef"`
}

type ParameterStoreRef struct {
	Name string `json:"name"`
	Path string `json:"path"`
	// +kubebuilder:default:=true
	Recursive bool `json:"recursive"`
}

type ParametersStoreRef struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// ParameterStoreStatus defines the observed state of ParameterStore
type ParameterStoreStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ParameterStore is the Schema for the parameterstores API
type ParameterStore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ParameterStoreSpec   `json:"spec,omitempty"`
	Status ParameterStoreStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ParameterStoreList contains a list of ParameterStore
type ParameterStoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ParameterStore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ParameterStore{}, &ParameterStoreList{})
}
