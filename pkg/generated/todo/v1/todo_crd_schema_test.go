package v1

import (
	"encoding/json"
	"os"
	"testing"
)

// crdDefinition represents the relevant parts of the CRD JSON for validation.
type crdDefinition struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name string `json:"name"`
	} `json:"metadata"`
	Spec struct {
		Group    string `json:"group"`
		Scope    string `json:"scope"`
		Names    struct {
			Kind   string `json:"kind"`
			Plural string `json:"plural"`
		} `json:"names"`
		Versions []struct {
			Name    string `json:"name"`
			Served  bool   `json:"served"`
			Storage bool   `json:"storage"`
			Schema  struct {
				OpenAPIV3Schema struct {
					Required   []string                          `json:"required"`
					Properties map[string]map[string]interface{} `json:"properties"`
				} `json:"openAPIV3Schema"`
			} `json:"schema"`
		} `json:"versions"`
	} `json:"spec"`
}

func loadCRD(t *testing.T) crdDefinition {
	t.Helper()
	data, err := os.ReadFile("../../../../definitions/todo.todo.grafana.app.json")
	if err != nil {
		t.Fatalf("failed to read CRD file: %v", err)
	}
	var crd crdDefinition
	if err := json.Unmarshal(data, &crd); err != nil {
		t.Fatalf("failed to unmarshal CRD: %v", err)
	}
	return crd
}

func TestCRDSchema_APIGroupAndKind(t *testing.T) {
	crd := loadCRD(t)

	if crd.APIVersion != "apiextensions.k8s.io/v1" {
		t.Errorf("expected apiVersion %q, got %q", "apiextensions.k8s.io/v1", crd.APIVersion)
	}
	if crd.Kind != "CustomResourceDefinition" {
		t.Errorf("expected kind %q, got %q", "CustomResourceDefinition", crd.Kind)
	}
	if crd.Spec.Group != "todo.grafana.app" {
		t.Errorf("expected group %q, got %q", "todo.grafana.app", crd.Spec.Group)
	}
	if crd.Spec.Names.Kind != "Todo" {
		t.Errorf("expected names.kind %q, got %q", "Todo", crd.Spec.Names.Kind)
	}
	if crd.Spec.Names.Plural != "todos" {
		t.Errorf("expected names.plural %q, got %q", "todos", crd.Spec.Names.Plural)
	}
	if crd.Spec.Scope != "Namespaced" {
		t.Errorf("expected scope %q, got %q", "Namespaced", crd.Spec.Scope)
	}
}

func TestCRDSchema_V1Served(t *testing.T) {
	crd := loadCRD(t)

	if len(crd.Spec.Versions) == 0 {
		t.Fatal("expected at least one version")
	}

	v1 := crd.Spec.Versions[0]
	if v1.Name != "v1" {
		t.Errorf("expected version name %q, got %q", "v1", v1.Name)
	}
	if !v1.Served {
		t.Error("expected v1 to be served")
	}
	if !v1.Storage {
		t.Error("expected v1 to be storage version")
	}
}

func TestCRDSchema_SpecFieldsValidation(t *testing.T) {
	// Re-read raw JSON to inspect nested schema properties
	data, err := os.ReadFile("../../../../definitions/todo.todo.grafana.app.json")
	if err != nil {
		t.Fatalf("failed to read CRD file: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Navigate to spec.versions[0].schema.openAPIV3Schema.properties.spec
	versions := raw["spec"].(map[string]interface{})["versions"].([]interface{})
	v1Schema := versions[0].(map[string]interface{})["schema"].(map[string]interface{})
	openAPI := v1Schema["openAPIV3Schema"].(map[string]interface{})
	props := openAPI["properties"].(map[string]interface{})
	specProps := props["spec"].(map[string]interface{})["properties"].(map[string]interface{})

	// Verify title field exists and is a string
	titleField := specProps["title"].(map[string]interface{})
	if titleField["type"] != "string" {
		t.Errorf("expected title type %q, got %q", "string", titleField["type"])
	}

	// Verify description field exists and is a string
	descField := specProps["description"].(map[string]interface{})
	if descField["type"] != "string" {
		t.Errorf("expected description type %q, got %q", "string", descField["type"])
	}

	// Verify status field has enum constraint
	statusField := specProps["status"].(map[string]interface{})
	if statusField["type"] != "string" {
		t.Errorf("expected status type %q, got %q", "string", statusField["type"])
	}
	enumValues, ok := statusField["enum"].([]interface{})
	if !ok {
		t.Fatal("expected status field to have enum constraint")
	}

	expectedEnums := map[string]bool{"open": false, "in_progress": false, "done": false}
	for _, v := range enumValues {
		s, ok := v.(string)
		if !ok {
			t.Errorf("expected enum value to be string, got %T", v)
			continue
		}
		if _, exists := expectedEnums[s]; !exists {
			t.Errorf("unexpected enum value %q", s)
		}
		expectedEnums[s] = true
	}
	for k, found := range expectedEnums {
		if !found {
			t.Errorf("missing expected enum value %q", k)
		}
	}

	// Verify required fields in spec
	specRequired := props["spec"].(map[string]interface{})["required"].([]interface{})
	requiredMap := make(map[string]bool)
	for _, r := range specRequired {
		requiredMap[r.(string)] = true
	}
	if !requiredMap["title"] {
		t.Error("title should be required in spec")
	}
	if !requiredMap["status"] {
		t.Error("status should be required in spec")
	}
	if requiredMap["description"] {
		t.Error("description should NOT be required in spec")
	}
	if requiredMap["priority"] {
		t.Error("priority should NOT be required in spec (optional with default)")
	}

	// Verify priority field has enum constraint and default
	priorityField, ok := specProps["priority"].(map[string]interface{})
	if !ok {
		t.Fatal("expected priority field in spec properties")
	}
	if priorityField["type"] != "string" {
		t.Errorf("expected priority type %q, got %q", "string", priorityField["type"])
	}
	if priorityField["default"] != "medium" {
		t.Errorf("expected priority default %q, got %v", "medium", priorityField["default"])
	}
	priorityEnums, ok := priorityField["enum"].([]interface{})
	if !ok {
		t.Fatal("expected priority field to have enum constraint")
	}
	expectedPriorityEnums := map[string]bool{"low": false, "medium": false, "high": false, "critical": false}
	for _, v := range priorityEnums {
		s, ok := v.(string)
		if !ok {
			t.Errorf("expected enum value to be string, got %T", v)
			continue
		}
		if _, exists := expectedPriorityEnums[s]; !exists {
			t.Errorf("unexpected priority enum value %q", s)
		}
		expectedPriorityEnums[s] = true
	}
	for k, found := range expectedPriorityEnums {
		if !found {
			t.Errorf("missing expected priority enum value %q", k)
		}
	}
}
