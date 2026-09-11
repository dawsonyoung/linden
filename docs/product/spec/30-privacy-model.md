# Privacy Model

This document outlines the privacy guarantees, data boundaries, and network behavior of the Linden appliance.

## Data categories

| Category | Stored Location | Retention | User-deletable | Leaves device |
|----------|-----------------|-----------|----------------|---------------|
| Message content | Local disk (`data/sessions`) | Host file store | Yes (clearing data dir) | No |
| Model selection | In-memory / browser state | Session duration | Yes (reset browser) | No |
| Operational logs | Host stdout / journald | Rotated by host OS | Yes (OS log purge) | No |

## Logging commitments

1. **Zero User Data in Logs:** No message content, prompt text, or model output ever appears in server logs at any log level (including `DEBUG`).
2. **Metadata Only:** Server logs record operational metadata only: HTTP method, path, response status, duration, and generated `X-Request-ID`.
3. **Query String Scrubbing:** URL query parameters are stripped and excluded from access logs to prevent accidental content leakage.
4. **Sanitized Error Payloads:** Error responses returned to clients never contain internal stack traces or database internals.

## Network behavior

1. **LAN-Bound Operation:** Linden binds to the configured local address (`0.0.0.0:8080` by default) and advertises service presence over local multicast DNS (`linden.local`). It does not open firewall ports or establish outbound internet tunnels.
2. **Local Inference Communication:** By default, Linden communicates exclusively across loopback with the local inference engine (`http://localhost:11434`).
3. **Zero Telemetry:** Linden contains no telemetry collectors, crash reporting daemons, or usage tracking modules. No phone-home requests are initiated.

## Extensibility boundary (External Engines)

Linden supports user-configured inference backend URLs. If an operator manually configures Linden to target a cloud-hosted inference provider or third-party proxy, data dispatched to that endpoint is subject to that external provider's policies. Linden remains strictly local by default.
