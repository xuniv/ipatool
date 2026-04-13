# Security

## Threat Model

### Assets
- Apple ID credentials and session artifacts.
- App metadata lookup results and downloaded IPA files.
- Local keychain and any credential material used for authentication.

### Trust Boundaries
- Client host running `ipatool`.
- Apple APIs and transport over HTTPS.
- Local storage (filesystem and keychain integrations).
- CI/runtime environments where non-interactive commands may run.

### Representative Threats
- Credential leakage through shell history, logs, CI logs, or environment dumps.
- Man-in-the-middle attempts when traffic is intercepted by untrusted proxies.
- Local privilege abuse where another user/process reads stored secrets.
- Accidental data exposure by writing output artifacts to world-readable paths.

### Security Controls
- Prefer non-interactive secure secret injection over plain command-line literals.
- Keep verbose logging disabled unless needed for diagnostics.
- Use OS-provided secret storage where available (e.g., keychain integration).
- Restrict filesystem permissions on output directories and credential material.
- Rotate compromised credentials and revoke App Store sessions when exposure is suspected.

## Secret Handling Policy

### Allowed Secret Sources
- Interactive prompt entry.
- Secure OS credential store.
- CI secret manager variables with masked log output.

### Prohibited Patterns
- Committing secrets to source control.
- Hardcoding secrets in scripts or checked-in config files.
- Printing raw credentials or session tokens to standard output.

### Operational Requirements
- Redact secrets from bug reports and support artifacts.
- Scope credentials to least privilege and shortest feasible lifetime.
- Revoke and replace credentials after any confirmed or suspected leak.
