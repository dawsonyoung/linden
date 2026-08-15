---
name: "Review Product Spec"
description: "Classify a change against the product spec, PRDs, and interface spec before implementation. Use at the Spec step of every cycle, before any contract test or code is written."
---
# Review Product Spec

Run this before writing contract tests or implementation. Its output satisfies the Spec Gate.

## Step 1: Classify

Answer all four. "No" must be a decision, not an omission.

| Question | If yes |
|----------|--------|
| Does this add a capability a user could name? | An accepted PRD is required in `docs/product/prd/` before implementation |
| Does this change behavior a user can observe? | `docs/product/spec/` must be updated in the same PR as the behavior |
| Does this change a layer interface, schema, error code, or config key? | `docs/architecture/layer-interface-spec.md` must be updated |
| Is this internal only? | Record `Spec impact: none` and why |

## Step 2: Locate the affected pages

- Capability behavior → `docs/product/spec/10-capabilities.md`
- Flow or error experience → `docs/product/spec/20-user-flows.md`
- Data handling, retention, or network egress → `docs/product/spec/30-privacy-model.md`
- Client-facing endpoint or event → `docs/product/spec/40-api-surface.md`
- New term the spec relies on → `docs/product/spec/90-glossary.md`

## Step 3: Update

- Spec is present tense and describes only what is real. Planned behavior stays in the PRD until it ships.
- No implementation detail: no file paths, function names, or struct fields.
- A new page must be added to `docs/product/SUMMARY.md` or it will not publish.
- Privacy-relevant changes always touch `30-privacy-model.md`, even when the answer is "nothing changed" — say so explicitly.

## Step 4: Report

Produce, for the PR body:

```
Spec impact: <none | spec updated | PRD required | interface updated>
Pages touched: <paths, or none>
PRD: <PRD-NNNN and status, or N/A with reason>
```

## Stop conditions

- A new capability with no PRD — stop and escalate.
- A PRD with unresolved Open questions — stop; it cannot move to Accepted.
- Behavior change that cannot be described without naming internals — the design is leaking; escalate to Design Agent.
