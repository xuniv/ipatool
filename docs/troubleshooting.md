# Troubleshooting

## 1) Web CORS Issues

### Symptoms
- Browser-based integrations report CORS preflight failures.
- Requests succeed with CLI tools but fail in browser runtime.

### Why It Happens
- CORS is enforced by browsers, not by CLI binaries.
- Cross-origin request headers/methods may not be allowed by intermediary services.

### Mitigation
- Prefer direct CLI usage when possible.
- If using a browser bridge/backend, route requests through same-origin backend endpoints.
- Ensure `OPTIONS` preflight responses include expected allow-origin/allow-method/allow-header values.
- Avoid wildcard CORS configurations for credentialed flows.

## 2) Runtime Compatibility Issues

### Symptoms
- Binary fails to start due to architecture mismatch.
- Unexpected behavior across older/newer runtime environments.

### Checks
- Validate OS/architecture (`amd64`, `arm64`) matches the binary you installed.
- Confirm supported Go toolchain for local rebuilds.
- Rebuild from source in your target environment when in doubt.

### Mitigation
- Use release artifacts that explicitly match your runtime.
- Pin CI runner images and Go versions to reduce drift.
- Track compatibility expectations in the version policy document.

## 3) Permission Errors

### Symptoms
- `permission denied` when executing the binary.
- Write failures when saving downloaded IPA artifacts.
- Keychain access failures in restricted shells/containers.

### Mitigation
- Mark binary executable (`chmod +x ./ipatool`).
- Write output to user-owned directories.
- Verify keychain access policies and unlock state on host OS.
- In CI, avoid privileged paths and run with least-required permissions.
