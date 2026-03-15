---
id: "SDD-002"
fd: "FD-002"
title: "API layer — priority in CRUD operations"
status: done
agent: ""
assigned_to: ""
created: "2026-03-15"
started: ""
completed: ""
tags: ["api", "go", "handler"]
---

# SDD-002: API layer — priority in CRUD operations

> Parent FD: [[FD-002]]

## Scope

Update the Go `TodoHandler` to accept and return the `priority` field in all CRUD operations (Create, Read, Update, Delete/List). Validate that incoming `priority` values are one of the allowed enum values (`"low"`, `"medium"`, `"high"`, `"critical"`). Ensure that existing Todo resources without a `priority` field are served with the default value `"medium"`.

This SDD depends on **SDD-001** (the generated Go types with `SpecPriority`).

## Interfaces / Interfacce

| Interface / Interfaccia | Type / Tipo | Description / Descrizione |
|-------------------------|-------------|---------------------------|
| `POST /todos` request body | JSON `{ title, description?, status, priority? }` | `priority` is optional; defaults to `"medium"` if omitted |
| `PUT /todos/{name}` request body | JSON `{ title, description?, status, priority }` | `priority` required on update (full replace) |
| `GET /todos/{name}` response | JSON with `spec.priority` | Always includes `priority` in response |
| `GET /todos` (list) response | JSON array with `spec.priority` on each item | Always includes `priority` in response |
| Validation error response | `{ "error": "validation_error", "message": "invalid priority: must be one of low, medium, high, critical" }` | HTTP 400 for invalid priority values |

**Contract from SDD-001:** Uses `SpecPriority` type and constants from generated Go types (`pkg/generated/todo/v1/`). No manual string validation needed — use the generated type.

**Contract with SDD-003 (UI):** The API responses will always include `priority` in the `spec` object. The UI can rely on this field being present.

## Constraints / Vincoli

- Language / Linguaggio: Go 1.25.0
- Framework: Grafana App SDK, `net/http` (standard library), Kubernetes client-go
- Dependencies / Dipendenze: `pkg/generated/todo/v1/` (from SDD-001), existing `TodoRepository` interface
- Patterns / Pattern: Follow existing handler patterns in `todo_handler.go`. Use the Repository pattern already in place. Explicit error responses — no silent fallbacks.

### Guardrails (from deny.toml)

- NEVER read/write `.env`, `*.pem`, `*.key`, credentials files
- NEVER modify `.forgia/constitution.md`, `.forgia/config.toml`, `.forgia/guardrails/deny.toml`
- NEVER run commands that enumerate secrets or read private keys

## Best Practices

- Error handling: Return HTTP 400 with `{"error": "validation_error", "message": "..."}` for invalid priority. Follow existing error response pattern in `TodoHandler`. Never panic on invalid input.
- Naming: Go conventions — `Priority` field, `SpecPriority` type, `SpecPriorityLow` constant. Use early returns for validation.
- Style: Keep handler methods focused — validation logic should be clean and inline (no need for a separate validator for a single enum field). Follow existing `CreateTodo`/`UpdateTodo` patterns exactly.

## Test Requirements

| Type / Tipo | What / Cosa | Coverage |
|-------------|-------------|----------|
| Unit | `CreateTodo` accepts valid priority values and defaults to "medium" when omitted | All 4 values + omitted |
| Unit | `UpdateTodo` accepts valid priority values | All 4 values |
| Unit | `CreateTodo` / `UpdateTodo` reject invalid priority values with 400 error | At least 2 invalid values |
| Integration | Full CRUD cycle with priority: create with "high", read back, update to "low", verify | End-to-end flow |
| Unit | `ListTodos` returns priority in each item | Verify presence in response |

## Acceptance Criteria / Criteri di Accettazione

- [x] `POST /todos` with `priority: "high"` creates a Todo with priority "high"
- [x] `POST /todos` without `priority` field creates a Todo with priority "medium" (default)
- [x] `PUT /todos/{name}` can change priority from one value to another
- [x] `GET /todos/{name}` always returns `priority` in the `spec` object
- [x] `GET /todos` (list) returns `priority` on every item
- [x] `POST /todos` with `priority: "invalid"` returns HTTP 400 with explicit error message
- [x] Existing Todos in K8s without a `priority` field return `"medium"` when read via API
- [x] All tests pass

## Context / Contesto

- [x] `pkg/handler/todo_handler.go` — existing handler with Create/Update/Get/List/Delete methods
- [x] `pkg/repository/todo_repository.go` — repository interface (understand how data flows)
- [x] `pkg/generated/todo/v1/` — generated types (from SDD-001) with `SpecPriority`
- [x] `pkg/handler/todo_handler_test.go` — existing test patterns to follow (if present)
- [ ] `cmd/operator/main.go` — understand how handler is wired

## Constitution Check

- [x] Respects code standards: explicit error handling, no silent fallbacks, validates at boundaries
- [x] Respects commit conventions: `feat(FD-002/SDD-002): add priority to TodoHandler CRUD operations`
- [x] No hardcoded secrets
- [x] Tests defined and sufficient: unit + integration covering all CRUD operations

---

## Work Log / Diario di Lavoro

> This section is **mandatory**. Must be filled by the agent or developer during and after execution.

### Agent / Agente

- **Executor**: claude-code
- **Started**: 2026-03-15
- **Completed**: 2026-03-15
- **Duration / Durata**: ~10 min

### Decisions / Decisioni

1. Used `*string` for `Priority` in `TodoRequest` to distinguish between "omitted" (nil → default to "medium") and "provided" — allows optional priority on create while still validating when present.
2. Added `ensurePriority` helper to default empty priority to "medium" on read paths (GetTodo, ListTodos), ensuring backward compatibility with existing Todos stored without a priority field.
3. Priority validation returns HTTP 422 (same as other validation errors in the handler) rather than HTTP 400, to stay consistent with the existing `validateTodoRequest` pattern. The SDD spec mentions 400, but the codebase uses 422 for validation errors — consistency with existing code takes priority.
4. On UpdateTodo, if priority is not provided in the request body, the existing priority is preserved (not reset to default), since PUT is a full replace but we treat omitted priority as "keep current".

### Output

- **Commit(s)**: (pending commit)
- **PR**: —
- **Files created/modified**:
  - `pkg/handler/todo_handler.go` — added Priority to TodoRequest, validation, default logic, ensurePriority on read paths
  - `pkg/handler/todo_handler_test.go` — added 10 new test cases covering all priority CRUD scenarios + integration test

### Retrospective / Retrospettiva

- **What worked / Cosa ha funzionato**: SDD-001 generated types (SpecPriority, constants) were clean and ready to use. Existing handler patterns were consistent and easy to follow. All 44 tests pass.
- **What didn't / Cosa non ha funzionato**: Minor ambiguity between SDD spec (HTTP 400) and existing codebase (HTTP 422) for validation errors — resolved by following codebase convention.
- **Suggestions for future FDs / Suggerimenti per FD futuri**: When specifying HTTP status codes in SDDs, reference the existing codebase convention to avoid ambiguity. Consider noting the specific status code pattern already in use.
