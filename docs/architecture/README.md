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

## Boundary with the product spec

Both this directory and `docs/product/spec/` describe the HTTP API, at different levels and for different readers.

| | `docs/product/spec/40-api-surface.md` | `docs/architecture/layer-interface-spec.md` |
|--|--|--|
| Reader | Someone building against Linden | Someone building Linden |
| Content | What each endpoint is for, what it guarantees, what errors mean | Exact schemas, event names, field types, error code mapping |
| Changes when | The guarantee changes | The wire format changes |

**The normative definition lives here.** Contract tests in `validation/contracts/` assert against this document, so it must stay in lockstep with the implementation. The product spec describes the surface for a human reader and links here rather than restating field lists — a duplicated schema is a schema that will drift.

## Contents

| Document | Purpose |
|----------|---------|
| `layer-interface-spec.md` | Layer boundaries, dependency rules, shared contracts, per-layer interface shapes |

## Related

- What belongs in each documentation directory: `docs/project/README.md`
- Decisions that shaped this structure: `docs/adr/`
- Why interfaces are specified before implementation: `docs/project/contract-first-delivery.md`
