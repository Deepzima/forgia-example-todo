---
id: "SDD-003"
fd: "FD-002"
title: "UI components — priority selector, badges, and filtering"
status: done
agent: "claude-code"
assigned_to: "claude-code"
created: "2026-03-15"
started: ""
completed: ""
tags: ["ui", "react", "typescript", "grafana"]
---

# SDD-003: UI components — priority selector, badges, and filtering

> Parent FD: [[FD-002]]

## Scope

Add priority support to the Grafana plugin frontend:
1. **TodoForm**: Add a `Select` dropdown for priority (low/medium/high/critical), defaulting to "medium"
2. **TodoList**: Display a color-coded `Badge` for each Todo's priority level
3. **TodoPage**: Add sort-by-priority and filter-by-priority controls

This SDD depends on **SDD-001** (generated TypeScript types) and **SDD-002** (API returns priority in responses).

## Interfaces / Interfacce

| Interface / Interfaccia | Type / Tipo | Description / Descrizione |
|-------------------------|-------------|---------------------------|
| `Spec.priority` (TS) | `"low" \| "medium" \| "high" \| "critical"` | From generated types (SDD-001). Used for form state and display. |
| Priority Select options | `SelectableValue<string>[]` | `[{label: "Low", value: "low"}, {label: "Medium", value: "medium"}, {label: "High", value: "high"}, {label: "Critical", value: "critical"}]` |
| Priority Badge color map | `Record<string, BadgeColor>` | `{low: "blue", medium: "yellow", high: "orange", critical: "red"}` — uses Grafana `Badge` component color prop |
| `todoApi.create/update` | Request body includes `priority` | Already handled by existing API client since it passes full spec |
| Sort order | `critical > high > medium > low` | Numeric mapping for sort: `{critical: 4, high: 3, medium: 2, low: 1}` |

**Contract from SDD-001:** TypeScript `Spec` interface includes optional `priority` field.

**Contract from SDD-002:** All API responses include `priority` in `spec`. Create without priority defaults to `"medium"`.

## Constraints / Vincoli

- Language / Linguaggio: TypeScript 5.6.0 (strict mode)
- Framework: React 18.3.0, Grafana UI components (`@grafana/ui ^11.0.0`)
- Dependencies / Dipendenze: No new dependencies. Use existing `@grafana/ui` components: `Select`, `Badge`, `RadioButtonGroup` (for filter).
- Patterns / Pattern: Follow existing component patterns in `TodoForm.tsx` and `TodoList.tsx`. Use `const` over `let`. Explicit return types on exported functions.

### Guardrails (from deny.toml)

- NEVER read/write `.env`, `*.pem`, `*.key`, credentials files
- NEVER modify `.forgia/constitution.md`, `.forgia/config.toml`, `.forgia/guardrails/deny.toml`
- NEVER run commands that enumerate secrets or read private keys

## Best Practices

- Error handling: Priority field should always be present in responses (SDD-002 guarantees this). If somehow missing, default to `"medium"` in display — but log a console warning (no silent fallback).
- Naming: `priorityOptions` for Select options array. `priorityColorMap` for badge colors. `sortByPriority` for sort function. Use `const` for all.
- Style: Keep color mapping as a simple object — no need for a class or factory. Place shared constants (options, colors, sort order) in a `priorityUtils.ts` file if they are used by more than one component; otherwise keep them co-located. Follow existing component structure.

## Test Requirements

| Type / Tipo | What / Cosa | Coverage |
|-------------|-------------|----------|
| Unit | TodoForm renders priority Select with 4 options and default "medium" | Render + interaction |
| Unit | TodoForm submits with selected priority value | Form submission |
| Unit | TodoList renders correct Badge color for each priority level | All 4 values |
| Unit | Sort by priority orders correctly: critical > high > medium > low | Sort function |
| Unit | Filter by priority shows only matching todos | Filter function |

## Acceptance Criteria / Criteri di Accettazione

- [x] TodoForm shows a "Priority" dropdown with options: Low, Medium, High, Critical
- [x] TodoForm defaults to "Medium" when creating a new Todo
- [x] TodoForm pre-selects the current priority when editing an existing Todo
- [x] TodoList shows a colored Badge next to each Todo: low=blue, medium=yellow, high=orange, critical=red
- [x] TodoPage has a sort control that can sort by priority (critical first or low first)
- [x] TodoPage has a filter control to show only specific priority levels
- [x] All existing tests still pass (no regressions)
- [x] New tests cover form interaction, badge rendering, sort, and filter

## Context / Contesto

- [x] `plugin/src/components/TodoForm.tsx` — existing form with title, description, status fields
- [x] `plugin/src/components/TodoList.tsx` — existing list with status Badge pattern to follow
- [x] `plugin/src/pages/TodoPage.tsx` — main page with view mode state management
- [x] `plugin/src/hooks/useTodos.ts` — state hook (may need sort/filter additions)
- [x] `plugin/src/api/todoApi.ts` — API client (should work without changes since it passes full spec)
- [x] `plugin/src/generated/todo/v1/` — generated TypeScript types (from SDD-001)

## Constitution Check

- [x] Respects code standards: TypeScript strict mode, `const` over `let`, explicit return types
- [x] Respects commit conventions: `feat(FD-002/SDD-003): add priority UI components`
- [x] No hardcoded secrets
- [x] Tests defined and sufficient: unit tests for form, badges, sort, and filter

---

## Work Log / Diario di Lavoro

> This section is **mandatory**. Must be filled by the agent or developer during and after execution.

### Agent / Agente

- **Executor**: claude-code
- **Started**: 2026-03-15
- **Completed**: 2026-03-15
- **Duration / Durata**: ~15 min

### Decisions / Decisioni

1. Created a shared `priorityUtils.ts` module for constants (options, color map, sort weights) and pure functions (`sortByPriority`, `filterByPriority`, `getPriority`) since they are used by multiple components (TodoForm, TodoList, TodoPage).
2. Used `RadioButtonGroup` for sort control (3 states: no sort, critical-first, low-first) and `Select` for priority filter — keeps controls compact and consistent with existing Grafana UI patterns.
3. `getPriority()` logs a `console.warn` when priority is missing rather than silently falling back — follows the "no silent fallbacks" constitution rule while still being resilient.
4. Filter uses a single-select dropdown (one priority at a time or "All") to keep the UI simple; the underlying `filterByPriority` function supports multi-select for future extensibility.

### Output

- **Commit(s)**: b8b3426
- **PR**: —
- **Files created/modified**:
  - `plugin/src/components/priorityUtils.ts` (new — shared constants and pure functions)
  - `plugin/src/components/priorityUtils.test.ts` (new — unit tests for sort, filter, getPriority)
  - `plugin/src/components/TodoForm.tsx` (modified — added priority Select dropdown)
  - `plugin/src/components/TodoForm.test.tsx` (modified — added priority tests)
  - `plugin/src/components/TodoList.tsx` (modified — added priority Badge)
  - `plugin/src/components/TodoList.test.tsx` (modified — added priority badge color tests)
  - `plugin/src/pages/TodoPage.tsx` (modified — added sort and filter controls)
  - `plugin/src/__mocks__/@grafana/ui.tsx` (modified — added RadioButtonGroup mock)

### Retrospective / Retrospettiva

- **What worked / Cosa ha funzionato**: Existing component patterns (Select, Badge) made it straightforward to add priority support. The mock-based test setup allowed fast iteration. All 47 tests pass with no regressions.
- **What didn't / Cosa non ha funzionato**: Nothing significant — the generated Spec type already had the optional priority field from SDD-001.
- **Suggestions for future FDs / Suggerimenti per FD futuri**: Consider specifying whether filter controls should be single-select or multi-select in the SDD to avoid implementation ambiguity.
