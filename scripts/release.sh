#!/usr/bin/env bash
set -euo pipefail

TAG=${1:-v1.0.0}
TITLE="ghprofile ${TAG}"
NOTES=${2:-"Release ${TAG}"}

# ensure release-bin exists
mkdir -p release-bin

# Build common targets if not present
./scripts/build-arch.sh linux amd64 || true
./scripts/build-arch.sh linux arm64 || true
./scripts/build-arch.sh darwin amd64 || true
./scripts/build-arch.sh darwin arm64 || true
./scripts/build-arch.sh windows amd64 || true

# Create annotated tag if doesn't exist
if ! git rev-parse ${TAG} >/dev/null 2>&1; then
  git tag -a ${TAG} -m "${TITLE}"
  git push origin ${TAG}
fi

# Create GitHub release and upload assets
gh release create ${TAG} release-bin/ghprofile-linux-amd64 release-bin/ghprofile-linux-arm64 release-bin/ghprofile-darwin-amd64 release-bin/ghprofile-darwin-arm64 release-bin/ghprofile-windows-amd64.exe --title "${TITLE}" --notes "${NOTES}"

echo "Release ${TAG} created with assets."
