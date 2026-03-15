package v1

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/grafana/grafana-app-sdk/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestTodo_ImplementsResourceObject(t *testing.T) {
	var _ resource.Object = &Todo{}
}

func TestTodoList_ImplementsResourceListObject(t *testing.T) {
	var _ resource.ListObject = &TodoList{}
}

func TestTodo_FullJSONStructure(t *testing.T) {
	desc := "Full test"
	todo := &Todo{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "todo.grafana.app/v1",
			Kind:       "Todo",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-todo",
			Namespace: "default",
		},
		Spec: Spec{
			Title:       "My TODO",
			Description: &desc,
			Status:      SpecStatusOpen,
			Priority:    SpecPriorityHigh,
		},
	}

	data, err := json.Marshal(todo)
	if err != nil {
		t.Fatalf("failed to marshal todo: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if parsed["apiVersion"] != "todo.grafana.app/v1" {
		t.Errorf("expected apiVersion 'todo.grafana.app/v1', got '%v'", parsed["apiVersion"])
	}
	if parsed["kind"] != "Todo" {
		t.Errorf("expected kind 'Todo', got '%v'", parsed["kind"])
	}

	metadata, ok := parsed["metadata"].(map[string]interface{})
	if !ok {
		t.Fatal("expected metadata to be a map")
	}
	if metadata["name"] != "test-todo" {
		t.Errorf("expected name 'test-todo', got '%v'", metadata["name"])
	}
	if metadata["namespace"] != "default" {
		t.Errorf("expected namespace 'default', got '%v'", metadata["namespace"])
	}

	spec, ok := parsed["spec"].(map[string]interface{})
	if !ok {
		t.Fatal("expected spec to be a map")
	}
	if spec["title"] != "My TODO" {
		t.Errorf("expected title 'My TODO', got '%v'", spec["title"])
	}
	if spec["description"] != "Full test" {
		t.Errorf("expected description 'Full test', got '%v'", spec["description"])
	}
	if spec["status"] != "open" {
		t.Errorf("expected status 'open', got '%v'", spec["status"])
	}
	if spec["priority"] != "high" {
		t.Errorf("expected priority 'high', got '%v'", spec["priority"])
	}
}

func TestJSONCodec_RoundTrip(t *testing.T) {
	desc := "Codec test"
	original := &Todo{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "todo.grafana.app/v1",
			Kind:       "Todo",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "codec-test",
			Namespace: "default",
		},
		Spec: Spec{
			Title:       "Codec Test",
			Description: &desc,
			Status:      SpecStatusDone,
		},
	}

	codec := &JSONCodec{}
	var buf bytes.Buffer

	if err := codec.Write(&buf, original); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	decoded := NewTodo()
	if err := codec.Read(&buf, decoded); err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if decoded.Spec.Title != "Codec Test" {
		t.Errorf("expected title 'Codec Test', got '%s'", decoded.Spec.Title)
	}
	if decoded.Spec.Status != SpecStatusDone {
		t.Errorf("expected status 'done', got '%s'", decoded.Spec.Status)
	}
	if decoded.Name != "codec-test" {
		t.Errorf("expected name 'codec-test', got '%s'", decoded.Name)
	}
}

func TestConstants_GroupAndVersion(t *testing.T) {
	if APIGroup != "todo.grafana.app" {
		t.Errorf("expected APIGroup 'todo.grafana.app', got '%s'", APIGroup)
	}
	if APIVersion != "v1" {
		t.Errorf("expected APIVersion 'v1', got '%s'", APIVersion)
	}
	expected := schema.GroupVersion{Group: "todo.grafana.app", Version: "v1"}
	if GroupVersion != expected {
		t.Errorf("expected GroupVersion %v, got %v", expected, GroupVersion)
	}
}

func TestSchema_ScopeIsNamespaced(t *testing.T) {
	s := Schema()
	if s.Scope() != resource.NamespacedScope {
		t.Errorf("expected NamespacedScope, got '%v'", s.Scope())
	}
}

func TestKind_HasJSONCodec(t *testing.T) {
	k := Kind()
	codec, ok := k.Codecs[resource.KindEncodingJSON]
	if !ok {
		t.Fatal("expected JSON codec in Kind")
	}
	if codec == nil {
		t.Fatal("JSON codec is nil")
	}
}

func TestSpecPriority_ConstantsMatchExpectedValues(t *testing.T) {
	tests := []struct {
		constant SpecPriority
		expected string
	}{
		{SpecPriorityLow, "low"},
		{SpecPriorityMedium, "medium"},
		{SpecPriorityHigh, "high"},
		{SpecPriorityCritical, "critical"},
	}
	for _, tt := range tests {
		if string(tt.constant) != tt.expected {
			t.Errorf("expected %q, got %q", tt.expected, string(tt.constant))
		}
	}
}

func TestSpec_PriorityOmittedWhenEmpty(t *testing.T) {
	// Backward compatibility: existing Todos without priority should serialize without the field
	spec := Spec{
		Title:  "No priority",
		Status: SpecStatusOpen,
	}
	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if _, exists := parsed["priority"]; exists {
		t.Error("priority should be omitted when empty (backward compatibility)")
	}
}

func TestSpec_PriorityIncludedWhenSet(t *testing.T) {
	spec := Spec{
		Title:    "With priority",
		Status:   SpecStatusOpen,
		Priority: SpecPriorityHigh,
	}
	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if parsed["priority"] != "high" {
		t.Errorf("expected priority %q, got %v", "high", parsed["priority"])
	}
}

func TestTodoList_SetAndGetItems(t *testing.T) {
	list := &TodoList{}
	todo1 := NewTodo()
	todo1.Name = "todo-1"
	todo1.Spec.Title = "First"
	todo1.Spec.Status = SpecStatusOpen

	todo2 := NewTodo()
	todo2.Name = "todo-2"
	todo2.Spec.Title = "Second"
	todo2.Spec.Status = SpecStatusDone

	list.SetItems([]resource.Object{todo1, todo2})

	items := list.GetItems()
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0].(*Todo).Spec.Title != "First" {
		t.Errorf("expected first item title 'First', got '%s'", items[0].(*Todo).Spec.Title)
	}
}

