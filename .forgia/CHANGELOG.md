# Forgia Changelog

## [FD-002] Add priority field to Todo — 2026-03-15

Aggiunto campo `priority` (low, medium, high, critical) al modello Todo su tutti i layer: schema CUE, tipi generati Go/TS, API handler, e componenti UI con badge colorati, ordinamento e filtro.

### SDDs Completati

- **SDD-001**: Data model & persistence — campo `priority` enum in CUE schema, tipi Go (`SpecPriority` + 4 costanti) e TypeScript generati, CRD JSON con validazione OpenAPI e default `"medium"`
- **SDD-002**: API layer — `TodoHandler` aggiornato per accettare/restituire/validare priority in tutte le operazioni CRUD, default `"medium"` per Todo esistenti, 10 nuovi test
- **SDD-003**: UI components — `Select` dropdown nel form, `Badge` colorato nella lista (blue/yellow/orange/red), controlli sort e filter nella pagina, modulo `priorityUtils.ts` condiviso

### Decisioni Chiave

- **Enum a 4 livelli** (low/medium/high/critical) invece di numerico 1-10 — segue il pattern `status` esistente, UX piu' semplice
- **`*string` per Priority in TodoRequest** — distingue "omesso" (nil → default medium) da "fornito" per validazione
- **HTTP 422** per errori di validazione (non 400) — coerenza con il pattern esistente nel codebase
- **`priorityUtils.ts` condiviso** — costanti, color map, funzioni pure per sort/filter usate da piu' componenti
- **`console.warn`** su priority mancante — nessun fallback silente, rispetta la constitution

### Retrospettiva Aggregata

**Cosa ha funzionato:**
- Il pattern `status` enum esistente era chiaro e facile da replicare per `priority` — tutti i 14 test Go passano
- I tipi generati da SDD-001 (SpecPriority, costanti) erano puliti e pronti all'uso in SDD-002
- I pattern esistenti dei componenti (Select, Badge) hanno reso semplice aggiungere il supporto priority nella UI
- Il setup di test basato su mock ha permesso iterazione veloce — tutti i 47 test frontend passano senza regressioni

**Cosa non ha funzionato:**
- Ambiguita' tra SDD spec (HTTP 400) e convenzione esistente nel codebase (HTTP 422) per errori di validazione — risolta seguendo la convenzione del codebase

**Suggerimenti per FD futuri:**
- Quando si specificano HTTP status code nelle SDD, fare riferimento alla convenzione gia' in uso nel codebase
- Specificare se i controlli filter devono essere single-select o multi-select per evitare ambiguita' implementativa
- Considerare un target Makefile per codegen (`make generate`) cosi' gli agent possono rigenerare i tipi automaticamente

---

## [FD-001] TODO Grafana App - grafana-app-sdk CRUD with zima-lab deploy — 2026-03-15

Applicazione completa per la gestione di TODO come Grafana app plugin, con backend operator Kubernetes e deploy su zima-lab.

### SDDs Completati

- **SDD-001**: TODO Custom Resource Definition — CUE schema, Go/TS codegen, CRD con validazione OpenAPI
- **SDD-002**: TODO Backend / Operator — CRUD handlers con repository pattern, watcher lifecycle, 98.8% test coverage
- **SDD-003**: TODO Frontend Plugin — React/TypeScript con @grafana/ui, 32 test (form, list, hook, API, E2E)
- **SDD-004**: Deployment to zima-lab — Dockerfile distroless, K8s manifests (RBAC, ResourceQuota), CI/CD script

### Decisioni Chiave

- **grafana-app-sdk v0.52.0** con CUE schema per definire il CRD Todo (Go >= 1.25.0 required)
- **Repository pattern** per separare handler HTTP da accesso dati K8s — facilita testing con mock
- **Webpack manuale** per il plugin frontend (`@grafana/plugin-configs` non disponibile su npm)
- **Distroless runtime image** (nonroot) per minimizzare superficie di attacco del container

### Retrospettiva Aggregata

**Cosa ha funzionato:**
- grafana-app-sdk codegen produce tipi Go, TS e CRD corretti dal CUE schema
- La separazione in layer (handler/repository/client e API/hook/componenti) rende il testing semplice
- Struttura modulare dei manifest K8s (un file per risorsa) facilita manutenzione e debug
- Mock di @grafana/ui funzionano bene per testing senza rendering reale

**Cosa non ha funzionato:**
- `grafana-app-sdk project kind add` va in panic senza prompt interattivi — richiede gestione manuale di manifest.cue
- `@grafana/plugin-configs` non esiste su npm — necessaria configurazione webpack manuale
- Coverage di `useTodos.ts` risulta bassa (53%) per artefatti source map SWC, nonostante branch testati
- Validazione YAML richiede pyyaml/kubectl — risolto con fallback a controlli strutturali

**Suggerimenti per FD futuri:**
- Documentare la versione esatta di grafana-app-sdk e tooling Grafana nel FD
- Includere diagrammi architetturali dei layer (handler/repository/watcher) nelle SDD
- Specificare se il CRUD e' esposto via HTTP custom o K8s API
- Aggiungere kubeconform come dipendenza obbligatoria (via mise.toml) per validazione schema K8s
- Considerare envtest per test CRD senza cluster reale
- Considerare Kustomize per varianti di deployment (dev/staging/prod)
- Considerare @grafana/scenes per plugin piu' complessi
