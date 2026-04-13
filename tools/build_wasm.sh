#!/usr/bin/env bash
set -euo pipefail

OUT_DIR="${1:-build/wasm}"
mkdir -p "$OUT_DIR"

GOOS=js GOARCH=wasm go build -o "$OUT_DIR/ipatool.wasm" ./
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" "$OUT_DIR/wasm_exec.js"

echo "WASM build complete at $OUT_DIR"
