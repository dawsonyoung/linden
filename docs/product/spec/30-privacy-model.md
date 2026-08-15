# Privacy Model

> **Status:** Stub. This document is a commitment surface — treat changes here as breaking.

## Data categories

TODO: Enumerate every category of data Linden handles, and for each: where it is stored, how long it persists, whether the user can delete it, and whether it can leave the device.

| Category | Stored | Retention | User-deletable | Leaves device |
|----------|--------|-----------|----------------|---------------|
| Message content | TODO | TODO | TODO | No |
| Model selection | TODO | TODO | TODO | No |
| Logs | TODO | TODO | TODO | No |

## Logging commitment

No user content appears in logs at any level, including debug. Logs record request metadata only: method, path, status, duration, and request identifier.

## Network behavior

TODO: Enumerate every outbound connection Linden can make, what triggers it, and how the user disables it. The expected steady state for the MVP is: connections to a local model backend only.

## Verification

TODO: Describe how a user or auditor can verify these claims without reading the source.
