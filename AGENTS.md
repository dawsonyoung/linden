# Linden

A privacy-first, safety optimized personal AI platform. Local-only capable on your home network for a secure chat experience, built on a foundation of trust that opens the door for future optional cloud-enhanced services.

## Architecture

### Source Layers (`src/`)

Each folder in `src/` is an architectural layer with a defined interface boundary.
Layers communicate only through their public interfaces — never by reaching into
another layer's internals. Every interface is testable in `validation/`.

```
src/
├── cmd/             # Entry point (main.go) — wires layers together, no logic
├── api/             # HTTP gateway — receives client requests, returns responses
│                      Interface: routes, middleware, SSE streaming, OpenAI compat
│                      Depends on: orchestrator
├── orchestrator/    # Request routing — manages conversation context, dispatches
│                      Interface: ChatService (accepts message, returns streamed response)
│                      Depends on: inference, storage
├── inference/       # LLM abstraction — talks to Ollama or other backends
│                      Interface: Client (ChatStream, ListModels)
│                      Depends on: errs (shared kernel)
├── storage/         # Data persistence — conversations, profiles, documents
│                      Interface: Store (Save, Load, Query)
│                      Depends on: errs (shared kernel)
├── discovery/       # Network discovery — mDNS advertisement
│                      Interface: Advertiser (Start, Stop)
│                      Depends on: errs (shared kernel)
├── mcp/             # Extensibility — Model Context Protocol layer boundary
│                      Depends on: errs (shared kernel)
├── errs/            # Shared kernel — the published error taxonomy
│                      Interface: Code, Error, New, Wrap, CodeOf, Is
│                      Depends on: nothing (stdlib only)
├── config/          # Environment configuration loading and validation
│                      Depends on: nothing (stdlib only)
└── web/             # SvelteKit frontend — builds to static files served by api
                       Interface: HTTP (static assets consumed by api layer)
                       Depends on: nothing (standalone build artifact)
```

**Layer rules:**
- `cmd/` depends on all layers (it wires them). No other layer imports `cmd/`.
- `api/` depends only on `orchestrator/`.
- `orchestrator/` depends on `inference/` and `storage/`.
- `inference/`, `storage/`, `discovery/`, `mcp/` are leaf layers — they depend on no other layer.
- `web/` is a separate build (Node/SvelteKit). Its output is static files consumed by `api/`.

**Shared kernels:** a dependency-free, logic-free package may be imported by any
layer, including leaf layers. Leaf layers exist for modularity — replaceability
behind a clear interface — not independence for its own sake. `errs/` is the only
such package; a second requires an ADR. See `docs/adr/0004-shared-error-taxonomy.md`.

### Validation (`validation/`)

```
validation/
├── contracts/       # Interface contract tests — does each layer honor its interface?
├── integration/     # Cross-layer tests — do layers work together correctly?
├── security/        # Security validation — prompt injection, data leak detection
└── smoke/           # Full E2E smoke tests — start server, hit endpoints, verify responses
```

### AI-First Directories

```
prompts/         # Reusable prompt templates for product features and dev tooling
evals/           # AI output evaluation datasets and test harnesses
tools/           # MCP-compatible tool definitions, sandboxed plugins, docgen
.agents/         # Agent workflow definitions and role configurations
.mcp/            # MCP server configuration and tool manifests
.context/        # Orientation index for AI agents (see .context/README.md)
llm.txt          # Top-level LLM orientation file (codebase summary for AI agents)
```

### Documentation (`docs/`)

```
docs/
├── product/      # Published: product spec and PRDs. Observable behavior only.
├── architecture/ # Layer interfaces and dependency rules
├── process/      # Branching, contract-first delivery, agent workflow
└── adr/          # Architecture decision records
```

`docs/product/` never contains file paths, function names, or struct fields. It is written for readers who have not seen the codebase.

### Planned (future directories at repo root)
- Deferred expansion directories (content-sdk, protocol, marketplace, clients/mobile, system) are tracked in external strategic roadmaps and should not be scaffolded in the current repository phase.

## Build & Run

```sh
# Prerequisites: Go 1.22+, Node 20+, Ollama with tinyllama pulled
make dev          # Build web UI + start server on :8080
make test         # Run all unit tests (Go + web)
make validate     # Run validation suite (contracts, integration, smoke)
make lint         # Static analysis (go vet, staticcheck, svelte-check)
make build        # Production build → bin/linden
```

## Conventions

### Go (src/cmd/, src/api/, src/orchestrator/, src/inference/, src/storage/, src/discovery/)
- Use stdlib. No web frameworks. Justify every external dependency.
- Each layer exports a Go interface. Consumers depend on the interface, not the implementation.
- Return errors with context wrapping. Never panic in library code.
- Table-driven unit tests per layer. Mock dependencies via interfaces.
- Test names: `Test_FunctionName_Scenario_Expected`
- No user data in logs. Ever. Not even at debug level.

### Web (src/web/)
- SvelteKit with static adapter. Output served by the api layer.
- TypeScript strict mode. No `any`.
- Tailwind CSS v4 and DaisyUI for clean utility styling; `@ai-sdk/svelte` for chat streaming.
- Fetch + Server-Sent Events for streaming.

### Validation (validation/)
- `contracts/` — verify each layer's interface contract in isolation (mocked dependencies).
- `integration/` — verify layer interactions with real (not mocked) wiring.
- `security/` — verify no user data leaks in logs, responses, or error messages.
- `smoke/` — start the full server, hit endpoints, verify streamed responses.

### Security (all code)
- Validate all inputs at the API boundary.
- No hardcoded secrets, tokens, or keys.
- No user content in error messages returned to clients.
- Sanitize any content before embedding in RAG (future).
- TLS for all network communication (future).

### Git
- Branch per PR, created from latest `main` when work starts.
- Branch types: `feat/`, `fix/`, `test/`, `refactor/`, `docs/`, `ci/`, `chore/`, `workflow/`
- Format: `<type>/<scope>-<short-description>` (e.g. `test/inference-contracts`)
- `workflow/` is reserved for `.agents/**`, `.context/**`, and workflow tooling changes. Never mix with product code.
- Commits: Conventional Commits (`feat:`, `fix:`, `docs:`, etc.)
- Every PR must pass CI (lint + test + build).
- Squash merge to main; delete branch after merge.
- See `docs/project/branching-strategy.md` for the full branching model and PR rules.

## Design Principles

1. **Provably private** — architecture enforces privacy, not just policy.
2. **Zero-config for users** — complexity hidden; simplicity exposed.
3. **Tool, not authority** — answers questions, cites sources, has no opinions.
4. **Owned, not rented** — works forever without company infrastructure.
5. **Offline-first** — full functionality without internet.

## Self-Correction and Learning

- `.context/` is an orientation index, not a source of truth. See `.context/README.md` for its role and rules.
- Read `.context/lessons_learned.md` at the start of every session to avoid known failure patterns.
- Before closing any task that involved a build error, test failure, lint issue, or repeated retry, append a structured lesson to `.context/draft_lessons.md`.
- Use this format for each entry: `[Rule-NNN]: <actionable rule> | Context: <trigger> | Negative: <what to avoid>`.
- Never write directly to `.context/lessons_learned.md`; that file is human-promoted only.

## Agent Guidelines

When working in this codebase:

- Consult `plans/` as a dynamic holding area for forward-looking design intent, architectural ideas, or when directed to implement specific future capabilities.
- Prefer simple, obvious solutions over clever ones.
- Do not add features, abstractions, or refactors beyond what is requested.
- Do not add comments to code you did not write or change.
- When creating new files, follow the patterns of adjacent existing files.
- When uncertain about a design decision, note it as a TODO rather than guessing.
- Never generate placeholder/example user data that looks real (names, emails, SSNs, etc).
- Test edge cases: empty input, missing Ollama, network timeout, malformed JSON.
- This is a privacy product. Treat every line of code as if an auditor is reviewing it for data leaks.
- Always strictly follow the full development process (contract-first delivery, implementation sequence stages, branching rules). Start by invoking the Workflow Orchestrator (`.agents/agents/workflow-orchestrator.agent.md`) as the primary entrypoint for any task. Never circumvent the process, and always run local validation gates to confirm health before pushing; do not rely on remote CI as the primary testing environment.
