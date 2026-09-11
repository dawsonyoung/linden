# Discovery Advertiser Technical Specification

## 1. Overview & Boundary

The `discovery` layer (`src/discovery/`) broadcasts Linden's presence across the local network using Multicast DNS (mDNS) and DNS-SD (DNS Service Discovery). This enables zero-configuration discovery so users on the local Wi-Fi can navigate to `http://linden.local:8080` without looking up IP addresses.

- **Package:** `github.com/dawsonyoung/linden/discovery`
- **Exported Interface:** `Advertiser`
- **Dependencies:** Leaf layer. Depends on `github.com/grandcat/zeroconf` and `errs`. Never imports `api`, `orchestrator`, or `cmd`.
- **Consumers:** Initialized and controlled by `src/cmd/main.go`.

---

## 2. Interface Definition

```go
type Advertiser interface {
    // Start advertises the Linden HTTP service on the given TCP port.
    Start(port int) error

    // Stop shuts down the mDNS responder and unregisters records.
    Stop()
}
```

---

## 3. mDNS Service Parameters & TXT Records

| Parameter | Value | Description |
|---|---|---|
| **Instance Name** | `Linden AI (<hex_id>)` | Human-readable service instance advertised on LAN. |
| **Service Type** | `_linden._tcp` | Standard DNS-SD service identifier. |
| **Domain** | `local.` | Standard mDNS top-level domain. |
| **Target Host** | `linden.local` | Hostname mapped to local network interface IPv4 addresses. |
| **Default Port** | `8080` | Port passed to `Start(port)`. |

### DNS-SD TXT Records Table
| Key | Example Value | Purpose |
|---|---|---|
| `version` | `0.1.0` | Informs network clients of the Linden runtime version. |
| `api_path` | `/` | Base HTTP endpoint path. |
| `id` | `a3f8` | Unique 4-character random instance identifier. |

---

## 4. Lifecycle & Registration Semantics

```
[ Uninitialized ] ──► Start(port) ──► [ Advertising: linden.local ] ──► Stop() ──► [ Stopped ]
```

1. **Interface Discovery:** At startup, `Start` enumerates all active local network interfaces (`net.Interfaces()`), filters for active non-loopback IPv4 addresses, and binds mDNS responses to those addresses.
2. **Proxy Registration with Fallback:** Linden attempts proxy registration (`zeroconf.RegisterProxy`) to bind `linden.local`. If proxy registration is rejected by the OS, it falls back to standard system hostname registration (`zeroconf.Register`).
3. **Graceful Teardown:** `Stop()` issues mDNS goodbye packets and frees underlying UDP multicast socket resources. `Stop()` is idempotent and safe to call multiple times.

---

## 5. Error Taxonomy & Mapping Table

| Failure Scenario | `errs.Code` | Notes |
|---|---|---|
| Port <= 0 or > 65535 | `errs.InvalidArgument` | Invalid TCP port parameter. |
| UDP 5353 multicast socket failure | `errs.Internal` | Network interface permission failure or firewall restriction. |
| Zeroconf server initialization fault | `errs.Internal` | Returned wrapped in `errs.Error`. |

---

## 6. Contract Test Verification Matrix

Every requirement in this specification is verified in `validation/contracts/`:

| Test Function | Contract Verified |
|---|---|
| `Test_Discovery_StartAndStop_Idempotent` | Verifies `Start` binds cleanly and `Stop` cleans up without panics. |
| `Test_Discovery_InvalidPort_Rejection` | Verifies invalid ports return `errs.InvalidArgument`. |
