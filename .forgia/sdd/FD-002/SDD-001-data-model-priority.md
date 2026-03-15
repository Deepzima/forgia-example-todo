---
id: "SDD-001"
fd: "FD-002"
title: "Data model & persistence — priority field"
status: planned
agent: ""
assigned_to: ""
created: "2026-03-15"
started: ""
completed: ""
tags: ["schema", "cue", "codegen"]
---

# SDD-001: Data model & persistence — priority field

> Parent FD: [[FD-002]]

## Scope

Add a `priority` enum field to the Todo CUE schema (`todo_v1.cue`) with values `"low" | "medium" | "high" | "critical"`. The field must be optional with a default of `"medium"` to ensure backward compatibility with existing Todo resources in Kubernetes. After updating the schema, regenerate the Go and TypeScript types, and update the CRD JSON definition.

This SDD touches only the schema layer — no handler or UI changes.

## Interfaces / Interfacce

| Interface / Interfaccia | Type / Tipo | Description / Descrizione |
|-------------------------|-------------|---------------------------|
| `Spec.priority` (CUE) | `*"medium" \| "low" \| "high" \| "critical"` | New optional enum field with default `"medium"`. Follows the same pattern as `Spec.status`. |
| `SpecPriority` (Go) | `type SpecPriority string` | Generated Go string type with constants: `SpecPriorityLow`, `SpecPriorityMedium`, `SpecPriorityHigh`, `SpecPriorityCritical` |
| `Spec.priority` (TS) | `"low" \| "medium" \| "high" \| "critical"` | Generated TypeScript union type in the `Spec` interface |
| CRD JSON schema | OpenAPI v3 property | `priority` field in `spec` with `enum` constraint and `default: "medium"` |

**Contract with SDD-002 (API layer):** The generated Go `Spec` struct will include a `Priority` field of type `SpecPriority`. SDD-002 will use this type for validation and serialization — no manual parsing needed.

**Contract with SDD-003 (UI):** The generated TypeScript `Spec` interface will include `priority` as an optional field. SDD-003 will use this type for form binding and display.

## Constraints / Vincoli

- Language / Linguaggio: CUE (schema), Go (generated types), TypeScript (generated types)
- Framework: Grafana App SDK (`github.com/grafana/grafana-app-sdk v0.52.0`), Kubernetes CRD
- Dependencies / Dipendenze: No new dependencies. Use existing code generation tooling.
- Patterns / Pattern: Follow the exact pattern used by `status: "open" | "in_progress" | "done"` — same CUE syntax, same generated type structure.

### Guardrails (from deny.toml)

- NEVER read/write `.env`, `*.pem`, `*.key`, credentials files
- NEVER modify `.forgia/constitution.md`, `.forgia/config.toml`, `.forgia/guardrails/deny.toml`
- NEVER run commands that enumerate secrets or read private keys

## Best Practices

- Error handling: CUE schema validation is declarative — invalid values are rejected at the schema level. The generated Go types must use the enum constants, never raw strings.
- Naming: Use `priority` (lowercase) in CUE/JSON. Use `Priority` (PascalCase) in Go struct. Use `SpecPriority` for the Go type name, following the `SpecStatus` pattern.
- Style: Keep the CUE definition adjacent to the `status` field for readability. The `*"medium"` syntax marks the default value in CUE.

## Test Requirements

| Type / Tipo | What / Cosa | Coverage |
|-------------|-------------|----------|
| Unit | CUE schema validation: valid values accepted, invalid values rejected | 100% of enum values |
| Unit | Go generated type: verify constants exist and match expected string values | All 4 priority constants |
| Unit | TypeScript generated type: verify the interface includes `priority` field | Type compilation check |

## Acceptance Criteria / Criteri di Accettazione

- [ ] `kinds/todo_v1.cue` contains `priority: *"medium" | "low" | "high" | "critical"` in the spec
- [ ] `pkg/generated/todo/v1/` Go types include `Priority SpecPriority` field in `Spec` struct
- [ ] `pkg/generated/todo/v1/` Go types include `SpecPriority` type with 4 constants
- [ ] `plugin/src/generated/todo/v1/` TypeScript types include `priority` in `Spec` interface
- [ ] `definitions/todo.todo.grafana.app.json` CRD includes `priority` property with enum constraint
- [ ] Existing Todo resources without `priority` field are valid (optional field with default)
- [ ] Code generation runs without errors

## Context / Contesto

- [ ] `kinds/todo_v1.cue` — current CUE schema with `status` enum pattern to follow
- [ ] `pkg/generated/todo/v1/` — current generated Go types (understand structure before regenerating)
- [ ] `plugin/src/generated/todo/v1/` — current generated TypeScript types
- [ ] `definitions/todo.todo.grafana.app.json` — current CRD JSON definition
- [ ] `go.mod` — verify Grafana App SDK version for codegen compatibility

## Constitution Check

- [ ] Respects code standards: enum field follows existing pattern, no silent fallbacks
- [ ] Respects commit conventions: `feat(FD-002/SDD-001): add priority field to Todo schema`
- [ ] No hardcoded secrets
- [ ] Tests defined and sufficient: schema validation + type verification

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
