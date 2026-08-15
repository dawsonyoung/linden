# Product Documentation

Human-readable product documentation. This directory is the published surface — it ships with the source and is rendered to HTML for reading outside the repository.

Everything here describes **what Linden does and why**. Nothing here describes how the repository is organized or how work gets done; that lives in `docs/project/`.

## Structure

| Path | Contains | Audience |
|------|----------|----------|
| `prd/` | Product requirement documents, one per capability | Product, engineering, reviewers |
| `spec/` | The product specification — the authoritative description of behavior | Everyone, including end users |

## PRD vs Spec

They answer different questions and have different lifecycles.

|  | PRD | Spec |
|--|-----|------|
| Question | Should we build this, and for whom? | What does the system do today? |
| Tense | Future — proposes | Present — describes |
| Lifecycle | Written once, then frozen with a status | Living; updated with every behavior change |
| Numbered | Yes, sequential and immutable | No; organized by topic |

A PRD is a decision record for a capability. Once accepted and shipped, it is not rewritten — the spec becomes the current truth and the PRD stands as the reason.

## Rules

1. Describe behavior, not implementation. No file paths, function names, or struct fields.
2. Interface-level detail belongs in `docs/architecture/layer-interface-spec.md`. Link to it.
3. Every user-visible behavior change updates `spec/` in the same PR that changes the behavior.
4. New capabilities require a PRD before implementation begins.
5. Write for a reader who has never seen the codebase.

## Publishing

Sources are Markdown. `make docs` renders this directory to HTML into `docs/.site/` (git-ignored). See `tools/docgen/README.md` and `docs/adr/0001-documentation-publishing-toolchain.md`.

Navigation is declared in `spec/SUMMARY.md`. A new page is not published until it appears there.
