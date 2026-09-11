# Troubleshooting Guide

A practical diagnostics and recovery guide for homeowners, operators, and developers running Linden on a local network.

---

## Quick Diagnostics

Before diving into specific issues, verify Linden's operational health using its standard diagnostic endpoints:

| Action | Command / URL | Expected Response | Meaning |
|---|---|---|---|
| **Service Liveness** | `curl http://localhost:8080/health`<br>or open `http://localhost:8080/health` | `{"status":"ok"}` (HTTP 200) | Linden server process is running and accepting HTTP connections. |
| **Version & Commit** | `curl http://localhost:8080/version`<br>or open `http://localhost:8080/version` | `{"version":"...","commit":"...","go_version":"..."}` | Confirms binary version, build commit, and Go compiler version. |
| **Model Connectivity** | `curl http://localhost:8080/models`<br>or open `http://localhost:8080/models` | `{"models":[{"name":"..."}]}` | Linden successfully communicates with the local Ollama inference daemon. |

---

## Network Discovery & Access (`linden.local`)

### Symptom: Cannot connect to `http://linden.local:8080`

Linden advertises itself over multicast DNS (mDNS) as `linden.local`. If your web browser cannot find the server, work through these common network causes:

1. **Client / Guest Wi-Fi Isolation:**
   - **Cause:** Many home routers have a "Guest Network" or "AP Isolation" feature that blocks multicast traffic and peer-to-peer communication between devices on the same Wi-Fi.
   - **Remedy:** Ensure both your browsing device and the Linden appliance are on the primary local Wi-Fi or connected via Ethernet.
2. **Active VPN Connection:**
   - **Cause:** Corporate or commercial VPNs often tunnel all DNS queries to remote resolvers, preventing the resolution of `.local` top-level domains.
   - **Remedy:** Temporarily disconnect the VPN or enable local network LAN access in your VPN client settings.
3. **OS mDNS Resolver Limitations:**
   - **Windows:** Ensure the "mDNS" protocol is active (standard on Windows 10/11). If `.local` fails to resolve, you can install Apple Bonjour for Windows or use the direct IP fallback.
   - **Android:** Some Android versions do not resolve `.local` mDNS hostnames by default in the Chrome browser.
4. **Direct IP Fallback:**
   - Find the device's local IP address (e.g. from your router's DHCP client list or running `ip addr` / `ipconfig` on the host).
   - Navigate directly to `http://<device-ip>:8080` (e.g., `http://192.168.1.150:8080`).

---

## Inference Engine & Model Readiness

Linden delegates local token generation to an underlying inference backend (defaulting to Ollama at `http://localhost:11434`).

### Symptom: "Inference service unavailable" or HTTP 502 Bad Gateway

- **Cause:** The Ollama inference daemon is not running or is unreachable from the Linden process.
- **Diagnostics:**
  - Test Ollama directly: `curl http://localhost:11434/api/tags`
  - In Docker environments, ensure `LINDEN_INFERENCE_URL` points to `http://host.docker.internal:11434` or the appropriate container network alias rather than `localhost`.
- **Remedy:**
  - Start the Ollama daemon: `ollama serve` (or start the Ollama system service).
  - Check Ollama logs for out-of-memory (OOM) crashes.

### Symptom: "Model not found" or HTTP 404

- **Cause:** The requested model (such as `tinyllama:latest` or a custom configured model) has not been downloaded to the local inference host.
- **Diagnostics:**
  - Run `curl http://localhost:8080/models` to view all models currently installed and recognized by Ollama.
- **Remedy:**
  - Download the missing model directly on the machine running Ollama:
    ```sh
    ollama pull tinyllama
    ```
  - Once the download completes, refresh the Linden web interface or re-query `GET /models`.

### Symptom: Sluggish token generation or frozen streaming responses

- **Cause:** The selected model is too large for the host's available physical RAM or GPU VRAM, forcing the operating system to swap weights to disk or fall back entirely to single-core CPU execution.
- **Remedy:**
  - Choose a lighter quantized model tailored to your hardware tier (e.g., `tinyllama` or `qwen2.5:0.5b` for 8GB RAM devices; `qwen2.5:7b` for 16GB+ RAM with dedicated GPU acceleration).
  - Close memory-intensive background applications on the inference host.

---

## Port & Startup Conflicts

### Symptom: Linden server exits on startup with `bind: address already in use`

- **Cause:** Another service on the host is already bound to port 8080 (the default Linden listening port).
- **Diagnostics:**
  - Find what is listening on port 8080:
    - **Linux/macOS:** `lsof -i :8080` or `ss -tulpn | grep 8080`
    - **Windows:** `netstat -ano | findstr :8080`
- **Remedy:**
  - Terminate the conflicting process, or
  - Configure Linden to bind to a different port by setting the `LINDEN_PORT` environment variable:
    ```sh
    export LINDEN_PORT=8085
    ./bin/linden
    ```

---

## Session Persistence & Storage

### Symptom: Chat history is not saved across restarts

- **Cause:** Linden persists conversation sessions to atomic JSON files inside the local storage directory (default: `data/sessions/`). If the process lacks write permissions or the disk is full, session writes will fail.
- **Diagnostics:**
  - Check disk space: `df -h`
  - Check directory ownership: `ls -ld data/sessions/`
- **Remedy:**
  - Ensure the user running Linden owns the `data/` directory:
    ```sh
    mkdir -p data/sessions
    chmod -R u+rwX data
    ```
  - Free up storage space if disk usage is at 100%.

---

## Privacy & Network Verification

Homeowners and auditors can easily verify Linden's privacy guarantees independently without trusting code comments:

1. **Verify Zero Outbound Network Traffic:**
   - Monitor the Linden process using network monitoring utilities (`tcpdump`, `iftop`, Little Snitch, or Wireshark):
     ```sh
     sudo tcpdump -i any host not 127.0.0.1 and not 192.168.0.0/16
     ```
   - Confirm that sending chat messages generates **zero** egress traffic to public internet IP addresses.
2. **Audit Application Logs:**
   - Review Linden's standard output logs.
   - Verify that log entries contain only operational events (e.g., HTTP status codes, request durations, and request IDs).
   - Observe that message bodies, chat prompts, and user responses are never written to log files, even in failure modes.
