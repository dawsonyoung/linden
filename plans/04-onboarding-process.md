# Linden Onboarding & Ownership Process (Future Plan)

This document outlines the proposed design for the initial device onboarding experience, enabling a user to claim ownership of a Linden appliance and configure it via a mobile or web app.

## Goals
1. **Zero-Config Discovery:** The app must discover the Linden device on the local network automatically without the user entering IP addresses.
2. **Secure Claiming:** The first user to connect should be able to establish themselves as the owner of the device.
3. **Customization:** The owner can assign a friendly name to the device (e.g., "Living Room AI"), which updates its mDNS broadcast.
4. **Offline First:** This entire process must happen without requiring an internet connection or external cloud accounts.

## Discovery Phase
As implemented in the MVP, the Linden server broadcasts an mDNS service (`_linden._tcp`) with a unique randomly-generated ID in the TXT records (e.g., `id=4F8A`).
1. The onboarding app (mobile or web client) scans the LAN for `_linden._tcp`.
2. It displays a list of discovered devices (e.g., `Linden AI (4F8A)`).
3. The user selects the device they want to configure.

## Claiming Phase (The "Ownership" Handshake)
To prevent unauthorized users on the same network from hijacking the configuration, the device must have a verifiable "claiming" step.
- **State:** The Linden device maintains an `IsClaimed` boolean in its persistent storage.
- **Unclaimed State:** If `IsClaimed` is false, the API allows an initial `POST /api/onboard` request.
- **The Handshake:**
  1. The app prompts the user to physically verify the device (e.g., matching the 4-digit ID printed on the bottom of the device to the `4F8A` broadcasted in mDNS).
  2. The app sends a `POST /api/onboard` request with an initial configuration payload (owner name, device name, etc.).
  3. The Linden device saves this configuration to persistent storage, sets `IsClaimed = true`, and generates a secure local token (or uses public-key cryptography) to authenticate future administrative actions from that specific app instance.
- **Claimed State:** Once claimed, the `/api/onboard` endpoint is locked. Any further configuration changes require the administrative token. (Note: Standard chat API endpoints might remain open to the LAN depending on the privacy policy, but admin endpoints are locked).

## Configuration Phase
During or after claiming, the user can configure:
1. **Custom Device Name:** The user renames the device from "Linden AI" to "Living Room AI".
2. **mDNS Update:** The server dynamically updates its Zeroconf broadcast. The Instance Name changes to `Living Room AI (4F8A)`.
3. **Network Persistence:** The unique ID (`4F8A`) remains constant, allowing the app to seamlessly track the device across IP changes or name changes.

## Resetting Ownership
If a user sells or moves the device, there must be a physical mechanism to reset it (e.g., holding a physical reset button for 10 seconds), which wipes the persistent storage and returns `IsClaimed` to false.
