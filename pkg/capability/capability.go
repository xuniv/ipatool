package capability

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	NameNetFetch   = "net.fetch"
	NameStorageGet = "storage.get"
	NameStorageSet = "storage.set"
	NameStorageDel = "storage.delete"
	NameSecretGet  = "secret.get"
	NameClockNow   = "clock.now"
)

var ErrUnsupportedEnv = errors.New("ERR_UNSUPPORTED_ENV")

type Manifest struct {
	AllowNetFetch bool
	AllowStorage  map[string]bool
	AllowSecret   bool
	AllowClock    bool
}

type Policy struct {
	manifest Manifest
}

func NewPolicy(manifest Manifest) Policy {
	if manifest.AllowStorage == nil {
		manifest.AllowStorage = map[string]bool{}
	}

	return Policy{manifest: manifest}
}

func (p Policy) IsAllowed(name string) bool {
	switch name {
	case NameNetFetch:
		return p.manifest.AllowNetFetch
	case NameStorageGet:
		return p.manifest.AllowStorage[NameStorageGet]
	case NameStorageSet:
		return p.manifest.AllowStorage[NameStorageSet]
	case NameStorageDel:
		return p.manifest.AllowStorage[NameStorageDel]
	case NameSecretGet:
		return p.manifest.AllowSecret
	case NameClockNow:
		return p.manifest.AllowClock
	default:
		return false
	}
}

func (p Policy) Check(name string) error {
	if !p.IsAllowed(name) {
		return fmt.Errorf("capability %q is denied", name)
	}

	return nil
}

type Core interface {
	NetFetch(ctx context.Context, req *http.Request) (*http.Response, error)
	StorageGet(ctx context.Context, key string) ([]byte, error)
	StorageSet(ctx context.Context, key string, value []byte) error
	StorageDelete(ctx context.Context, key string) error
	SecretGet(ctx context.Context, key string) ([]byte, error)
	ClockNow(ctx context.Context) (time.Time, error)
}
