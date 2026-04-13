package capability

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type NativeStorage interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Delete(key string) error
}

type NativeSecret interface {
	Get(key string) ([]byte, error)
}

type NativeAdapterArgs struct {
	Policy   Policy
	Client   *http.Client
	Storage  NativeStorage
	Secret   NativeSecret
	ClockNow func() time.Time
}

type NativeAdapter struct {
	policy   Policy
	client   *http.Client
	storage  NativeStorage
	secret   NativeSecret
	clockNow func() time.Time
}

func NewNativeAdapter(args NativeAdapterArgs) *NativeAdapter {
	return &NativeAdapter{
		policy:   args.Policy,
		client:   args.Client,
		storage:  args.Storage,
		secret:   args.Secret,
		clockNow: args.ClockNow,
	}
}

func (a *NativeAdapter) NetFetch(_ context.Context, req *http.Request) (*http.Response, error) {
	if err := a.policy.Check(NameNetFetch); err != nil {
		return nil, err
	}
	if a.client == nil {
		return nil, ErrUnsupportedEnv
	}

	res, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("net.fetch failed: %w", err)
	}

	return res, nil
}

func (a *NativeAdapter) StorageGet(_ context.Context, key string) ([]byte, error) {
	if err := a.policy.Check(NameStorageGet); err != nil {
		return nil, err
	}
	if a.storage == nil {
		return nil, ErrUnsupportedEnv
	}

	data, err := a.storage.Get(key)
	if err != nil {
		return nil, fmt.Errorf("storage.get failed: %w", err)
	}

	return data, nil
}

func (a *NativeAdapter) StorageSet(_ context.Context, key string, value []byte) error {
	if err := a.policy.Check(NameStorageSet); err != nil {
		return err
	}
	if a.storage == nil {
		return ErrUnsupportedEnv
	}

	if err := a.storage.Set(key, value); err != nil {
		return fmt.Errorf("storage.set failed: %w", err)
	}

	return nil
}

func (a *NativeAdapter) StorageDelete(_ context.Context, key string) error {
	if err := a.policy.Check(NameStorageDel); err != nil {
		return err
	}
	if a.storage == nil {
		return ErrUnsupportedEnv
	}

	if err := a.storage.Delete(key); err != nil {
		return fmt.Errorf("storage.delete failed: %w", err)
	}

	return nil
}

func (a *NativeAdapter) SecretGet(_ context.Context, key string) ([]byte, error) {
	if err := a.policy.Check(NameSecretGet); err != nil {
		return nil, err
	}
	if a.secret == nil {
		return nil, ErrUnsupportedEnv
	}

	data, err := a.secret.Get(key)
	if err != nil {
		return nil, fmt.Errorf("secret.get failed: %w", err)
	}

	return data, nil
}

func (a *NativeAdapter) ClockNow(_ context.Context) (time.Time, error) {
	if err := a.policy.Check(NameClockNow); err != nil {
		return time.Time{}, err
	}
	if a.clockNow == nil {
		return time.Time{}, ErrUnsupportedEnv
	}

	return a.clockNow(), nil
}
