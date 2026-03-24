#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
NPM_DIR="$ROOT_DIR/npm"
MAIN_PACKAGE_JSON="$NPM_DIR/ccv/package.json"

usage() {
  cat <<'EOF'
Usage:
  bash scripts/bump-npm-version.sh <version>

Example:
  bash scripts/bump-npm-version.sh 0.1.1
EOF
}

if [[ $# -ne 1 ]]; then
  usage >&2
  exit 1
fi

NEW_VERSION="$1"

if [[ ! "$NEW_VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+([-.][0-9A-Za-z.-]+)?$ ]]; then
  echo "invalid version: $NEW_VERSION" >&2
  exit 1
fi

mapfile -t PACKAGE_JSON_FILES < <(find "$NPM_DIR" -name package.json | sort)

if [[ ${#PACKAGE_JSON_FILES[@]} -eq 0 ]]; then
  echo "no package.json files found under $NPM_DIR" >&2
  exit 1
fi

for pkg_json in "${PACKAGE_JSON_FILES[@]}"; do
  node -e '
    const fs = require("fs");
    const path = process.argv[1];
    const version = process.argv[2];
    const data = JSON.parse(fs.readFileSync(path, "utf8"));

    data.version = version;

    if (path.endsWith("/npm/ccv/package.json") || path.endsWith("\\npm\\ccv\\package.json")) {
      if (data.optionalDependencies && typeof data.optionalDependencies === "object") {
        for (const name of Object.keys(data.optionalDependencies)) {
          data.optionalDependencies[name] = version;
        }
      }
    }

    fs.writeFileSync(path, JSON.stringify(data, null, 2) + "\n");
  ' "$pkg_json" "$NEW_VERSION"

  rel_path="${pkg_json#$ROOT_DIR/}"
  echo "updated ${rel_path} -> $NEW_VERSION"
done

echo "npm package versions are now synced to $NEW_VERSION"
