# Linden Documentation

Welcome to the documentation for **Linden**, a privacy-first, zero-cloud personal AI platform for your home network.

Linden runs directly on hardware you own. It requires no user accounts, makes no external network calls to cloud LLMs, and ensures all conversations, documents, and vectors remain strictly within your local environment.

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
- **[User Flows & First Run](spec/20-user-flows.md):** Step-by-step walkthroughs of discovering `linden.local`, starting local chat sessions, selecting models, and handling offline inference recovery.
- **[Privacy & Trust Model](spec/30-privacy-model.md):** Clear explanations of data boundaries, local file storage, and the strict zero-logging policy.

---

### 2. Technical Specifications
Authoritative, living technical references describing Linden's observable software surface, API endpoints, and system behaviors as they exist today.

- **[Capabilities Matrix](spec/10-capabilities.md):** The current status of all platform capabilities (Local Streaming Chat, OpenAI Compatibility, mDNS LAN Discovery).
- **[API Surface](spec/40-api-surface.md):** Exact HTTP contracts and schemas for client integration, including `POST /chat`, `GET /models`, `GET /health`, and the OpenAI-compatible `POST /v1/chat/completions`.
- **[Glossary](spec/90-glossary.md):** Standard terminology used across Linden's interfaces and documentation.

*(Note: Internal layer boundaries, Go interfaces, and architectural guidelines are maintained in [`docs/architecture/`](https://github.com/dawsonyoung/linden/tree/main/docs/architecture)).*

---

### 3. Product Requirements (PRDs)
High-level, problem-centric requirement documents establishing user pain points, target personas, scope boundaries (Must-Have vs. Non-Goals), and privacy constraints.

- **Role:** PRDs serve as the primary **behavioral anchor** for engineering and autonomous AI agents before code or contract tests are written.
- **Agent Workflow:** Agents are not blocked by missing low-level specs; they reference the PRD tree for intent and propose spec extensions for review.
- **[How We Design Linden](prd/README.md):** Design philosophy and guide to how PRDs bound development.
- **[PRD-0001: Local Streaming Chat](prd/0001-local-chat.md):** Shipped MVP capability for multi-turn local chat with SSE token streaming.
- **[PRD-0002: Document Ingestion & Retrieval](prd/0002-document-retrieval.md):** Accepted requirement specification for Track C (Private Document Ingestion & Local RAG).
- **[PRD Template](prd/TEMPLATE.md):** Standardized format for authoring future capability PRDs.
