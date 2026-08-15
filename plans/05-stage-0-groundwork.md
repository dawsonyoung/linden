# Stage 0: Groundwork

## Status
Ready to implement.

## Purpose

Make the repository truthful and buildable. Today the Makefile and CI reference source files, manifests, and scripts that do not exist, and the CI validation job masks missing suites with pass-through echoes so it cannot fail. Nothing downstream can be trusted until this is corrected.

Branching rules for all work below: `docs/process/branching-strategy.md`.

## Contract-First Exemption

Stage 0 is exempt from the contract-first ordering rule. It introduces no cross-layer interface — only process bootstrap, build tooling, and CI. `/health` and `/version` are self-contained infrastructure endpoints with no orchestrator dependency.

The exemption ends at Stage A.1. Every interface method from that point requires a contract test first.

## Shared Decisions

| Decision | Value |
|----------|-------|
| Go module path | `github.com/dawsonyoung/linden` |
| `go.mod` location | `src/` (all tooling already does `cd src`) |
| Go version | 1.22 |
| Node version | 20 |
| Default listen address | `:8080` |
| Logging | `log/slog`, JSON handler, no user data at any level |
| Config source | Environment variables only in Stage 0 |

Environment variables introduced:

| Variable | Default | Purpose |
|----------|---------|---------|
| `LINDEN_ADDR` | `:8080` | Listen address |
| `LINDEN_LOG_LEVEL` | `info` | slog level |
| `OLLAMA_URL` | `http://localhost:11434` | Reserved for Stage A |

## Sequencing

```
0.1 chore/repo-hygiene
     ↓
0.2 feat/cmd-server-bootstrap
     ↓
     ├──→ 0.3 ci/linux-gate-honesty
     └──→ 0.4 ci/docker-linux-parity
                ↓
           0.5 chore/cross-platform-make
```

0.3 and 0.4 are independent of each other and may run in parallel once 0.2 is merged.

---

## 0.1 `chore/repo-hygiene`

Intent: satisfy Phase 0 baseline hygiene and stop the README from claiming capability that does not exist.

Agent: Document Agent. No other specialist required.

### Tasks

1. `CONTRIBUTING.md` — prerequisites, setup, branch and commit conventions, test-first rule, PR checklist, required status checks, Windows-via-WSL note.
2. `SECURITY.md` — disclosure process, response expectation, no-user-data-in-logs rule, no secrets in repo.
3. `.github/CODEOWNERS` — default owner for all paths; `.agents/**` and `.context/**` called out explicitly.
4. `README.md` — replace the false status line, add support matrix and MVP checklist linked to the MVP plan.

`.github/pull_request_template.md` already exists. CONTRIBUTING links the PR standard in `docs/process/branching-strategy.md` rather than restating it.

### Acceptance Criteria

1. README status reflects actual repository state: scaffold plus agent framework, no runnable server yet.
2. A new contributor can read CONTRIBUTING and know the branch name, commit order, and gate expectations.
3. CONTRIBUTING links `docs/process/branching-strategy.md` rather than restating it.
4. PR template renders on a test PR.

### Gate

Lint only. No code changes in this PR.

### Done When

`git diff --stat` shows documentation files only, and README contains no claim that is not currently true.

Estimated size: 200–300 lines.

---

## 0.2 `feat/cmd-server-bootstrap`

Intent: produce the first binary that builds and boots, so every later gate has something real to run against.

Agents: Implementation Agent, then Unit Test Agent, then Validation Agent.

### Tasks

| File | Contents |
|------|----------|
| `src/go.mod` | Module `github.com/dawsonyoung/linden`, go 1.22, stdlib only |
| `src/config/config.go` | Env parsing with defaults and validation |
| `src/cmd/main.go` | Config load, logger init, server construction, listen, graceful shutdown |
| `src/api/server.go` | `http.ServeMux` router, timeouts, exported `NewServer` |
| `src/api/health.go` | `GET /health` → 200 `{"status":"ok"}` |
| `src/api/version.go` | `GET /version` → build version, commit, Go version |
| `src/api/middleware.go` | Request ID generation and injection, structured access logging |
| `src/api/*_test.go` | Unit tests for handlers and middleware |
| `src/config/config_test.go` | Unit tests for env parsing and defaults |

### Behavior Requirements

1. Server sets `ReadHeaderTimeout` and `WriteTimeout`. No unbounded reads.
2. Request ID is generated if absent, echoed in `X-Request-ID`, and attached to every log line.
3. Graceful shutdown drains in-flight requests with a bounded timeout on SIGINT and SIGTERM.
4. Version values are injected at build time via `-ldflags`, defaulting to `dev` when unset.
5. Access logs contain method, path, status, duration, request ID. Never query strings, headers, or bodies.

### Acceptance Criteria

```sh
cd src && go build ./cmd/          # exits 0
cd src && go vet ./...             # exits 0
cd src && go test -race ./...      # exits 0
```

Manual verification:

```sh
./bin/linden &
curl -i http://localhost:8080/health    # 200, {"status":"ok"}, X-Request-ID present
curl -s http://localhost:8080/version   # valid JSON
kill -INT %1                            # clean shutdown, no dropped request
```

### Unit Test Coverage Required

- Health handler: 200 and correct body.
- Version handler: valid JSON, `dev` default when ldflags unset.
- Middleware: generates request ID when absent, preserves when supplied, echoes in response header.
- Config: defaults applied, invalid values rejected, all three env vars parsed.

### Gate

Lint + Unit.

### Risks

`src/web` does not exist yet, so `make dev` and `make build` still fail on `build-web`. Do not add a static file handler in this PR — web wiring is out of scope until C.3.

Estimated size: 400–550 lines.

---

## 0.3 `ci/linux-gate-honesty`

Intent: make CI fail when something is broken. The validation job currently cannot fail.

Agents: Implementation Agent (workflow YAML), Validation Agent (verification).

### Tasks

1. Delete `2>/dev/null || echo "No contract tests yet"` from both validation steps in `.github/workflows/ci.yml`.
2. Split monolithic jobs into one job per gate: `lint`, `unit`, `contract`, `integration`, `web`. Failures must be attributable to a single gate.
3. Where a suite has no packages yet, guard with an explicit visible step condition. Never swallow a non-zero exit.
4. Add `concurrency` with in-progress cancellation for superseded pushes.
5. Add `timeout-minutes` per job.
6. Add Go build/module cache and npm cache.
7. Keep `ubuntu-latest` as the only runner.
8. Document required status checks for branch protection in `CONTRIBUTING.md`.

### Acceptance Criteria

1. A deliberately broken test causes a red check. Verify by pushing a temporary failing commit, confirming red, then reverting within the same PR. Link both runs in the PR body.
2. No step in the workflow can exit non-zero and still report success.
3. Job names map one-to-one onto the quality gates in `docs/process/agentic-workflow-framework.md`.

### Gate

CI green on the PR itself, plus the recorded red run from the deliberate breakage.

### Note

`contract` and `integration` jobs will legitimately have zero packages until A.1 and B.2. Use a guarded step that reports "no packages yet" as a visible skip, not as a passing test run.

Estimated size: 150–250 lines.

---

## 0.4 `ci/docker-linux-parity`

Intent: get a Linux container running early, since strict parity is the standing requirement.

Agents: Implementation Agent, then Validation Agent.

### Tasks

| File | Contents |
|------|----------|
| `Dockerfile` | Multi-stage build, see below |
| `.dockerignore` | Exclude `.git`, `node_modules`, `bin`, `plans`, `docs`, `.agents`, `.context` |
| `docker-compose.yml` | `linden` service plus optional `ollama` service on a shared network |
| `.github/workflows/ci.yml` | Add `docker` job: build, run, poll `/health`, assert 200, print logs on failure |
| `Makefile` | Add `docker-build` and `docker-smoke` targets |

Dockerfile stages:

1. `node:20-alpine` — build web assets. Keep the stage but tolerate absent input in Stage 0.
2. `golang:1.22-alpine` — build static binary, `CGO_ENABLED=0`, `-ldflags` version injection.
3. `gcr.io/distroless/static-debian12` — copy binary only.

### Runtime Requirements

1. Container runs as a non-root user.
2. `HEALTHCHECK` invokes `/health`.
3. No secrets baked into the image. No build args carrying credentials.
4. Final image contains the binary and nothing else — no compiler, no shell.

### Acceptance Criteria

```sh
docker build -t linden:ci .                              # exits 0 on Linux CI
docker run -d -p 8080:8080 --name linden-smoke linden:ci
curl -sf http://localhost:8080/health                    # 200 within 10s
docker inspect linden-smoke --format '{{.Config.User}}'  # non-root
```

Record baseline image size in the PR body.

### Gate

Docker build + container health check.

### Risk

The Linux/Docker parity instruction requires an SSE smoke check, but no chat endpoint exists yet. Defer the SSE portion of container smoke to C.4 and state the deferral explicitly in the PR body.

Estimated size: 200–300 lines.

---

## 0.5 `chore/cross-platform-make`

Intent: remove Windows-hostile constructs from developer entry points while keeping Linux authoritative.

Agent: Implementation Agent, then Validation Agent.

### Tasks

1. Replace `rm -rf bin/ src/web/build/` in `clean` with a portable form.
2. `validate-smoke` calls `./run.sh`, which does not exist. Drop the target rather than ship one that cannot run. C.4 reintroduces it.
3. Guard `build-web`, `test`, and `lint` web steps so they no-op cleanly while `src/web` has no `package.json`.
4. Confirm `docker-build` and `docker-smoke` targets from 0.4 are present.
5. Add `scripts/make.ps1` wrapper, or document in CONTRIBUTING that Windows contributors run targets through WSL or Git Bash.

### Acceptance Criteria

1. Every Makefile target either succeeds or fails with an actionable message on a clean Linux checkout.
2. No target references a file that does not exist.
3. Windows path documented and verified once manually.
4. CI continues to gate on the Linux form only.

### Verification

Run every target on a clean Linux checkout and record results:

```sh
make build; make test; make lint; make clean; make docker-build; make docker-smoke
```

### Gate

Lint, plus manual verification of each target on Linux.

Estimated size: 100–200 lines.

---

## Stage 0 Exit Criteria

All five branches merged, and:

1. `git clone` → `make build` produces `bin/linden` on Linux with no manual steps.
2. `docker build` and container `/health` both pass in CI.
3. Every CI job can fail. Verified by deliberate breakage during 0.3.
4. README describes the repository as it actually is.
5. No Makefile target references a nonexistent file.
6. Contract, integration, and smoke jobs exist and report visible skips rather than false passes.

Only then does Stage A begin.

---

## 0.6 `ci/docgen-publishing`

Intent: make `docs/product/` publishable so the spec can be read outside the repository.

Agents: Design Agent (resolve the ADR), then Implementation Agent, then Validation Agent.
### Tasks

1. Add `docs/book.toml` configuring mdBook with `docs/product/` as the book source and `docs/.site/` as output.
2. Decide how `docs/product/prd/` is rendered — listed in `SUMMARY.md` as a section, or built as a second book. Record the choice in the ADR consequences.
3. Add the mdBook fetch to CI with a pinned version.
4. Replace the placeholder `docs` and `docs-serve` Makefile targets with real ones.
5. Add a CI job that builds the docs and fails on a broken link or a page missing from `SUMMARY.md`.

### Acceptance Criteria

1. `make docs` renders `docs/product/` to `docs/.site/` on a clean Linux checkout.
2. Output is browsable from `file://` with no server.
3. A page added without a `SUMMARY.md` entry fails the build. Verify deliberately.
4. A broken relative link fails the build. Verify deliberately.

### Gate

Docs build job green, plus the two deliberate failure verifications recorded in the PR body.

Depends on: 0.3 merged, so the new job slots into an honest CI.

Estimated size: 250–400 lines.


## Post-Stage Action

Run the Reflective Learning Agent against the completed stage. Promote accepted lessons on a `workflow/promote-lessons-<date>` branch.
