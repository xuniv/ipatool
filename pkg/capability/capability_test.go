package capability

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

type hostStub struct {
	result interface{}
	err    error
}

func (s hostStub) Invoke(_ context.Context, _ string, _ interface{}) (interface{}, error) {
	return s.result, s.err
}

func TestPolicyDenyAllByDefault(t *testing.T) {
	policy := NewPolicy(Manifest{})

	if policy.IsAllowed(NameNetFetch) {
		t.Fatalf("expected %s to be denied", NameNetFetch)
	}
	if policy.IsAllowed(NameStorageGet) {
		t.Fatalf("expected %s to be denied", NameStorageGet)
	}
	if policy.IsAllowed(NameStorageSet) {
		t.Fatalf("expected %s to be denied", NameStorageSet)
	}
	if policy.IsAllowed(NameStorageDel) {
		t.Fatalf("expected %s to be denied", NameStorageDel)
	}
	if policy.IsAllowed(NameSecretGet) {
		t.Fatalf("expected %s to be denied", NameSecretGet)
	}
	if policy.IsAllowed(NameClockNow) {
		t.Fatalf("expected %s to be denied", NameClockNow)
	}
}

func TestNativeAdapterUnsupported(t *testing.T) {
	policy := NewPolicy(Manifest{
		AllowClock:    true,
		AllowNetFetch: true,
		AllowSecret:   true,
		AllowStorage: map[string]bool{
			NameStorageGet: true,
			NameStorageSet: true,
			NameStorageDel: true,
		},
	})
	adapter := NewNativeAdapter(NativeAdapterArgs{Policy: policy})

	if _, err := adapter.NetFetch(context.Background(), &http.Request{}); !errors.Is(err, ErrUnsupportedEnv) {
		t.Fatalf("expected unsupported error for net.fetch, got %v", err)
	}
	if _, err := adapter.StorageGet(context.Background(), "a"); !errors.Is(err, ErrUnsupportedEnv) {
		t.Fatalf("expected unsupported error for storage.get, got %v", err)
	}
	if err := adapter.StorageSet(context.Background(), "a", []byte("b")); !errors.Is(err, ErrUnsupportedEnv) {
		t.Fatalf("expected unsupported error for storage.set, got %v", err)
	}
	if err := adapter.StorageDelete(context.Background(), "a"); !errors.Is(err, ErrUnsupportedEnv) {
		t.Fatalf("expected unsupported error for storage.delete, got %v", err)
	}
	if _, err := adapter.SecretGet(context.Background(), "a"); !errors.Is(err, ErrUnsupportedEnv) {
		t.Fatalf("expected unsupported error for secret.get, got %v", err)
	}
	if _, err := adapter.ClockNow(context.Background()); !errors.Is(err, ErrUnsupportedEnv) {
		t.Fatalf("expected unsupported error for clock.now, got %v", err)
	}
}

func TestWasmAdapterUsesSameCoreAPI(t *testing.T) {
	now := time.Now().UTC()
	policy := NewPolicy(Manifest{
		AllowClock:    true,
		AllowNetFetch: true,
		AllowSecret:   true,
		AllowStorage: map[string]bool{
			NameStorageGet: true,
			NameStorageSet: true,
			NameStorageDel: true,
		},
	})

	clockAdapter := NewWasmAdapter(WasmAdapterArgs{Policy: policy, Host: hostStub{result: now}})
	gotNow, err := clockAdapter.ClockNow(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !gotNow.Equal(now) {
		t.Fatalf("expected %v got %v", now, gotNow)
	}

	storageAdapter := NewWasmAdapter(WasmAdapterArgs{Policy: policy, Host: hostStub{result: []byte("ok")}})
	data, err := storageAdapter.StorageGet(context.Background(), "k")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(data) != "ok" {
		t.Fatalf("expected ok, got %s", string(data))
	}
}
