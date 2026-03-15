package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"log/slog"

	todov1 "github.com/zima/forgia-example-todo/pkg/generated/todo/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// mockTodoRepository is a test double for TodoRepository.
type mockTodoRepository struct {
	todos     map[string]*todov1.Todo
	lastErr   error
	getErr    error
	createErr error
	updateErr error
	deleteErr error
}

func newMockRepo() *mockTodoRepository {
	return &mockTodoRepository{
		todos: make(map[string]*todov1.Todo),
	}
}

func (m *mockTodoRepository) Get(_ context.Context, namespace, name string) (*todov1.Todo, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	if m.lastErr != nil {
		return nil, m.lastErr
	}
	key := namespace + "/" + name
	todo, ok := m.todos[key]
	if !ok {
		return nil, fmt.Errorf("todo %s not found", key)
	}
	return todo.DeepCopy(), nil
}

func (m *mockTodoRepository) List(_ context.Context, namespace string) (*todov1.TodoList, error) {
	if m.lastErr != nil {
		return nil, m.lastErr
	}
	var items []todov1.Todo
	for _, t := range m.todos {
		if t.Namespace == namespace {
			items = append(items, *t.DeepCopy())
		}
	}
	return &todov1.TodoList{Items: items}, nil
}

func (m *mockTodoRepository) Create(_ context.Context, todo *todov1.Todo) (*todov1.Todo, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	if m.lastErr != nil {
		return nil, m.lastErr
	}
	created := todo.DeepCopy()
	if created.Name == "" {
		created.Name = "todo-abc123"
	}
	created.UID = "uid-123"
	created.ResourceVersion = "1"
	key := created.Namespace + "/" + created.Name
	m.todos[key] = created
	return created.DeepCopy(), nil
}

func (m *mockTodoRepository) Update(_ context.Context, todo *todov1.Todo) (*todov1.Todo, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	if m.lastErr != nil {
		return nil, m.lastErr
	}
	key := todo.Namespace + "/" + todo.Name
	if _, ok := m.todos[key]; !ok {
		return nil, fmt.Errorf("todo %s not found", key)
	}
	updated := todo.DeepCopy()
	m.todos[key] = updated
	return updated.DeepCopy(), nil
}

func (m *mockTodoRepository) Delete(_ context.Context, namespace, name string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if m.lastErr != nil {
		return m.lastErr
	}
	key := namespace + "/" + name
	delete(m.todos, key)
	return nil
}

func (m *mockTodoRepository) addTodo(namespace, name, title string, status todov1.SpecStatus) {
	key := namespace + "/" + name
	m.todos[key] = &todov1.Todo{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Todo",
			APIVersion: "todo.grafana.app/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Namespace:       namespace,
			ResourceVersion: "1",
			UID:             types.UID("uid-" + name),
		},
		Spec: todov1.Spec{
			Title:  title,
			Status: status,
		},
	}
}

func newTestHandler() (*TodoHandler, *mockTodoRepository) {
	repo := newMockRepo()
	logger := slog.Default()
	h := NewTodoHandler(repo, logger)
	return h, repo
}

func makeRequest(method, body string) *http.Request {
	var r *http.Request
	if body != "" {
		r = httptest.NewRequest(method, "/", bytes.NewBufferString(body))
	} else {
		r = httptest.NewRequest(method, "/", nil)
	}
	r.Header.Set("Content-Type", "application/json")
	return r
}

func decodeResponse(t *testing.T, w *httptest.ResponseRecorder, v any) {
	t.Helper()
	if err := json.NewDecoder(w.Body).Decode(v); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
}

// --- CreateTodo tests ---

func TestCreateTodo_Success(t *testing.T) {
	h, _ := newTestHandler()
	w := httptest.NewRecorder()
	body := `{"title": "Buy groceries", "status": "open"}`
	r := makeRequest(http.MethodPost, body)

	h.CreateTodo(w, r, "default")

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}
	var todo todov1.Todo
	decodeResponse(t, w, &todo)
	if todo.Spec.Title != "Buy groceries" {
		t.Errorf("expected title 'Buy groceries', got %q", todo.Spec.Title)
	}
	if todo.Spec.Status != todov1.SpecStatusOpen {
		t.Errorf("expected status 'open', got %q", todo.Spec.Status)
	}
	if todo.Name == "" {
		t.Error("expected generated name, got empty")
	}
}

func TestCreateTodo_WithDescription(t *testing.T) {
	h, _ := newTestHandler()
	w := httptest.NewRecorder()
	body := `{"title": "Buy groceries", "description": "Milk, eggs, bread", "status": "open"}`
	r := makeRequest(http.MethodPost, body)

	h.CreateTodo(w, r, "default")

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}
	var todo todov1.Todo
	decodeResponse(t, w, &todo)
	if todo.Spec.Description == nil || *todo.Spec.Description != "Milk, eggs, bread" {
		t.Errorf("expected description 'Milk, eggs, bread', got %v", todo.Spec.Description)
	}
}

func TestCreateTodo_MissingTitle(t *testing.T) {
	h, _ := newTestHandler()
	w := httptest.NewRecorder()
	body := `{"status": "open"}`
	r := makeRequest(http.MethodPost, body)

	h.CreateTodo(w, r, "default")

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d", w.Code)
	}
	var errResp ErrorResponse
	decodeResponse(t, w, &errResp)
	if errResp.Error != "validation_error" {
		t.Errorf("expected error 'validation_error', got %q", errResp.Error)
	}
}

func TestCreateTodo_EmptyTitle(t *testing.T) {
	h, _ := newTestHandler()
	w := httptest.NewRecorder()
	body := `{"title": "   ", "status": "open"}`
	r := makeRequest(http.MethodPost, body)

	h.CreateTodo(w, r, "default")

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d", w.Code)
	}
}

func TestCreateTodo_InvalidStatus(t *testing.T) {
	h, _ := newTestHandler()
	w := httptest.NewRecorder()
	body := `{"title": "Test", "status": "invalid"}`
	r := makeRequest(http.MethodPost, body)

	h.CreateTodo(w, r, "default")

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d", w.Code)
	}
	var errResp ErrorResponse
	decodeResponse(t, w, &errResp)
	if errResp.Error != "validation_error" {
		t.Errorf("expected error 'validation_error', got %q", errResp.Error)
	}
}

func TestCreateTodo_InvalidJSON(t *testing.T) {
	h, _ := newTestHandler()
	w := httptest.NewRecorder()
	body := `{invalid json}`
	r := makeRequest(http.MethodPost, body)

	h.CreateTodo(w, r, "default")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestCreateTodo_RepoError(t *testing.T) {
	h, repo := newTestHandler()
	repo.lastErr = fmt.Errorf("connection refused")
	w := httptest.NewRecorder()
	body := `{"title": "Test", "status": "open"}`
	r := makeRequest(http.MethodPost, body)

	h.CreateTodo(w, r, "default")

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestCreateTodo_AllStatuses(t *testing.T) {
	statuses := []string{"open", "in_progress", "done"}
	for _, status := range statuses {
		t.Run(status, func(t *testing.T) {
			h, _ := newTestHandler()
			w := httptest.NewRecorder()
			body := fmt.Sprintf(`{"title": "Test", "status": %q}`, status)
			r := makeRequest(http.MethodPost, body)

			h.CreateTodo(w, r, "default")

			if w.Code != http.StatusCreated {
				t.Fatalf("expected status 201 for status %q, got %d", status, w.Code)
			}
		})
	}
}

// --- GetTodo tests ---

func TestGetTodo_Success(t *testing.T) {
	h, repo := newTestHandler()
	repo.addTodo("default", "todo-1", "My Todo", todov1.SpecStatusOpen)
	w := httptest.NewRecorder()
	r := makeRequest(http.MethodGet, "")

	h.GetTodo(w, r, "default", "todo-1")

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var todo todov1.Todo
	decodeResponse(t, w, &todo)
	if todo.Spec.Title != "My Todo" {
		t.Errorf("expected title 'My Todo', got %q", todo.Spec.Title)
	}
}

func TestGetTodo_NotFound(t *testing.T) {
	h, _ := newTestHandler()
	w := httptest.NewRecorder()
	r := makeRequest(http.MethodGet, "")

	h.GetTodo(w, r, "default", "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestGetTodo_InternalError(t *testing.T) {
	h, repo := newTestHandler()
	repo.getErr = fmt.Errorf("connection timeout")
	w := httptest.NewRecorder()
	r := makeRequest(http.MethodGet, "")

	h.GetTodo(w, r, "default", "todo-1")

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

// --- ListTodos tests ---

func TestListTodos_Empty(t *testing.T) {
	h, _ := newTestHandler()
	w := httptest.NewRecorder()
	r := makeRequest(http.MethodGet, "")

	h.ListTodos(w, r, "default")

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var list todov1.TodoList
	decodeResponse(t, w, &list)
	if len(list.Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(list.Items))
	}
}

func TestListTodos_WithItems(t *testing.T) {
	h, repo := newTestHandler()
	repo.addTodo("default", "todo-1", "First", todov1.SpecStatusOpen)
	repo.addTodo("default", "todo-2", "Second", todov1.SpecStatusDone)
	w := httptest.NewRecorder()
	r := makeRequest(http.MethodGet, "")

	h.ListTodos(w, r, "default")

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var list todov1.TodoList
	decodeResponse(t, w, &list)
	if len(list.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(list.Items))
	}
}

func TestListTodos_NamespaceIsolation(t *testing.T) {
	h, repo := newTestHandler()
	repo.addTodo("ns-a", "todo-1", "NS A Todo", todov1.SpecStatusOpen)
	repo.addTodo("ns-b", "todo-2", "NS B Todo", todov1.SpecStatusOpen)
	w := httptest.NewRecorder()
	r := makeRequest(http.MethodGet, "")

	h.ListTodos(w, r, "ns-a")

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var list todov1.TodoList
	decodeResponse(t, w, &list)
	if len(list.Items) != 1 {
		t.Errorf("expected 1 item for ns-a, got %d", len(list.Items))
	}
}

func TestListTodos_RepoError(t *testing.T) {
	h, repo := newTestHandler()
	repo.lastErr = fmt.Errorf("database error")
	w := httptest.NewRecorder()
	r := makeRequest(http.MethodGet, "")

	h.ListTodos(w, r, "default")

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

// --- UpdateTodo tests ---

func TestUpdateTodo_Success(t *testing.T) {
	h, repo := newTestHandler()
	repo.addTodo("default", "todo-1", "Old Title", todov1.SpecStatusOpen)
	w := httptest.NewRecorder()
	body := `{"title": "New Title", "status": "in_progress"}`
	r := makeRequest(http.MethodPut, body)

	h.UpdateTodo(w, r, "default", "todo-1")

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var todo todov1.Todo
	decodeResponse(t, w, &todo)
	if todo.Spec.Title != "New Title" {
		t.Errorf("expected title 'New Title', got %q", todo.Spec.Title)
	}
	if todo.Spec.Status != todov1.SpecStatusInProgress {
		t.Errorf("expected status 'in_progress', got %q", todo.Spec.Status)
	}
}

func TestUpdateTodo_NotFound(t *testing.T) {
	h, _ := newTestHandler()
	w := httptest.NewRecorder()
	body := `{"title": "Test", "status": "open"}`
	r := makeRequest(http.MethodPut, body)

	h.UpdateTodo(w, r, "default", "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestUpdateTodo_InvalidStatus(t *testing.T) {
	h, repo := newTestHandler()
	repo.addTodo("default", "todo-1", "Test", todov1.SpecStatusOpen)
	w := httptest.NewRecorder()
	body := `{"title": "Test", "status": "invalid"}`
	r := makeRequest(http.MethodPut, body)

	h.UpdateTodo(w, r, "default", "todo-1")

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d", w.Code)
	}
}

func TestUpdateTodo_MissingTitle(t *testing.T) {
	h, repo := newTestHandler()
	repo.addTodo("default", "todo-1", "Test", todov1.SpecStatusOpen)
	w := httptest.NewRecorder()
	body := `{"status": "open"}`
	r := makeRequest(http.MethodPut, body)

	h.UpdateTodo(w, r, "default", "todo-1")

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d", w.Code)
	}
}

func TestUpdateTodo_GetInternalError(t *testing.T) {
	h, repo := newTestHandler()
	repo.getErr = fmt.Errorf("connection timeout")
	w := httptest.NewRecorder()
	body := `{"title": "Test", "status": "open"}`
	r := makeRequest(http.MethodPut, body)

	h.UpdateTodo(w, r, "default", "todo-1")

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestUpdateTodo_RepoUpdateError(t *testing.T) {
	h, repo := newTestHandler()
	repo.addTodo("default", "todo-1", "Test", todov1.SpecStatusOpen)
	repo.updateErr = fmt.Errorf("conflict: resource version changed")
	w := httptest.NewRecorder()
	body := `{"title": "Updated", "status": "done"}`
	r := makeRequest(http.MethodPut, body)

	h.UpdateTodo(w, r, "default", "todo-1")

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestUpdateTodo_InvalidJSON(t *testing.T) {
	h, repo := newTestHandler()
	repo.addTodo("default", "todo-1", "Test", todov1.SpecStatusOpen)
	w := httptest.NewRecorder()
	body := `{broken`
	r := makeRequest(http.MethodPut, body)

	h.UpdateTodo(w, r, "default", "todo-1")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

// --- DeleteTodo tests ---

func TestDeleteTodo_Success(t *testing.T) {
	h, repo := newTestHandler()
	repo.addTodo("default", "todo-1", "Test", todov1.SpecStatusOpen)
	w := httptest.NewRecorder()
	r := makeRequest(http.MethodDelete, "")

	h.DeleteTodo(w, r, "default", "todo-1")

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestDeleteTodo_NotFound(t *testing.T) {
	h, _ := newTestHandler()
	w := httptest.NewRecorder()
	r := makeRequest(http.MethodDelete, "")

	h.DeleteTodo(w, r, "default", "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestDeleteTodo_RepoDeleteError(t *testing.T) {
	h, repo := newTestHandler()
	repo.addTodo("default", "todo-1", "Test", todov1.SpecStatusOpen)
	repo.deleteErr = fmt.Errorf("connection timeout")
	w := httptest.NewRecorder()
	r := makeRequest(http.MethodDelete, "")

	h.DeleteTodo(w, r, "default", "todo-1")

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestDeleteTodo_GetInternalError(t *testing.T) {
	h, repo := newTestHandler()
	repo.getErr = fmt.Errorf("connection timeout")
	w := httptest.NewRecorder()
	r := makeRequest(http.MethodDelete, "")

	h.DeleteTodo(w, r, "default", "todo-1")

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

// --- validateTodoRequest tests ---

func TestValidateTodoRequest_ValidOpen(t *testing.T) {
	req := &TodoRequest{Title: "Test", Status: "open"}
	if msg := validateTodoRequest(req); msg != "" {
		t.Errorf("expected no error, got %q", msg)
	}
}

func TestValidateTodoRequest_ValidInProgress(t *testing.T) {
	req := &TodoRequest{Title: "Test", Status: "in_progress"}
	if msg := validateTodoRequest(req); msg != "" {
		t.Errorf("expected no error, got %q", msg)
	}
}

func TestValidateTodoRequest_ValidDone(t *testing.T) {
	req := &TodoRequest{Title: "Test", Status: "done"}
	if msg := validateTodoRequest(req); msg != "" {
		t.Errorf("expected no error, got %q", msg)
	}
}

func TestValidateTodoRequest_EmptyTitle(t *testing.T) {
	req := &TodoRequest{Title: "", Status: "open"}
	if msg := validateTodoRequest(req); msg == "" {
		t.Error("expected validation error for empty title")
	}
}

func TestValidateTodoRequest_WhitespaceTitle(t *testing.T) {
	req := &TodoRequest{Title: "   ", Status: "open"}
	if msg := validateTodoRequest(req); msg == "" {
		t.Error("expected validation error for whitespace-only title")
	}
}

func TestValidateTodoRequest_InvalidStatus(t *testing.T) {
	req := &TodoRequest{Title: "Test", Status: "unknown"}
	if msg := validateTodoRequest(req); msg == "" {
		t.Error("expected validation error for invalid status")
	}
}

func TestValidateTodoRequest_EmptyStatus(t *testing.T) {
	req := &TodoRequest{Title: "Test", Status: ""}
	if msg := validateTodoRequest(req); msg == "" {
		t.Error("expected validation error for empty status")
	}
}

// --- isNotFound tests ---

func TestIsNotFound_True(t *testing.T) {
	cases := []string{
		"todo default/x not found",
		"resource NotFound",
		"HTTP 404 error",
	}
	for _, msg := range cases {
		if !isNotFound(fmt.Errorf("%s", msg)) {
			t.Errorf("expected isNotFound=true for %q", msg)
		}
	}
}

func TestIsNotFound_False(t *testing.T) {
	if isNotFound(fmt.Errorf("connection refused")) {
		t.Error("expected isNotFound=false for 'connection refused'")
	}
}
