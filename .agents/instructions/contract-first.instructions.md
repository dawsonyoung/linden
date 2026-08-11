---
applyTo: "validation/**,src/**"
---
# Contract-First Development

Contract tests in `validation/contracts/` must exist and compile before any implementation is merged.

- Write the contract test first. Commit it. Then implement.
- Contract tests use build tag `//go:build contracts` and assert interface shape, not internals.
- Every new interface method requires at minimum: happy path, error path, and empty input case.
- Implementation PRs that lack a corresponding contract test are not ready for review.
