# Architecture

`ipatool` is organized around a simple layered structure:

- `internal/core`: business use-cases and orchestration boundary.
- `internal/adapters/native`: native CLI adapter.
- `internal/adapters/wasm`: browser/embedded wasm adapter.
- `web/ui`: frontend UI layer that talks to wasm bridge.
- `cmd`: native command definitions kept thin and focused on invoking core.

## Call flow diagrams

### 1) Web flow

```text
web/ui
  -> internal/adapters/wasm/entry
    -> internal/core
      -> pkg/appstore + pkg/http + pkg/keychain
```

### 2) Restricted CLI flow

```text
main.go
  -> internal/adapters/native.Execute
    -> cmd/* (flag parsing / minimal command wiring)
      -> internal/core
        -> pkg/appstore
```

### 3) Embedded flow

```text
host process (embedding wasm)
  -> internal/adapters/wasm/entry bridge
    -> internal/core
      -> domain operations (search/purchase/download/metadata)
```

## Notes

- Adapters should remain transport-focused (CLI/wasm/input-output translation).
- Business rules should converge in `internal/core` over time.
