---
applyTo: "**"
---
# PR Size and History

- Target 250–600 net lines changed per PR. Max 900 for cross-cutting wiring PRs.
- One logical concern per PR. Split if review would take more than 30 minutes.
- Commit order within a PR: contract tests → implementation → unit tests.
- Commit messages use Conventional Commits: `feat:`, `fix:`, `test:`, `docs:`, `refactor:`.
- PR body must include: Intent, Interface changes, Tests added, Linux/Docker evidence, Risks.
