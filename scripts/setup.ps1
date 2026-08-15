# Linden Development Environment Setup (Windows)
# Run: .\scripts\setup.ps1
#
# Convenience only. The canonical setup script is scripts/setup.sh, and the
# recommended environment is the devcontainer, which matches CI exactly.
# See docs/adr/0003-target-platform-policy.md for when this file is removed.

$ErrorActionPreference = "Stop"

function Write-Step($msg) { Write-Host "`n==> $msg" -ForegroundColor Cyan }
function Write-Ok($msg)   { Write-Host "    OK: $msg" -ForegroundColor Green }
function Write-Skip($msg) { Write-Host "    SKIP: $msg" -ForegroundColor Yellow }
function Write-Fail($msg) { Write-Host "    FAIL: $msg" -ForegroundColor Red }

# --- Go ---
Write-Step "Checking Go"
if (Get-Command go -ErrorAction SilentlyContinue) {
    $goVer = go version
    Write-Ok $goVer
} else {
    Write-Step "Installing Go via winget"
    winget install GoLang.Go --accept-source-agreements --accept-package-agreements
    if ($LASTEXITCODE -ne 0) { Write-Fail "Go install failed"; exit 1 }
    Write-Ok "Go installed — restart your terminal to use it"
}

# --- Node.js ---
Write-Step "Checking Node.js"
if (Get-Command node -ErrorAction SilentlyContinue) {
    $nodeVer = node --version
    Write-Ok "Node $nodeVer"
} else {
    Write-Step "Installing Node.js LTS via winget"
    winget install OpenJS.NodeJS.LTS --accept-source-agreements --accept-package-agreements
    if ($LASTEXITCODE -ne 0) { Write-Fail "Node install failed"; exit 1 }
    Write-Ok "Node installed — restart your terminal to use it"
}

# --- Ollama ---
Write-Step "Checking Ollama"
if (Get-Command ollama -ErrorAction SilentlyContinue) {
    $ollamaVer = ollama --version
    Write-Ok "Ollama $ollamaVer"
} else {
    Write-Step "Installing Ollama via winget"
    winget install Ollama.Ollama --accept-source-agreements --accept-package-agreements
    if ($LASTEXITCODE -ne 0) { Write-Fail "Ollama install failed"; exit 1 }
    Write-Ok "Ollama installed — restart your terminal to use it"
}

# --- Pull model ---
Write-Step "Checking for tinyllama model"
if (Get-Command ollama -ErrorAction SilentlyContinue) {
    $models = ollama list 2>&1
    if ($models -match "tinyllama") {
        Write-Ok "tinyllama already pulled"
    } else {
        Write-Step "Pulling tinyllama (this may take a few minutes)"
        ollama pull tinyllama
        if ($LASTEXITCODE -ne 0) { Write-Fail "Model pull failed"; exit 1 }
        Write-Ok "tinyllama pulled"
    }
} else {
    Write-Skip "Ollama not in PATH yet — pull tinyllama after restarting terminal: ollama pull tinyllama"
}

# --- Summary ---
Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "Setup complete." -ForegroundColor Cyan
Write-Host "If any tools were just installed, restart your terminal then run:"
Write-Host "  ollama pull tinyllama    (if skipped above)"
Write-Host "  make dev                (to start developing)"
Write-Host "========================================`n" -ForegroundColor Cyan
