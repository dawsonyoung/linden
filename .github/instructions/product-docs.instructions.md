---
applyTo: "docs/product/**"
---
# Product Documentation

`docs/product/` is the published, human-readable surface. It ships with the source and is rendered to HTML.

- Describe observable behavior only. No file paths, function names, struct fields, or package names.
- Write for a reader who has never seen the codebase.
- The spec is present tense and describes what is real today. Planned behavior belongs in a PRD until it ships.
- A PRD is frozen once shipped. Do not rewrite it to match what was built; update the spec instead.
- Layer-internal detail belongs in `docs/architecture/layer-interface-spec.md`. Link to it rather than restating it.
- A new page must be listed in `docs/product/SUMMARY.md` or it is not published.
- Any change touching data handling, retention, or network egress must update `spec/30-privacy-model.md`, even if only to state that nothing changed.
- Every user-visible behavior change updates the spec in the same PR that changes the behavior.
