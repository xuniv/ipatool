package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeTransport struct {
	responses []InvocationResponse
	calls     int
}

func (f *fakeTransport) Call(ctx context.Context, operation string, payload []byte) InvocationResponse {
	_ = operation
	_ = payload
	select {
	case <-ctx.Done():
		return InvocationResponse{Err: ctx.Err()}
	default:
	}

	if f.calls >= len(f.responses) {
		return InvocationResponse{}
	}

	response := f.responses[f.calls]
	f.calls++
	return response
}

func TestEmbeddedHost_DeniesByManifest(t *testing.T) {
	host := NewEmbeddedHost(CapabilityManifest{}, HostOptions{}, &fakeTransport{})

	if err := host.NetworkCall(context.Background(), "https://example.com"); !errors.Is(err, ErrNetworkDenied) {
		t.Fatalf("expected network denied, got %v", err)
	}
	if err := host.StorageCall(context.Background(), "k", []byte("v")); !errors.Is(err, ErrStorageDenied) {
		t.Fatalf("expected storage denied, got %v", err)
	}
	if err := host.SecretCall(context.Background(), "selector"); !errors.Is(err, ErrSecretDenied) {
		t.Fatalf("expected secret denied, got %v", err)
	}
}

func TestEmbeddedHost_DeserializesStandardError(t *testing.T) {
	transport := &fakeTransport{responses: []InvocationResponse{{
		Stderr: []byte(`{"domain":"network","code":"forbidden","message":"blocked"}`),
		Err:    errors.New("failure"),
	}}}
	host := NewEmbeddedHost(CapabilityManifest{AllowNetwork: true}, HostOptions{}, transport)

	err := host.NetworkCall(context.Background(), "https://example.com")
	var stdErr *StandardError
	if !errors.As(err, &stdErr) {
		t.Fatalf("expected standard error, got %T %v", err, err)
	}
	if stdErr.Code != "forbidden" {
		t.Fatalf("unexpected code %s", stdErr.Code)
	}
}

func TestEmbeddedHost_RetryOption(t *testing.T) {
	transport := &fakeTransport{responses: []InvocationResponse{{Err: errors.New("temporary")}, {}}}
	host := NewEmbeddedHost(
		CapabilityManifest{AllowNetwork: true},
		HostOptions{Retries: 1, Timeout: time.Second},
		transport,
	)

	if err := host.NetworkCall(context.Background(), "https://example.com"); err != nil {
		t.Fatalf("expected retry to recover, got %v", err)
	}
	if transport.calls != 2 {
		t.Fatalf("expected 2 calls, got %d", transport.calls)
	}
}

func TestEmbeddedHost_TimeoutOption(t *testing.T) {
	host := NewEmbeddedHost(
		CapabilityManifest{AllowNetwork: true},
		HostOptions{Timeout: 10 * time.Millisecond},
		HostTransportFunc(func(ctx context.Context, operation string, payload []byte) InvocationResponse {
			_ = operation
			_ = payload
			time.Sleep(20 * time.Millisecond)
			select {
			case <-ctx.Done():
				return InvocationResponse{Err: ctx.Err()}
			default:
				return InvocationResponse{}
			}
		}),
	)

	err := host.NetworkCall(context.Background(), "https://example.com")
	if err == nil {
		t.Fatalf("expected timeout error")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded, got %v", err)
	}
}

type HostTransportFunc func(ctx context.Context, operation string, payload []byte) InvocationResponse

func (f HostTransportFunc) Call(ctx context.Context, operation string, payload []byte) InvocationResponse {
	return f(ctx, operation, payload)
}
