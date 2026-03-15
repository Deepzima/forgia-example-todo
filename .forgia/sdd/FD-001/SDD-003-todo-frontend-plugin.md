---
id: "SDD-003"
fd: "FD-001"
title: "TODO Frontend Plugin (Grafana app plugin UI, pages, components)"
status: planned
agent: ""
assigned_to: ""
created: "2026-03-15"
started: ""
completed: ""
tags: [grafana, plugin, react, typescript, frontend]
---

# SDD-003: TODO Frontend Plugin

> Parent FD: [[FD-001]]

## Scope

Implementare il frontend del plugin Grafana app per la gestione dei TODO. Questo include:

- Scaffolding del plugin Grafana app (usando `@grafana/create-plugin` o struttura manuale)
- Pagina lista TODO: mostra tutti i TODO con stato, titolo, data di creazione
- Form creazione/modifica TODO: campi `title`, `description`, `status` (dropdown)
- Azioni inline: cambia stato, modifica, elimina
- Integrazione con le API REST del backend (SDD-002)
- Plugin registration in Grafana (`plugin.json`)

L'interfaccia deve essere semplice, funzionale e allineata con il design system Grafana (`@grafana/ui`).

## Interfaces / Interfacce

| Interface / Interfaccia | Type / Tipo | Description / Descrizione |
|-------------------------|-------------|---------------------------|
| REST API (consumed) | HTTP REST | Chiama le API di SDD-002: `GET/POST/PUT/DELETE /apis/todo.grafana.app/v1/namespaces/{ns}/todos` |
| Plugin registration | `plugin.json` | Registrazione del plugin in Grafana con routes, includes, pages |
| Todo data model (frontend) | TypeScript interface | `{metadata: {uid, name, namespace, creationTimestamp}, spec: {title, description, status}}` |

### Contract con SDD-002 (Backend)

Il frontend consuma le API REST esposte dal backend. Interfaccia TypeScript attesa:

```typescript
interface TodoSpec {
  title: string;
  description?: string;
  status: 'open' | 'in_progress' | 'done';
}

interface TodoMetadata {
  uid: string;
  name: string;
  namespace: string;
  creationTimestamp: string;
}

interface TodoResource {
  metadata: TodoMetadata;
  spec: TodoSpec;
}

interface TodoList {
  items: TodoResource[];
  metadata: {
    continue?: string;
    remainingItemCount?: number;
  };
}
```

### API Calls

```typescript
// List todos
GET /apis/todo.grafana.app/v1/namespaces/{ns}/todos → TodoList

// Get single todo
GET /apis/todo.grafana.app/v1/namespaces/{ns}/todos/{name} → TodoResource

// Create todo
POST /apis/todo.grafana.app/v1/namespaces/{ns}/todos
Body: { spec: TodoSpec } → TodoResource (201)

// Update todo
PUT /apis/todo.grafana.app/v1/namespaces/{ns}/todos/{name}
Body: { spec: TodoSpec } → TodoResource (200)

// Delete todo
DELETE /apis/todo.grafana.app/v1/namespaces/{ns}/todos/{name} → 200
```

## Constraints / Vincoli

- Language / Linguaggio: TypeScript (strict mode)
- Framework: React, `@grafana/ui`, `@grafana/runtime`, `@grafana/data`
- Dependencies / Dipendenze: Grafana Plugin SDK (frontend), `@grafana/ui` components
- Patterns / Pattern: Grafana app plugin pattern (AppPlugin, AppRootPage)
- Build: webpack (Grafana plugin tooling)
- **Dipende da SDD-002**: le API REST devono essere disponibili (ma lo sviluppo puo' procedere con mock)
- Nessun segreto hardcoded

### Guardrails (deny.toml)

L'agent NON deve:
- Leggere file `.env`, `*.pem`, `*.key`, `kubeconfig`, `credentials.json`
- Scrivere file `.env`, `.forgia/constitution.md`, `.forgia/config.toml`, `.forgia/guardrails/deny.toml`
- Eseguire comandi che leggono chiavi SSH, GPG, o enumerano variabili d'ambiente con segreti

## Best Practices

- Error handling: mostrare errori all'utente con `@grafana/ui` Alert component, mai fallback silenti. Gestire errori API (network, 4xx, 5xx) con messaggi chiari
- Naming: TypeScript conventions — `TodoList`, `TodoForm`, `useTodos`, `todoApi`; `const` over `let`; explicit return types su funzioni exported
- Style: componenti `@grafana/ui` per coerenza con il design system Grafana (Button, Input, Select, Card, Alert)
- Separation of concerns: componenti UI separati dalla logica di data fetching (custom hooks)
- Validazione: validare input utente nel form (title required, status enum) prima di inviare all'API

## Test Requirements

| Type / Tipo | What / Cosa | Coverage |
|-------------|-------------|----------|
| Unit | Componenti React (TodoList, TodoForm), custom hooks (useTodos) | 80%+ sui componenti |
| Integration | API calls con mock server | Tutti gli endpoint |
| E2E | Flusso completo: lista → crea → modifica → cancella | Happy path |

## Acceptance Criteria / Criteri di Accettazione

- [ ] Il plugin Grafana si carica correttamente (`plugin.json` valido, nessun errore in console)
- [ ] La pagina lista TODO mostra tutti i TODO con titolo, stato e data di creazione
- [ ] Il form di creazione permette di inserire titolo, descrizione e stato
- [ ] Il form di modifica permette di aggiornare titolo, descrizione e stato di un TODO esistente
- [ ] L'azione di cancellazione elimina un TODO con conferma
- [ ] Il dropdown dello stato mostra solo `open`, `in_progress`, `done`
- [ ] Gli errori API vengono mostrati all'utente con un messaggio chiaro
- [ ] I test passano con `npm test` / `yarn test`
- [ ] Nessun segreto hardcoded

## Context / Contesto

- [ ] Grafana Plugin Development documentation
- [ ] `@grafana/create-plugin` scaffolding tool
- [ ] `@grafana/ui` component library documentation
- [ ] API contract da SDD-002: `.forgia/sdd/FD-001/SDD-002-todo-backend-operator.md`
- [ ] File FD: `.forgia/fd/FD-001-todo-grafana-app.md`

## Constitution Check

- [ ] Rispetta le code standards (TypeScript strict, explicit return types, no silent fallbacks)
- [ ] Rispetta le commit conventions (`feat(FD-001/SDD-003): ...` + `Co-Authored-By`)
- [ ] Nessun secret hardcoded — env vars o secret manager
- [ ] Test definiti e sufficienti (unit + integration + E2E)

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
