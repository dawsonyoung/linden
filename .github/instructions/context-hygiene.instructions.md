---
applyTo: ".context/**"
---
# .context Hygiene

`.context` is an orientation index for agents, not a source of truth. See `.context/README.md`.

- Record only system boundaries, architectural patterns, conventions, and module relationships.
- Never record implementation details, code snippets, function signatures, or payload shapes. If a source edit would falsify the line, it does not belong here.
- Point at the authoritative file instead of reproducing its content. One line and a path.
- Code-level documentation stays beside the code. Architecture stays in `docs/`. `.context` links to both.
- When adding or moving a doc in `docs/` or `plans/`, update the index in `.context/README.md` in the same PR.
- Lessons capture process rules and failure patterns, not domain facts already documented elsewhere.
- Never write to `.context/lessons_learned.md` directly; append to `.context/draft_lessons.md` for human promotion.
- All `.context` changes go on a `workflow/` branch.
