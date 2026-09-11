# Architecture

How the system is decomposed internally: what the parts are, what depends on what, and the shape of each seam between them.

## What belongs here

Four tests, in order. A document belongs in `docs/architecture/` if it fails the first two and passes the third.

| Question | If yes |
|----------|--------|
| Could a user observe this without reading code? | `docs/product/` |
| Does it describe how the project works? | `docs/project/` |
| Does it describe how the system is built internally? | here |
| Is it a single decision with alternatives and consequences? | `docs/adr/` |

Concretely, this directory owns:

- Layer decomposition and responsibilities
- Dependency direction and forbidden imports
- The interface shape at each internal seam
- Normative wire contracts — request and response schemas, event names, error codes
- Cross-cutting technical requirements: timeouts, logging, request identity

It does not own capability descriptions, user flows, privacy commitments, roadmap, or test policy.

## Boundary with the product documentation

Client-facing HTTP wire contracts, internal Go interfaces, and subsystem specifications live here in `docs/architecture/interfaces/` rather than in `docs/product/`.

| | `docs/product/` | `docs/architecture/interfaces/api-gateway.md` |
|--|--|--|
| Reader | Homeowners, operators, and evaluators | Developers and autonomous agents building or integrating with Linden |
| Content | User guides, privacy guarantees, capability matrices, and troubleshooting | Exact HTTP routes, JSON/SSE wire schemas, headers, status codes, and error taxonomies |
| Changes when | The user experience or capability changes | The wire format or interface changes |

**The normative definitions live here.** Contract tests in `validation/contracts/` assert against these documents (enforced by `tools/docgen -verify`), so they stay in lockstep with the implementation. The product documentation focuses on user-observable behavior and troubleshooting without duplicating low-level API schemas.

## Contents

| Document | Purpose |
|----------|---------|
| `layer-interface-spec.md` | Layer boundaries, dependency rules, shared contracts, per-layer interface shapes |
| `system_design.md` | Software & runtime stack, inference daemon, MCP tooling, and update pipelines |
| `interfaces/` | Dedicated, table-driven technical specifications for each layer interface |

## Related

- What belongs in each documentation directory: `docs/project/README.md`
- Decisions that shaped this structure: `docs/adr/`
- Why interfaces are specified before implementation: `docs/project/contract-first-delivery.md`
