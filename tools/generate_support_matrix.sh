#!/usr/bin/env bash
set -euo pipefail

OUT_FILE="${1:-build/support-matrix.md}"
mkdir -p "$(dirname "$OUT_FILE")"

cat > "$OUT_FILE" <<'MD'
## 지원/제한 기능 매트릭스

| 영역 | 지원 상태 | 비고 |
| --- | --- | --- |
| Native CLI (Windows/Linux/macOS) | ✅ 지원 | 공식 릴리스 바이너리 제공 |
| WASM (`GOOS=js`, `GOARCH=wasm`) | ✅ 지원 | `ipatool.wasm` 아티팩트 포함 |
| WebUI 정적 번들 | ✅ 지원 | `index.html`, `wasm_exec.js`, `.wasm` 포함 |
| 제한 CLI 스모크 (`--non-interactive --help`) | ✅ 지원 | E2E 자동 검증 |
| 임베디드 실행 (`cmd.Execute`) | ✅ 지원 | E2E 자동 검증 |
| 네트워크 의존 명령 실동작 (`auth/login/download`) | ⚠️ 제한 | CI에서는 외부 자격 증명 미사용 |
MD

echo "Support matrix generated at $OUT_FILE"
