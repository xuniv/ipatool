package protocol

const (
	Version = "1.0.0"
)

var SupportedSchemas = []string{
	"appstore.lookup.v1",
	"appstore.search.v1",
	"error.v1",
}

var ErrorCodes = []string{
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
