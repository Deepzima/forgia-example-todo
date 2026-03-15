---
id: "SDD-001"
fd: "FD-001"
title: "TODO Custom Resource Definition (CRD schema, grafana-app-sdk kind)"
status: done
agent: "claude-code"
assigned_to: "claude-code"
created: "2026-03-15"
started: "2026-03-15"
completed: "2026-03-15"
tags: [grafana-app-sdk, crd, kubernetes, go]
---

# SDD-001: TODO Custom Resource Definition

> Parent FD: [[FD-001]]

## Scope

Definire la Custom Resource `Todo` utilizzando `grafana-app-sdk`. Questo include:

- Definizione del kind `Todo` con `grafana-app-sdk` (Go types + CUE schema)
- Campi spec: `title` (string, required), `description` (string, optional), `status` (enum: `open` | `in_progress` | `done`)
- Campi metadata gestiti automaticamente: `uid`, `createdAt`, `updatedAt`
- Validazione OpenAPI dello schema CRD
- Generazione del codice Go tramite `grafana-app-sdk` codegen

Il CRD risultante deve essere installabile su un cluster Kubernetes e deve esporre l'API group `todo.grafana.app/v1` con kind `Todo`.

## Interfaces / Interfacce

| Interface / Interfaccia | Type / Tipo | Description / Descrizione |
|-------------------------|-------------|---------------------------|
| `Todo` Kind definition | Go struct + CUE schema | Definizione del tipo `Todo` per grafana-app-sdk, usato da SDD-002 per i CRUD handlers |
| CRD YAML | Kubernetes CRD manifest | `apiVersion: apiextensions.k8s.io/v1`, `kind: CustomResourceDefinition`, group `todo.grafana.app` |
| Todo spec fields | Go struct | `Title string`, `Description string`, `Status TodoStatus` (enum) |
| Todo response shape | JSON | `{uid, title, description, status, createdAt, updatedAt}` — consumato da SDD-003 (frontend) |

### Contract con SDD-002 (Backend/Operator)

Il backend importera' i tipi Go generati da questo SDD per registrare i CRUD handlers. L'interfaccia e':

```go
// Kind registration — SDD-002 usa questo per registrare l'operator
type Todo struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`
    Spec              TodoSpec   `json:"spec"`
    Status            TodoStatus `json:"status,omitempty"`
}

type TodoSpec struct {
    Title       string `json:"title"`
    Description string `json:"description,omitempty"`
    Status      string `json:"status"` // "open" | "in_progress" | "done"
}
```

### Contract con SDD-003 (Frontend)

Il frontend consumera' le risorse Todo tramite REST API. Il formato JSON atteso e':

```json
{
  "metadata": {
    "uid": "uuid",
    "name": "todo-name",
    "namespace": "default",
    "creationTimestamp": "2026-03-15T00:00:00Z"
  },
  "spec": {
    "title": "My TODO",
    "description": "Description text",
    "status": "open"
  }
}
```

## Constraints / Vincoli

- Language / Linguaggio: Go
- Framework: `grafana-app-sdk` (latest stable)
- Dependencies / Dipendenze: `k8s.io/apimachinery`, `grafana-app-sdk`
- Patterns / Pattern: grafana-app-sdk kind definition pattern (CUE + Go codegen)
- Il campo `status` DEVE accettare solo: `open`, `in_progress`, `done` — validato tramite schema OpenAPI (`enum`)
- Nessun segreto hardcoded

### Guardrails (deny.toml)

L'agent NON deve:
- Leggere file `.env`, `*.pem`, `*.key`, `kubeconfig`, `credentials.json`
- Scrivere file `.env`, `.forgia/constitution.md`, `.forgia/config.toml`, `.forgia/guardrails/deny.toml`
- Eseguire comandi che leggono chiavi SSH, GPG, o enumerano variabili d'ambiente con segreti

## Best Practices

- Error handling: usare `fmt.Errorf` con wrapping (`%w`) per errori espliciti, mai fallback silenti
- Naming: Go conventions — `TodoSpec`, `TodoStatus`, `TodoList`, exported types in PascalCase
- Style: `gofmt` + `golint`, commenti in inglese nel codice
- Validazione: usare marker annotations di kubebuilder (`+kubebuilder:validation:Enum`) o equivalenti grafana-app-sdk per la validazione dello schema

## Test Requirements

| Type / Tipo | What / Cosa | Coverage |
|-------------|-------------|----------|
| Unit | Validazione struct TodoSpec (campi obbligatori, enum status) | 100% sui campi spec |
| Integration | CRD installabile su cluster K8s (kind, minikube, o envtest) | CRD apply + get |
| E2E | — | Coperto da SDD-002 |

## Acceptance Criteria / Criteri di Accettazione

- [ ] Il CRD `Todo` (apiVersion: `todo.grafana.app/v1`) e' installabile su un cluster Kubernetes
- [ ] Lo schema OpenAPI del CRD valida correttamente i campi `title` (required), `description` (optional), `status` (enum)
- [ ] Il campo `status` accetta solo i valori `open`, `in_progress`, `done` — valori non validi vengono rifiutati
- [ ] I tipi Go sono generati tramite `grafana-app-sdk` codegen e compilano senza errori
- [ ] I test unitari passano con `go test ./...`
- [ ] Nessun segreto hardcoded nel codice o nei manifest

## Context / Contesto

- [ ] Documentazione `grafana-app-sdk`: https://github.com/grafana/grafana-app-sdk
- [ ] Esempio kind definition in grafana-app-sdk
- [ ] CUE schema language reference
- [ ] Kubernetes CRD documentation: apiextensions.k8s.io/v1
- [ ] File FD: `.forgia/fd/FD-001-todo-grafana-app.md`

## Constitution Check

- [ ] Rispetta le code standards (Go conventions, explicit errors, no silent fallbacks)
- [ ] Rispetta le commit conventions (`feat(FD-001/SDD-001): ...` + `Co-Authored-By`)
- [ ] Nessun secret hardcoded — env vars o secret manager
- [ ] Test definiti e sufficienti (unit + integration)

---

## Work Log / Diario di Lavoro

> Questa sezione e' **obbligatoria**. Deve essere compilata dall'agent o dallo sviluppatore durante e dopo l'esecuzione.

### Agent / Agente

- **Executor**: claude-code
- **Started**: 2026-03-15
- **Completed**: 2026-03-15
- **Duration / Durata**: ~15 min

### Decisions / Decisioni

1. CUE schema gia' definito in `kinds/todo.cue` e `kinds/todo_v1.cue` con i campi `title` (required), `description` (optional), `status` (enum). Codegen grafana-app-sdk gia' eseguita e funzionante.
2. Aggiunto test di validazione dello schema CRD (`todo_crd_schema_test.go`) che verifica la struttura OpenAPI del CRD JSON senza necessita' di un cluster Kubernetes — controlla i campi required, i tipi, e i vincoli enum.
3. I test di integrazione (`tests/integration/crd_test.go`) usano il build tag `integration` e richiedono un cluster K8s; verificano install CRD, creazione risorse, e rifiuto di valori status non validi.

### Output

- **Commit(s)**: pending
- **PR**: —
- **Files created/modified**:
  - `kinds/todo.cue` — kind definition (scope, pluralName, codegen flags)
  - `kinds/todo_v1.cue` — v1 schema (spec: title, description?, status enum)
  - `kinds/config.cue` — codegen configuration (paths Go/TS)
  - `kinds/manifest.cue` — manifest config (appName, group override)
  - `pkg/generated/todo/v1/` — Go types generati (Todo, Spec, Status, Client, Codec, Schema)
  - `pkg/generated/manifestdata/` — manifest data generato
  - `definitions/todo.todo.grafana.app.json` — CRD manifest generato
  - `definitions/todo-manifest.json` — app manifest generato
  - `plugin/src/generated/todo/v1/` — TypeScript types generati
  - `pkg/generated/todo/v1/todo_spec_test.go` — unit test per Spec (marshal, unmarshal, enum, defaults)
  - `pkg/generated/todo/v1/todo_object_test.go` — unit test per Todo object (metadata, deep copy, schema)
  - `pkg/generated/todo/v1/todo_crd_schema_test.go` — unit test validazione schema CRD/OpenAPI
  - `tests/integration/crd_test.go` — integration test (CRD install, create resource, reject invalid status)

### Retrospective / Retrospettiva

- **Cosa ha funzionato**: grafana-app-sdk codegen ha prodotto tipi Go e TS corretti dal CUE schema. Il CRD generato include validazione OpenAPI con enum constraint per status. La struttura del progetto e' pulita e segue le convenzioni grafana-app-sdk.
- **Cosa non ha funzionato**: nulla di significativo — il progetto era gia' ben impostato.
- **Suggerimenti per FD futuri**: includere un test di validazione dello schema CRD come unit test (senza cluster) per verificare le invarianti OpenAPI in CI senza dipendenze esterne.
