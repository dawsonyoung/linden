# Linden: Repository Setup and MVP Plan

## Objective
Create a working MVP that proves Linden can run local AI reliably, expose a secure LAN-accessible server, and deliver a usable web client with clear privacy controls.

## MVP Scope
In scope:

1. Local server in Go.
2. Local model inference via Ollama.
3. Web client for chat and settings.
4. Local document ingestion and retrieval (basic RAG).
5. LAN pairing/authentication for second device access (via explicit IP; mDNS deferred).
6. Structured logging and baseline observability.
7. CI checks and validation workflows in the repository.

Out of scope (for MVP):

1. Mobile native apps.
2. Large plugin marketplace.
3. Distributed multi-node inference.
4. Hardware appliance manufacturing readiness.

## Repository Setup Plan

### Phase 0: Baseline Hygiene
Deliverables:

1. Confirm canonical root structure:
   - src/cmd
   - src/api
   - src/orchestrator
   - src/inference
   - src/storage
   - src/discovery
   - src/web
   - validation/contracts
   - validation/integration
   - validation/security
   - validation/smoke
   - docs/adr
   - scripts
   - .github/workflows
2. Add CONTRIBUTING.md, SECURITY.md, CODEOWNERS.
3. Expand README with architecture, run instructions, and support matrix.
4. Add reproducible local setup script support for Windows first; keep Linux/macOS stubs documented.

Acceptance criteria:

1. New contributor can clone, run setup, run tests, and launch local dev stack.
2. CI validates lint, unit tests, and smoke start-up checks.

### Phase 1: Go Service Foundation
Deliverables:

1. Initialize Go module and package boundaries per layer.
2. Build HTTP server bootstrap under src/cmd.
3. Add health endpoint and version endpoint.
4. Add structured config loading (env + local config file).
5. Add structured logging with request IDs.

Acceptance criteria:

1. Service boots consistently with deterministic config behavior.
2. Health check passes locally and in CI.

### Phase 2: Inference Integration (Ollama)
Deliverables:

1. Define inference provider interface in src/inference.
2. Implement Ollama adapter for text generation and model listing.
3. Add timeout, retry, and graceful error mapping.
4. Add model availability preflight check at startup.

Acceptance criteria:

1. Chat request can complete against local model.
2. User-facing errors are actionable when Ollama/model is unavailable.

### Phase 3: API + Orchestration
Deliverables:

1. Define API contracts for:
   - chat completion
   - model selection
   - session history
   - settings
2. Implement orchestrator pipeline:
   - request validation
   - context assembly
   - inference call
   - response formatting
3. Add unit tests for orchestration decisions and failure modes.

Acceptance criteria:

1. End-to-end request flow works from HTTP to model and back.
2. Contract tests enforce response schema stability.

### Phase 4: Storage + Basic RAG
Deliverables:

1. Local metadata store for sessions, settings, and document references.
2. Document ingestion pipeline:
   - file parsing (initially txt/md/pdf if practical)
   - chunking
   - embedding generation (local-compatible path)
   - vector index persistence
3. Retrieval integration in orchestrator with simple relevance threshold.

Acceptance criteria:

1. User can add documents and get grounded responses.
2. Retrieval is optional and can be disabled in settings.

### Phase 5: Discovery and LAN Access
Deliverables:

1. Service discovery strategy (explicit host flow; mDNS deferred to post-MVP).
2. Device pairing flow with one-time approval token.
3. Auth/session token model for local network clients.
4. Basic allow-list of paired devices.

Acceptance criteria:

1. Second device can discover/connect with explicit approval.
2. Unpaired devices are denied by default.

### Phase 6: Web Client MVP
Deliverables:

1. Chat interface (history, retries, stop generation).
2. Settings panel:
   - model selection
   - retrieval toggles
   - device management
   - data controls (clear history/index)
3. Privacy UX copy for local-first behavior and permission prompts.

Acceptance criteria:

1. Non-technical user can complete first-use flow in under 10 minutes.
2. Core features work on desktop and mobile browsers.

### Phase 7: Validation and Hardening
Deliverables:

1. validation/contracts tests for API schemas.
2. validation/integration tests for server + inference + storage.
3. validation/security checks for auth boundaries, input constraints, and safe defaults.
4. validation/smoke tests for install and launch.

Acceptance criteria:

1. CI gates on all validation suites relevant to MVP.
2. Known high-risk failure modes are documented with mitigations.

## MVP Release Criteria

1. Install script succeeds on target Windows environments.
2. Single-host chat is stable for repeated sessions.
3. LAN client pairing and usage are functional and secure by default.
4. Basic retrieval improves grounded answer quality for local docs.
5. Privacy controls are visible, understandable, and test-covered.

## Suggested 10-Week MVP Timeline

1. Weeks 1-2: Repo hygiene, Go service foundation, CI baseline.
2. Weeks 3-4: Ollama integration and core API/orchestrator flow.
3. Weeks 5-6: Storage and retrieval pipeline.
4. Weeks 7-8: Discovery, pairing, and web client completion.
5. Weeks 9-10: Validation hardening, docs, and release candidate stabilization.

## Immediate Next Actions

1. Finalize Go module bootstrap and minimum server startup.
2. Implement inference interface + Ollama adapter first.
3. Lock initial API contracts and write contract tests before UI expansion.
4. Add MVP definition checklist to README and track weekly against it.