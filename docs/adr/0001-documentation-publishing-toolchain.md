# ADR-0001: Documentation Publishing Toolchain

| Field | Value |
|-------|-------|
| Status | Proposed |
| Date | 2026-08-15 |
| Deciders | TODO |

## Context

`docs/product/` must be readable outside the repository — as a published HTML site — while remaining plain Markdown in version control. The toolchain must:

1. Run on Linux in CI without a manual step.
2. Work offline, consistent with the product's own commitments.
3. Add no language toolchain beyond Go and Node, which CI already installs.
4. Produce output that is usable from a local filesystem, not only from a web server.
5. Stay small enough that nobody needs to learn it to edit a page.

## Options

### A. mdBook

A single static Rust binary, purpose-built for specifications and technical books. Declarative navigation via `SUMMARY.md`, built-in client-side search, no configuration beyond a small TOML file.

- Distributed as a prebuilt binary; CI fetches a release rather than installing Rust.
- Output works from `file://`.
- Adds a third-party binary to the build.

### B. Hugo

A single static Go binary with a large feature surface and a theme ecosystem.

- Same distribution model as mdBook.
- Considerably more capability than a spec site needs; themes carry their own upgrade burden.

### C. Custom Go generator

A small program under `tools/docgen/` using a Markdown library.

- No external binary; `go run ./tools/docgen` and nothing else.
- Full control over output and zero theme churn.
- We own search, navigation, and cross-link checking — all of which mdBook provides for free.

## Decision

TODO: Not yet decided.

Recommendation is **Option A, mdBook**. It matches the shape of the problem, `SUMMARY.md` is already how the spec declares navigation, and the effort saved on search and cross-linking outweighs adding one fetched binary. Option C becomes preferable if the published output needs to be generated from the running binary itself, or if fetching a third-party artifact in CI is unacceptable.

## Consequences

If A is chosen:

- CI fetches a pinned mdBook release; the version is recorded here and in the workflow.
- `docs/product/spec/SUMMARY.md` remains the navigation source.
- `book.toml` is added at the repository root or under `docs/`.

If C is chosen:

- `SUMMARY.md` stays as-is; the generator reads it.
- Cross-link validation must be implemented, since nothing else will catch a broken relative link.

Either way:

- Output goes to `docs/.site/`, which is git-ignored.
- `make docs` builds; `make docs-serve` previews.
- A broken link or a page missing from `SUMMARY.md` fails the build.
