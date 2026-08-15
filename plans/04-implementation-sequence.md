# Implementation Sequence

## Status
Active roadmap.

## Scope

Stage-by-stage branch and PR sequence from groundwork to MVP hardening.

- Branching rules: `docs/branching-strategy.md`
- Stage 0 detail: `plans/05-stage-0-groundwork.md`
- Interface definitions: `docs/architecture/01-layer-interface-spec.md`
- Delivery policy: `docs/architecture/02-contract-first-delivery-plan.md`

Each row is one branch and one PR. Create the branch from `main` when the work starts, not before.

## Stage 0: Groundwork

Fully specified in `plans/05-stage-0-groundwork.md`.

| # | Branch | Purpose | Gates |
|---|--------|---------|-------|
| 0.1 | `chore/repo-hygiene` | CONTRIBUTING, SECURITY, CODEOWNERS, README status accuracy | Lint |
| 0.2 | `feat/cmd-server-bootstrap` | Go module init, HTTP server, `/health`, `/version`, config, request-ID logging | Unit |
| 0.3 | `ci/linux-gate-honesty` | Remove pass-through echoes, split CI jobs per gate, Linux authoritative | CI green |
| 0.4 | `ci/docker-linux-parity` | Dockerfile, .dockerignore, compose, container health check | Docker + health |
| 0.5 | `chore/cross-platform-make` | Replace POSIX-only Makefile targets with portable forms | Lint |

## Stage A: Contract-first

| # | Branch | Purpose | Gates |
|---|--------|---------|-------|
| A.1 | `test/contract-harness` | Contract test harness, shared assertion helpers, build tags | Contract |
| A.2 | `test/inference-contracts` | `ListModels`, `ChatStream` contract tests | Contract |
| A.3 | `feat/inference-ollama-adapter` | Ollama adapter implementation + unit tests | Contract + Unit |
| A.4 | `test/orchestrator-contracts` | `ChatService` contract tests with doubles | Contract |
| A.5 | `feat/orchestrator-chat-service` | ChatService implementation + unit tests | Contract + Unit |
| A.6 | `test/api-contracts-sse` | HTTP schema + SSE framing, ordering, terminal event contracts | Contract |
| A.7 | `feat/api-chat-handlers` | Chat handler (JSON + SSE), model list, validation + unit tests | Contract + Unit |

## Stage B: Vertical slice and integration

| # | Branch | Purpose | Gates |
|---|--------|---------|-------|
| B.1 | `feat/cmd-wire-chat-path` | Wire cmd → api → orchestrator → inference | All + Docker |
| B.2 | `test/integration-chat-path` | End-to-end chat happy path and error propagation | Integration |
| B.3 | `test/storage-contracts` | Session save/load/list contract tests | Contract |
| B.4 | `feat/storage-session-store` | Session store implementation + unit tests | Contract + Unit |
| B.5 | `test/integration-session-path` | Session persistence across chat turns | Integration |

## Stage C: Functional and hardening

| # | Branch | Purpose | Gates |
|---|--------|---------|-------|
| C.1 | `test/discovery-contracts` | Start/Stop/Status lifecycle contracts | Contract |
| C.2 | `feat/discovery-mdns` | mDNS advertisement + unit tests | Contract + Unit |
| C.3 | `feat/web-chat-ui` | SvelteKit chat view consuming SSE | Web check |
| C.4 | `test/smoke-container` | Container startup, `/health`, SSE stream smoke | Smoke + Docker |
| C.5 | `test/security-baseline` | Input validation, log leakage, error message safety | Security |

## Continuous: Workflow

| Branch | Trigger |
|--------|---------|
| `workflow/promote-lessons-<date>` | After any stage completes, promote accepted draft lessons |
| `workflow/agent-<change>` | Reflective Learning Agent proposes an agent definition change |

## Merge Order Dependencies

1. Stage 0 completes before Stage A begins. CI must be honest before it gates anything.
2. Within Stage A, contract branch merges before its paired implementation branch.
3. Stage B.1 requires A.3, A.5, and A.7 merged.
4. Stage C.4 requires B.1 merged and Docker gate green.
5. `workflow/` branches merge independently at any time.

## Stage Detail Documents

Detailed, implementation-ready plans are written one stage ahead:

| Stage | Document | Status |
|-------|----------|--------|
| 0 | `plans/05-stage-0-groundwork.md` | Ready |
| A | `plans/06-stage-a-contract-first.md` | Not written |
| B | — | Not written |
| C | — | Not written |

Write the next stage detail during the Reflect step of the preceding stage. In the same change, update the "Current stage detail" row in `.context/README.md` to point at the new document.
