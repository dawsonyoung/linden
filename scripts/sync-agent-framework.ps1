# Sync agent framework definitions from .agents to provider-specific paths.
# Source of truth is always .agents — never edit provider mirrors directly.
# Usage: .\scripts\sync-agent-framework.ps1 [-DryRun] [-Provider github]

param(
    [switch]$DryRun,
    [string[]]$Provider = @("github")
)

$ErrorActionPreference = "Stop"
$RepoRoot = Split-Path -Parent $PSScriptRoot
$AgentsRoot = Join-Path $RepoRoot ".agents"

function Write-Step($msg) { Write-Host "`n==> $msg" -ForegroundColor Cyan }
function Write-Ok($msg)   { Write-Host "    OK: $msg" -ForegroundColor Green }
function Write-Skip($msg) { Write-Host "    SKIP: $msg" -ForegroundColor Yellow }
function Write-Fail($msg) { Write-Host "    FAIL: $msg" -ForegroundColor Red; exit 1 }

# Maps .agents subdirectory -> provider destination
$ProviderMaps = @{
    "github" = @{
        "agents"       = ".github/agents"
        "instructions" = ".github/instructions"
        "prompts"      = ".github/prompts"
        # skills are not mirrored to github — consumed directly from .agents
    }
}

function Sync-Directory($srcDir, $destDir) {
    if (-not (Test-Path $srcDir)) { Write-Skip "Source not found: $srcDir"; return }

    if (-not (Test-Path $destDir)) {
        if ($DryRun) { Write-Host "  [dry-run] Would create: $destDir" }
        else         { New-Item -ItemType Directory -Path $destDir -Force | Out-Null }
    }

    foreach ($file in Get-ChildItem -Path $srcDir -File) {
        $dest = Join-Path $destDir $file.Name
        $srcHash  = (Get-FileHash $file.FullName -Algorithm SHA256).Hash
        $destHash = if (Test-Path $dest) { (Get-FileHash $dest -Algorithm SHA256).Hash } else { "" }

        if ($srcHash -eq $destHash) {
            Write-Skip "Up to date: $($file.Name)"
        } elseif ($DryRun) {
            Write-Host "  [dry-run] Would copy: $($file.FullName) -> $dest"
        } else {
            Copy-Item $file.FullName $dest -Force
            Write-Ok "Synced: $($file.Name)"
        }
    }
}

function Assert-InSync($srcDir, $destDir) {
    # Returns $true if all files in srcDir match destDir (used in CI)
    if (-not (Test-Path $srcDir)) { return $true }
    foreach ($file in Get-ChildItem -Path $srcDir -File) {
        $dest = Join-Path $destDir $file.Name
        if (-not (Test-Path $dest)) { return $false }
        $srcHash  = (Get-FileHash $file.FullName -Algorithm SHA256).Hash
        $destHash = (Get-FileHash $dest -Algorithm SHA256).Hash
        if ($srcHash -ne $destHash) { return $false }
    }
    return $true
}

foreach ($prov in $Provider) {
    if (-not $ProviderMaps.ContainsKey($prov)) {
        Write-Fail "Unknown provider: $prov. Supported: $($ProviderMaps.Keys -join ', ')"
    }

    Write-Step "Syncing to provider: $prov"
    $map = $ProviderMaps[$prov]

    foreach ($subdir in $map.Keys) {
        $src  = Join-Path $AgentsRoot $subdir
        $dest = Join-Path $RepoRoot $map[$subdir]
        Write-Host "  $subdir -> $($map[$subdir])"
        Sync-Directory $src $dest
    }
}

if (-not $DryRun) {
    Write-Host "`n========================================" -ForegroundColor Cyan
    Write-Host "Sync complete. Provider mirrors are up to date." -ForegroundColor Cyan
    Write-Host "DO NOT edit files under provider paths directly." -ForegroundColor Yellow
    Write-Host "========================================`n" -ForegroundColor Cyan
}
