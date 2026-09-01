<#
.SYNOPSIS
Runs the Linden validation suite inside a temporary Docker container.

.DESCRIPTION
This script spins up a Linux container (golang:1.22-alpine), installs the required dependencies 
(make, nodejs, npm), mounts the repository root, and executes the validation targets.
This guarantees Linux/Docker parity and avoids Windows-specific execution constraints.

.EXAMPLE
.\scripts\docker-validate.ps1
#>

$ErrorActionPreference = "Stop"

Write-Host "Starting Docker validation suite..." -ForegroundColor Cyan

# Use golang:1.22-alpine as the base image since it's lightweight and we already know Go 1.22+ is required.
# We mount the current working directory to /workspace and run the Makefile targets.
docker run --rm -v "${PWD}:/workspace" -w /workspace golang:1.22-alpine sh -c "apk add --no-cache make nodejs npm && make lint test validate"

if ($LASTEXITCODE -eq 0) {
    Write-Host "Validation suite passed successfully!" -ForegroundColor Green
} else {
    Write-Host "Validation suite failed." -ForegroundColor Red
    exit $LASTEXITCODE
}
