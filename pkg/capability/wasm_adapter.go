package capability

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type HostInvoker interface {
	Invoke(ctx context.Context, capability string, payload interface{}) (interface{}, error)
}

type WasmAdapterArgs struct {
	Policy Policy
	Host   HostInvoker
}

type WasmAdapter struct {
	policy Policy
	host   HostInvoker
}

func NewWasmAdapter(args WasmAdapterArgs) *WasmAdapter {
	return &WasmAdapter{policy: args.Policy, host: args.Host}
}

func (a *WasmAdapter) NetFetch(ctx context.Context, req *http.Request) (*http.Response, error) {
	if err := a.policy.Check(NameNetFetch); err != nil {
		return nil, err
	}
	if a.host == nil {
		return nil, ErrUnsupportedEnv
	}

	result, err := a.host.Invoke(ctx, NameNetFetch, req)
	if err != nil {
		return nil, fmt.Errorf("net.fetch failed: %w", err)
	}

	res, ok := result.(*http.Response)
	if !ok {
		return nil, ErrUnsupportedEnv
	}

	return res, nil
}

func (a *WasmAdapter) StorageGet(ctx context.Context, key string) ([]byte, error) {
	if err := a.policy.Check(NameStorageGet); err != nil {
		return nil, err
	}
	if a.host == nil {
		return nil, ErrUnsupportedEnv
	}

	result, err := a.host.Invoke(ctx, NameStorageGet, map[string]string{"key": key})
	if err != nil {
		return nil, fmt.Errorf("storage.get failed: %w", err)
	}

	data, ok := result.([]byte)
	if !ok {
		return nil, ErrUnsupportedEnv
	}

	return data, nil
}

func (a *WasmAdapter) StorageSet(ctx context.Context, key string, value []byte) error {
	if err := a.policy.Check(NameStorageSet); err != nil {
		return err
	}
	if a.host == nil {
		return ErrUnsupportedEnv
	}

	_, err := a.host.Invoke(ctx, NameStorageSet, map[string]interface{}{"key": key, "value": value})
	if err != nil {
		return fmt.Errorf("storage.set failed: %w", err)
	}

	return nil
}

func (a *WasmAdapter) StorageDelete(ctx context.Context, key string) error {
	if err := a.policy.Check(NameStorageDel); err != nil {
		return err
	}
	if a.host == nil {
		return ErrUnsupportedEnv
	}

	_, err := a.host.Invoke(ctx, NameStorageDel, map[string]string{"key": key})
	if err != nil {
		return fmt.Errorf("storage.delete failed: %w", err)
	}

	return nil
}

func (a *WasmAdapter) SecretGet(ctx context.Context, key string) ([]byte, error) {
	if err := a.policy.Check(NameSecretGet); err != nil {
		return nil, err
	}
	if a.host == nil {
		return nil, ErrUnsupportedEnv
	}

	result, err := a.host.Invoke(ctx, NameSecretGet, map[string]string{"key": key})
	if err != nil {
		return nil, fmt.Errorf("secret.get failed: %w", err)
	}

	data, ok := result.([]byte)
	if !ok {
		return nil, ErrUnsupportedEnv
	}

	return data, nil
}

func (a *WasmAdapter) ClockNow(ctx context.Context) (time.Time, error) {
	if err := a.policy.Check(NameClockNow); err != nil {
		return time.Time{}, err
	}
	if a.host == nil {
		return time.Time{}, ErrUnsupportedEnv
	}

	result, err := a.host.Invoke(ctx, NameClockNow, nil)
	if err != nil {
		return time.Time{}, fmt.Errorf("clock.now failed: %w", err)
	}

	now, ok := result.(time.Time)
	if !ok {
		return time.Time{}, ErrUnsupportedEnv
	}

	return now, nil
}
