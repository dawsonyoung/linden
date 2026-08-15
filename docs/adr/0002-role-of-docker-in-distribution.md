# ADR-0002: Role of Docker in Distribution

| Field | Value |
|-------|-------|
| Status | Accepted |
| Date | 2026-08-15 |
| Deciders | dawsonyoung |

## Context

Stage 0.4 introduces a Dockerfile, a compose file, and a CI job that builds and
smoke-tests an image. Before more is built on top of it, the image's role needs
to be settled: is it a test fixture, or a product artifact?

The question is not cosmetic. It determines whether multi-architecture builds,
registry publishing, signing, SBOMs, and a documented data-volume contract are
requirements or waste.

Forces:

1. `plans/03` Track E1 describes distribution as cross-platform installers,
   portable packages, signed binaries, and update channels with rollback. It
   does not mention a container registry.
2. Track F describes the appliance as a device image with first-boot
   provisioning, not a container runtime.
3. The MVP release criterion is that an install script succeeds on Windows.
4. The product premise is a non-technical user plugging a device into a home
   network. That user does not install Docker.
5. `plans/01` names "builders and tinkerers who want local AI infrastructure"
   as a target audience. That audience expects a compose file, not an installer.
6. Linux CI is the authoritative gate, and the developer machine is Windows.
   Some behavior — SIGTERM handling, non-root execution, static linking — cannot
   be verified anywhere else.

## Options

### A. Verification vehicle only

The image exists to prove Linux runtime behavior in CI. It is never published.

Cheapest. Forfeits a distribution channel that a named audience would prefer,
and risks the image rotting into a CI-only artifact that diverges from how the
software actually runs.

### B. Primary distribution for all users

Docker becomes the supported way to run Linden.

Contradicts the product premise and Track E. Would require non-technical users
to install a container runtime, and would make the Ollama topology a user
problem rather than an implementation detail.

### C. Verification now, secondary distribution later

The image is a CI verification vehicle today. It is built to production
standards where those standards are free, and may later be published for
technical users. It is never the primary consumer path and never the appliance
mechanism.

## Decision

**Option C.**

The image is a verification vehicle today. It is not a supported distribution
artifact, and no support commitment attaches to it until it is published.

Production-grade properties that cost nothing are kept from the start:
distroless base, non-root user, static binary, self-probing healthcheck. These
are good CI hygiene regardless, and they are the expensive part to retrofit.

Properties that only matter for publication are deferred to Track E:
multi-architecture builds, registry publishing, signing, and SBOM generation.

GitHub Container Registry (`ghcr.io`) is the intended target if and when
publication happens. It authenticates in Actions with `GITHUB_TOKEN` under
`packages: write`, supports OCI multi-arch manifests, and inherits repository
visibility, so adopting it later is a small change rather than a new system.

## Consequences

Immediate:

1. The image stays unpublished. CI builds and smoke-tests it; nothing pushes it.
2. `docker-compose.yml` is a development convenience, not a supported topology.
3. Image size is recorded as a CI notice to establish a baseline, not enforced
   as a budget.

Deferred, with the trigger that forces each:

| Concern | Trigger |
|---------|---------|
| arm64 multi-arch build | Appliance hardware selection in Track F1, or the first user request |
| GHCR publishing and versioned tags | Track E1 packaging work |
| Image signing and provenance attestation | First published tag |
| Vulnerability scanning in CI | First published tag |

Forced independently of this decision:

1. **Data volume contract.** Stage B.4 introduces persistence. A container
   without a declared data directory and volume loses conversations on restart.
   The data directory must be decided at B.4 whether or not the image is ever
   published.
2. **Ollama topology.** Stage A.3 introduces a real dependency on a model
   backend. The supported container topologies — Ollama on the host versus
   Ollama as a sibling service — must be documented then.

Reversal: if the appliance work in Track F concludes that a container runtime is
the right base for the device image, this ADR is superseded rather than edited.

## Notes

This ADR governs the role of the image, not its contents. Base image choice,
healthcheck mechanism, and build layout are implementation decisions recorded in
the Dockerfile itself.

## Amendment, 2026-08-15

`docs/adr/0003-target-platform-policy.md` establishes that short-term deployment
on a Windows host runs the Linux container, with Ollama native on the host.

This does not change the decision: the image is still not a published artifact
and carries no external support commitment. It does mean the image is running
real workloads sooner than Track E assumed, so image quality is load-bearing
now rather than at first publication. The deferred items and their triggers are
unchanged.
