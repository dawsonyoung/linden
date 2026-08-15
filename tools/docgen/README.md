# docgen

Renders `docs/product/` to a browsable HTML site.

> **Status:** Not yet implemented. Toolchain is **mdBook**, decided in `docs/adr/0001-documentation-publishing-toolchain.md`. Implementation is Stage 0.6 in `plans/05-stage-0-groundwork.md`.

mdBook is a prebuilt binary fetched in CI — no Rust toolchain is installed. This directory holds the configuration and any wrapper scripting needed to invoke it; there is no generator source to write.

## Intended usage

```sh
make docs         # render to docs/.site/
make docs-serve   # render and serve locally with reload
```

## Requirements the implementation must satisfy

1. Input is `docs/product/`; navigation comes from `docs/product/spec/SUMMARY.md`.
2. Output is `docs/.site/`, git-ignored, and browsable from the local filesystem without a server.
3. A page present on disk but absent from `SUMMARY.md` fails the build. Unpublished pages are invisible, and silent invisibility is worse than a build error.
4. A broken relative link fails the build.
5. Runs on Linux in CI with no interactive step.
6. No network access required at render time.

## Not in scope

- Rendering `docs/process/`, `docs/architecture/`, or `plans/`. Those are contributor-facing and read in the repository.
- Publishing or deployment. Rendering only; where the output goes is a separate decision.
