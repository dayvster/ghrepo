#!/usr/bin/env bash
set -euo pipefail

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

echo "Building $OUT..."
GOOS=${OS} GOARCH=${ARCH} go build -o ${OUT} ./cmd/main.go
chmod +x ${OUT}
echo "Built ${OUT}" 
