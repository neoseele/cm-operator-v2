package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CustomMetricSpec defines the desired state of CustomMetric
type CustomMetricSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	Project  string   `json:"project"`
	Cluster  string   `json:"cluster"`
	Location string   `json:"location"`
	Metrics  []string `json:"metrics"`
}

// CustomMetricStatus defines the observed state of CustomMetric
type CustomMetricStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CustomMetric is the Schema for the custommetrics API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=custommetrics,scope=Namespaced
type CustomMetric struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CustomMetricSpec   `json:"spec,omitempty"`
	Status CustomMetricStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CustomMetricList contains a list of CustomMetric
type CustomMetricList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CustomMetric `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CustomMetric{}, &CustomMetricList{})
}
