---
id: "SDD-002"
fd: "FD-001"
title: "TODO Backend / Operator (CRUD handlers, grafana-app-sdk operator logic)"
status: planned
agent: ""
assigned_to: ""
created: "2026-03-15"
started: ""
completed: ""
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

- [ ] Le operazioni CRUD (create, list, get, update, delete) funzionano via REST API
- [ ] `POST` crea un Todo e restituisce `201 Created` con il resource creato
- [ ] `GET` (list) restituisce tutti i Todo nel namespace
- [ ] `GET` (singolo) restituisce il Todo o `404 Not Found`
- [ ] `PUT` aggiorna un Todo esistente e restituisce `200 OK`
- [ ] `DELETE` cancella un Todo e restituisce `200 OK`
- [ ] Input invalido (title mancante, status non valido) restituisce `400 Bad Request` o `422`
- [ ] L'operator gestisce correttamente il ciclo di vita delle risorse (creazione, aggiornamento, cancellazione)
- [ ] I test passano con `go test ./...`
- [ ] Nessun segreto hardcoded

## Context / Contesto

- [ ] Tipi Go da SDD-001: `pkg/apis/todo/v1/` (generati)
- [ ] Documentazione `grafana-app-sdk` operator/watcher
- [ ] Grafana plugin backend SDK documentation
- [ ] File FD: `.forgia/fd/FD-001-todo-grafana-app.md`
- [ ] File SDD-001: `.forgia/sdd/FD-001/SDD-001-todo-custom-resource.md`

## Constitution Check

- [ ] Rispetta le code standards (Go conventions, explicit errors, no silent fallbacks)
- [ ] Rispetta le commit conventions (`feat(FD-001/SDD-002): ...` + `Co-Authored-By`)
- [ ] Nessun secret hardcoded — env vars o secret manager
- [ ] Test definiti e sufficienti (unit + integration)

---

## Work Log / Diario di Lavoro

> Questa sezione e' **obbligatoria**. Deve essere compilata dall'agent o dallo sviluppatore durante e dopo l'esecuzione.

### Agent / Agente

- **Executor**: <!-- openhands | claude-code | manual | name -->
- **Started**: <!-- timestamp -->
- **Completed**: <!-- timestamp -->
- **Duration / Durata**: <!-- total time -->

### Decisions / Decisioni

1. <!-- decisione 1: cosa e perche' -->

### Output

- **Commit(s)**: <!-- hash -->
- **PR**: <!-- link -->
- **Files created/modified**:
  - `path/to/file`

### Retrospective / Retrospettiva

- **Cosa ha funzionato**:
- **Cosa non ha funzionato**:
- **Suggerimenti per FD futuri**:
