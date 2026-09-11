# PRD-Driven Agentic Development Workflow

## Status
Active.

## Objective
Establish a formal, repeatable development process where high-level Product Requirements Documents (PRDs) serve as the primary behavioral anchor for autonomous and pair-programming AI agents. This workflow connects product intent directly to contract-first engineering and local quality gates.

---

## 1. Why PRDs Matter in an Agentic Codebase

Autonomous AI agents excel at code generation, refactoring, and test writing. However, without rigid boundary conditions, agents frequently suffer from common failure patterns:
1. **Scope Creep:** Inventing ad-hoc UI features, configuration flags, or endpoints not requested by the user.
2. **Context Drift:** Guessing at user requirements or making unverified assumptions about error recovery.
3. **Leaking Internals:** Conflating internal Go structs or database schemas with public, user-observable product behavior.
4. **Premature Blocking:** Halting progress because an exhaustive low-level spec is not yet written.

By placing structured PRDs in [`docs/product/prd/`](../product/prd/) at the top of the engineering hierarchy, Linden provides agents with a high-level **behavioral anchor**. PRDs define what the feature is, who it is for, what it must never do (non-goals), and what privacy constraints strictly bind it.

---

## 2. The 3-Tier Traceability Chain

Linden structures all development around a three-tier traceability chain that flows from user problem statements to executable machine verification:

```
┌────────────────────────────────────────────────────────────────────────┐
│             Tier 1: Product Requirements Document (PRD)                │
│                         "The Why & Intent"                             │
│   • Problem Statement & Target Persona                                 │
│   • Explicit Boundaries (Must-Have vs. Non-Goals)                      │
│   • Privacy, Safety, & Trust Guarantees                                │
│   ↳ docs/product/prd/0001-local-chat.md                                │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │ Bounding Context & Scope
                                    ▼
┌────────────────────────────────────────────────────────────────────────┐
│             Tier 2: Product Specifications (Living Spec)               │
│                     "The What & Observable Reality"                    │
│   • Present-tense description of real observable behavior              │
│   • Capabilities Matrix, User Flows, API Endpoints, Glossary           │
│   ↳ docs/product/spec/10-capabilities.md                               │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │ Implemented & Verified By
                                    ▼
┌────────────────────────────────────────────────────────────────────────┐
│          Tier 3: Technical Architecture & Interface Contracts          │
│                       "The How & Verification"                         │
│   • Go interface boundaries, package layering, error taxonomy          │
│   • Contract tests (validation/contracts/) & unit suites               │
│   ↳ docs/architecture/layer-interface-spec.md                          │
└────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Workflow Evolution: Current State vs. Future Vision

### Current State: Non-Blocking Reference
- **No Hard Blocking:** Agents must **never** be blocked by the lack of an existing, exhaustive specification in `docs/product/spec/`.
- **PRD Tree as Reference:** Agents consult the PRD tree (`docs/product/prd/`) to understand problem context, user jobs, non-goals, and privacy constraints.
- **Drafting Spec Extensions:** When a change touches user-observable behavior, the agent uses the PRD as reference to draft or extend the corresponding spec pages for human review alongside the implementation.

### Future State: Autonomous Spec Tree Synthesis
- **Automated Decomposition:** Agents will ingest high-level PRD problem statements and autonomously synthesize full, multi-layered specification trees (user flows, edge-case matrices, error taxonomies, and API schemas).
- **Pre-Implementation Spec Review:** The synthesized spec tree will undergo explicit human review and approval before contract tests and Go implementation begin.

---

## 4. The Step-by-Step Agentic Cycle

Every task executed by Linden workflow agents follows this disciplined lifecycle:

```
Intake ──► Spec (PRD Check) ──► Contract ──► Implement ──► Unit ──► Integrate ──► Review ──► Reflect
```

### Step 1: Intake & PRD Classification (Spec Gate)
- **Agent:** Workflow Orchestrator Agent (`workflow-orchestrator.agent.md`) using `review-product-spec.prompt.md`.
- **Procedure:**
  1. Determine whether the requested change introduces a capability a user could name or observe.
  2. Locate the corresponding PRD in `docs/product/prd/` (e.g. PRD-0001 for chat, PRD-0002 for document retrieval).
  3. Verify scope against the PRD's **Out of Scope (Non-Goals)** list to ensure no prohibited features are being added.
  4. Verify compliance with the PRD's **Privacy, Safety, & Trust Constraints**.
- **Exit Condition (Spec Gate):** Change is classified with explicit PRD reference and spec impact stated (`none`, `spec updated`, or `spec drafted`).

### Step 2: Specification Synthesis / Extension
- **Agent:** Document Agent (`document.agent.md`).
- **Procedure:**
  - If the capability introduces observable behavior, draft or update the corresponding page in `docs/product/spec/` (`10-capabilities.md`, `20-user-flows.md`, `30-privacy-model.md`, or `40-api-surface.md`).
  - Ensure spec language is in the **present tense** and strictly describes observable behavior without referencing internal Go struct names, file paths, or private functions.

### Step 3: Contract Test Authoring (Contract Gate)
- **Agent:** Contract Test Agent (`contract-tests.agent.md`) using `create-contract-tests.prompt.md`.
- **Procedure:**
  - Define or update the Go interface in `src/` or `docs/architecture/layer-interface-spec.md`.
  - Author contract tests in `validation/contracts/` using the `//go:build contracts` tag.
  - Contract tests must cover the happy path, error path, empty input, and timeout/cancellation.
- **Exit Condition (Contract Gate):** `go test -tags=contracts ./validation/contracts/...` exits 0.

### Step 4: Layer Implementation & Unit Testing
- **Agents:** Implementation Agent (`implementation.agent.md`) and Unit Test Agent (`unit-tests.agent.md`).
- **Procedure:**
  - Implement layer code in `src/` satisfying the contract tests.
  - Write table-driven unit tests (`Test_FunctionName_Scenario_Expected`) covering edge cases and boundary conditions.
- **Exit Condition (Unit Gate):** `go test -race ./src/...` exits 0.

### Step 5: Validation & Gate Enforcement
- **Agent:** Validation Agent (`validation.agent.md`) using `validation-gates` skill.
- **Procedure:**
  - Run all quality gates locally: Lint (`go vet`, `staticcheck`), Unit, Contract, Integration, Documentation (`mdbook build`), and Container build (`docker build`).
- **Exit Condition:** All CI gates pass locally before PR creation.

### Step 6: Post-Cycle Reflective Learning
- **Agent:** Reflective Learning Agent (`reflective-learning.agent.md`).
- **Procedure:**
  - Review the cycle for friction, retries, or spec ambiguities.
  - Append structured lessons to `.context/draft_lessons.md` using the format:
    `[Rule-NNN]: <actionable rule> | Context: <trigger> | Negative: <what to avoid>`.
  - Verify that the published documentation accurately reflects what was merged.

---

## 5. Agent Constraints & Guardrails

When operating under this PRD-driven workflow, agents must strictly adhere to the following rules:

1. **Never Invent Outside the PRD:** If a feature or capability is listed under "Out of Scope" in the relevant PRD, the agent must reject or de-scope it immediately.
2. **Never Mention Internals in Specs:** Product specifications and PRDs must be readable by a user or evaluator who has never opened the repository. Do not include Go package names, struct field definitions, or database table names.
3. **Keep PRDs Frozen Once Shipped:** When a capability is marked `Shipped`, its PRD is frozen as historical rationale. Subsequent behavioral changes are updated directly in `docs/product/spec/`.
4. **Honest Spec Currency:** An endpoint or capability listed as working that does not exist in the code is a critical defect. Specs must accurately reflect current system reality.
