---
id: "SDD-002"
fd: "FD-002"
title: "API layer — priority in CRUD operations"
status: planned
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

- [ ] `POST /todos` with `priority: "high"` creates a Todo with priority "high"
- [ ] `POST /todos` without `priority` field creates a Todo with priority "medium" (default)
- [ ] `PUT /todos/{name}` can change priority from one value to another
- [ ] `GET /todos/{name}` always returns `priority` in the `spec` object
- [ ] `GET /todos` (list) returns `priority` on every item
- [ ] `POST /todos` with `priority: "invalid"` returns HTTP 400 with explicit error message
- [ ] Existing Todos in K8s without a `priority` field return `"medium"` when read via API
- [ ] All tests pass

## Context / Contesto

- [ ] `pkg/handler/todo_handler.go` — existing handler with Create/Update/Get/List/Delete methods
- [ ] `pkg/repository/todo_repository.go` — repository interface (understand how data flows)
- [ ] `pkg/generated/todo/v1/` — generated types (from SDD-001) with `SpecPriority`
- [ ] `pkg/handler/todo_handler_test.go` — existing test patterns to follow (if present)
- [ ] `cmd/operator/main.go` — understand how handler is wired

## Constitution Check

- [ ] Respects code standards: explicit error handling, no silent fallbacks, validates at boundaries
- [ ] Respects commit conventions: `feat(FD-002/SDD-002): add priority to TodoHandler CRUD operations`
- [ ] No hardcoded secrets
- [ ] Tests defined and sufficient: unit + integration covering all CRUD operations

---

## Work Log / Diario di Lavoro

> This section is **mandatory**. Must be filled by the agent or developer during and after execution.

### Agent / Agente

- **Executor**: <!-- openhands | claude-code | manual | name -->
- **Started**: <!-- timestamp -->
- **Completed**: <!-- timestamp -->
- **Duration / Durata**: <!-- total time -->

### Decisions / Decisioni

1. <!-- decision 1: what and why -->

### Output

- **Commit(s)**: <!-- hash -->
- **PR**: <!-- link -->
- **Files created/modified**:
  - `path/to/file`

### Retrospective / Retrospettiva

- **What worked / Cosa ha funzionato**:
- **What didn't / Cosa non ha funzionato**:
- **Suggestions for future FDs / Suggerimenti per FD futuri**:
