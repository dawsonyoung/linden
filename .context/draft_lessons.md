# Draft Lessons — Pending Human Review

Agent-generated candidates. Review and promote accepted rules to `.context/lessons_learned.md` during PR.
Reject or delete rules that are incorrect, too narrow, or already covered.

## Pending Rules

<!-- Agents append here during task completion when they encounter unexpected failures or retries. -->
<!-- Format: [Rule-NNN]: <actionable rule> | Context: <trigger> | Negative: <what to avoid> -->
<!-- Prefer a structural fix (agent definition, instruction, script) over a rule when one would enforce the behavior directly. -->

[Rule-002]: Always verify local build/lint health before concluding a task | Context: Changing function signatures (e.g. NewService) or pushing code | Negative: Do not rely solely on GitHub Actions to catch unformatted code (gofmt) or missed unit test compilation errors (go vet). Run them locally.
