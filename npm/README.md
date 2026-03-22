# npm packaging

This directory contains the npm distribution layout for `@timerzz/ccv`.

Packages:

- `@timerzz/ccv`: main launcher package
- `@timerzz/ccv-linux-x64`
- `@timerzz/ccv-linux-arm64`
- `@timerzz/ccv-darwin-x64`
- `@timerzz/ccv-darwin-arm64`
- `@timerzz/ccv-win32-x64`
- `@timerzz/ccv-win32-arm64`

Build and sync binaries:

```bash
make npm-packages
```

Recommended publish order:

1. Publish all platform packages under `npm/platforms/*`
2. Publish the main package under `npm/ccv`

Publish with the helper script:

```bash
bash scripts/publish-npm.sh
```

Publish with an explicit token:

```bash
NPM_TOKEN=your_token_here bash scripts/publish-npm.sh
```

Dry run:

```bash
bash scripts/publish-npm.sh --dry-run
```

Local validation on Linux x64:

```bash
packdir="$(mktemp -d)"
npm pack ./npm/platforms/ccv-linux-x64 --pack-destination "$packdir"
npm pack ./npm/ccv --pack-destination "$packdir"
tmpdir="$(mktemp -d)"
npm install --prefix "$tmpdir" "$packdir"/timerzz-ccv-linux-x64-*.tgz "$packdir"/timerzz-ccv-*.tgz
"$tmpdir/node_modules/.bin/ccv" --help
```
