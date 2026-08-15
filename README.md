# Linden

A privacy-first local AI platform. Runs on hardware you own, on your own network, with content you own. Ask your Linden AI anything with safety, privacy, and ownership.

See [AGENTS.md](AGENTS.md) for architecture and conventions.

## Status

Early development. The server runs and serves operational endpoints; there is no
chat capability yet.

Current stage and what remains: [plans/04-implementation-sequence.md](plans/04-implementation-sequence.md).

This file deliberately does not track per-stage progress. That belongs in the
plan, which is updated as the work is done.

## Quick Start

Linux, or Windows and macOS via the devcontainer.

```sh
./scripts/setup.sh     # check the toolchain
make build             # produces bin/linden
./bin/linden           # serves on :8080
```

```sh
curl http://localhost:8080/health    # {"status":"ok"}
curl http://localhost:8080/version   # build metadata
```

Or run it as a container:

```sh
make docker-smoke
```

Contributors should read [CONTRIBUTING.md](CONTRIBUTING.md) first.

## Support Matrix

| Platform | Role | Notes |
|----------|------|-------|
| Linux | Deployment target | The only supported runtime; all CI gates run here |
| Docker (Linux) | Deployment vehicle | Required CI gate |
| Windows | Development host | Via devcontainer; deploys by running the Linux container |
| macOS | Development host | Expected to work; not gated |

See [docs/adr/0003-target-platform-policy.md](docs/adr/0003-target-platform-policy.md).

| Dependency | Version |
|------------|---------|
| Go | 1.22+ |
| Node.js | 20+ |
| Ollama | current, with a pulled model |

## MVP Checklist

Tracked against [plans/02-repository-setup-and-mvp-plan.md](plans/02-repository-setup-and-mvp-plan.md).

- [x] Server builds and boots with health and version endpoints
- [x] Linux CI gates can fail, and gate on every PR
- [x] Docker image builds and passes a container health check
- [ ] Chat request completes against a local model
- [ ] Replies stream to the client
- [ ] Conversations persist across restarts
- [ ] A second device on the LAN can pair and connect
- [ ] Privacy controls are visible and test-covered

## Documentation

**Product** — what Linden does

1. [Product Spec](docs/product/spec/00-overview.md)
2. [Product Requirements (PRDs)](docs/product/prd/)

**Engineering** — how it is built

3. [Layer Interface Specification](docs/architecture/layer-interface-spec.md)
4. [Contract-First Delivery](docs/project/contract-first-delivery.md)
5. [Agentic Workflow Framework](docs/project/agentic-workflow-framework.md)
6. [Branching Strategy](docs/project/branching-strategy.md)
7. [Implementation Sequence](plans/04-implementation-sequence.md)
8. [Stage 0: Groundwork](plans/05-stage-0-groundwork.md)

## Security

Report vulnerabilities privately. See [SECURITY.md](SECURITY.md).

## License

See [LICENSE](LICENSE).

