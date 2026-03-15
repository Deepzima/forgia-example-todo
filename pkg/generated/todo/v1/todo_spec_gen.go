// Code generated - EDITING IS FUTILE. DO NOT EDIT.

package v1

// +k8s:openapi-gen=true
type Spec struct {
	// Title of the TODO item (required)
	Title string `json:"title"`
	// Optional description of the TODO item
	Description *string `json:"description,omitempty"`
	// Current status of the TODO item
	Status SpecStatus `json:"status"`
	// Priority level of the TODO item (optional, defaults to "medium")
	Priority SpecPriority `json:"priority,omitempty"`
}

// NewSpec creates a new Spec object.
func NewSpec() *Spec {
	return &Spec{}
}

// OpenAPIModelName returns the OpenAPI model name for Spec.
func (Spec) OpenAPIModelName() string {
	return "com.github.zima.forgia-example-todo.pkg.generated.todo.v1.Spec"
}

// +k8s:openapi-gen=true
type SpecStatus string

const (
	SpecStatusOpen       SpecStatus = "open"
	SpecStatusInProgress SpecStatus = "in_progress"
	SpecStatusDone       SpecStatus = "done"
)

// OpenAPIModelName returns the OpenAPI model name for SpecStatus.
func (SpecStatus) OpenAPIModelName() string {
	return "com.github.zima.forgia-example-todo.pkg.generated.todo.v1.SpecStatus"
}

// +k8s:openapi-gen=true
type SpecPriority string

const (
	SpecPriorityLow      SpecPriority = "low"
	SpecPriorityMedium   SpecPriority = "medium"
	SpecPriorityHigh     SpecPriority = "high"
	SpecPriorityCritical SpecPriority = "critical"
)

// OpenAPIModelName returns the OpenAPI model name for SpecPriority.
func (SpecPriority) OpenAPIModelName() string {
	return "com.github.zima.forgia-example-todo.pkg.generated.todo.v1.SpecPriority"
}
