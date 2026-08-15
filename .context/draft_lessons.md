# Draft Lessons — Pending Human Review

Agent-generated candidates. Review and promote accepted rules to `.context/lessons_learned.md` during PR.
Reject or delete rules that are incorrect, too narrow, or already covered.

## Pending Rules

[Rule-001]: When extracting a standard into its own document, grep the repo for existing statements of that standard and replace each with a link.
Context: Branching rules were moved to `docs/branching-strategy.md` while `docs/architecture/02` kept its own prefix list, which was already missing `workflow/`.
Negative: Do not leave the original copy in place "for convenience." Two copies of a rule diverge on the first edit, and the stale one is indistinguishable from the current one.

[Rule-002]: When adding, moving, or deleting a file under `docs/` or `plans/`, update the `.context/README.md` index in the same change.
Context: The index shipped in the same PR that defined it, already missing `plans/01` through `plans/03`.
Negative: Do not assume an index is complete because it was just authored. Enumerate the directory and diff against the index.

[Rule-003]: An index entry whose value changes over time needs a named owner step that refreshes it.
Context: The "Current stage detail" row points at one stage plan and becomes wrong the moment the next stage begins.
Negative: Do not add time-varying pointers without stating which recurring step updates them.

[Rule-004]: When repository layout changes, update the directory listings in `llm.txt` and `AGENTS.md`, not just the doc index.
Context: `docs/architecture/` and `docs/branching-strategy.md` existed for several changes while `llm.txt` still listed only `docs/adr/`.
Negative: Do not treat the `.context` index as the sole map. Orientation files are read first and are trusted as complete.

[Rule-005]: Review any PR that establishes a policy against that policy itself before requesting review.
Context: This PR codified single-source-of-truth while shipping two duplicated rule sets.
Negative: Do not assume authorship implies compliance. The document defining a rule is the most likely place to violate it.
