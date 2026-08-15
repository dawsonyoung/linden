# Linden

A privacy-first local AI platform. See [AGENTS.md](AGENTS.md) for architecture and conventions.

## Quick Start

```powershell
# Install prerequisites
.\scripts\setup.ps1

# Start development
make dev
```

## Status

**MVP in progress** — hello world inference loop (Go server + Ollama + SvelteKit chat UI).

## Documentation

**Product** — what Linden does

1. [Product Spec](docs/product/spec/00-overview.md)
2. [Product Requirements (PRDs)](docs/product/prd/)

**Engineering** — how it is built

3. [Layer Interface Specification](docs/architecture/layer-interface-spec.md)
4. [Contract-First Delivery](docs/process/contract-first-delivery.md)
5. [Agentic Workflow Framework](docs/process/agentic-workflow-framework.md)
6. [Branching Strategy](docs/process/branching-strategy.md)
7. [Implementation Sequence](plans/04-implementation-sequence.md)
8. [Stage 0: Groundwork](plans/05-stage-0-groundwork.md)
