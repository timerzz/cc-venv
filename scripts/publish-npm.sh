#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
NPM_DIR="$ROOT_DIR/npm"
MAIN_PACKAGE_DIR="$NPM_DIR/ccv"
PLATFORMS_DIR="$NPM_DIR/platforms"

BUILD_FIRST=1
DRY_RUN=0
DIST_TAG=""
TMP_NPMRC=""

usage() {
  cat <<'EOF'
Usage:
  bash scripts/publish-npm.sh [--dry-run] [--skip-build] [--tag <dist-tag>]

Environment:
  NPM_TOKEN     Optional npm token used explicitly for publish/auth operations

Options:
  --dry-run     Run npm publish with --dry-run
  --skip-build  Skip `make npm-packages`
  --tag <tag>   Publish using an npm dist-tag such as next or beta
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --dry-run)
      DRY_RUN=1
      shift
      ;;
    --skip-build)
      BUILD_FIRST=0
      shift
      ;;
    --tag)
      if [[ $# -lt 2 ]]; then
        echo "missing value for --tag" >&2
        exit 1
      fi
      DIST_TAG="$2"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "unknown argument: $1" >&2
      usage >&2
      exit 1
      ;;
  esac
done

cleanup() {
  if [[ -n "$TMP_NPMRC" && -f "$TMP_NPMRC" ]]; then
    rm -f "$TMP_NPMRC"
  fi
}

trap cleanup EXIT

pkg_name() {
  local pkg_dir="$1"
  node -p "require(process.argv[1]).name" "$pkg_dir/package.json"
}

pkg_version() {
  local pkg_dir="$1"
  node -p "require(process.argv[1]).version" "$pkg_dir/package.json"
}

published_version() {
  local pkg_name="$1"
  local version

  if version="$(npm view "$pkg_name" version 2>/dev/null)"; then
    printf '%s\n' "$version"
    return 0
  fi

  return 1
}

main_pkg_name="$(pkg_name "$MAIN_PACKAGE_DIR")"
main_version="$(pkg_version "$MAIN_PACKAGE_DIR")"

platform_dirs=(
  "$PLATFORMS_DIR/ccv-linux-x64"
  "$PLATFORMS_DIR/ccv-linux-arm64"
  "$PLATFORMS_DIR/ccv-darwin-x64"
  "$PLATFORMS_DIR/ccv-darwin-arm64"
  "$PLATFORMS_DIR/ccv-win32-x64"
  "$PLATFORMS_DIR/ccv-win32-arm64"
)

for dir in "${platform_dirs[@]}"; do
  version="$(pkg_version "$dir")"
  if [[ "$version" != "$main_version" ]]; then
    echo "version mismatch: $(pkg_name "$dir") has $version, expected $main_version" >&2
    exit 1
  fi
done

if [[ "$BUILD_FIRST" -eq 1 ]]; then
  echo "==> Building binaries and syncing npm package contents"
  make -C "$ROOT_DIR" npm-packages
fi

if [[ -n "${NPM_TOKEN:-}" ]]; then
  TMP_NPMRC="$(mktemp /tmp/ccv-npmrc.XXXXXX)"
  printf '%s\n' "//registry.npmjs.org/:_authToken=${NPM_TOKEN}" "access=public" > "$TMP_NPMRC"
  export NPM_CONFIG_USERCONFIG="$TMP_NPMRC"
  echo "==> Using npm token from NPM_TOKEN"
fi

publish_args=(publish --access public)
if [[ -n "$DIST_TAG" ]]; then
  publish_args+=(--tag "$DIST_TAG")
fi
if [[ "$DRY_RUN" -eq 1 ]]; then
  publish_args+=(--dry-run)
fi

echo "==> Publishing platform packages"
for dir in "${platform_dirs[@]}"; do
  name="$(pkg_name "$dir")"
  version="$(pkg_version "$dir")"
  if existing="$(published_version "$name")" && [[ "$existing" == "$version" ]]; then
    echo "   -> ${name}@${version} already published, skipping"
    continue
  fi

  echo "   -> ${name}@${version}"
  (
    cd "$dir"
    npm "${publish_args[@]}"
  )
done

echo "==> Publishing main package"
if existing="$(published_version "$main_pkg_name")" && [[ "$existing" == "$main_version" ]]; then
  echo "   -> ${main_pkg_name}@${main_version} already published, skipping"
else
  echo "   -> ${main_pkg_name}@${main_version}"
  (
    cd "$MAIN_PACKAGE_DIR"
    npm "${publish_args[@]}"
  )
fi

echo "==> Done"
