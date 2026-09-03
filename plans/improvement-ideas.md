# Improvement Ideas

This document serves as a backlog for architectural, structural, or engineering improvements that we want to consider for the future, but have decided to defer in order to maintain velocity or current conventions. 

When an idea is ready to be implemented, it should be moved to a formal plan or ADR (Architecture Decision Record).

## 1. Move Test Doubles to Implementation Packages

**Context**: Currently, test doubles (fakes/simulators) live in the `validation/contracts/` directory alongside the contract tests. 
**Proposal**: Move test doubles to live alongside the actual implementations they simulate (e.g., `src/inference/double.go` or `src/storage/memory.go`). 
**Rationale**: It is good engineering practice for the team that owns the production code to also own the simulator for that code. This makes it easier to keep the double in sync with the real implementation, and allows other depending modules to easily consume the simulator for their own testing needs.
