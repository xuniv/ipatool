#!/usr/bin/env bash
set -euo pipefail

OUT_DIR="${1:-build/webui}"
mkdir -p "$OUT_DIR"

cat > "$OUT_DIR/index.html" <<'HTML'
<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>ipatool webui smoke bundle</title>
  </head>
  <body>
    <h1>ipatool webui smoke bundle</h1>
    <p>This static bundle is generated in CI for release packaging.</p>
    <script src="./wasm_exec.js"></script>
    <script>
      const go = new Go();
      WebAssembly.instantiateStreaming(fetch("./ipatool.wasm"), go.importObject)
        .then((result) => go.run(result.instance));
    </script>
  </body>
</html>
HTML

echo "WebUI build complete at $OUT_DIR"
