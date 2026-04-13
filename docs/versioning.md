# Version Compatibility Policy

## Scope
This policy defines compatibility expectations for:
- Protocol-level interactions (network request/response behavior with upstream services).
- WASM-based embeddings/integrations (where `ipatool` capabilities are wrapped).
- Native CLI usage on supported OS targets.

## Policy Baseline

### 1) Semantic Versioning
- `MAJOR`: breaking behavior or compatibility changes.
- `MINOR`: backward-compatible functionality additions.
- `PATCH`: backward-compatible fixes/docs/internal corrections.

### 2) Protocol Compatibility
- Patch/minor releases should preserve protocol behavior used by existing automations.
- Any unavoidable protocol break requires a major release note and migration guidance.
- Deprecated protocol fields/flows should be announced before removal whenever feasible.

### 3) WASM Compatibility Window
- Public WASM-facing contracts should remain backward-compatible for all patch/minor releases within the same major line.
- Contract-breaking changes (imports/exports, expected payload schema) require a major version bump.
- Host integrations should pin major versions and may float patch updates.

### 4) Native CLI Compatibility Window
- CLI flags and output formats are stable within a major line unless explicitly marked experimental.
- Removing or changing semantics of existing stable flags requires a major release.
- Default behavior changes that can affect automation should include release notes and an opt-in transition period when possible.

## Support Expectations
- Latest major release line receives active updates.
- Older major lines may receive limited or no fixes unless explicitly stated in release notes.
