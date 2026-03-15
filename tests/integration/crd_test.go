//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

func getKubeconfig(t *testing.T) string {
	t.Helper()
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			t.Skipf("skipping: cannot determine home dir: %v", err)
		}
		kubeconfig = home + "/.kube/config"
	}
	if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
		t.Skip("skipping: no kubeconfig found")
	}
	return kubeconfig
}

func TestCRD_Install(t *testing.T) {
	kubeconfig := getKubeconfig(t)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		t.Fatalf("failed to build config: %v", err)
	}

	apiextClient, err := apiextensionsclient.NewForConfig(config)
	if err != nil {
		t.Fatalf("failed to create apiextensions client: %v", err)
	}

	// Read CRD definition
	crdJSON, err := os.ReadFile("../../definitions/todo.todo.grafana.app.json")
	if err != nil {
		t.Fatalf("failed to read CRD file: %v", err)
	}

	var crd apiextensionsv1.CustomResourceDefinition
	if err := json.Unmarshal(crdJSON, &crd); err != nil {
		t.Fatalf("failed to unmarshal CRD: %v", err)
	}

	ctx := context.Background()

	// Clean up if exists
	_ = apiextClient.ApiextensionsV1().CustomResourceDefinitions().Delete(
		ctx, crd.Name, metav1.DeleteOptions{},
	)

	// Apply CRD
	created, err := apiextClient.ApiextensionsV1().CustomResourceDefinitions().Create(
		ctx, &crd, metav1.CreateOptions{},
	)
	if err != nil {
		t.Fatalf("failed to create CRD: %v", err)
	}
	t.Cleanup(func() {
		_ = apiextClient.ApiextensionsV1().CustomResourceDefinitions().Delete(
			ctx, created.Name, metav1.DeleteOptions{},
		)
	})

	// Verify CRD was created
	got, err := apiextClient.ApiextensionsV1().CustomResourceDefinitions().Get(
		ctx, "todos.todo.grafana.app", metav1.GetOptions{},
	)
	if err != nil {
		t.Fatalf("failed to get CRD: %v", err)
	}

	if got.Spec.Group != "todo.grafana.app" {
		t.Errorf("expected group %q, got %q", "todo.grafana.app", got.Spec.Group)
	}
	if got.Spec.Names.Kind != "Todo" {
		t.Errorf("expected kind %q, got %q", "Todo", got.Spec.Names.Kind)
	}

	t.Log("CRD installed successfully")
}

func TestCRD_CreateTodoResource(t *testing.T) {
	kubeconfig := getKubeconfig(t)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		t.Fatalf("failed to build config: %v", err)
	}

	apiextClient, err := apiextensionsclient.NewForConfig(config)
	if err != nil {
		t.Fatalf("failed to create apiextensions client: %v", err)
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		t.Fatalf("failed to create dynamic client: %v", err)
	}

	// Ensure CRD is installed
	crdJSON, err := os.ReadFile("../../definitions/todo.todo.grafana.app.json")
	if err != nil {
		t.Fatalf("failed to read CRD file: %v", err)
	}

	var crd apiextensionsv1.CustomResourceDefinition
	if err := json.Unmarshal(crdJSON, &crd); err != nil {
		t.Fatalf("failed to unmarshal CRD: %v", err)
	}

	ctx := context.Background()
	_ = apiextClient.ApiextensionsV1().CustomResourceDefinitions().Delete(
		ctx, crd.Name, metav1.DeleteOptions{},
	)
	_, err = apiextClient.ApiextensionsV1().CustomResourceDefinitions().Create(
		ctx, &crd, metav1.CreateOptions{},
	)
	if err != nil {
		t.Fatalf("failed to create CRD: %v", err)
	}
	t.Cleanup(func() {
		_ = apiextClient.ApiextensionsV1().CustomResourceDefinitions().Delete(
			ctx, crd.Name, metav1.DeleteOptions{},
		)
	})

	// Create a Todo resource
	todoGVR := schema.GroupVersionResource{
		Group:    "todo.grafana.app",
		Version:  "v1",
		Resource: "todos",
	}

	todo := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "todo.grafana.app/v1",
			"kind":       "Todo",
			"metadata": map[string]interface{}{
				"name":      "integration-test-todo",
				"namespace": "default",
			},
			"spec": map[string]interface{}{
				"title":       "Integration Test Todo",
				"description": "Created by integration test",
				"status":      "open",
			},
		},
	}

	created, err := dynClient.Resource(todoGVR).Namespace("default").Create(
		ctx, todo, metav1.CreateOptions{},
	)
	if err != nil {
		t.Fatalf("failed to create Todo resource: %v", err)
	}
	t.Cleanup(func() {
		_ = dynClient.Resource(todoGVR).Namespace("default").Delete(
			ctx, "integration-test-todo", metav1.DeleteOptions{},
		)
	})

	// Verify
	spec, ok := created.Object["spec"].(map[string]interface{})
	if !ok {
		t.Fatal("spec not found in created resource")
	}
	if spec["title"] != "Integration Test Todo" {
		t.Errorf("expected title %q, got %q", "Integration Test Todo", spec["title"])
	}
	if spec["status"] != "open" {
		t.Errorf("expected status %q, got %q", "open", spec["status"])
	}

	t.Log("Todo resource created and verified successfully")
}

func TestCRD_RejectInvalidStatus(t *testing.T) {
	kubeconfig := getKubeconfig(t)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		t.Fatalf("failed to build config: %v", err)
	}

	apiextClient, err := apiextensionsclient.NewForConfig(config)
	if err != nil {
		t.Fatalf("failed to create apiextensions client: %v", err)
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		t.Fatalf("failed to create dynamic client: %v", err)
	}

	// Ensure CRD is installed
	crdJSON, err := os.ReadFile("../../definitions/todo.todo.grafana.app.json")
	if err != nil {
		t.Fatalf("failed to read CRD file: %v", err)
	}

	var crd apiextensionsv1.CustomResourceDefinition
	if err := json.Unmarshal(crdJSON, &crd); err != nil {
		t.Fatalf("failed to unmarshal CRD: %v", err)
	}

	ctx := context.Background()
	_ = apiextClient.ApiextensionsV1().CustomResourceDefinitions().Delete(
		ctx, crd.Name, metav1.DeleteOptions{},
	)
	_, err = apiextClient.ApiextensionsV1().CustomResourceDefinitions().Create(
		ctx, &crd, metav1.CreateOptions{},
	)
	if err != nil {
		t.Fatalf("failed to create CRD: %v", err)
	}
	t.Cleanup(func() {
		_ = apiextClient.ApiextensionsV1().CustomResourceDefinitions().Delete(
			ctx, crd.Name, metav1.DeleteOptions{},
		)
	})

	// Try to create a Todo with invalid status
	todoGVR := schema.GroupVersionResource{
		Group:    "todo.grafana.app",
		Version:  "v1",
		Resource: "todos",
	}

	invalidTodo := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "todo.grafana.app/v1",
			"kind":       "Todo",
			"metadata": map[string]interface{}{
				"name":      "invalid-todo",
				"namespace": "default",
			},
			"spec": map[string]interface{}{
				"title":  "Invalid Status Todo",
				"status": "invalid_status",
			},
		},
	}

	_, err = dynClient.Resource(todoGVR).Namespace("default").Create(
		ctx, invalidTodo, metav1.CreateOptions{},
	)
	if err == nil {
		// Clean up if it was accidentally created
		_ = dynClient.Resource(todoGVR).Namespace("default").Delete(
			ctx, "invalid-todo", metav1.DeleteOptions{},
		)
		t.Fatal("expected error when creating Todo with invalid status, but got none")
	}

	t.Logf("correctly rejected invalid status: %v", err)
}
