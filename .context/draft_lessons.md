# Draft Lessons — Pending Human Review

Agent-generated candidates. Review and promote accepted rules to `.context/lessons_learned.md` during PR.
Reject or delete rules that are incorrect, too narrow, or already covered.

## Pending Rules

<!-- Agents append here during task completion when they encounter unexpected failures or retries. -->
<!-- Format: [Rule-NNN]: <actionable rule> | Context: <trigger> | Negative: <what to avoid> -->
<!-- Prefer a structural fix (agent definition, instruction, script) over a rule when one would enforce the behavior directly. -->

[Rule-002]: Always verify local build/lint health before concluding a task | Context: Changing function signatures (e.g. NewService) or pushing code | Negative: Do not rely solely on GitHub Actions to catch unformatted code (gofmt) or missed unit test compilation errors (go vet). Run them locally.

[Rule-003]: Always run all CI validation gates locally (lint, unit, contract, integration, web) via manual commands or make targets prior to pushing to verify changes. | Context: Implementing new features or fixing bugs. | Negative: Do not push code assuming it works without local verification, and do not rely on remote CI pipelines as the primary debugging environment.
[Rule-004]: Actively monitor remote CI via 'gh pr checks --watch' or polling after pushing. | Context: Waiting for CI completion. | Negative: Do not stop and wait for the user to report CI failures; automatically poll and resolve them.
