# tui-images installer for Windows (PowerShell)
# Usage: .\install.ps1 (from repo root)
#   or:  powershell -ExecutionPolicy Bypass -File install.ps1

$ErrorActionPreference = "Stop"

$BinaryName = "tui-images.exe"
$InstallDir = Join-Path $env:USERPROFILE "go\bin"

function Write-Info    { param($msg) Write-Host "[tui-images] $msg" -ForegroundColor Green }
function Write-Warn    { param($msg) Write-Host "[tui-images] $msg" -ForegroundColor Yellow }
function Write-Err     { param($msg) Write-Host "[tui-images] $msg" -ForegroundColor Red }

# Check Go installation
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Err "Go is not installed or not in PATH."
    Write-Err "Install it from: https://go.dev/doc/install"
    Write-Err "After installing, restart your terminal and run this script again."
    exit 1
}

$goVersion = go version
Write-Info "Using $goVersion"

# Create install directory if it doesn't exist
if (-not (Test-Path $InstallDir)) {
    Write-Info "Creating install directory: $InstallDir"
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

# Check if install directory is in PATH
$pathDirs = $env:PATH -split ';'
if ($pathDirs -notcontains $InstallDir) {
    Write-Warn "Adding $InstallDir to your PATH..."
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($currentPath -notlike "*$InstallDir*") {
        [Environment]::SetEnvironmentVariable("Path", "$currentPath;$InstallDir", "User")
        $env:PATH = "$env:PATH;$InstallDir"
        Write-Info "PATH updated. You may need to restart your terminal for changes to take effect."
    }
}

# Build
Write-Info "Building tui-images..."
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
Set-Location $ScriptDir

go build -ldflags="-s -w" -o $BinaryName ./cmd/main.go

if (-not (Test-Path $BinaryName)) {
    Write-Err "Build failed. Make sure you're running this from the repository root."
    exit 1
}

# Install
Write-Info "Installing to $InstallDir\$BinaryName..."
Copy-Item $BinaryName (Join-Path $InstallDir $BinaryName) -Force

# Cleanup
Remove-Item $BinaryName -ErrorAction SilentlyContinue

# Verify
if (Test-Path (Join-Path $InstallDir $BinaryName)) {
    Write-Info "Installation complete!"
    Write-Info ""
    Write-Info "If 'tui-images' is not recognized, restart your terminal or run:"
    Write-Info "  `$env:PATH += `";$InstallDir`""
    Write-Info ""
    Write-Info "Then run: tui-images"
} else {
    Write-Err "Installation failed. Please check the errors above."
    exit 1
}
