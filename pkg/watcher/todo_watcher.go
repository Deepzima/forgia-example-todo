package watcher

import (
	"context"
	"log/slog"
	"time"

	"github.com/grafana/grafana-app-sdk/resource"
	todov1 "github.com/zima/forgia-example-todo/pkg/generated/todo/v1"
	"github.com/zima/forgia-example-todo/pkg/repository"
)

const operatorID = "todo-operator"

// TodoWatcher watches Todo resources and manages their lifecycle.
type TodoWatcher struct {
	repo   repository.TodoRepository
	logger *slog.Logger
}

// NewTodoWatcher creates a new watcher for Todo lifecycle management.
func NewTodoWatcher(repo repository.TodoRepository, logger *slog.Logger) *TodoWatcher {
	return &TodoWatcher{
		repo:   repo,
		logger: logger,
	}
}

// Add is called when a new Todo resource is created.
func (w *TodoWatcher) Add(ctx context.Context, obj resource.Object) error {
	todo, ok := obj.(*todov1.Todo)
	if !ok {
		w.logger.Error("received non-Todo object in Add")
		return nil
	}

	w.logger.Info("todo created",
		"name", todo.Name,
		"namespace", todo.Namespace,
		"title", todo.Spec.Title,
		"status", todo.Spec.Status,
	)

	return w.updateOperatorStatus(ctx, todo, todov1.StatusOperatorStateStateSuccess, "resource created")
}

// Update is called when a Todo resource is modified.
func (w *TodoWatcher) Update(ctx context.Context, old, new resource.Object) error {
	oldTodo, ok := old.(*todov1.Todo)
	if !ok {
		return nil
	}
	newTodo, ok := new.(*todov1.Todo)
	if !ok {
		return nil
	}

	w.logger.Info("todo updated",
		"name", newTodo.Name,
		"namespace", newTodo.Namespace,
		"oldStatus", oldTodo.Spec.Status,
		"newStatus", newTodo.Spec.Status,
	)

	return w.updateOperatorStatus(ctx, newTodo, todov1.StatusOperatorStateStateSuccess, "resource updated")
}

// Delete is called when a Todo resource is deleted.
func (w *TodoWatcher) Delete(ctx context.Context, obj resource.Object) error {
	todo, ok := obj.(*todov1.Todo)
	if !ok {
		return nil
	}

	w.logger.Info("todo deleted",
		"name", todo.Name,
		"namespace", todo.Namespace,
	)

	// No status update needed on delete - resource is being removed
	return nil
}

// Sync is called on operator restart for pre-existing resources.
func (w *TodoWatcher) Sync(ctx context.Context, obj resource.Object) error {
	todo, ok := obj.(*todov1.Todo)
	if !ok {
		return nil
	}

	w.logger.Info("todo synced",
		"name", todo.Name,
		"namespace", todo.Namespace,
		"status", todo.Spec.Status,
	)

	return nil
}

// updateOperatorStatus sets the operator state on a Todo's status subresource.
func (w *TodoWatcher) updateOperatorStatus(ctx context.Context, todo *todov1.Todo, state todov1.StatusOperatorStateState, description string) error {
	if todo.Status.OperatorStates == nil {
		todo.Status.OperatorStates = make(map[string]todov1.StatusOperatorState)
	}

	desc := description
	todo.Status.OperatorStates[operatorID] = todov1.StatusOperatorState{
		LastEvaluation:   todo.ResourceVersion,
		State:            state,
		DescriptiveState: &desc,
		Details: map[string]interface{}{
			"lastUpdated": time.Now().UTC().Format(time.RFC3339),
		},
	}

	return nil
}
