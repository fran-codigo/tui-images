#!/usr/bin/env bash
set -euo pipefail

# tui-images installer for Linux/macOS
# Usage: curl -fsSL https://raw.githubusercontent.com/fran-codigo/tui-images/main/install.sh | bash
#   or:  ./install.sh (from repo root)

BINARY_NAME="tui-images"
REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INSTALL_DIR="/usr/local/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info()    { echo -e "${GREEN}[tui-images]${NC} $1"; }
warn()    { echo -e "${YELLOW}[tui-images]${NC} $1"; }
error()   { echo -e "${RED}[tui-images]${NC} $1"; }

# Check if running as root for system-wide install
if [ "$EUID" -ne 0 ]; then
    # Try user-local install
    if [ -d "$HOME/.local/bin" ]; then
        INSTALL_DIR="$HOME/.local/bin"
        info "Installing to $INSTALL_DIR (user-local)"
    else
        mkdir -p "$HOME/.local/bin"
        INSTALL_DIR="$HOME/.local/bin"
        info "Installing to $INSTALL_DIR (user-local)"
    fi
else
    info "Installing system-wide to $INSTALL_DIR"
fi

# Check Go installation
if ! command -v go &> /dev/null; then
    error "Go is not installed. Install it first: https://go.dev/doc/install"
    exit 1
fi

GO_VERSION=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+')
info "Using Go $(go version | grep -oP 'go\K[0-9.]+')"

# Build
info "Building $BINARY_NAME..."
cd "$REPO_DIR"
go build -ldflags="-s -w" -o "$BINARY_NAME" ./cmd/main.go

# Install
info "Installing to $INSTALL_DIR/$BINARY_NAME..."
cp "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Verify PATH
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    warn "$INSTALL_DIR is not in your PATH."
    warn "Add it by running: export PATH=\"\$PATH:$INSTALL_DIR\""
    warn "Add this line to your ~/.bashrc or ~/.zshrc to make it permanent."
else
    info "Installation complete! Run '$BINARY_NAME' to start."
fi

# Cleanup
rm -f "$REPO_DIR/$BINARY_NAME"
info "Done!"
