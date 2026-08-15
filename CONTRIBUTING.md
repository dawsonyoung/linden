# Contributing to Linden

## Prerequisites

| Tool | Version | Notes |
|------|---------|-------|
| Go | 1.22+ | Server |
| Node.js | 20+ | Web client |
| Ollama | current | Local model backend |
| Docker | current | Required for the runtime gate |
| make | any | Linux and macOS native; see Windows below |

## Setup

```powershell
.\scripts\setup.ps1
```

Installs Go, Node, and Ollama if missing and pulls the `tinyllama` model. Restart your terminal afterward if anything was newly installed.

## Windows

Linux is the authoritative platform. Windows is supported for development, but CI gates on Linux only.

`make` is not available on Windows by default. Run targets through WSL or Git Bash. PowerShell scripts under `scripts/` work natively.

If a command behaves differently across platforms, the Linux form is correct.

## Branching and commits

Follow [docs/project/branching-strategy.md](docs/project/branching-strategy.md). In short:

- One branch per PR, created from latest `main` when work starts.
- Format: `<type>/<scope>-<short-description>`.
- Conventional Commits, ordered contract tests → implementation → unit tests.
- `workflow/` branches are reserved for `.agents/**`, `.context/**`, and workflow tooling. Never mix with product code.

## Test-first rule

Contract tests come before implementation. This is not a preference.

1. Write the contract test in `validation/contracts/`. Commit it.
2. Implement until it passes.
3. Add unit tests alongside the implementation.

An implementation PR without a corresponding contract test is not ready for review. See [docs/project/contract-first-delivery.md](docs/project/contract-first-delivery.md).

Stage 0 is the one exception — it introduces no cross-layer interface. The exemption ends at Stage A.

## Before opening a PR

1. Run the Spec step. Classify the change against the product spec, PRDs, and interface spec, and state the result in the PR body — including `Spec impact: none` when nothing applies.
2. Run the gates below.
3. Fill in every section of the PR template. Unused sections say `N/A` with a reason; do not delete them.

## Gates

```sh
make lint                                                    # go vet + svelte-check
cd src && go test -race ./...                                # unit
cd src && go test -tags=contracts -race ../validation/contracts/...
cd src && go test -tags=integration -race ../validation/integration/...
make docker-build && make docker-smoke                       # runtime
```

Not every gate applies to every change. Documentation-only PRs need lint alone.

## Required status checks

`main` is protected. These checks must pass before merge:

| Check | Gate |
|-------|------|
| `lint` | `go vet` and `gofmt` |
| `unit` | `go test -race` across `src/` |
| `build` | Server binary compiles |
| `contract` | Contract suite, when present |
| `integration` | Integration suite, when present |
| `web` | Type check and build, when present |

`contract`, `integration`, and `web` report a visible skip until the
corresponding suite or project exists. A skip is recorded as a notice in the
job log; it never reports a passing test run.

The Docker runtime check joins this list in Stage 0.4.

## Documentation

Four directories, each with an admission test in its README:

| Directory | Holds |
|-----------|-------|
| `docs/product/` | What Linden does — published, observable behavior only |
| `docs/architecture/` | How the system is built internally |
| `docs/project/` | How the project works |
| `docs/adr/` | Why a decision was made |

Roadmap and stage plans live in `plans/`, because they expire.

Do not restate a rule that exists elsewhere. Link to it. Duplicated standards diverge on the first edit.

## Agent workflow

Development is driven by the agent framework in `.agents/`, with the Workflow Orchestrator as the entrypoint. `.agents/` is canonical; `.github/agents`, `.github/instructions`, and `.github/prompts` are generated mirrors.

After changing anything under `.agents/`, regenerate the mirrors in the same PR:

```powershell
.\scripts\sync-agent-framework.ps1 -Provider github
```

## Privacy rules

Non-negotiable, and reviewed on every PR:

- No user content in logs, at any level, including debug.
- No user content in error messages returned to clients.
- Validate all input at the API boundary.
- No hardcoded secrets, tokens, or keys.

Treat every change as if an auditor is reading it for data leaks.
