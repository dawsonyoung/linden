---
applyTo: "**"
---
# Linux and Docker Parity

Linux CI is the authoritative gate. Windows support is additive.

- All commands used in Makefile, CI, and scripts must run on Linux without modification.
- Docker build and health check must pass before any runtime-affecting PR is merged.
- Never use shell constructs, path separators, or tools that are Windows-only in CI-facing scripts.
- If a command differs by platform, provide both forms and gate CI on the Linux form.
- Container smoke checks must include: `GET /health` returns 200, SSE stream delivers at least one `message` event.
