---
id: "SDD-004"
fd: "FD-001"
title: "Deployment to zima-lab (Kubernetes manifests, CI/CD, Tanka/Helm config)"
status: planned
agent: ""
assigned_to: ""
created: "2026-03-15"
started: ""
completed: ""
tags: [kubernetes, deployment, ci-cd, zima-lab, helm]
---

# SDD-004: Deployment to zima-lab

> Parent FD: [[FD-001]]

## Scope

Configurare il deployment dell'applicazione TODO (operator + plugin) nell'ambiente Kubernetes zima-lab. Questo include:

- Kubernetes manifests per il deployment dell'operator (Deployment, Service, ServiceAccount, RBAC)
- Installazione del CRD (da SDD-001) nel cluster
- Configurazione del plugin Grafana (provisioning o sidecar)
- Dockerfile per il backend operator
- CI/CD pipeline o script di deploy per zima-lab
- Namespace e resource quotas

**NON include**: la configurazione di kubeconfig (gestito dal CD pipeline), l'infrastruttura del cluster (gia' esistente).

## Interfaces / Interfacce

| Interface / Interfaccia | Type / Tipo | Description / Descrizione |
|-------------------------|-------------|---------------------------|
| CRD manifest | YAML | Da SDD-001: `kubectl apply -f` del CRD Todo |
| Operator Deployment | YAML | Deployment K8s per il backend operator (container image) |
| Plugin provisioning | YAML/JSON | Configurazione per caricare il plugin in Grafana (env vars o provisioning file) |
| Dockerfile (operator) | Dockerfile | Build multi-stage per il backend Go |
| CI/CD pipeline | Shell/YAML | Script o pipeline per build + deploy su zima-lab |

### Contract con SDD-001 (CRD)

Il CRD YAML generato da SDD-001 deve essere applicato prima del deployment dell'operator:

```bash
kubectl apply -f deploy/crd/todo-crd.yaml
kubectl apply -f deploy/operator/
```

### Contract con SDD-002 (Backend) e SDD-003 (Frontend)

- L'immagine Docker del backend include l'operator compilato (Go binary)
- Il plugin frontend e' incluso come asset statico o distribuito tramite Grafana plugin provisioning
- Entrambi devono essere buildati prima del deploy

## Constraints / Vincoli

- Language / Linguaggio: YAML (K8s manifests), Dockerfile, Shell (deploy scripts)
- Framework: Kubernetes, Docker
- Dependencies / Dipendenze: Output di SDD-001 (CRD), SDD-002 (binary Go), SDD-003 (plugin dist)
- Patterns / Pattern: K8s Deployment + Service + RBAC, multi-stage Dockerfile
- **Dipende da SDD-001, SDD-002, SDD-003**: tutti devono essere completati prima del deploy
- Target cluster: zima-lab (Kubernetes, namespace dedicato)
- Nessun segreto hardcoded nei manifest — usare K8s Secrets referenziati tramite `secretKeyRef`
- Kubeconfig NON deve essere incluso nei file — gestito dal CD pipeline

### Guardrails (deny.toml)

L'agent NON deve:
- Leggere file `.env`, `*.pem`, `*.key`, `kubeconfig`, `kubeconfig.*`, `credentials.json`
- Scrivere file `.env`, `.forgia/constitution.md`, `.forgia/config.toml`, `.forgia/guardrails/deny.toml`
- Eseguire comandi che leggono chiavi SSH, GPG, o enumerano variabili d'ambiente con segreti
- **IMPORTANTE**: NON includere valori reali di secret nei manifest. Usare placeholder `<SECRET_PLACEHOLDER>` dove necessario

## Best Practices

- Error handling (shell): `set -euo pipefail` in ogni script, exit codes espliciti, messaggi di errore su stderr
- Naming: K8s resource names in kebab-case (`todo-operator`, `todo-app-plugin`), shell functions con `prefix_action` pattern
- Style (shell): shellcheck mandatory, 2 spaces indent, `command -v` per tool check
- Style (YAML): indent 2 spaces, commenti esplicativi per ogni risorsa K8s
- Dockerfile: multi-stage build (builder + runtime), immagine runtime minimale (distroless o alpine), no root user
- RBAC: principio del minimo privilegio — l'operator deve avere solo i permessi necessari per gestire le risorse Todo

## Test Requirements

| Type / Tipo | What / Cosa | Coverage |
|-------------|-------------|----------|
| Unit | Validazione YAML con `kubeval` o `kubeconform` | Tutti i manifest |
| Integration | Deploy su cluster locale (kind/minikube) con verifica pod running | Deploy completo |
| E2E | Deploy su zima-lab (o staging) con smoke test CRUD | Pipeline completa |

## Acceptance Criteria / Criteri di Accettazione

- [ ] Il CRD e' applicabile con `kubectl apply` senza errori
- [ ] Il Deployment dell'operator parte correttamente (pod in stato Running)
- [ ] L'RBAC concede solo i permessi necessari (verbs: get, list, watch, create, update, delete su todos.todo.grafana.app)
- [ ] Il Dockerfile builda con successo e produce un'immagine funzionante
- [ ] Il plugin Grafana e' caricato nell'istanza Grafana su zima-lab
- [ ] Il deploy su zima-lab avviene con successo tramite pipeline/script CD
- [ ] Nessun segreto e' hardcoded nei manifest o negli script
- [ ] Gli script di deploy usano `set -euo pipefail` e controllano i prerequisiti con `command -v`

## Context / Contesto

- [ ] Struttura cluster zima-lab (namespace, risorse disponibili)
- [ ] CRD manifest da SDD-001
- [ ] Binary Go da SDD-002
- [ ] Plugin dist da SDD-003
- [ ] Grafana plugin provisioning documentation
- [ ] File FD: `.forgia/fd/FD-001-todo-grafana-app.md`
- [ ] Shell conventions: `.forgia/dev-guide/lang/shell.md`

## Constitution Check

- [ ] Rispetta le code standards (shell conventions, YAML standards, explicit errors)
- [ ] Rispetta le commit conventions (`feat(FD-001/SDD-004): ...` + `Co-Authored-By`)
- [ ] Nessun secret hardcoded — K8s Secrets con `secretKeyRef`, placeholder dove necessario
- [ ] Test definiti e sufficienti (YAML validation + integration deploy)

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
