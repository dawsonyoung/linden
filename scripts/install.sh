#!/bin/sh
# Linden Linux Installer
# Single-command turnkey installer for barebones Linux systems.
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/dawsonyoung/linden/main/scripts/install.sh | bash
#   ./scripts/install.sh [--yes-ollama | --skip-ollama] [--model <model_name>] [--version <tag>]

set -eu

REPO="dawsonyoung/linden"
INSTALL_DIR="/usr/local/bin"
SERVICE_DIR="/etc/systemd/system"
STATE_DIR="/var/lib/linden"

YES_OLLAMA=0
SKIP_OLLAMA=0
CHOSEN_MODEL=""
VERSION="latest"

# Parse arguments
while [ $# -gt 0 ]; do
  case "$1" in
    --yes-ollama)
      YES_OLLAMA=1
      shift
      ;;
    --skip-ollama)
      SKIP_OLLAMA=1
      shift
      ;;
    --model)
      CHOSEN_MODEL="$2"
      shift 2
      ;;
    --version)
      VERSION="$2"
      shift 2
      ;;
    -h|--help)
      echo "Linden Linux Installer"
      echo "Usage: install.sh [options]"
      echo ""
      echo "Options:"
      echo "  --yes-ollama       Install Ollama automatically without prompting"
      echo "  --skip-ollama      Do not install or configure Ollama"
      echo "  --model <name>     Specify model to pull (default: auto based on RAM)"
      echo "  --version <tag>    Linden release version to install (default: latest)"
      exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      exit 1
      ;;
  esac
done

info()  { printf '\n\033[1;32m==>\033[0m \033[1m%s\033[0m\n' "$1"; }
warn()  { printf '\033[1;33m[!] %s\033[0m\n' "$1"; }
error() { printf '\033[1;31m[ERROR] %s\033[0m\n' "$1" >&2; exit 1; }

# Verify root/sudo privileges
SUDO=""
if [ "$(id -u)" -ne 0 ]; then
  if command -v sudo >/dev/null 2>&1; then
    SUDO="sudo"
  else
    error "Root privileges or sudo are required to install Linden system-wide."
  fi
fi

# Detect architecture
info "Detecting system architecture"
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)
    TARGET_ARCH="amd64"
    ;;
  aarch64|arm64)
    TARGET_ARCH="arm64"
    ;;
  *)
    error "Unsupported architecture: $ARCH (Linden supports x86_64 and aarch64/arm64)"
    ;;
esac
echo "    Architecture: $TARGET_ARCH"

# Detect memory for model recommendation
TOTAL_RAM_KB=$(awk '/MemTotal/ {print $2}' /proc/meminfo 2>/dev/null || echo 0)
TOTAL_RAM_GB=$((TOTAL_RAM_KB / 1024 / 1024))
echo "    Detected RAM: ${TOTAL_RAM_GB} GB"

if [ -z "$CHOSEN_MODEL" ]; then
  if [ "$TOTAL_RAM_GB" -ge 16 ]; then
    CHOSEN_MODEL="qwen2.5:7b"
  else
    CHOSEN_MODEL="qwen2.5:3b"
  fi
fi

# Install Linden Binary
info "Installing Linden single-binary runtime"
BINARY_NAME="linden-linux-${TARGET_ARCH}"
TMP_BIN="/tmp/linden-bin"

if [ "$VERSION" = "latest" ]; then
  DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME}"
else
  DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}"
fi

echo "    Fetching: $DOWNLOAD_URL"
if command -v curl >/dev/null 2>&1; then
  if ! curl -fsSL -o "$TMP_BIN" "$DOWNLOAD_URL" 2>/dev/null; then
    # If binary release doesn't exist yet (e.g. running from local clone), check local bin
    if [ -f "bin/linden" ]; then
      warn "Remote release not found. Using locally built bin/linden."
      cp "bin/linden" "$TMP_BIN"
    else
      error "Failed to download $DOWNLOAD_URL and no local bin/linden was found."
    fi
  fi
elif command -v wget >/dev/null 2>&1; then
  wget -qO "$TMP_BIN" "$DOWNLOAD_URL" || error "Failed to download $DOWNLOAD_URL"
else
  error "Neither curl nor wget is available. Please install curl first."
fi

chmod +x "$TMP_BIN"
$SUDO mkdir -p "$INSTALL_DIR"
$SUDO mv "$TMP_BIN" "${INSTALL_DIR}/linden"
echo "    Installed to: ${INSTALL_DIR}/linden"

# Ollama Detection & Prompting
info "Checking Ollama LLM backend"
if command -v ollama >/dev/null 2>&1; then
  echo "    OK: Ollama is already installed."
elif [ "$SKIP_OLLAMA" -eq 1 ]; then
  warn "Skipping Ollama installation as requested (--skip-ollama)."
  warn "Remember to set OLLAMA_URL to a reachable Ollama host."
else
  INSTALL_OLLAMA=0
  if [ "$YES_OLLAMA" -eq 1 ]; then
    INSTALL_OLLAMA=1
  else
    printf "\n"
    warn "Ollama was not detected on this system."
    printf "    Linden uses Ollama to execute privacy-first AI models locally.\n"
    printf "    Estimated download size: ~1.5GB to ~4GB (Ollama runtime + initial model).\n\n"
    printf "    Would you like Linden to install Ollama now? [y/N]: "
    # Read from /dev/tty to support pipe-to-bash execution (curl ... | bash)
    if [ -t 0 ]; then
      read -r response || response="n"
    elif [ -e /dev/tty ]; then
      read -r response < /dev/tty || response="n"
    else
      response="n"
    fi
    case "$response" in
      [yY][eE][sS]|[yY])
        INSTALL_OLLAMA=1
        ;;
      *)
        INSTALL_OLLAMA=0
        ;;
    esac
  fi

  if [ "$INSTALL_OLLAMA" -eq 1 ]; then
    info "Installing Ollama"
    curl -fsSL https://ollama.com/install.sh | sh
    if command -v systemctl >/dev/null 2>&1; then
      $SUDO systemctl enable --now ollama || true
    fi
  else
    warn "Ollama installation skipped."
    warn "Configure OLLAMA_URL if connecting to a remote or existing Ollama instance."
  fi
fi

# Pull recommended model if Ollama is available
if command -v ollama >/dev/null 2>&1 && [ "$SKIP_OLLAMA" -eq 0 ]; then
  info "Preparing recommended document model: ${CHOSEN_MODEL}"
  echo "    Checking if ${CHOSEN_MODEL} is downloaded..."
  if ! ollama list 2>/dev/null | grep -q "$CHOSEN_MODEL"; then
    echo "    Pulling ${CHOSEN_MODEL} (this may take a few minutes depending on your connection)..."
    ollama pull "$CHOSEN_MODEL" || warn "Could not pull ${CHOSEN_MODEL} immediately. You can pull it later using: ollama pull ${CHOSEN_MODEL}"
  else
    echo "    Model ${CHOSEN_MODEL} is already available."
  fi
fi

# Configure Systemd Service
if command -v systemctl >/dev/null 2>&1; then
  info "Configuring Linden systemd service"
  $SUDO mkdir -p "$STATE_DIR"
  $SUDO mkdir -p "$SERVICE_DIR"

  TMP_SVC="/tmp/linden.service"
  cat << 'EOF' > "$TMP_SVC"
[Unit]
Description=Linden: A Home AI Solution
Documentation=https://github.com/dawsonyoung/linden
After=network-online.target ollama.service
Wants=network-online.target
After=ollama.service

[Service]
Type=simple
ExecStart=/usr/local/bin/linden
Restart=on-failure
RestartSec=5s

# Security & Sandboxing
DynamicUser=yes
StateDirectory=linden
Environment=LINDEN_HOST=0.0.0.0
Environment=LINDEN_PORT=8080
Environment=LINDEN_DATA_DIR=/var/lib/linden
Environment=OLLAMA_URL=http://127.0.0.1:11434

NoNewPrivileges=true
ProtectSystem=full
ProtectHome=read-only
PrivateTmp=true

[Install]
WantedBy=multi-user.target
EOF

  $SUDO mv "$TMP_SVC" "${SERVICE_DIR}/linden.service"
  $SUDO systemctl daemon-reload
  $SUDO systemctl enable --now linden.service
  echo "    Service installed and started: linden.service"
fi

# Determine LAN IP
LAN_IP="localhost"
if command -v ip >/dev/null 2>&1; then
  DETECTED_IP=$(ip -4 route get 8.8.8.8 2>/dev/null | awk '{print $7}' | tr -d '\n')
  if [ -n "$DETECTED_IP" ]; then
    LAN_IP="$DETECTED_IP"
  fi
fi

# Completion Banner
printf '\n'
printf '\033[1;32m===================================================================\033[0m\n'
printf '\033[1;32m   Linden is successfully installed and running!                   \033[0m\n'
printf '\033[1;32m===================================================================\033[0m\n'
printf '\n'
printf '  Access your local home AI:\n'
printf '    • On this host:      \033[1mhttp://localhost:8080\033[0m\n'
printf '    • From home Wi-Fi:   \033[1mhttp://linden.local:8080\033[0m\n'
printf '    • Direct LAN IP:     \033[1mhttp://%s:8080\033[0m\n' "$LAN_IP"
printf '\n'
printf '  Manage service:\n'
printf '    • Check status:      \033[2msudo systemctl status linden\033[0m\n'
printf '    • View logs:         \033[2msudo journalctl -u linden -f\033[0m\n'
printf '    • Restart:           \033[2msudo systemctl restart linden\033[0m\n'
printf '\n'
