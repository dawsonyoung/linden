# Product Requirements (PRDs)

## How We Design Linden

In Linden, every major capability begins with a high-level **Product Requirements Document (PRD)**.

Unlike technical layer specifications or low-level API contracts, a PRD is a problem-centric, observable document written from the user's perspective. It does not describe internal Go structs, database queries, or function names. Instead, it defines:
1. **The Real-World Problem:** Why this feature matters to a homeowner, privacy auditor, or local AI enthusiast.
2. **User Stories & Scenarios:** What the user does and what value they receive.
3. **Explicit Boundaries & Non-Goals:** What the capability deliberately will *not* do in its initial release.
4. **Privacy & Trust Guarantees:** Strict constraints ensuring user data remains protected, local, and auditable.

---

## The Role of PRDs in Our Agentic Workflow

Linden is developed using advanced AI-assisted and autonomous agent workflows. In this architecture, PRDs serve as the primary **behavioral anchor**:

* **Bounding Agent Scope:** AI coding agents reference the PRD tree to understand the overarching design intent without hallucinating unnecessary features or over-engineering solutions.
* **Non-Blocking Reference:** During early development cycles, agents are not blocked if an exhaustive technical specification is not yet written; they reference the PRD tree for high-level guidance and propose spec extensions for human review.
* **Evolution Toward Full Spec Trees:** Over time, agents will autonomously generate complete, fine-grained specification trees based on these high-level PRDs for review before implementation begins.

---

## Capability Index

| PRD | Title | Status | Primary Capability |
|---|---|---|---|
| [PRD-0001](0001-local-chat.md) | Local Streaming Chat | **Shipped** | Multi-turn conversational chat with SSE token streaming |
| [PRD-0002](0002-document-retrieval.md) | Local Document Ingestion & Retrieval | **Accepted** | Private document parsing, vector indexing, and citations |
