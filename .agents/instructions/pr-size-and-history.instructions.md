---
applyTo: "**"
---
# PR Size and History

- Every PR gets its own branch, created from latest `main` when work starts.
- Branch format: `<type>/<scope>-<short-description>` using `feat/`, `fix/`, `test/`, `refactor/`, `docs/`, `ci/`, `chore/`, `workflow/`.
- Agent framework and `.context/` learning file changes go on `workflow/` branches only, never mixed with product code.
- Target 250–600 net lines changed per PR. Max 900 for cross-cutting wiring PRs.
- One logical concern per PR. Split if review would take more than 30 minutes.
- Commit order within a PR: contract tests → implementation → unit tests.
- Commit messages use Conventional Commits: `feat:`, `fix:`, `test:`, `docs:`, `refactor:`.
- PR titles and bodies follow the standard in `docs/project/branching-strategy.md`. The body template is `.github/pull_request_template.md`; every section is required and unused sections say `N/A` with a reason.
- Full branching model: `docs/project/branching-strategy.md`.
