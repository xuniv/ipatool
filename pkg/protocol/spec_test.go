package protocol

import "testing"

func TestProtocolVersionCompatibility(t *testing.T) {
	const expectedVersion = "1.0.0"
	if Version != expectedVersion {
		t.Fatalf("protocol version changed: got %q, want %q", Version, expectedVersion)
	}
}

func TestProtocolSchemaCompatibility(t *testing.T) {
	expected := []string{
		"appstore.lookup.v1",
		"appstore.search.v1",
		"error.v1",
	}

	if len(SupportedSchemas) != len(expected) {
		t.Fatalf("schema set changed: got %v, want %v", SupportedSchemas, expected)
	}

	for idx, value := range expected {
		if SupportedSchemas[idx] != value {
			t.Fatalf("schema mismatch at %d: got %q, want %q", idx, SupportedSchemas[idx], value)
		}
	}
}

func TestProtocolErrorCodeCompatibility(t *testing.T) {
	expected := []string{
		"AUTH_INVALID_CREDENTIALS",
		"AUTH_2FA_REQUIRED",
		"AUTH_SESSION_EXPIRED",
		"APP_NOT_FOUND",
		"APP_VERSION_NOT_FOUND",
		"DOWNLOAD_FAILED",
		"NETWORK_UNAVAILABLE",
		"RATE_LIMITED",
		"INTERNAL_ERROR",
	}

	if len(ErrorCodes) != len(expected) {
		t.Fatalf("error code set changed: got %v, want %v", ErrorCodes, expected)
	}

	for idx, value := range expected {
		if ErrorCodes[idx] != value {
			t.Fatalf("error code mismatch at %d: got %q, want %q", idx, ErrorCodes[idx], value)
		}
	}
}
