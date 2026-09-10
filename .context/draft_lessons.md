# Draft Lessons — Pending Human Review

Agent-generated candidates. Review and promote accepted rules to `.context/lessons_learned.md` during PR.
Reject or delete rules that are incorrect, too narrow, or already covered.

## Pending Rules

<!-- Agents append here during task completion when they encounter unexpected failures or retries. -->
<!-- Format: [Rule-NNN]: <actionable rule> | Context: <trigger> | Negative: <what to avoid> -->
<!-- Prefer a structural fix (agent definition, instruction, script) over a rule when one would enforce the behavior directly. -->

[Rule-005]: Parse and unpack error payloads at UI boundaries before surfacing to users or asserting in E2E tests.
Context: Surfacing HTTP/SSE failure responses in frontend components.
Negative: Do not pass raw response body text directly into Error constructors or alert elements, which leaks raw JSON syntax into user-facing banners and fails validation assertions.