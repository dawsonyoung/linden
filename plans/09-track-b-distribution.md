# Track B: Distribution and Operations

## Status
Active implementation.

## Purpose

Establish effortless, single-command installation and distribution for Linden on Linux hosts and low-cost home appliances (<32GB RAM). This milestone transitions Linden from a source-built repository to a turnkey product installable on barebones Linux hardware with zero local development toolchain requirements.

---

## Target Audience & Hardware Persona

- **Primary Target**: Low-cost, always-on home appliance hardware (mini-PCs such as Intel N100/N97, AMD Ryzen 5000/7000 mini-PCs, or 8GB/16GB Raspberry Pi 5).
- **Memory Ceiling**: `< 32GB RAM` (sweet spot: 8GB to 16GB).
- **Core Use Case**: High-privacy text document ingestion and local inference (sensitive financial records, tax documents, health histories, and educational notes) with zero cloud telemetry.
- **Recommended OS**: **Debian 12 (Bookworm) Minimal** (~150MB base RAM) or **Ubuntu Server 24.04 LTS** (~300MB base RAM).

---

## Sequencing & Deliverables

```
B.1 ci/binary-release-workflow ──→ B.2 feat/turnkey-linux-installer ──→ B.3 docs/hardware-and-distro-guidance
```

---

## B.1 `ci/binary-release-workflow`

Intent: Automate reproducible, statically linked Linux binary builds with embedded SvelteKit assets for distribution via GitHub Releases.

### Deliverables
1. `.github/workflows/release.yml` with triggers on semantic version tags (`v*.*.*`) and `workflow_dispatch`.
2. Cross-compilation matrix for `linux/amd64` and `linux/arm64` using `CGO_ENABLED=0`.
3. SvelteKit frontend compilation embedded directly into the Go binary (`embed.FS`).
4. Automated SHA256 checksum generation (`checksums.txt`).
5. Draft/final GitHub release publishing.

### Gates
- Workflow syntax validation.
- Successful artifact compilation of both architectures.

---

## B.2 `feat/turnkey-linux-installer`

Intent: Provide a single-command installer script (`curl ... | bash`) that provisions Linden, systemd service units, and local LLM backends on barebones Linux hosts.

### Deliverables
1. `scripts/install.sh`:
   - System architecture detection (`x86_64` vs `aarch64`).
   - Binary download and verification from GitHub Releases into `/usr/local/bin/linden`.
   - Explicit user prompt before downloading or installing Ollama (`--yes-ollama` and `--skip-ollama` flags supported).
   - Dynamic model recommendation based on host RAM (`qwen2.5:3b` for <16GB; `qwen2.5:7b` for &ge;16GB).
   - Systemd service configuration (`packaging/systemd/linden.service`) with sandbox protections (`DynamicUser=yes`, `ProtectSystem=full`).
   - LAN and mDNS endpoint display (`http://linden.local:8080`).
2. `packaging/systemd/linden.service`:
   - Systemd unit definition with auto-restart and state directory isolation (`/var/lib/linden`).

### Gates
- Clean shell syntax check (`sh -n scripts/install.sh`).
- Smoke execution in clean Docker container (`ubuntu:24.04` and `debian:12`).

---

## B.3 `docs/hardware-and-distro-guidance`

Intent: Update primary project documentation to reflect the new turnkey install method, recommended Linux distributions, and targeted hardware capabilities.

### Deliverables
1. Root `README.md` updates:
   - Turnkey installation command prominently featured in Quick Start.
   - Recommended Linux distributions (Debian 12 / Ubuntu Server 24.04).
   - Target Hardware Specification & Document Ingestion use case profile.
2. Architecture & plan cross-links in `.context/README.md`.

---

## Track B Exit Criteria

1. Tagging a release publishes valid Linux `amd64` and `arm64` binaries to GitHub Releases with checksums.
2. Running `curl -fsSL .../install.sh | bash` on a fresh Linux server sets up Linden as an active systemd service.
3. The installer respects user consent prior to initiating Ollama installation.
4. Documentation clearly articulates hardware recommendations (<32GB RAM) and recommended Linux distributions.
