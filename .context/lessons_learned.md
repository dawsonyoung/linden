# System Memory — Learned Patterns

Human-promoted rules only. To add a rule, promote it from `.context/draft_lessons.md` via PR review.

## Domain Rules

[Rule-001]: Verify a change against the policy it establishes before requesting review.
Context: Any change that defines, extracts, or tightens a standard.
Negative: Do not assume authorship implies compliance. The change that defines a rule is the most likely place to violate it.

[Rule-002]: A gate is unproven until it has been observed failing.
Context: Adding or modifying any check — CI job, test, script, or build setting.
Negative: Do not treat a passing check as evidence that it works. The original validation job passed for weeks while being incapable of failing.

[Rule-003]: When a tool claims a guarantee, verify the tool enforces it before documenting the guarantee.
Context: Writing an ADR or spec that depends on third-party tooling behavior.
Negative: Do not describe a protection on the assumption that the chosen tool provides it. mdBook silently ignores unlisted pages and never checks links, both of which had been documented as build failures.

[Rule-004]: State current state at the coarsest useful granularity, or only where the work already updates it.
Context: Any status table, progress list, or "current state" section.
Negative: Do not add one that duplicates something tracked in plans/. A per-stage README table went stale within five PRs of being written to fix exactly that problem.

