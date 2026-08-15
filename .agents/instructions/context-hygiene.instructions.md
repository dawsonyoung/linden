---
applyTo: ".context/**,docs/**,plans/**"
---
# Documentation and Context Hygiene

`.context` is an orientation index for agents, not a source of truth. See `.context/README.md`.

## Writing to `.context`

- Record only system boundaries, architectural patterns, conventions, and module relationships.
- Never record implementation details, code snippets, function signatures, or payload shapes. If a source edit would falsify the line, it does not belong here.
- Point at the authoritative file instead of reproducing its content. One line and a path.
- Lessons capture process rules and failure patterns, not domain facts already documented elsewhere.
- Never write to `.context/lessons_learned.md` directly; append to `.context/draft_lessons.md` for human promotion.

## Adding or moving documentation

- Adding, moving, or deleting a file under `docs/` or `plans/` requires updating the index in `.context/README.md` in the same change.
- Changing repository layout requires updating the directory listings in `llm.txt` and `AGENTS.md`, not just the index.
- Before stating a rule, search for an existing statement of it. If one exists, link to it instead of restating it.
- When extracting a standard into its own document, replace every prior statement of that standard with a link in the same change.

## Scope

- Code-level documentation stays beside the code. Architecture stays in `docs/`. `.context` links to both.
- All `.context` changes go on a `workflow/` branch.
