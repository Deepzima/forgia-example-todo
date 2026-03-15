package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	todov1 "github.com/zima/forgia-example-todo/pkg/generated/todo/v1"
	"github.com/zima/forgia-example-todo/pkg/repository"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TodoRequest is the JSON body for creating or updating a Todo.
type TodoRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Status      string  `json:"status"`
	Priority    *string `json:"priority,omitempty"`
}

// ErrorResponse is a standard error response body.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// TodoHandler handles HTTP requests for Todo resources.
type TodoHandler struct {
	repo   repository.TodoRepository
	logger *slog.Logger
}

// NewTodoHandler creates a new handler with the given repository and logger.
func NewTodoHandler(repo repository.TodoRepository, logger *slog.Logger) *TodoHandler {
	return &TodoHandler{
		repo:   repo,
		logger: logger,
	}
}

// CreateTodo handles POST requests to create a new Todo.
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request, namespace string) {
	req, err := decodeRequest(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", fmt.Sprintf("failed to decode request body: %v", err))
		return
	}

	if validationErr := validateTodoRequest(req); validationErr != "" {
		writeError(w, http.StatusUnprocessableEntity, "validation_error", validationErr)
		return
	}

	priority := todov1.SpecPriorityMedium
	if req.Priority != nil {
		priority = todov1.SpecPriority(*req.Priority)
	}

	todo := &todov1.Todo{
		TypeMeta: metav1.TypeMeta{
			Kind:       todov1.Kind().Kind(),
			APIVersion: todov1.GroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    namespace,
			GenerateName: "todo-",
		},
		Spec: todov1.Spec{
			Title:       req.Title,
			Description: req.Description,
			Status:      todov1.SpecStatus(req.Status),
			Priority:    priority,
		},
	}

	created, err := h.repo.Create(r.Context(), todo)
	if err != nil {
		h.logger.Error("failed to create todo", "error", err)
		writeError(w, http.StatusInternalServerError, "create_failed", fmt.Sprintf("failed to create todo: %v", err))
		return
	}

	h.logger.Info("todo created", "name", created.Name, "namespace", namespace)
	writeJSON(w, http.StatusCreated, created)
}

// GetTodo handles GET requests for a single Todo.
func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request, namespace, name string) {
	todo, err := h.repo.Get(r.Context(), namespace, name)
	if err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, "not_found", fmt.Sprintf("todo %s/%s not found", namespace, name))
			return
		}
		h.logger.Error("failed to get todo", "namespace", namespace, "name", name, "error", err)
		writeError(w, http.StatusInternalServerError, "get_failed", fmt.Sprintf("failed to get todo: %v", err))
		return
	}

	ensurePriority(todo)
	writeJSON(w, http.StatusOK, todo)
}

// ListTodos handles GET requests to list all Todos in a namespace.
func (h *TodoHandler) ListTodos(w http.ResponseWriter, r *http.Request, namespace string) {
	list, err := h.repo.List(r.Context(), namespace)
	if err != nil {
		h.logger.Error("failed to list todos", "namespace", namespace, "error", err)
		writeError(w, http.StatusInternalServerError, "list_failed", fmt.Sprintf("failed to list todos: %v", err))
		return
	}

	for i := range list.Items {
		ensurePriority(&list.Items[i])
	}
	writeJSON(w, http.StatusOK, list)
}

// UpdateTodo handles PUT requests to update an existing Todo.
func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request, namespace, name string) {
	req, err := decodeRequest(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", fmt.Sprintf("failed to decode request body: %v", err))
		return
	}

	if validationErr := validateTodoRequest(req); validationErr != "" {
		writeError(w, http.StatusUnprocessableEntity, "validation_error", validationErr)
		return
	}

	// Fetch current to get resourceVersion for optimistic concurrency
	existing, err := h.repo.Get(r.Context(), namespace, name)
	if err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, "not_found", fmt.Sprintf("todo %s/%s not found", namespace, name))
			return
		}
		h.logger.Error("failed to get todo for update", "namespace", namespace, "name", name, "error", err)
		writeError(w, http.StatusInternalServerError, "get_failed", fmt.Sprintf("failed to get todo: %v", err))
		return
	}

	priority := existing.Spec.Priority
	if req.Priority != nil {
		priority = todov1.SpecPriority(*req.Priority)
	}

	existing.Spec = todov1.Spec{
		Title:       req.Title,
		Description: req.Description,
		Status:      todov1.SpecStatus(req.Status),
		Priority:    priority,
	}

	updated, err := h.repo.Update(r.Context(), existing)
	if err != nil {
		h.logger.Error("failed to update todo", "namespace", namespace, "name", name, "error", err)
		writeError(w, http.StatusInternalServerError, "update_failed", fmt.Sprintf("failed to update todo: %v", err))
		return
	}

	h.logger.Info("todo updated", "name", name, "namespace", namespace)
	writeJSON(w, http.StatusOK, updated)
}

// DeleteTodo handles DELETE requests to remove a Todo.
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request, namespace, name string) {
	// Verify the todo exists before deleting
	_, err := h.repo.Get(r.Context(), namespace, name)
	if err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, "not_found", fmt.Sprintf("todo %s/%s not found", namespace, name))
			return
		}
		h.logger.Error("failed to get todo for delete", "namespace", namespace, "name", name, "error", err)
		writeError(w, http.StatusInternalServerError, "get_failed", fmt.Sprintf("failed to get todo: %v", err))
		return
	}

	if err := h.repo.Delete(r.Context(), namespace, name); err != nil {
		h.logger.Error("failed to delete todo", "namespace", namespace, "name", name, "error", err)
		writeError(w, http.StatusInternalServerError, "delete_failed", fmt.Sprintf("failed to delete todo: %v", err))
		return
	}

	h.logger.Info("todo deleted", "name", name, "namespace", namespace)
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// validateTodoRequest validates a TodoRequest and returns an error message or empty string.
func validateTodoRequest(req *TodoRequest) string {
	if strings.TrimSpace(req.Title) == "" {
		return "title is required and cannot be empty"
	}

	switch todov1.SpecStatus(req.Status) {
	case todov1.SpecStatusOpen, todov1.SpecStatusInProgress, todov1.SpecStatusDone:
		// valid
	default:
		return fmt.Sprintf("invalid status %q: must be one of: open, in_progress, done", req.Status)
	}

	if req.Priority != nil {
		if msg := validatePriority(*req.Priority); msg != "" {
			return msg
		}
	}

	return ""
}

func validatePriority(value string) string {
	switch todov1.SpecPriority(value) {
	case todov1.SpecPriorityLow, todov1.SpecPriorityMedium,
		todov1.SpecPriorityHigh, todov1.SpecPriorityCritical:
		return ""
	default:
		return fmt.Sprintf("invalid priority: must be one of low, medium, high, critical")
	}
}

// ensurePriority defaults empty priority to "medium" for existing Todos.
func ensurePriority(todo *todov1.Todo) {
	if todo.Spec.Priority == "" {
		todo.Spec.Priority = todov1.SpecPriorityMedium
	}
}

func decodeRequest(body io.Reader) (*TodoRequest, error) {
	var req TodoRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

func writeError(w http.ResponseWriter, status int, errType, message string) {
	writeJSON(w, status, ErrorResponse{
		Error:   errType,
		Message: message,
	})
}

// isNotFound checks if an error indicates a resource was not found.
func isNotFound(err error) bool {
	return strings.Contains(err.Error(), "not found") ||
		strings.Contains(err.Error(), "NotFound") ||
		strings.Contains(err.Error(), "404")
}
