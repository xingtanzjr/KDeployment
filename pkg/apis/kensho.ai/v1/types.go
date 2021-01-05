package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KDeploymentSpec struct {
	TotalReplicas *int32 `json:"totalReplicas"`
	ReplicaPolicy string `json:"replicaPolicy"`

	DeploymentTemplate appsv1.DeploymentSpec `json:"deploymentSpec"`
}

type KDeploymentStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KDeploymentSpec   `json:"spec,omitempty"`
	Status KDeploymentStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDeployment `json:"items"`
}
