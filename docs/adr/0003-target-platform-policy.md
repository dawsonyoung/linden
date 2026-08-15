# ADR-0003: Target Platform Policy

| Field | Value |
|-------|-------|
| Status | Accepted |
| Date | 2026-08-15 |
| Deciders | dawsonyoung |

## Context

Development happens on Windows. Deployment is intended to be Linux, and Windows
support is expected to be dropped entirely. In the short term the software must
run on a Windows machine.

Three separate questions have been getting conflated:

1. What operating system do developers work on?
2. What operating system does the software run on?
3. What platforms do we commit to supporting?

Treating these as one question produces half-finished platform support that
nobody asked for and nobody removes.

The cost of ambiguity is already visible. The Windows development host cannot
run `go test -race` without a C compiler, cannot run a Docker daemon, and
produces CRLF checkouts that made a `gofmt` gate fail locally while passing in
CI. Stage 0.5 is currently scoped to make the Makefile work on Windows, which
is effort spent on a platform that is planned for removal.

The forcing decision is near. Stage B.4 introduces persistence and needs a data
directory. `%APPDATA%` versus `$XDG_DATA_HOME` behind a `runtime.GOOS` switch is
the kind of branch that outlives the platform that motivated it.

## Options

### A. Support Windows and Linux as peer platforms

Native binaries for both, platform-conditional code where they differ, CI on
both.

Doubles the test surface and puts platform divergence into the source
permanently, for a platform slated for removal.

### B. Linux only, immediately

Drop Windows entirely now.

Blocks the short-term need to deploy on a Windows machine.

### C. Linux is the only deployment target; Windows is a development host, and
short-term Windows deployment runs the Linux container

Platform divergence never enters the source, because the container is the
compatibility layer. Windows remains viable as both a development host and an
interim deployment host without becoming a supported target.

## Decision

**Option C.**

1. **Linux is the only deployment target.** All runtime behavior is specified,
   tested, and supported on Linux.
2. **Windows is a supported development host**, not a target. It is expected to
   remain one indefinitely; that is an editor and toolchain concern, not a
   product commitment.
3. **Short-term Windows deployment runs the Linux container**, with Ollama
   native on the host for GPU access. See ADR-0002.
4. **No platform-conditional code.** No `runtime.GOOS` branches, no `_windows.go`
   build-tagged files, no Windows service integration, no registry or
   `%APPDATA%` handling. Filesystem paths follow Linux conventions.
5. **No Windows CI job.** Adding one converts a free option into an obligation.

A native Windows binary remains a one-command fallback, because the code is pure
standard library with no syscalls. `GOOS=windows go build` works today and will
keep working as long as rule 4 holds. That is an unsupported convenience
artifact, not a platform.

## Consequences

Immediate:

1. Stage 0.5 is rescoped. The Makefile stays POSIX-only; no PowerShell wrapper
   is written. Windows contributors run targets inside the devcontainer.
2. A devcontainer becomes the recommended development environment, making the
   local toolchain identical to CI. This closes the `-race`, Docker, and line
   ending gaps in one step.
3. `scripts/setup.sh` becomes the canonical setup script. `scripts/setup.ps1`
   remains as a convenience for developers not using the devcontainer, and is
   deleted when the devcontainer is the only documented path.
4. Data directories, when introduced at B.4, follow XDG and FHS conventions with
   no Windows branch.

Ongoing:

1. Case sensitivity bugs are caught only by CI, since Windows filesystems are
   case-insensitive. This is an accepted risk, reduced by developing inside the
   devcontainer.
2. `docs/product/` describes Linux behavior. It does not document Windows
   installation.

Sunset criteria for Windows as a development host — none currently apply, and
support ends when any of these becomes true:

1. A dependency is required that cannot be built or run on Windows, and no
   devcontainer workaround exists.
2. Maintaining `scripts/setup.ps1` and `scripts/sync-agent-framework.ps1`
   requires effort disproportionate to their use.
3. All contributors use the devcontainer, making the native Windows path
   untested in practice.

Removal, when it happens, is a documentation and script deletion. It does not
touch the source, because rule 4 keeps the source platform-neutral.

## Notes

This ADR governs platform commitments. The role of the container image in
distribution is ADR-0002.
