// +kubebuilder:object:generate=true
package deepcopy

// TestConfig is a sample struct for testing DeepCopyInto generation.
// +kubebuilder:object:generate=true
type TestConfig struct {
	Name     string            `json:"name"`
	Replicas int               `json:"replicas"`
	Labels   map[string]string `json:"labels,omitempty"`
	Ports    []int             `json:"ports,omitempty"`
	Nested   *NestedConfig     `json:"nested,omitempty"`
}

// NestedConfig demonstrates nested struct deep copying.
// +kubebuilder:object:generate=true
type NestedConfig struct {
	Host    string   `json:"host"`
	Port    int      `json:"port"`
	Options []string `json:"options,omitempty"`
}
