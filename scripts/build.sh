#!/bin/sh
# Build packages in place from the repository root. End users run the bundled binaries.
set -eu
cd "$(dirname "$0")/.."
win=plugins/seek-windows
mkdir -p "$win/.codex-plugin" "$win/bin"
cp -R plugins/seek/assets plugins/seek/licenses plugins/seek/skills "$win/"
sed -e 's/"name": "seek"/"name": "seek-windows"/' -e 's/Trim Seek (macOS \/ Linux)/Trim Seek (Windows)/' plugins/seek/.codex-plugin/plugin.json > "$win/.codex-plugin/plugin.json"
for os in darwin windows linux; do
  for arch in arm64 amd64; do
    out="plugins/seek/bin/seek-mcp-$os-$arch"
    if [ "$os" = windows ]; then out="$win/bin/seek-mcp-$os-$arch.exe"; fi
    CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" go build -trimpath -buildvcs=false -ldflags='-s -w' -o "$out" ./cmd/seek-mcp
  done
done
