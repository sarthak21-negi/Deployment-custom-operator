/*
Copyright 2026.

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

// DeploymentScaleTarget defines one deployment scaling target
// +kubebuilder:object:generate=true

type DeploymentScaleTarget struct {
	// Deployment name
	Name string `json:"name"`

	// Namespace of deployment
	Namespace string `json:"namespace"`

	// Desired replicas during active schedule
	Replicas int32 `json:"replicas"`
}

// ScalingSchedule defines time window
// +kubebuilder:object:generate=true

type ScalingSchedule struct {
	// Start hour in 24-hour format
	// Example: 9 means 9 AM
	StartHour int `json:"startHour"`

	// End hour in 24-hour format
	// Example: 18 means 6 PM
	EndHour int `json:"endHour"`

	// Replicas outside active period
	DefaultReplicas int32 `json:"defaultReplicas"`
}

// DeploymentCustomOperatorSpec defines the desired state of DeploymentCustomOperator.
// +kubebuilder:object:generate=true

type DeploymentCustomOperatorSpec struct {
	// List of deployments to scale
	Targets []DeploymentScaleTarget `json:"targets"`

	// Schedule configuration
	Schedule ScalingSchedule `json:"schedule"`
}

// DeploymentCustomOperatorStatus defines the observed state of DeploymentCustomOperator.
// +kubebuilder:object:generate=true

type DeploymentCustomOperatorStatus struct {
	LastScaleTime string `json:"lastScaleTime,omitempty"`

	Active bool `json:"active,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// DeploymentCustomOperator is the Schema for the deploymentcustomoperators API.
type DeploymentCustomOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeploymentCustomOperatorSpec   `json:"spec,omitempty"`
	Status DeploymentCustomOperatorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DeploymentCustomOperatorList contains a list of DeploymentCustomOperator.
type DeploymentCustomOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DeploymentCustomOperator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DeploymentCustomOperator{}, &DeploymentCustomOperatorList{})
}
