# Linden Interface Specifications

This directory contains the normative technical specifications for each architectural layer interface in Linden.

Each document defines the exact Go interface declaration, input/output parameter tables, execution semantics, concurrency rules, error mappings, and the corresponding contract tests in `validation/contracts/`.

---

## Interface Index

| Document | Layer | Interface / Boundary | Primary Consumer | Contract Test Suite |
|---|---|---|---|---|
| [api-gateway.md](api-gateway.md) | `api` | HTTP Wire Handlers & OpenAI Compat | Web UI, External LAN Clients | `validation/contracts/api_contract_test.go` |
| [orchestrator-service.md](orchestrator-service.md) | `orchestrator` | `ChatService` | `api` | `validation/contracts/orchestrator_contract_test.go` |
| [inference-client.md](inference-client.md) | `inference` | `Client` | `orchestrator` | `validation/contracts/inference_contract_test.go` |
| [storage-store.md](storage-store.md) | `storage` | `Store` | `orchestrator` | `validation/contracts/storage_contract_test.go` |
| [discovery-advertiser.md](discovery-advertiser.md) | `discovery` | `Advertiser` | `cmd` | `validation/contracts/discovery_contract_test.go` |
| [mcp-boundary.md](mcp-boundary.md) | `mcp` | Tool & Protocol Boundary | `orchestrator` | `validation/contracts/` |
| [error-taxonomy.md](error-taxonomy.md) | `errs` | Shared Kernel Taxonomy | All Layers | Stdlib unit tests |

---

## Normative Rules

1. **Exact Parity with Code:** The tables in these documents represent the normative standard. Any changes to interface method signatures, struct fields, or error return codes must update the corresponding specification in the same commit.
2. **Automated Verification:** Interface definitions are verified against the Go AST using the `tools/docgen` verification utility (`go run ./tools/docgen -verify`).
3. **Behavioral Conformance:** Implementations must satisfy the isolated interface conformance test suites located in `validation/contracts/` prior to merging.
