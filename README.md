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

## Architecture Specs

1. [Layer Interface Specification](docs/architecture/01-layer-interface-spec.md)
2. [Contract-First Delivery Plan](docs/architecture/02-contract-first-delivery-plan.md)
3. [Agentic Workflow Framework](docs/architecture/03-agentic-workflow-framework.md)
