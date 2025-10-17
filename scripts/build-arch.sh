#!/usr/bin/env bash
set -euo pipefail

# Colors
RED="\033[0;31m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
BLUE="\033[0;34m"
BOLD="\033[1m"
RESET="\033[0m"

info() { echo -e "${BLUE}[build-arch]${RESET} $*"; }
success() { echo -e "${GREEN}[build-arch]${RESET} $*"; }
warn() { echo -e "${YELLOW}[build-arch]${RESET} $*"; }
error() { echo -e "${RED}[build-arch]${RESET} $*"; }

# On any error, print a helpful message with context
trap 'error "Build failed for ${OS:-unknown}:${ARCH:-unknown} (exit $? )"; exit 1' ERR

if [ $# -lt 2 ]; then
  echo "Usage: $0 <os> <arch>"
  echo "Example: $0 linux amd64"
  exit 1
fi

OS=$1
ARCH=$2
OUT=release-bin/ghprofile-${OS}-${ARCH}
if [ "$OS" = "windows" ]; then
  OUT=${OUT}.exe
fi

info "Starting build for OS=${OS} ARCH=${ARCH}"
info "Output file: ${OUT}"
info "Go version: $(go version)"
info "Running: GOOS=${OS} GOARCH=${ARCH} go build -o ${OUT} ./cmd/main.go"
if GOOS=${OS} GOARCH=${ARCH} go build -o ${OUT} ./cmd/main.go; then
  success "Go build completed for ${OS}/${ARCH}"
else
  error "Go build failed for ${OS}/${ARCH}"
  exit 1
fi
chmod +x ${OUT} || true
ls -lh ${OUT} || warn "Could not list output file: ${OUT}"
success "Build finished: ${OUT}"
