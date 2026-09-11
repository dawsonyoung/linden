---
name: workflow-orchestration
description: "Step-by-step procedures for the Workflow Orchestrator. Use when planning a development cycle, classifying a task, building a delegation graph, or deciding which agent to invoke next."
---
# Workflow Orchestration

## Task Classification

Classify every incoming request before building a plan.

| Request type | First agent | Gate before advancing |
|-------------|-------------|----------------------|
| New interface | Design | Spec Gate |
| New contract tests | Contract Tests | Contract Gate |
| Implementation | Implementation | Contract Gate must already pass |
| Bug fix | Debug | Unit Gate |
| Docs update | Document | — |
| Review | Review | Review Gate |
| PR packaging | Git Ops | Review Gate |
| Post-cycle learning | Reflective Learning | — |

## Building a Delegation Graph

1. List all work types present in the request.
2. Order them: Design → Contract → Implement → Unit → Integrate → Functional → Review → Package → Reflect.
3. Remove stages not applicable to this request.
4. For each stage, write a delegation contract (Goal, Scope, Inputs, Output, Validation, Stop condition).

## Gate Enforcement

Never advance past a stage without confirming its gate:

- Spec Gate: Classify change against product spec, PRD tree (`docs/product/prd/`), and interface spec. Agents are not blocked by a lack of exhaustive specs; they use the PRD tree as reference and draft or extend the spec tree for review before implementation.
- Contract Gate: `go test -tags=contracts` exits 0.
- Unit Gate: `go test -race` exits 0.
- Integration Gate: `go test -tags=integration` exits 0.
- Runtime Gate: Docker build + health check pass.

## Escalation

If a specialist returns a failure or ambiguity:
1. Hand the failure output to Debug Agent (for gate failures).
2. Hand ambiguous design decisions to Design Agent.
3. If unresolved after two rounds, surface to human with a clear question.
