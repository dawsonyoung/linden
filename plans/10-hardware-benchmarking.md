# Track 10: Hardware Benchmarking

## Current Limitations

Linden is engineered for low-cost, dedicated home appliance hardware operating 24/7. However, absolute performance benchmarking (e.g., measuring Tokens Per Second for specific LLMs) is impossible to run reliably in standard CI/CD environments (like GitHub-hosted Ubuntu runners) because:
1.  Standard runners are heavily virtualized, multi-tenant environments with noisy neighbors, causing wildly fluctuating CPU performance.
2.  Standard runners lack GPU acceleration.
3.  Pulling massive GGUF model blobs on every PR is prohibitively expensive and slow.

## The Goal

To provide users with accurate expectations, Linden must publish and automatically verify its inference performance across its "Target Hardware Profiles" (e.g., Intel N100 mini-PCs, Raspberry Pi 5 16GB).

## The Plan: Self-Hosted Benchmarking Fleet

To solve this, Linden will implement a self-hosted hardware benchmarking pipeline.

### 1. Dedicated Hardware Nodes
We will procure and set up dedicated, isolated bare-metal hardware that matches our target deployment profiles:
*   Node A: Intel N100 (16GB RAM)
*   Node B: Raspberry Pi 5 (8GB RAM)
*   Node C: AMD Ryzen 5000 series w/ Integrated Graphics

### 2. GitHub Actions Self-Hosted Runners
Each node will be registered to the Linden GitHub repository as a `self-hosted` runner. They will be labeled accordingly (e.g., `runs-on: [self-hosted, n100, linux, x64]`).

### 3. The Benchmarking Workflow
We will introduce a new workflow (`.github/workflows/benchmark-hardware.yml`) that runs on a weekly schedule (or manually dispatched) across all hardware nodes.
This workflow will:
1.  Pull the latest `main` branch and compile Linden.
2.  Pre-load specific GGUF models (e.g., `qwen2.5:3b`, `qwen2.5:7b-q4_K_M`) stored persistently on the hardware node.
3.  Start Linden and use a benchmarking tool (like `ghz` or a custom script targeting the `/v1/chat/completions` endpoint) to send a standardized prompt payload.
4.  Calculate exact Tokens Per Second (TPS), Time To First Token (TTFT), and Peak Memory Consumption.
5.  Automatically update a `BENCHMARKS.md` file in the repository or a published dashboard if the results are stable.

This guarantees that performance regressions introduced by new orchestrator logic or model engine updates are caught on actual production-tier hardware.
