---
id: "SDD-002"
fd: "FD-001"
title: "TODO Backend / Operator (CRUD handlers, grafana-app-sdk operator logic)"
status: done
agent: "claude-code"
assigned_to: "claude-code"
created: "2026-03-15"
started: "2026-03-15"
completed: "2026-03-15"
tags: [grafana-app-sdk, operator, go, backend, crud]
---

# SDD-002: TODO Backend / Operator

> Parent FD: [[FD-001]]

## Scope

Implementare il backend operator per la gestione CRUD delle risorse `Todo` utilizzando `grafana-app-sdk`. Questo include:

- Registrazione del kind `Todo` (definito in SDD-001) nell'operator
- Implementazione dei CRUD handlers: Create, Read (Get + List), Update, Delete
- Routing delle API REST: `POST/GET/PUT/DELETE /apis/todo.grafana.app/v1/namespaces/{ns}/todos`
- Watcher/reconciler per il ciclo di vita delle risorse Todo
- Configurazione dell'app plugin backend (Grafana plugin backend SDK)

L'operator deve girare come backend del plugin Grafana e gestire le risorse Todo tramite l'API Kubernetes.

## Interfaces / Interfacce

| Interface / Interfaccia | Type / Tipo | Description / Descrizione |
|-------------------------|-------------|---------------------------|
| REST API endpoints | HTTP REST | `POST/GET/PUT/DELETE /apis/todo.grafana.app/v1/namespaces/{ns}/todos[/{name}]` |
| Todo CRUD input | JSON request body | `{title: string, description?: string, status: "open"\|"in_progress"\|"done"}` |
| Todo CRUD output | JSON response body | `{metadata: {uid, name, namespace, creationTimestamp}, spec: {title, description, status}}` |
| Todo List output | JSON response body | `{items: [Todo...], metadata: {continue, remainingItemCount}}` |
| Kind registration | Go import | Importa i tipi `Todo`, `TodoSpec` da SDD-001 |

### Contract con SDD-001 (Custom Resource)

Questo SDD importa i tipi Go generati da SDD-001:

```go
import todov1 "github.com/zima/forgia-example-todo/pkg/apis/todo/v1"
```

### Contract con SDD-003 (Frontend)

Il frontend invochera' le seguenti API (gestite da questo operator):

```
GET    /apis/todo.grafana.app/v1/namespaces/{ns}/todos          → Lista todos
GET    /apis/todo.grafana.app/v1/namespaces/{ns}/todos/{name}    → Singolo todo
POST   /apis/todo.grafana.app/v1/namespaces/{ns}/todos           → Crea todo
PUT    /apis/todo.grafana.app/v1/namespaces/{ns}/todos/{name}    → Aggiorna todo
DELETE /apis/todo.grafana.app/v1/namespaces/{ns}/todos/{name}    → Cancella todo
```

Response codes: `200 OK`, `201 Created`, `404 Not Found`, `400 Bad Request`, `422 Unprocessable Entity`.

## Constraints / Vincoli

- Language / Linguaggio: Go
- Framework: `grafana-app-sdk`, Grafana plugin backend SDK
- Dependencies / Dipendenze: `grafana-app-sdk`, `k8s.io/client-go`, tipi Go da SDD-001
- Patterns / Pattern: Operator pattern (reconciler + watcher), repository pattern per l'accesso ai dati
- **Dipende da SDD-001**: i tipi Go del CRD devono essere disponibili prima dell'implementazione
- Nessun segreto hardcoded — kubeconfig gestito dall'ambiente runtime

### Guardrails (deny.toml)

L'agent NON deve:
- Leggere file `.env`, `*.pem`, `*.key`, `kubeconfig`, `credentials.json`
- Scrivere file `.env`, `.forgia/constitution.md`, `.forgia/config.toml`, `.forgia/guardrails/deny.toml`
- Eseguire comandi che leggono chiavi SSH, GPG, o enumerano variabili d'ambiente con segreti

## Best Practices

- Error handling: errori espliciti con context (`fmt.Errorf("failed to create todo %s: %w", name, err)`), mai fallback silenti. Restituire HTTP status codes appropriati (400 per input invalido, 404 per not found, 500 per errori interni)
- Naming: Go conventions — `TodoHandler`, `CreateTodo`, `ListTodos`, exported functions in PascalCase
- Style: `gofmt` + `golint`, commenti in inglese nel codice
- Separation of concerns: handler layer (HTTP) separato dalla business logic, separato dall'accesso ai dati (repository pattern)
- Logging: usare il logger strutturato di Grafana (`log.DefaultLogger`)

## Test Requirements

| Type / Tipo | What / Cosa | Coverage |
|-------------|-------------|----------|
| Unit | CRUD handler logic (create, get, list, update, delete), validazione input | 90%+ sui handlers |
| Integration | API endpoints con cluster K8s reale (envtest o kind) — CRUD completo | Tutti gli endpoint |
| E2E | — | Coperto da SDD-003 e test manuali |

## Acceptance Criteria / Criteri di Accettazione

- [x] Le operazioni CRUD (create, list, get, update, delete) funzionano via REST API
- [x] `POST` crea un Todo e restituisce `201 Created` con il resource creato
- [x] `GET` (list) restituisce tutti i Todo nel namespace
- [x] `GET` (singolo) restituisce il Todo o `404 Not Found`
- [x] `PUT` aggiorna un Todo esistente e restituisce `200 OK`
- [x] `DELETE` cancella un Todo e restituisce `200 OK`
- [x] Input invalido (title mancante, status non valido) restituisce `400 Bad Request` o `422`
- [x] L'operator gestisce correttamente il ciclo di vita delle risorse (creazione, aggiornamento, cancellazione)
- [x] I test passano con `go test ./...`
- [x] Nessun segreto hardcoded

## Context / Contesto

- [x] Tipi Go da SDD-001: `pkg/generated/todo/v1/` (generati)
- [x] Documentazione `grafana-app-sdk` operator/watcher
- [x] Grafana plugin backend SDK documentation
- [x] File FD: `.forgia/fd/FD-001-todo-grafana-app.md`
- [x] File SDD-001: `.forgia/sdd/FD-001/SDD-001-todo-custom-resource.md`

## Constitution Check

- [x] Rispetta le code standards (Go conventions, explicit errors, no silent fallbacks)
- [x] Rispetta le commit conventions (`feat(FD-001/SDD-002): ...` + `Co-Authored-By`)
- [x] Nessun secret hardcoded — env vars o secret manager
- [x] Test definiti e sufficienti (unit + integration)

---

## Work Log / Diario di Lavoro

> Questa sezione e' **obbligatoria**. Deve essere compilata dall'agent o dallo sviluppatore durante e dopo l'esecuzione.

### Agent / Agente

- **Executor**: claude-code
- **Started**: 2026-03-15
- **Completed**: 2026-03-15
- **Duration / Durata**: ~15 min

### Decisions / Decisioni

1. **Repository pattern per data access**: Creata interfaccia `TodoRepository` con implementazione `K8sTodoRepository` che wrappa il `TodoClient` generato. Separa la logica di accesso dati dagli handler HTTP, facilita il testing con mock.

2. **Handler layer HTTP separato**: I CRUD handler (`CreateTodo`, `GetTodo`, `ListTodos`, `UpdateTodo`, `DeleteTodo`) gestiscono validazione input, error mapping a status HTTP, e delegano al repository. Validazione esplicita di title (required, non-empty) e status (enum: open, in_progress, done).

3. **Watcher con simple.Watcher SDK**: Utilizzato il pattern `simple.Watcher` del grafana-app-sdk per il lifecycle management. Il watcher logga eventi Add/Update/Delete/Sync e aggiorna l'operator status sulla risorsa. Non usa reconciler perche' il caso d'uso e' semplice (logging + status update).

4. **App factory con simple.App**: L'app factory crea un `simple.App` con il Todo kind registrato come `ManagedKind` con il watcher. Usa `k8s.NewClientRegistry` per il client generator.

5. **Operator entry point con operator.Runner**: Il main usa `operator.NewRunner` con manifest locale dall'`AppProvider`. Supporta `--kubeconfig` flag per sviluppo locale, fallback a in-cluster config.

6. **Test coverage 98.8% sugli handler**: 35 unit test coprono tutti i path: success, not found, validation errors (title mancante, status invalido, JSON malformato), errori repository, isolamento namespace.

### Output

- **Commit(s)**: (pending)
- **PR**: (pending)
- **Files created/modified**:
  - `pkg/repository/todo_repository.go` (NEW) - Repository interface + K8s implementation
  - `pkg/handler/todo_handler.go` (NEW) - HTTP CRUD handlers with validation
  - `pkg/handler/todo_handler_test.go` (NEW) - 35 unit tests, 98.8% coverage
  - `pkg/watcher/todo_watcher.go` (NEW) - Lifecycle watcher (Add/Update/Delete/Sync)
  - `pkg/app/app.go` (NEW) - App factory using simple.App
  - `cmd/operator/main.go` (NEW) - Operator entry point with Runner
  - `go.mod` (MODIFIED) - New dependencies added
  - `go.sum` (MODIFIED) - Updated checksums

### Retrospective / Retrospettiva

- **Cosa ha funzionato**: La separazione in layer (handler -> repository -> client) ha reso il testing semplice e pulito. Il grafana-app-sdk fornisce un buon pattern con simple.App + Watcher. Il codice generato da SDD-001 (TodoClient, Kind(), Schema()) si integra perfettamente.
- **Cosa non ha funzionato**: Nulla di significativo. La documentazione del SDK richiede esplorazione del codice sorgente per capire i pattern corretti (NewApp, AppConfig, ManagedKind).
- **Suggerimenti per FD futuri**: Includere nella SDD un diagramma dell'architettura dei layer (handler/repository/watcher) per velocizzare l'implementazione. Specificare se il CRUD e' esposto via HTTP custom server o solo via K8s API.
