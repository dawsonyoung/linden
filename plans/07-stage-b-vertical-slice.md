# Stage B: Vertical Slice and Integration

## Status
Ready to implement.

## Purpose

Wire the completed Stage A layers together, prove they operate collectively via integration tests, and then implement the persistence tier (storage) for conversational sessions. 

Branching rules: `docs/project/branching-strategy.md`.
Delivery policy: `docs/project/contract-first-delivery.md`.

## Sequencing

```
B.1 feat/cmd-wire-chat-path ──→ B.2 test/integration-chat-path
                                        ↓
B.3 test/storage-contracts ───→ B.4 feat/storage-session-store
                                        ↓
B.5 test/integration-session-path
```

Each contract branch merges before its paired implementation branch. The implementation PR is not ready for review until the contract tests it satisfies are already on `main`.

---

## B.1 `feat/cmd-wire-chat-path`

Intent: Wire the executable (`src/cmd/main.go`) to connect the HTTP API, orchestrator, and inference layers together into a functional vertical slice.

Agents: Implementation Agent, then Review Agent.

### Tasks

1. Update `src/cmd/main.go` to construct the `inference.Client` and `orchestrator.ChatService` implementations.
2. Inject these instances into `api.NewServer`.
3. Provide the actual runtime configuration (e.g., Ollama URL) via environment variables loaded by the `config` package.

### Gate

Build and unit tests must remain green. Docker validate passes.

Estimated size: 50–100 lines.

---

## B.2 `test/integration-chat-path`

Intent: Prove that the connected layers interact correctly.

Agents: Contract Test Agent (for integration), then Review Agent.

### Tasks

1. Create `validation/integration/chat_integration_test.go`.
2. Construct the real API server, Orchestrator, and Inference layers in the test harness.
3. Test end-to-end chat requests against a test double for the actual Ollama network boundary (using `httptest.Server` at the inference layer boundary).

### Coverage Required

- Sending a chat request successfully propagates through the stack and streams SSE back.
- Network timeouts at the Ollama boundary propagate out as 503 HTTP responses from the API layer.
- Ensure no log leakage of user inputs occurs during full request transit.

### Gate

Integration suite passes (`go test -tags=integration ./integration/...`).

Estimated size: 150–250 lines.

---

## B.3 `test/storage-contracts`

Intent: Define the interface for session storage and its contract.

Agents: Design Agent (for interface shape), Contract Test Agent, then Review Agent.

### Tasks

1. Define `src/storage.Store` interface with methods for: saving a session turn, loading a session history, and listing active sessions.
2. Create `validation/contracts/storage_contract_test.go` to enforce storage behavior constraints.

### Coverage Required

- A saved turn is durably retrieved.
- Listing sessions returns chronologically ordered metadata.
- Unknown sessions map to `errs.NotFound`.

### Gate

Contract tests turn green (using an in-memory double).

Estimated size: 200–300 lines.

---

## B.4 `feat/storage-session-store`

Intent: Implement the storage layer with a local SQLite or file-system database.

Agents: Implementation Agent, then Unit Test Agent, then Review Agent.

### Tasks

1. Implement the `storage.Store` interface in `src/storage/`.
2. Map IO and database errors to the shared taxonomy in `errs`.
3. Provide unit testing around file locking and edge cases.

### Obligation

Carried from `docs/adr/0002-role-of-docker-in-distribution.md`: must declare a data directory and a corresponding volume in the `docker-compose.yml` or container configuration so that restarts do not wipe conversations.

### Gate

Contract and Unit tests.

Estimated size: 300–400 lines.

---

## B.5 `test/integration-session-path`

Intent: Integrate the storage layer into the orchestrator so that chat sessions persist and restore historical context.

Agents: Implementation Agent, Integration Agent, Review Agent.

### Tasks

1. Update `orchestrator` to accept a `storage.Store` dependency.
2. Modify `ChatStream` to persist inputs and fetch previous history using a provided session ID.
3. Write `validation/integration/session_integration_test.go` proving multi-turn memory.

### Gate

Integration suite.

Estimated size: 200 lines.

## Stage B Exit Criteria

1. The API can process chat streams through a fully wired `cmd` entrypoint.
2. Conversations have a persistent identifier.
3. Chat history is durably stored to disk and reloaded for multi-turn context.
4. Container configuration mounts a persistent data volume.
5. Integration suites cover the happy path and network error boundaries.
