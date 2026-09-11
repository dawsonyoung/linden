# Linden Documentation

Welcome to the documentation for **Linden**, a privacy-first, zero-cloud local AI platform.

Linden runs directly on hardware you own and control on your local network. It requires no accounts, makes no external network calls to cloud LLMs, and stores all conversations and documents locally with complete privacy.

---

## The Three Tiers of Linden Documentation

To ensure clarity, prevent scope creep, and provide explicit guardrails for human contributors and autonomous AI agents alike, Linden's documentation is organized into three distinct tiers:

```
                  ┌────────────────────────────────────────┐
                  │   Tier 1: Product Requirements (PRDs)  │
                  │        "The Why & Design Intent"       │
                  └───────────────────┬────────────────────┘
                                      │ Informs & Bounds
                                      ▼
                  ┌────────────────────────────────────────┐
                  │    Tier 2: Product Specification       │
                  │      "The What & Current Reality"      │
                  └───────────────────┬────────────────────┘
                                      │ Implemented By
                                      ▼
                  ┌────────────────────────────────────────┐
                  │    Tier 3: Technical Architecture      │
                  │     "The How & System Contracts"       │
                  └────────────────────────────────────────┘
```

---

### Tier 1: Product Requirements (PRDs) — *The "Why"*
**Purpose:** High-level problem statements, target user personas, user stories, explicit boundaries (Must-Have vs. Out of Scope), and privacy/trust constraints.

- **Role:** Answers *"Why should this exist, who is it for, and what must it never do?"* PRDs provide a durable behavioral anchor that bounds scope before implementation begins.
- **Tense & Lifecycle:** Future-facing proposals that are reviewed, accepted, and frozen once shipped.
- **Key Documents:**
  - [Design Philosophy & Capability Index](prd/README.md)
  - [PRD-0001: Local Streaming Chat](prd/0001-local-chat.md) (`Shipped`)
  - [PRD-0002: Local Document Ingestion & Retrieval](prd/0002-document-retrieval.md) (`Accepted` - Track C)
  - [PRD Template](prd/TEMPLATE.md)

---

### Tier 2: Product Specifications — *The "What"*
**Purpose:** The authoritative, living description of observable behavior that users and client applications can experience and rely on today.

- **Role:** Answers *"What does Linden do right now?"* Written strictly from the user's perspective—never mentioning internal Go struct names, file paths, or private implementation details.
- **Tense & Lifecycle:** Present tense. Updated in the exact same pull request that introduces or modifies observable behavior.
- **Key Documents:**
  - [Overview & Core Commitments](spec/00-overview.md)
  - [Capabilities Matrix](spec/10-capabilities.md)
  - [User Flows & Experience](spec/20-user-flows.md)
  - [Privacy & Data Handling Model](spec/30-privacy-model.md)
  - [API Surface & Endpoints](spec/40-api-surface.md)
  - [Glossary of Terms](spec/90-glossary.md)

---

### Tier 3: Technical Architecture & Contracts — *The "How"*
**Purpose:** Engineering blueprints, layer boundaries, dependency rules, interface contracts, and architectural decisions.

- **Role:** Answers *"How is Linden constructed and tested?"* Enforces the architectural separation between the HTTP gateway, request orchestrator, inference engine, local persistence, network discovery, and SvelteKit frontend layers.
- **Key Documents:**
  - **Layer Interface Specification:** The formal boundaries and Go interfaces defined in [`docs/architecture/layer-interface-spec.md`](https://github.com/dawsonyoung/linden/blob/main/docs/architecture/layer-interface-spec.md).
  - **Architecture Decision Records (ADRs):** Historic rationale for key technical choices located in [`docs/adr/`](https://github.com/dawsonyoung/linden/tree/main/docs/adr).
  - **Contract-First Testing:** Interface verification suites located in `validation/contracts/`.

---

## How Autonomous Agents Use These Tiers

Linden is built with pair-programming AI agents (Workflow Orchestrator and specialist agents):
1. **PRD Tree as Behavioral Anchor:** Agents consult Tier 1 (PRDs) to understand feature intent, user journeys, non-goals, and privacy constraints.
2. **Non-Blocking Specification Workflow:** Agents are never blocked by a lack of exhaustive specifications; they reference the PRD tree and draft/propose Tier 2 specification extensions for review before implementation.
3. **Contract-First Implementation:** Agents define and satisfy Tier 3 Go interface contracts in `validation/contracts/` before writing production code in `src/`.
