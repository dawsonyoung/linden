# Project

How this project works: the practices that govern the work, and the rationale behind how the repository itself is arranged.

## What belongs here

Four tests, in order. A document belongs in `docs/project/` if it fails the first three and passes the fourth.

| Question | If yes |
|----------|--------|
| Could a user observe this without reading code? | `docs/product/` |
| Does it describe how the system is built internally? | `docs/architecture/` |
| Is it a single decision with alternatives and consequences? | `docs/adr/` |
| Does it describe how the project works? | here |

Concretely, this directory owns:

- Working practices — branching, review, delivery policy, quality gates
- Tooling that governs the work, including the agent workflow
- The rationale for repository structure: why a directory exists and what it is for
- Conventions that outlive any single change

## What does not belong here

- Roadmap, stage plans, and sequencing. Those are `plans/`, because they expire.
- A single technical decision with alternatives. That is an ADR.
- Anything a user of Linden would read.

The distinction from `plans/` is durability. A document here describes how work is done in general; a plan describes what is being done now and stops being true once it is finished.

## Contents

| Document | Purpose |
|----------|---------|
| `branching-strategy.md` | Branching model, naming, commits, PR titles and bodies |
| `contract-first-delivery.md` | Test sequencing policy and PR-by-PR delivery approach |
| `agentic-workflow-framework.md` | Agent topology, delegation contracts, quality gates |

## Admission test

Before adding a document here, state which of the four questions above it passes and why the preceding three do not apply. If that is hard to answer, the document probably belongs somewhere else — or does not need to exist.
