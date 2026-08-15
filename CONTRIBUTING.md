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

**Recommended: the devcontainer.** Open the repository in VS Code and reopen in
container. The toolchain then matches CI exactly, which means `go test -race`,
Docker, and line endings all behave the same locally as they do in the gates.

Without a devcontainer:

```sh
./scripts/setup.sh              # check the toolchain
./scripts/setup.sh --pull-model # also fetch the model used from Stage A.3
```

`scripts/setup.ps1` exists as a Windows convenience and is not the canonical
path. See `docs/adr/0003-target-platform-policy.md`.

## Platforms

Linux is the only deployment target. Windows and macOS are development hosts.

The Makefile and CI are POSIX-only by policy. Windows contributors run targets
inside the devcontainer rather than through a wrapper. There is deliberately no
Windows CI job, and no platform-conditional code is accepted in `src/`.

Short-term deployment on a Windows machine runs the Linux container with Ollama
native on the host. See `docs/adr/0002-role-of-docker-in-distribution.md`.

## Reaching Ollama

The server reads `OLLAMA_URL`. Two topologies are supported; both are the same
to the server, which only makes HTTP requests to that address.

**Ollama on the host, Linden in a container.** The default, and the only option
on Windows and macOS, where GPU passthrough into a Linux container is
impractical. The compose file resolves `host.docker.internal` for this, which on
Linux requires the `host-gateway` mapping it already declares.

```sh
OLLAMA_URL=http://host.docker.internal:11434 docker compose up --build
```

**Ollama as a sibling container.** Linux only, and preferable there: GPU access
works through `nvidia-container-toolkit` without a VM boundary.

```sh
OLLAMA_URL=http://ollama:11434 docker compose --profile ollama up --build
```

Ollama is never part of the Linden image. Models are multi-gigabyte and have a
different lifecycle from the binary, so they live in their own volume.

Running the server directly on the host needs neither: the default
`http://localhost:11434` is correct.

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
| `docs` | Product docs render; no unpublished page or broken link |
| `docker` | Image builds, runs non-root, serves health, shuts down gracefully |

`contract`, `integration`, and `web` report a visible skip until the
corresponding suite or project exists. A skip is recorded as a notice in the
job log; it never reports a passing test run.

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
