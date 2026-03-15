package v1

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/grafana/grafana-app-sdk/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestSpec_TitleRequired(t *testing.T) {
	spec := Spec{
		Title:  "Buy groceries",
		Status: SpecStatusOpen,
	}
	if spec.Title != "Buy groceries" {
		t.Errorf("expected title 'Buy groceries', got '%s'", spec.Title)
	}
}

func TestSpec_DescriptionOptional(t *testing.T) {
	// Description is a pointer, so nil means absent
	spec := Spec{
		Title:  "Test",
		Status: SpecStatusOpen,
	}
	if spec.Description != nil {
		t.Error("expected nil description for unset field")
	}

	desc := "Some description"
	spec.Description = &desc
	if *spec.Description != "Some description" {
		t.Errorf("expected 'Some description', got '%s'", *spec.Description)
	}
}

func TestSpecStatus_ValidValues(t *testing.T) {
	tests := []struct {
		name   string
		status SpecStatus
		want   string
	}{
		{"open status", SpecStatusOpen, "open"},
		{"in_progress status", SpecStatusInProgress, "in_progress"},
		{"done status", SpecStatusDone, "done"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.want {
				t.Errorf("expected '%s', got '%s'", tt.want, tt.status)
			}
		})
	}
}

func TestSpecStatus_AllConstantsDefined(t *testing.T) {
	valid := map[SpecStatus]bool{
		SpecStatusOpen:       true,
		SpecStatusInProgress: true,
		SpecStatusDone:       true,
	}
	if len(valid) != 3 {
		t.Errorf("expected 3 status constants, got %d", len(valid))
	}
}

func TestSpec_JSONSerialization(t *testing.T) {
	desc := "A test todo"
	spec := Spec{
		Title:       "Test TODO",
		Description: &desc,
		Status:      SpecStatusOpen,
	}

	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("failed to marshal spec: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal spec JSON: %v", err)
	}

	if parsed["title"] != "Test TODO" {
		t.Errorf("expected title 'Test TODO', got '%v'", parsed["title"])
	}
	if parsed["description"] != "A test todo" {
		t.Errorf("expected description 'A test todo', got '%v'", parsed["description"])
	}
	if parsed["status"] != "open" {
		t.Errorf("expected status 'open', got '%v'", parsed["status"])
	}
}

func TestSpec_JSONOmitsNilDescription(t *testing.T) {
	spec := Spec{
		Title:  "No description",
		Status: SpecStatusDone,
	}

	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("failed to marshal spec: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal spec JSON: %v", err)
	}

	if _, exists := parsed["description"]; exists {
		t.Error("expected description to be omitted when nil")
	}
}

func TestSpec_JSONDeserialization(t *testing.T) {
	input := `{"title":"From JSON","status":"in_progress"}`
	var spec Spec
	if err := json.Unmarshal([]byte(input), &spec); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if spec.Title != "From JSON" {
		t.Errorf("expected title 'From JSON', got '%s'", spec.Title)
	}
	if spec.Status != SpecStatusInProgress {
		t.Errorf("expected status 'in_progress', got '%s'", spec.Status)
	}
	if spec.Description != nil {
		t.Error("expected nil description")
	}
}

func TestTodo_ImplementsResourceObject(t *testing.T) {
	var _ resource.Object = &Todo{}
}

func TestTodoList_ImplementsListObject(t *testing.T) {
	var _ resource.ListObject = &TodoList{}
}

func TestTodo_NewTodo(t *testing.T) {
	todo := NewTodo()
	if todo == nil {
		t.Fatal("NewTodo() returned nil")
	}
	if todo.Spec.Title != "" {
		t.Errorf("expected empty title, got '%s'", todo.Spec.Title)
	}
	if todo.Spec.Status != "" {
		t.Errorf("expected empty status, got '%s'", todo.Spec.Status)
	}
}

func TestTodo_GetSetSpec(t *testing.T) {
	todo := NewTodo()
	desc := "Test description"
	spec := Spec{
		Title:       "Test",
		Description: &desc,
		Status:      SpecStatusOpen,
	}

	if err := todo.SetSpec(spec); err != nil {
		t.Fatalf("SetSpec failed: %v", err)
	}

	got := todo.GetSpec()
	gotSpec, ok := got.(Spec)
	if !ok {
		t.Fatal("GetSpec() did not return Spec type")
	}
	if gotSpec.Title != "Test" {
		t.Errorf("expected title 'Test', got '%s'", gotSpec.Title)
	}
}

func TestTodo_SetSpec_InvalidType(t *testing.T) {
	todo := NewTodo()
	err := todo.SetSpec("invalid")
	if err == nil {
		t.Error("expected error when setting invalid spec type")
	}
}

func TestTodo_Subresources(t *testing.T) {
	todo := NewTodo()
	subs := todo.GetSubresources()
	if _, ok := subs["status"]; !ok {
		t.Error("expected 'status' subresource")
	}

	status, ok := todo.GetSubresource("status")
	if !ok {
		t.Error("expected GetSubresource('status') to return true")
	}
	if _, ok := status.(Status); !ok {
		t.Error("expected status to be of type Status")
	}

	_, ok = todo.GetSubresource("nonexistent")
	if ok {
		t.Error("expected GetSubresource('nonexistent') to return false")
	}
}

func TestTodo_StaticMetadata(t *testing.T) {
	todo := NewTodo()
	todo.Name = "my-todo"
	todo.Namespace = "default"
	todo.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   APIGroup,
		Version: APIVersion,
		Kind:    "Todo",
	})

	sm := todo.GetStaticMetadata()
	if sm.Name != "my-todo" {
		t.Errorf("expected name 'my-todo', got '%s'", sm.Name)
	}
	if sm.Namespace != "default" {
		t.Errorf("expected namespace 'default', got '%s'", sm.Namespace)
	}
	if sm.Group != "todo.grafana.app" {
		t.Errorf("expected group 'todo.grafana.app', got '%s'", sm.Group)
	}
	if sm.Version != "v1" {
		t.Errorf("expected version 'v1', got '%s'", sm.Version)
	}
	if sm.Kind != "Todo" {
		t.Errorf("expected kind 'Todo', got '%s'", sm.Kind)
	}
}

func TestTodo_DeepCopy(t *testing.T) {
	desc := "Original"
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
			Title:       "Original Title",
			Description: &desc,
			Status:      SpecStatusOpen,
		},
	}

	cpy := todo.DeepCopy()
	if cpy.Spec.Title != "Original Title" {
		t.Errorf("expected copied title 'Original Title', got '%s'", cpy.Spec.Title)
	}

	// Mutating copy should not affect original
	cpy.Spec.Title = "Modified"
	if todo.Spec.Title != "Original Title" {
		t.Error("modifying copy affected original")
	}
}

func TestTodo_FullJSON(t *testing.T) {
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
		},
	}

	data, err := json.Marshal(todo)
	if err != nil {
		t.Fatalf("failed to marshal todo: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal todo JSON: %v", err)
	}

	// Verify top-level structure
	if parsed["apiVersion"] != "todo.grafana.app/v1" {
		t.Errorf("expected apiVersion 'todo.grafana.app/v1', got '%v'", parsed["apiVersion"])
	}
	if parsed["kind"] != "Todo" {
		t.Errorf("expected kind 'Todo', got '%v'", parsed["kind"])
	}

	// Verify metadata
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

	// Verify spec
	spec, ok := parsed["spec"].(map[string]interface{})
	if !ok {
		t.Fatal("expected spec to be a map")
	}
	if spec["title"] != "My TODO" {
		t.Errorf("expected title 'My TODO', got '%v'", spec["title"])
	}
	if spec["status"] != "open" {
		t.Errorf("expected status 'open', got '%v'", spec["status"])
	}
}

func TestSchema_GroupVersionKind(t *testing.T) {
	s := Schema()
	if s.Group() != "todo.grafana.app" {
		t.Errorf("expected group 'todo.grafana.app', got '%s'", s.Group())
	}
	if s.Version() != "v1" {
		t.Errorf("expected version 'v1', got '%s'", s.Version())
	}
	if s.Kind() != "Todo" {
		t.Errorf("expected kind 'Todo', got '%s'", s.Kind())
	}
}

func TestSchema_Plural(t *testing.T) {
	s := Schema()
	if s.Plural() != "todos" {
		t.Errorf("expected plural 'todos', got '%s'", s.Plural())
	}
}

func TestSchema_Scope(t *testing.T) {
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

func TestJSONCodec_ReadWrite(t *testing.T) {
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
}

func TestConstants_APIGroupAndVersion(t *testing.T) {
	if APIGroup != "todo.grafana.app" {
		t.Errorf("expected APIGroup 'todo.grafana.app', got '%s'", APIGroup)
	}
	if APIVersion != "v1" {
		t.Errorf("expected APIVersion 'v1', got '%s'", APIVersion)
	}
}

func TestGroupVersion(t *testing.T) {
	expected := schema.GroupVersion{Group: "todo.grafana.app", Version: "v1"}
	if GroupVersion != expected {
		t.Errorf("expected GroupVersion %v, got %v", expected, GroupVersion)
	}
}

func TestTodoList_GetSetItems(t *testing.T) {
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

	first := items[0].(*Todo)
	if first.Spec.Title != "First" {
		t.Errorf("expected first item title 'First', got '%s'", first.Spec.Title)
	}
}
