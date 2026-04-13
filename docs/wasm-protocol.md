# WASM Host Protocol

`ipatool` wasm integration uses a JSON message protocol with explicit command input/output schemas.

## 1. Envelope schema

### Request

```json
{
  "protocol_version": "1.0",
  "id": "req-001",
  "command": "search",
  "input": {
    "term": "chatgpt",
    "limit": 5,
    "storefront": "143441"
  }
}
```

### Response (success)

```json
{
  "protocol_version": "1.0",
  "id": "req-001",
  "command": "search",
  "output": {
    "results": [
      {
        "app_id": 6448311069,
        "bundle_id": "com.openai.chat",
        "name": "ChatGPT",
        "version": "1.2026.0401"
      }
    ]
  }
}
```

### Response (error)

```json
{
  "protocol_version": "1.0",
  "id": "req-001",
  "command": "search",
  "error": {
    "code": "invalid_request",
    "message": "term is required"
  }
}
```

## 2. Protocol version

All messages MUST include `protocol_version`.

- Current version: `1.0`
- If a peer cannot process a version, return `unsupported_protocol`.

## 3. Standard error codes

- `invalid_request`
- `unsupported_protocol`
- `unauthorized`
- `forbidden`
- `not_found`
- `rate_limited`
- `conflict`
- `internal`
- `unavailable`
- `timeout`

## 4. Explicit command schemas

### Search

Input: `SearchInput`

```json
{
  "term": "chatgpt",
  "limit": 5,
  "storefront": "143441"
}
```

Output: `SearchOutput`

```json
{
  "results": [
    {
      "app_id": 6448311069,
      "bundle_id": "com.openai.chat",
      "name": "ChatGPT",
      "version": "1.2026.0401"
    }
  ]
}
```

### Download

Input: `DownloadInput`

```json
{
  "app_id": 6448311069,
  "version": "latest",
  "output_path": "/tmp/chatgpt.ipa"
}
```

Output: `DownloadOutput`

```json
{
  "file_path": "/tmp/chatgpt.ipa",
  "size_bytes": 74512000,
  "blob_handle": "blob-87f6"
}
```

### Auth

Input: `AuthInput`

```json
{
  "email": "user@example.com",
  "password": "***",
  "two_factor_code": "123456"
}
```

Output: `AuthOutput`

```json
{
  "account_name": "Example User",
  "dsid": "123456789",
  "session_token": "st-abc",
  "expires_at": "2026-04-14T10:00:00Z"
}
```

## 5. Large payload standardization

For large payloads, use one of the following two mechanisms:

1. **Chunk events** (`Event.type = "chunk"`)
   - Emit `ChunkEvent` payload with `transfer_id`, chunk indexes, and base64 data.
   - Receiver reassembles all chunks in order.

Example chunk event:

```json
{
  "protocol_version": "1.0",
  "request_id": "req-002",
  "type": "chunk",
  "data": {
    "transfer_id": "xfer-001",
    "chunk_index": 0,
    "total_chunks": 3,
    "data_base64": "AAECAw==",
    "encoding": "base64",
    "content_type": "application/octet-stream"
  }
}
```

2. **Host blob-handle** (`Event.type = "blob_handle"` or `Response.blob_handle`)
   - Emit a `HostBlobHandle` with host-resolvable `handle` and metadata.
   - Receiver calls host-side blob API to read bytes out-of-band.

Example blob-handle response:

```json
{
  "protocol_version": "1.0",
  "id": "req-003",
  "command": "download",
  "blob_handle": {
    "handle": "blob-87f6",
    "size_bytes": 74512000,
    "content_type": "application/octet-stream",
    "checksum_sha256": "..."
  }
}
```
