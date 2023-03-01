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

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DeploySpec defines the desired state of Deploy
type DeploySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// 镜像地址
	Image string `json:"image"`
	// 副本数
	Replicas int32 `json:"replicas"`
	// 环境变量
	Environments []corev1.EnvVar `json:"environments,omitempty"`
	// service要暴露的端口
	Expose *Expose `json:"expose"`
	// service和container端口
	Port int32 `json:"port,omitempty"`
}

// ExposeMode Expose模式
type ExposeMode string

const (
	ExposeModeNodePort ExposeMode = "nodePort"
	ExposeModeIngress  ExposeMode = "ingress"
)

type Expose struct {
	// 模式，如果为ExposeModeNodePort，忽略IngressDomain
	Mode ExposeMode `json:"mode"`
	// ingress域名
	//+optional
	IngressDomain string `json:"ingressDomain,omitempty"`
	// nodePort端口
	NodePort int32 `json:"nodePort,omitempty"`
	// Service端口
	ServicePort int32 `json:"servicePort,omitempty"`
	// ingress pre path
	Path string `json:"path,omitempty"`
}

// DeployStatus defines the observed state of Deploy
type DeployStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// 所处阶段
	Phase      string      `json:"phase,omitempty"`
	Message    string      `json:"message,omitempty"`
	Reason     string      `json:"reason,omitempty"`
	Conditions []Condition `json:"conditions,omitempty"`
}

type Condition struct {
	// 子资源类型
	Type string `json:"type"`
	// 子资源状态信息
	Message string `json:"message"`
	// 子资源状态名称
	Status string `json:"status"`
	// 处于这个状态的原因
	Reason string `json:"reason"`
	// 最后创建更新时间
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Deploy is the Schema for the deploys API
type Deploy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeploySpec   `json:"spec,omitempty"`
	Status DeployStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DeployList contains a list of Deploy
type DeployList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Deploy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Deploy{}, &DeployList{})
}
