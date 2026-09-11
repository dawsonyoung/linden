# Linden Documentation

Welcome to the documentation for **Linden**, a privacy-first, safety optimized personal AI platform featuring streamlined deployment of the local-only chat experience.

Linden runs directly on hardware you own with complete local-only capability. By default, it operates fully offline without requiring external cloud LLM connections, keeping all conversations, documents, and vectors within your local environment while establishing the trust needed for future optional cloud-enhanced services.

---

## Documentation Structure

To serve users, developers, and autonomous agents effectively, Linden's documentation is divided into three primary sections:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              Linden Docs                                │
├─────────────────────────┬─────────────────────────┬─────────────────────┤
│      1. User Guide      │  2. Technical Specs     │      3. PRDs        │
│  "How to Use & Trust"   │ "What Exists & APIs"    │ "Why & What's Next" │
└─────────────────────────┴─────────────────────────┴─────────────────────┘
```

---

### 1. User Guide
Written for homeowners, operators, and everyday users who interact with Linden on their home network.

- **[Overview & Commitments](spec/00-overview.md):** What Linden is, core privacy guarantees, and current runtime status.
- **[User Flows & First Run](spec/20-user-flows.md):** Step-by-step walkthroughs of discovering `linden.local`, starting local chat sessions, and handling offline inference recovery.
- **[Privacy & Trust Model](spec/30-privacy-model.md):** Clear explanations of data boundaries, local file storage, and the strict zero-logging policy.
- **[Troubleshooting Guide](spec/40-troubleshooting.md):** Diagnostics, network mDNS recovery, Ollama model setup, and port configuration.

---

### 2. Technical Specifications
Authoritative, living technical references describing Linden's observable capabilities and terminology.

- **[Capabilities Matrix](spec/10-capabilities.md):** The current status of all platform capabilities (Local Streaming Chat, OpenAI Compatibility, mDNS LAN Discovery).
- **[Glossary](spec/90-glossary.md):** Standard terminology used across Linden's interfaces and documentation.

*(Note: Normative API contracts, HTTP wire schemas, Go interfaces, and layer specifications are maintained in [`docs/architecture/interfaces/`](https://github.com/dawsonyoung/linden/tree/main/docs/architecture/interfaces), including [`api-gateway.md`](https://github.com/dawsonyoung/linden/tree/main/docs/architecture/interfaces/api-gateway.md) for external client integration).*

---

### 3. Product Requirements (PRDs)
High-level, problem-centric requirement documents establishing user pain points, target personas, scope boundaries (Must-Have vs. Non-Goals), and privacy constraints.

- **Role:** PRDs serve as the primary **behavioral anchor** for engineering and autonomous AI agents before code or contract tests are written.
- **Agent Workflow:** Agents are not blocked by missing low-level specs; they reference the PRD tree for intent and propose spec extensions for review.
- **[How We Design Linden](prd/README.md):** Design philosophy and guide to how PRDs bound development.
- **[PRD-0001: Local Streaming Chat](prd/0001-local-chat.md):** Shipped MVP capability for multi-turn local chat with SSE token streaming.
- **[PRD-0002: Document Ingestion & Retrieval](prd/0002-document-retrieval.md):** Accepted requirement specification for Track C (Private Document Ingestion & Local RAG).
- **[PRD Template](prd/TEMPLATE.md):** Standardized format for authoring future capability PRDs.
