#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
BUILD_DIR="$ROOT_DIR/build/e2e"
mkdir -p "$BUILD_DIR"

# 1) 제한 CLI 스모크 테스트
GO111MODULE=on go build -o "$BUILD_DIR/ipatool" "$ROOT_DIR"
"$BUILD_DIR/ipatool" --non-interactive --help >/dev/null

# 2) 웹 스모크 테스트
"$ROOT_DIR/tools/build_wasm.sh" "$BUILD_DIR/web"
"$ROOT_DIR/tools/build_webui.sh" "$BUILD_DIR/web"
test -f "$BUILD_DIR/web/index.html"
test -f "$BUILD_DIR/web/ipatool.wasm"

# 3) 임베디드 스모크 테스트
cat > "$BUILD_DIR/embedded_smoke.go" <<'GO'
package main

import (
	"fmt"
	"os"

	"github.com/majd/ipatool/v2/cmd"
)

func main() {
	os.Args = []string{"ipatool", "--help"}
	exitCode := cmd.Execute()
	if exitCode != 0 {
		panic(fmt.Sprintf("unexpected exit code: %d", exitCode))
	}
}
GO

go run "$BUILD_DIR/embedded_smoke.go" >/dev/null

echo "E2E smoke suite passed"
