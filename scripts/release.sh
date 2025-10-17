gh release create ${TAG} release-bin/ghprofile-linux-amd64 release-bin/ghprofile-linux-arm64 release-bin/ghprofile-darwin-amd64 release-bin/ghprofile-darwin-arm64 release-bin/ghprofile-windows-amd64.exe --title "${TITLE}" --notes "${NOTES}"
#!/usr/bin/env bash
set -euo pipefail

# Colors
RED="\033[0;31m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
BLUE="\033[0;34m"
BOLD="\033[1m"
RESET="\033[0m"

info() { echo -e "${BLUE}[release]${RESET} $*"; }
success() { echo -e "${GREEN}[release]${RESET} $*"; }
warn() { echo -e "${YELLOW}[release]${RESET} $*"; }
error() { echo -e "${RED}[release]${RESET} $*"; }

TAG=${1:-v1.0.0}
TITLE="ghprofile ${TAG}"
NOTES=${2:-"Release ${TAG}"}

info "Starting release script"
info "Tag: ${TAG}"
info "Title: ${TITLE}"

# Trap errors to provide a helpful message
trap 'error "Release script failed (exit $?)"; exit 1' ERR

# ensure release-bin exists
mkdir -p release-bin

# Build common targets if not present
for arch in "linux:amd64" "linux:arm64" "darwin:amd64" "darwin:arm64" "windows:amd64"; do
  OS=${arch%%:*}
  ARCH=${arch##*:}
  OUT=release-bin/ghprofile-${OS}-${ARCH}
  if [ "$OS" = "windows" ]; then
    OUT=${OUT}.exe
  fi
  if [ ! -f "$OUT" ]; then
    echo "[release] Building missing binary: $OUT"
    ./scripts/build-arch.sh $OS $ARCH
  else
    echo "[release] Binary exists: $OUT"
  fi
done

echo "[release] Preparing git tag and push"
if git rev-parse ${TAG} >/dev/null 2>&1; then
  warn "Tag ${TAG} already exists locally"
else
  info "Creating tag ${TAG}"
  git tag -a ${TAG} -m "${TITLE}"
  git push origin ${TAG}
  success "Pushed tag ${TAG} to origin"
fi

info "Creating GitHub release with assets"
gh release create ${TAG} \
  release-bin/ghprofile-linux-amd64 \
  release-bin/ghprofile-linux-arm64 \
  release-bin/ghprofile-darwin-amd64 \
  release-bin/ghprofile-darwin-arm64 \
  release-bin/ghprofile-windows-amd64.exe \
  --title "${TITLE}" --notes "${NOTES}"

success "Release ${TAG} created with assets."
