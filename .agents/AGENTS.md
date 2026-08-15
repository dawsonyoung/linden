# Linden Agent Framework

This directory is the canonical source for all workflow agent definitions.
Provider-specific paths are generated mirrors — edit here, not there.

Only `.github/agents/`, `.github/instructions/`, and `.github/prompts/` are generated.
Other files under `.github/` are hand-authored and are not touched by the sync script.

## Agents

| File | Role | User-Invocable |
|------|------|---------------|
| agents/workflow-orchestrator.agent.md | Primary entrypoint; classifies, plans, delegates | Yes |
| agents/document.agent.md | Docs, ADRs, architecture, changelogs | No |
| agents/contract-tests.agent.md | Contract tests before implementation | No |
| agents/implementation.agent.md | Layer implementation against contracts | No |
| agents/unit-tests.agent.md | Table-driven unit tests | No |
| agents/integration-tests.agent.md | Cross-layer integration tests | No |
| agents/validation.agent.md | Runs all CI gates locally | No |
| agents/review.agent.md | Risk-first code review | No |
| agents/debug.agent.md | Reproduce and fix failures | No |
| agents/git-ops.agent.md | Branch, commit, PR preparation | No |
| agents/design.agent.md | Interface and API design | No |
| agents/reflective-learning.agent.md | Post-cycle learning and improvement proposals | No |

## Instructions

| File | Applies To |
|------|-----------|
| instructions/contract-first.instructions.md | All implementation and test work |
| instructions/linux-docker-parity.instructions.md | All build and CI work |
| instructions/pr-size-and-history.instructions.md | All PR preparation |
| instructions/context-hygiene.instructions.md | Changes under `.context/`, `docs/`, `plans/` |
| instructions/product-docs.instructions.md | Changes under `docs/product/` |

## Prompts

| File | Purpose |
|------|---------|
| prompts/create-contract-tests.prompt.md | Guided contract test authoring |
| prompts/implement-from-contract.prompt.md | Guided implementation from contract |
| prompts/run-validation-gates.prompt.md | Run and interpret all quality gates |
| prompts/review-product-spec.prompt.md | Classify a change against spec and PRDs at the Spec step |

## Skills

| Folder | Purpose |
|--------|---------|
| skills/workflow-orchestration/ | Orchestrator decision and delegation procedures |
| skills/contract-first-implementation/ | Step-by-step contract-first dev workflow |
| skills/validation-gates/ | Gate execution and pass/fail interpretation |

## Self-Correction

- `.context/` is an orientation index, not a source of truth. See `.context/README.md`.
- Read `.context/lessons_learned.md` at session start.
- Append to `.context/draft_lessons.md` when tasks encounter failures or retries.
- Lessons link to authoritative docs; they never restate them.
- Structural change proposals go to `plans/` first; lesson entries reference the plan file.
