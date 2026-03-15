package repository

import (
	"context"
	"fmt"

	"github.com/grafana/grafana-app-sdk/resource"
	todov1 "github.com/zima/forgia-example-todo/pkg/generated/todo/v1"
)

// TodoRepository defines the interface for Todo data access.
type TodoRepository interface {
	Get(ctx context.Context, namespace, name string) (*todov1.Todo, error)
	List(ctx context.Context, namespace string) (*todov1.TodoList, error)
	Create(ctx context.Context, todo *todov1.Todo) (*todov1.Todo, error)
	Update(ctx context.Context, todo *todov1.Todo) (*todov1.Todo, error)
	Delete(ctx context.Context, namespace, name string) error
}

// K8sTodoRepository implements TodoRepository using the generated TodoClient.
type K8sTodoRepository struct {
	client *todov1.TodoClient
}

// NewK8sTodoRepository creates a new repository backed by the Kubernetes API.
func NewK8sTodoRepository(client *todov1.TodoClient) *K8sTodoRepository {
	return &K8sTodoRepository{client: client}
}

func (r *K8sTodoRepository) Get(ctx context.Context, namespace, name string) (*todov1.Todo, error) {
	todo, err := r.client.Get(ctx, resource.Identifier{
		Namespace: namespace,
		Name:      name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get todo %s/%s: %w", namespace, name, err)
	}
	return todo, nil
}

func (r *K8sTodoRepository) List(ctx context.Context, namespace string) (*todov1.TodoList, error) {
	list, err := r.client.List(ctx, namespace, resource.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list todos in namespace %s: %w", namespace, err)
	}
	return list, nil
}

func (r *K8sTodoRepository) Create(ctx context.Context, todo *todov1.Todo) (*todov1.Todo, error) {
	created, err := r.client.Create(ctx, todo, resource.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create todo %s: %w", todo.Name, err)
	}
	return created, nil
}

func (r *K8sTodoRepository) Update(ctx context.Context, todo *todov1.Todo) (*todov1.Todo, error) {
	updated, err := r.client.Update(ctx, todo, resource.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to update todo %s: %w", todo.Name, err)
	}
	return updated, nil
}

func (r *K8sTodoRepository) Delete(ctx context.Context, namespace, name string) error {
	err := r.client.Delete(ctx, resource.Identifier{
		Namespace: namespace,
		Name:      name,
	}, resource.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete todo %s/%s: %w", namespace, name, err)
	}
	return nil
}

// Interface compliance check
var _ TodoRepository = &K8sTodoRepository{}
