# .context

Orientation layer for AI agents. **Not** a source of truth.

## Role

`.context` answers "where does this live and what rules apply" — never "what does the code do."
Every factual claim about behavior belongs in `docs/`, in the code, or beside the code.

Two file categories, governed by the same rules:

| Category | Files | Purpose |
|----------|-------|---------|
| Index | `README.md` (this file) | Map concerns to authoritative sources |
| Memory | `lessons_learned.md`, `draft_lessons.md` | Accumulated process rules from past failures |

## Rules

### 1. High-level and architecture-focused

Record system boundaries, architectural patterns, conventions, and module relationships.

Do not record implementation details, code snippets, function signatures, API payloads, or anything that changes when code changes. If an edit to a source file would falsify a line here, that line does not belong here.

### 2. Index, not manual

Point at the primary source. Do not reproduce it.

```
Good:  SSE event contract → docs/architecture/layer-interface-spec.md
Bad:   SSE emits `event: message` then `event: done` with fields ...
```

A `.context` entry should be one line and a path.

### 3. Single source of truth

| Content | Lives in |
|---------|----------|
| Function and type behavior | Godoc beside the code |
| Observable product behavior | `docs/product/spec/` |
| Capability proposals and rationale | `docs/product/prd/` |
| API request/response schemas | `docs/architecture/layer-interface-spec.md` |
| Architecture and layer rules | `docs/architecture/` |
| Process standards | `docs/project/`, `AGENTS.md` |
| Roadmap and stage plans | `plans/` |
| Agent behavior definitions | `.agents/` |
| Orientation and process lessons | `.context/` |

Duplication is drift. When content exists in `docs/`, `.context` links to it.

## Index

### Product

| Concern | Authoritative source |
|---------|---------------------|
| What Linden does, observable behavior | `docs/product/spec/` |
| Published navigation manifest | `docs/product/spec/SUMMARY.md` |
| Capability proposals and rationale | `docs/product/prd/` |
| Privacy commitments | `docs/product/spec/30-privacy-model.md` |
| Publishing toolchain decision | `docs/adr/0001-documentation-publishing-toolchain.md` |

### Architecture

| Concern | Authoritative source |
|---------|---------------------|
| What belongs in architecture docs | `docs/architecture/README.md` |
| Layer boundaries and dependency rules | `docs/architecture/layer-interface-spec.md` |
| Interface shapes, error taxonomy, SSE contract | `docs/architecture/layer-interface-spec.md` |
| Repository layout and conventions | `AGENTS.md` |

### Process

| Concern | Authoritative source |
|---------|---------------------|
| What belongs in each project doc | `docs/project/README.md` |
| Branching, commits, PR titles and bodies | `docs/project/branching-strategy.md` |
| Contract-first delivery policy | `docs/project/contract-first-delivery.md` |
| Agent workflow, gates, delegation | `docs/project/agentic-workflow-framework.md` |
| Product purpose and positioning | `plans/01-high-level-purpose-overview.md` |
| MVP scope and phase plan | `plans/02-repository-setup-and-mvp-plan.md` |
| Long-range roadmap and tracks | `plans/03-exhaustive-full-implementation-plan.md` |
| Implementation roadmap | `plans/04-implementation-sequence.md` |
| Current stage detail | `plans/05-stage-0-groundwork.md` |
| Agent definitions and tool scopes | `.agents/AGENTS.md` |

### Memory

| File | Lifecycle |
|------|-----------|
| `lessons_learned.md` | Human-promoted only. Read at session start. |
| `draft_lessons.md` | Agent-appended. Reviewed and promoted on `workflow/` branches. |

## Maintenance

- Adding a doc under `docs/` or `plans/`? Add its row to the index above.
- Moving or deleting a doc? Update the index in the same PR.
- An index row whose value changes over time must name the recurring step that refreshes it.
- Changes to `.context` go on a `workflow/` branch.
- A lesson that restates a doc is rejected in review — link instead.
