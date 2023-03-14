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

type CharmedK8sConfigTemplateResource struct {
	Spec CharmedK8sConfigSpec `json:"spec"`
}

// CharmedK8sConfigTemplateSpec defines the desired state of CharmedK8sConfigTemplate
type CharmedK8sConfigTemplateSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Template CharmedK8sConfigConfigTemplateResource `json:"template"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:path=charmedk8sconfigtemplates,scope=Namespaced,categories=cluster-api,shortName=ckct
//+kubebuilder:storageversion

// CharmedK8sConfigTemplate is the Schema for the charmedk8sconfigtemplates API
type CharmedK8sConfigTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CharmedK8sConfigTemplateSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// CharmedK8sConfigTemplateList contains a list of CharmedK8sConfigTemplate
type CharmedK8sConfigTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CharmedK8sConfigTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CharmedK8sConfigTemplate{}, &CharmedK8sConfigTemplateList{})
}
