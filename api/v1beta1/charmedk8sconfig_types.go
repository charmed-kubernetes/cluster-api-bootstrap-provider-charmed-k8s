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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CharmedK8sConfigSpec defines the desired state of CharmedK8sConfig
type CharmedK8sConfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// List of Juju applications to deploy to control plane machines
	ControlPlaneApplications []string `json:"controlPlaneApplications,omitempty"`
	// List of Juju applications to deploy to worker machines
	WorkerApplications []string `json:"workerApplications,omitempty"`
}

// CharmedK8sConfigStatus defines the observed state of CharmedK8sConfig
type CharmedK8sConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Required fields for bootstrap providers

	Ready          bool   `json:"ready,omitempty"`
	DataSecretName string `json:"dataSecretName,omitempty"`

	// Optional fields for bootstrap providers

	FailureReason  string `json:"failureReason,omitempty"`  // error string for programs
	FailureMessage string `json:"failureMessage,omitempty"` // error string for humans
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CharmedK8sConfig is the Schema for the charmedk8sconfigs API
type CharmedK8sConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CharmedK8sConfigSpec   `json:"spec,omitempty"`
	Status CharmedK8sConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CharmedK8sConfigList contains a list of CharmedK8sConfig
type CharmedK8sConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CharmedK8sConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CharmedK8sConfig{}, &CharmedK8sConfigList{})
}
