package protocol

import "encoding/json"

// ProtocolVersion identifies the wire protocol schema.
type ProtocolVersion string

const (
	// ProtocolVersionV1 is the initial wasm host/tool protocol.
	ProtocolVersionV1 ProtocolVersion = "1.0"
)

// StandardErrorCode is a normalized machine-readable error code set.
type StandardErrorCode string

const (
	ErrorCodeInvalidRequest      StandardErrorCode = "invalid_request"
	ErrorCodeUnsupportedProtocol StandardErrorCode = "unsupported_protocol"
	ErrorCodeUnauthorized        StandardErrorCode = "unauthorized"
	ErrorCodeForbidden           StandardErrorCode = "forbidden"
	ErrorCodeNotFound            StandardErrorCode = "not_found"
	ErrorCodeRateLimited         StandardErrorCode = "rate_limited"
	ErrorCodeConflict            StandardErrorCode = "conflict"
	ErrorCodeInternal            StandardErrorCode = "internal"
	ErrorCodeUnavailable         StandardErrorCode = "unavailable"
	ErrorCodeTimeout             StandardErrorCode = "timeout"
)

// CommandType is the high-level command namespace carried by a request.
type CommandType string

const (
	CommandSearch   CommandType = "search"
	CommandDownload CommandType = "download"
	CommandAuth     CommandType = "auth"
)

// Request is a host -> wasm command invocation envelope.
type Request struct {
	ProtocolVersion ProtocolVersion `json:"protocol_version"`
	ID              string          `json:"id"`
	Command         CommandType     `json:"command"`
	Input           json.RawMessage `json:"input"`
}

// Response is a wasm -> host command completion envelope.
type Response struct {
	ProtocolVersion ProtocolVersion `json:"protocol_version"`
	ID              string          `json:"id"`
	Command         CommandType     `json:"command"`
	Output          json.RawMessage `json:"output,omitempty"`
	Error           *ResponseError  `json:"error,omitempty"`
	BlobHandle      *HostBlobHandle `json:"blob_handle,omitempty"`
}

// ResponseError describes a standard protocol failure.
type ResponseError struct {
	Code    StandardErrorCode `json:"code"`
	Message string            `json:"message"`
	Details json.RawMessage   `json:"details,omitempty"`
}

// EventType represents a non-terminal stream event.
type EventType string

const (
	EventProgress   EventType = "progress"
	EventChunk      EventType = "chunk"
	EventBlobHandle EventType = "blob_handle"
)

// Event is a wasm -> host streaming envelope.
type Event struct {
	ProtocolVersion ProtocolVersion `json:"protocol_version"`
	RequestID       string          `json:"request_id"`
	Type            EventType       `json:"type"`
	Data            json.RawMessage `json:"data"`
}

// ChunkEvent provides an inline chunk transport for large payloads.
type ChunkEvent struct {
	TransferID     string `json:"transfer_id"`
	ChunkIndex     int    `json:"chunk_index"`
	TotalChunks    int    `json:"total_chunks"`
	DataBase64     string `json:"data_base64"`
	Encoding       string `json:"encoding,omitempty"`
	ContentType    string `json:"content_type,omitempty"`
	ChecksumSHA256 string `json:"checksum_sha256,omitempty"`
}

// HostBlobHandle provides host-owned out-of-band payload access.
type HostBlobHandle struct {
	Handle         string `json:"handle"`
	SizeBytes      int64  `json:"size_bytes"`
	ContentType    string `json:"content_type,omitempty"`
	ChecksumSHA256 string `json:"checksum_sha256,omitempty"`
}

// SearchInput defines request input for the search command.
type SearchInput struct {
	Term       string `json:"term"`
	Limit      int    `json:"limit,omitempty"`
	Storefront string `json:"storefront,omitempty"`
}

// SearchOutput defines response output for the search command.
type SearchOutput struct {
	Results []SearchResult `json:"results"`
}

// SearchResult represents one search hit.
type SearchResult struct {
	AppID    int64  `json:"app_id"`
	BundleID string `json:"bundle_id"`
	Name     string `json:"name"`
	Version  string `json:"version,omitempty"`
}

// DownloadInput defines request input for the download command.
type DownloadInput struct {
	AppID      int64  `json:"app_id"`
	Version    string `json:"version,omitempty"`
	OutputPath string `json:"output_path,omitempty"`
}

// DownloadOutput defines response output for the download command.
type DownloadOutput struct {
	FilePath   string `json:"file_path,omitempty"`
	SizeBytes  int64  `json:"size_bytes,omitempty"`
	BlobHandle string `json:"blob_handle,omitempty"`
}

// AuthInput defines request input for the auth command.
type AuthInput struct {
	Email         string `json:"email"`
	Password      string `json:"password,omitempty"`
	TwoFactorCode string `json:"two_factor_code,omitempty"`
}

// AuthOutput defines response output for the auth command.
type AuthOutput struct {
	AccountName  string `json:"account_name,omitempty"`
	DSID         string `json:"dsid,omitempty"`
	SessionToken string `json:"session_token,omitempty"`
	ExpiresAt    string `json:"expires_at,omitempty"`
}
