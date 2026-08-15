# Linden Agentic Workflow Framework

## Status
Proposed.

## Objective
Install a tool-agnostic multi-agent development workflow with one primary entrypoint agent that delegates to specialized autonomous activities while preserving auditability, test discipline, and strict Linux parity.

## Design Goals

1. Single operator entrypoint for consistency.
2. Deterministic delegation based on task type.
3. Least-privilege tool access per specialist.
4. Contract-first testing order enforced by workflow.
5. Reviewable PR-sized changes with stable history.
6. Strong observability of agent decisions and outputs.
7. Vendor-neutral configuration with .agents as the source of truth.

## Architectural Pattern

Use an orchestrator-worker model:

1. Workflow Orchestrator Agent is user-facing and manages plan state.
2. Specialist Agents execute bounded work in parallel or sequence.
3. Quality Gates evaluate outputs before merge.
4. Human Approval remains required for merge and any high-risk operation.

## Framework Placement and Portability

Canonical location:

1. Keep workflow framework definitions in .agents.
2. Treat .agents as the source of truth for agent behavior, prompts, and reusable skills.

Compatibility strategy:

1. Mirror or generate provider-specific files from .agents only when required by a tool runtime.
2. Never maintain parallel hand-edited copies in multiple provider folders.
3. If a provider requires its own path conventions, use a sync step that copies from .agents to provider-specific directories.

## Agent Topology

### Primary Entrypoint

Workflow Orchestrator Agent

Responsibilities:

1. Parse request and classify work type.
2. Build task graph with dependencies.
3. Delegate to specialists with scoped prompts.
4. Enforce testing sequence and policy gates.
5. Aggregate outputs into PR-ready changes.

Inputs:

1. User request
2. Repository state
3. Architecture and interface specifications

Outputs:

1. Execution plan
2. Ordered change sets
3. Validation evidence
4. PR summary and review checklist

### Specialist Agents

1. Document Agent
   - Updates plans, ADRs, architecture docs, changelogs.
2. Contract Test Agent
   - Adds or updates validation/contracts before implementation.
3. Implementation Agent
   - Implements layer code against approved interface contracts.
4. Unit Test Agent
   - Adds table-driven unit tests for new or changed behavior.
5. Integration Test Agent
   - Adds validation/integration once required modules are wired.
6. Validation Agent
   - Runs lint, unit, contract, integration, smoke, and container checks.
7. Review Agent
   - Performs risk-first code review with findings and coverage gaps.
8. Debug Agent
   - Reproduces failures and applies root-cause fixes.
9. Git Operations Agent
   - Prepares branch strategy, commit grouping, PR metadata.
10. Design Agent
   - Handles system/API/UI design and interface refinements.
11. Reflective Learning Agent
   - Evaluates completed workflow cycles.
   - Identifies inefficiencies, repeated failures, and gate bypass patterns.
   - Proposes concrete improvements to agent definitions, delegation contracts, and quality gates.
   - Persists lessons to `.context/lessons_learned.md` (promoted) and `.context/draft_lessons.md` (pending human review).
   - When proposing changes to .agents structure, workflow scripts, or CI configuration, generates a plan document in `plans/` and references it in the draft lesson entry.
   - Outputs findings as diffs to .agents definitions or as ADR candidates in docs/adr.

### Reflective Learning Memory Model

Layer 1: Developer-agent reflection (static, version-controlled).

1. When any agent encounters a build error, test failure, lint issue, or repeated retry, it appends a structured rule to `.context/draft_lessons.md` before closing the task.
2. Humans review draft lessons in PRs and promote accepted rules to `.context/lessons_learned.md`.
3. All agents read `.context/lessons_learned.md` at session start to avoid known failure patterns.
4. Rules include both positive patterns and explicit negative examples of what not to do.

Rule format:

```
[Rule-NNN]: <one-line actionable rule>
Context: <trigger or scenario>
Negative: <what to avoid>
```

Layer 2: Runtime reflection (structured, per-execution).

1. After any retry or error during execution, record a structured insight with trigger, error, insight text, and tags.
2. On future tasks with matching tags, retrieve relevant past insights and inject them into the agent context before planning.
3. Use tag-based retrieval to keep injection targeted and cheap.

### Reflective Learning Guardrails

1. Deduplication: before writing a new lesson, check for semantic overlap with existing rules. Merge rather than append when duplicate intent is detected.
2. Human audit gate: auto-generated lessons go to `.context/draft_lessons.md` first; only human-promoted rules enter `.context/lessons_learned.md`.
3. Negative memory: every lesson entry must include an explicit "what not to do" alongside the positive rule.
4. Scope limit: lessons describe observed behavior and constraints only. Proposed changes to agent structure, workflow scripts, or CI configuration require a plan document in `plans/` before any edits are made; the lesson entry must include a `Plan:` reference to that file.

Plan reference format:

```
[Rule-NNN]: <actionable rule>
Context: <trigger or scenario>
Negative: <what to avoid>
Plan: plans/<filename>.md
```

## Delegation Contract

Each delegated task must include:

1. Goal statement.
2. Scope boundaries.
3. Explicit inputs and files in scope.
4. Expected output schema.
5. Validation to run.
6. Stop conditions and escalation criteria.

Delegated tasks are rejected if required fields are missing.

## Workflow State Machine

1. Intake
   - Classify request: docs, tests, implementation, debug, review.
2. Plan
   - Build chunked PR plan aligned to interface-first policy.
3. Spec
   - Classify the change: does it alter user-visible behavior, add a capability, or change an interface?
   - New capability requires an accepted PRD in `docs/product/prd/` before work continues.
   - User-visible behavior change requires the affected `docs/product/spec/` pages updated in the same cycle.
   - Interface change requires `docs/architecture/layer-interface-spec.md` updated.
   - Record `Spec impact: none` explicitly when the change is internal only.
4. Contract
   - Add contract tests first.
5. Implement
   - Add implementation in smallest coherent unit.
6. Unit Validate
   - Add and run unit tests.
7. Integrate
   - Add integration tests when slice is wireable.
8. Functional Validate
   - Add smoke and security checks when vertical path exists.
9. Review
   - Run risk-first review and resolve findings.
10. Package
   - Prepare commits and PR body with evidence.
11. Reflect
   - Evaluate cycle outcomes against quality gate evidence.
   - Append new lessons to `.context/draft_lessons.md` for human review.
   - Inject promoted lessons from `.context/lessons_learned.md` into next cycle context.
   - If proposing changes to .agents, scripts, or CI: generate a plan document in `plans/` first, then include a `Plan:` reference in the lesson entry.
   - Surface agent definition improvement candidates as diffs or ADR stubs only after the plan document exists.

## Quality Gates

### Gate 1: Spec Gate

Required:

1. Change classified against product spec, PRD, and interface spec. A classification of "no impact" is stated, not assumed.
2. New capability has an accepted PRD before implementation begins.
3. User-visible behavior change has the corresponding `docs/product/spec/` pages updated in the same PR as the behavior.
4. Interface changes documented.
5. Dependency direction remains valid.
6. Any new spec page appears in `docs/product/SUMMARY.md`.7. **Currency check.** If this PR changes what the software can do, every statement of current state is re-read and corrected in the same PR: `README.md`, `docs/product/spec/00-overview.md`, and the endpoint table in `docs/product/spec/40-api-surface.md`.

A new endpoint, a changed response, or a new runtime dependency is a user-visible
change. "It was already listed as planned" does not satisfy item 3; a planned
entry must be changed to reflect that it shipped.
### Gate 2: Contract Gate

Required:

1. Contract tests exist before implementation merge.
2. Contract tests pass in Linux CI.

### Gate 3: Unit Gate

Required:

1. Unit tests for all changed behavior.
2. Edge cases covered: empty input, malformed input, timeout, unavailable dependency.

### Gate 4: Integration Gate

Required when modules are wired:

1. validation/integration coverage for affected path.

### Gate 5: Runtime Gate

Required:

1. Docker image build passes.
2. Container health endpoint passes.
3. SSE stream smoke check passes.

### Gate 6: Review Gate

Required:

1. Risk findings resolved or explicitly accepted.
2. Test gaps documented.

## Tool Access Policy

Principle: least privilege by default.

1. Document Agent: read, edit, search.
2. Contract Test Agent: read, edit, search, execute.
3. Implementation Agent: read, edit, search, execute.
4. Unit Test Agent: read, edit, search, execute.
5. Validation Agent: read, search, execute.
6. Review Agent: read, search.
7. Debug Agent: read, search, execute, edit.
8. Git Operations Agent: read, execute, search.
9. Design Agent: read, search, edit.
10. Reflective Learning Agent: read, search, edit.
    - edit is scoped to .context/draft_lessons.md, docs/adr, and plans/ only.

## Git and PR Strategy

1. One concern per PR.
2. Preferred PR size: 250 to 600 net lines changed.
3. Branch naming follows existing conventions.
4. Commits grouped by test-first then implementation.
5. PR template includes:
   - Intent
   - Interface changes
   - Tests added
   - Linux and Docker evidence
   - Risks and mitigations

## Test Sequencing Policy

Mandatory order:

1. Contract tests first.
2. Implementation second.
3. Unit tests with implementation.
4. Integration tests when multi-layer wiring exists.
5. Functional and smoke tests after vertical slice stability.

Merge is blocked if contract-first ordering is violated.

## Linux and Docker Policy

1. Linux CI is authoritative.
2. All required gates must pass on Linux runners.
3. Docker checks run in CI for runtime parity.
4. Windows support is additive and cannot bypass Linux failures.

## Observability and Auditability

1. Keep execution summaries per PR in repository docs.
2. Record gate results in PR checklist.
3. Preserve failure artifacts for debugging.
4. Require deterministic commands for validation steps.

## Recommended Repository Layout

1. .agents/
   - AGENTS.md
2. .agents/agents/
   - workflow-orchestrator.agent.md
   - document.agent.md
   - contract-tests.agent.md
   - implementation.agent.md
   - unit-tests.agent.md
   - integration-tests.agent.md
   - validation.agent.md
   - review.agent.md
   - debug.agent.md
   - git-ops.agent.md
   - design.agent.md
   - reflective-learning.agent.md
3. .agents/instructions/
   - contract-first.instructions.md
   - linux-docker-parity.instructions.md
   - pr-size-and-history.instructions.md
4. .agents/prompts/
   - create-contract-tests.prompt.md
   - implement-from-contract.prompt.md
   - run-validation-gates.prompt.md
5. .agents/skills/
   - workflow-orchestration/
   - contract-first-implementation/
   - validation-gates/
6. scripts/
   - sync-agent-framework.ps1
7. docs/architecture/
   - this document
   - layer and delivery specs
8. .context/
   - README.md           — role definition, hygiene rules, and index of authoritative sources
   - lessons_learned.md  — human-promoted rules injected at session start
   - draft_lessons.md    — agent-generated candidates pending human review

Provider-specific mirrors are optional outputs, not canonical inputs.

## Rollout Plan

### Phase 1: Framework bootstrap

1. Add Workflow Orchestrator Agent and three specialists:
   - Document
   - Contract Test
   - Implementation
2. Add contract-first and Linux parity instructions.
3. Add sync script to mirror .agents definitions to provider-required paths when needed.
4. Validate delegation on a docs-only pilot task.

### Phase 2: Quality and validation expansion

1. Add Unit Test, Validation, and Review agents.
2. Add CI-linked validation prompt flows.
3. Add CI check that provider-specific mirror files are in sync with .agents source.
4. Pilot on first API plus orchestrator contract PR.

### Phase 3: Full delivery chain

1. Add Debug, Git Ops, Design, and Reflective Learning agents.
2. Enforce all quality gates in PR template and workflow.
3. Run Reflective Learning Agent after each completed PR sequence.
4. Use findings to iterate on agent definitions and gate criteria before next cycle.

## Success Metrics

1. Median PR review time decreases.
2. Rework from late-found defects decreases.
3. Contract violations caught before integration increases.
4. Linux plus Docker gate pass rate stabilizes.
5. Time from request to merge becomes more predictable.

## Risks and Mitigations

1. Over-delegation risk
   - Mitigation: strict delegation contract and small task scopes.
2. Tool misuse risk
   - Mitigation: least-privilege tool sets and gate checks.
3. Process overhead risk
   - Mitigation: lightweight templates and bounded PR sizes.
4. Inconsistent outputs risk
   - Mitigation: required output schemas and review gate.

## Decision

Proceed with Phase 1 bootstrap first, then expand agent set only after first two PRs succeed with measurable quality improvements.
