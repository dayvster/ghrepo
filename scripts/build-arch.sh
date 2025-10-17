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
  echo "Examples: $0 linux amd64    # standard 64-bit"
  echo "          $0 linux 386      # 32-bit x86"
  echo "          $0 linux armv7     # 32-bit ARM (GOARCH=arm GOARM=7)"
  exit 1
fi

OS=$1
ARCH=$2

# Support arch names like armv6/armv7 -> GOARCH=arm GOARM=6/7
GOARCH="$ARCH"
GOARM=""
if [[ "$ARCH" =~ ^armv([0-9]+)$ ]]; then
  GOARCH="arm"
  GOARM="${BASH_REMATCH[1]}"
  OUT=release-bin/ghprofile-${OS}-armv${GOARM}
else
  OUT=release-bin/ghprofile-${OS}-${ARCH}
fi

if [ "$OS" = "windows" ]; then
  OUT=${OUT}.exe
fi

info "Starting build for OS=${OS} ARCH=${ARCH}"
info "Output file: ${OUT}"
info "Go version: $(go version)"
# Build with optimizations: disable cgo, strip debug/symbols via linker flags, and trim paths
LDFLAGS="-s -w"

info "Running: CGO_ENABLED=0 GOOS=${OS} GOARCH=${GOARCH} ${GOARM:+GOARM=${GOARM}} go build -trimpath -ldflags \"${LDFLAGS}\" -o ${OUT} ./cmd/main.go"
if [ -n "$GOARM" ]; then
  if CGO_ENABLED=0 GOOS=${OS} GOARCH=${GOARCH} GOARM=${GOARM} go build -trimpath -ldflags "${LDFLAGS}" -o ${OUT} ./cmd/main.go; then
    success "Go build completed for ${OS}/${ARCH} (GOARCH=${GOARCH} GOARM=${GOARM})"
  else
    error "Go build failed for ${OS}/${ARCH}"
    exit 1
  fi
else
  if CGO_ENABLED=0 GOOS=${OS} GOARCH=${GOARCH} go build -trimpath -ldflags "${LDFLAGS}" -o ${OUT} ./cmd/main.go; then
    success "Go build completed for ${OS}/${ARCH}"
  else
    error "Go build failed for ${OS}/${ARCH}"
    exit 1
  fi
fi

# Attempt to reduce size further: run strip if available and supports the binary format
if command -v strip >/dev/null 2>&1; then
  if strip --version >/dev/null 2>&1 || true; then
    if file ${OUT} | grep -Ei "ELF|PE32|Mach-O" >/dev/null 2>&1; then
      warn "Attempting to strip symbols from ${OUT} (may fail for cross targets)"
      if strip "${OUT}" >/dev/null 2>&1; then
        success "strip succeeded for ${OUT}"
      else
        warn "strip failed or is not compatible for ${OUT} — skipping"
      fi
    fi
  fi
fi

# Optionally compress with upx if available (install upx to enable). This reduces size but
# may be undesirable for some environments. We try and continue on failure.
if command -v upx >/dev/null 2>&1; then
  warn "upx found — attempting to compress ${OUT} with upx -9 (this is optional)"
  if upx -9 "${OUT}" >/dev/null 2>&1; then
    success "upx compression succeeded for ${OUT}"
  else
    warn "upx compression failed for ${OUT} — continuing without compression"
  fi
fi
  
chmod +x ${OUT} || true
ls -lh ${OUT} || warn "Could not list output file: ${OUT}"
success "Build finished: ${OUT}"
