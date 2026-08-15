# ADR-0004: Shared Error Taxonomy and Streaming Callback Signature

| Field | Value |
|-------|-------|
| Status | Accepted |
| Date | 2026-08-15 |
| Deciders | dawsonyoung |

## Context

Stage A.2 specifies the inference interface before any provider exists. Two
decisions must be settled first, because the contract tests assert both.

**Where the error taxonomy lives.** `docs/architecture/layer-interface-spec.md`
publishes nine error codes. They exist only as prose. Nothing in the code
represents them, so nothing enforces that a layer uses `unavailable` to mean what
another layer means by it. A.2's coverage list includes mapping a provider
timeout to `deadline_exceeded` and an unknown model to `not_found`; those
assertions need a code representation to assert against.

**How a streaming consumer stops.** `ChatStream` delivers chunks through a
callback. Whether that callback returns an error determines how a consumer
aborts a stream, and A.7 will hit this immediately: SSE writes fail when a
client disconnects, which is routine rather than exceptional.

## Options considered: error taxonomy

### A. Per-layer error types

Each layer defines its own errors; adjacent layers map at the boundary.

Preserves layer independence. Costs a mapping at every boundary, and those
mappings are where meanings drift — two layers can each define `unavailable`
and disagree about it. Leaves the published taxonomy with no representation in
code.

### B. Shared error kernel

One package defines the taxonomy. Every layer uses it directly.

Makes the published contract executable, removes N×M boundary mappings, and
lets `errors.Is` work across layers. Introduces a package that every layer
imports, which the previous reading of the leaf-layer rule forbade.

## Decision: shared error kernel

**Option B**, with a narrow charter.

`src/errs` contains the nine codes, a wrapping error type, and accessors.
Nothing else. It has no dependencies beyond the standard library, performs no
I/O, and knows nothing about HTTP. Transport mapping stays in `api`, which is
the layer that owns transport.

### The leaf-layer rule is clarified, not broken

`AGENTS.md` described `inference`, `storage`, and `discovery` as having zero
internal dependencies. That phrasing conflated two goals.

The purpose of leaf layers is **modularity**, not independence for its own sake.
A layer is modular when it is replaceable behind a clear interface. Sharing a
dependency-free vocabulary package does not make a layer less replaceable; it
makes the boundaries between layers describable in the same terms.

The rule now reads: leaf layers depend on no other *layer*. A shared kernel is
permitted when it is dependency-free, logic-free, and has an explicit charter.
`src/errs` is currently the only such package.

### Charter for `src/errs`

1. Standard library only. No third-party dependency, ever.
2. No I/O, no logging, no transport knowledge.
3. Only the taxonomy and the machinery to carry it.
4. Adding a code requires updating `docs/architecture/layer-interface-spec.md`
   in the same change; the document and the package are one contract. No test
   can enforce this, because a test cannot know what the document says.
5. It is not a utility package. Anything that is not the error taxonomy does not
   belong, and a second shared kernel needs its own ADR.
6. **No formatting variant.** `Wrap` takes a plain message rather than a format
   string. A `Wrapf` would make interpolating a prompt, a model name, or a
   session identifier into a logged string the natural thing to do. The
   non-formatting signature makes embedding user content a visible act of
   concatenation, which a reviewer can see.
7. **Exempt from the contract-test rule.** The spec requires contract tests for
   layer interfaces before implementation. `errs` is a kernel, not a layer; it
   has no consumers to isolate from and is covered by unit tests.

### Adding context without reclassifying

`Wrap` assigns a code, so it reclassifies. A layer that only wants to add
context uses `fmt.Errorf("...: %w", err)`, which preserves the original code
through the chain.

This matters because `CodeOf` returns the outermost code. Wrapping in `Internal`
simply to add a note would silently downgrade a precise `unavailable` to a 500.

## Decision: callback returns an error

`ChatStream(ctx, Request, func(Chunk) error) (Result, error)`.

The alternative, `func(Chunk)`, forces every consumer to smuggle failures out
through a closure and cancel the context by hand. The error return removes that
from each call site rather than adding to it.

The two mechanisms are distinct and both are needed:

- **Context cancellation** is user-initiated stop, required by PRD-0001 R3.
- **A callback error** is consumer-side failure, such as an SSE write to a
  disconnected client.

When the callback returns an error, `ChatStream` stops delivering, performs no
further provider work, and returns that error unwrapped, so the caller can match
on its own sentinel.

## Consequences

1. Every layer may import `src/errs`. No layer may import another layer's
   package for error definitions.
2. `api` owns the mapping from `errs.Code` to HTTP status and to the
   client-facing body. Internal detail never crosses that boundary.
3. A code added to the taxonomy is added in one place and is immediately
   available everywhere, which is the extensibility this was chosen for.
4. `AGENTS.md` and `docs/architecture/layer-interface-spec.md` are updated to
   state the clarified rule.
5. The charter is the guard against `errs` becoming a utility package. If it is
   violated, this ADR is the thing that was violated.

## Notes

If a second shared kernel is ever proposed, this ADR is precedent for the shape
of the argument but not approval for the package. Each needs its own charter.
