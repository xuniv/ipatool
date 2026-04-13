package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

var (
	ErrNetworkDenied = errors.New("network capability is denied by manifest")
	ErrStorageDenied = errors.New("storage capability is denied by manifest")
	ErrSecretDenied  = errors.New("secret capability is denied by manifest")
)

type CapabilityManifest struct {
	AllowNetwork bool `json:"allow_network"`
	AllowStorage bool `json:"allow_storage"`
	AllowSecret  bool `json:"allow_secret"`
}

func LoadCapabilityManifest(path string) (CapabilityManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return CapabilityManifest{}, err
	}

	var manifest CapabilityManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return CapabilityManifest{}, err
	}

	return manifest, nil
}

type HostOptions struct {
	RoundTripLatency time.Duration
	Timeout          time.Duration
	Retries          int
}

type InvocationResponse struct {
	Stderr []byte
	Err    error
}

type HostTransport interface {
	Call(ctx context.Context, operation string, payload []byte) InvocationResponse
}

type EmbeddedHost struct {
	manifest  CapabilityManifest
	options   HostOptions
	transport HostTransport
}

func NewEmbeddedHost(manifest CapabilityManifest, options HostOptions, transport HostTransport) *EmbeddedHost {
	if options.Timeout <= 0 {
		options.Timeout = time.Second
	}

	if options.Retries < 0 {
		options.Retries = 0
	}

	return &EmbeddedHost{
		manifest:  manifest,
		options:   options,
		transport: transport,
	}
}

func (h *EmbeddedHost) ExportedHostFunctions() []string {
	hostFunctions := []string{}
	if h.manifest.AllowNetwork {
		hostFunctions = append(hostFunctions, "host.network")
	}
	if h.manifest.AllowStorage {
		hostFunctions = append(hostFunctions, "host.storage")
	}
	if h.manifest.AllowSecret {
		hostFunctions = append(hostFunctions, "host.secret")
	}
	return hostFunctions
}

func (h *EmbeddedHost) NetworkCall(ctx context.Context, endpoint string) error {
	if !h.manifest.AllowNetwork {
		return ErrNetworkDenied
	}

	return h.invoke(ctx, "network", []byte(endpoint))
}

func (h *EmbeddedHost) StorageCall(ctx context.Context, key string, value []byte) error {
	if !h.manifest.AllowStorage {
		return ErrStorageDenied
	}

	payload := append([]byte(key+":"), value...)
	return h.invoke(ctx, "storage", payload)
}

func (h *EmbeddedHost) SecretCall(ctx context.Context, selector string) error {
	if !h.manifest.AllowSecret {
		return ErrSecretDenied
	}

	return h.invoke(ctx, "secret", []byte(selector))
}

func (h *EmbeddedHost) invoke(ctx context.Context, operation string, payload []byte) error {
	var lastErr error
	attempts := h.options.Retries + 1

	for attempt := 0; attempt < attempts; attempt++ {
		time.Sleep(h.options.RoundTripLatency)

		callCtx, cancel := context.WithTimeout(ctx, h.options.Timeout)
		response := h.transport.Call(callCtx, operation, payload)
		cancel()

		if response.Err == nil {
			return nil
		}

		stdErr, err := DeserializeStandardError(response.Stderr)
		if err == nil {
			lastErr = stdErr
		} else {
			lastErr = fmt.Errorf("operation %s failed: %w", operation, response.Err)
		}

		if errors.Is(response.Err, context.DeadlineExceeded) {
			lastErr = fmt.Errorf("operation %s timeout: %w", operation, response.Err)
		}
	}

	return lastErr
}

type DemoTransport struct{}

func (d *DemoTransport) Call(ctx context.Context, operation string, payload []byte) InvocationResponse {
	_ = payload

	select {
	case <-ctx.Done():
		return InvocationResponse{Err: ctx.Err()}
	default:
	}

	switch operation {
	case "network":
		return InvocationResponse{}
	case "storage":
		return InvocationResponse{}
	case "secret":
		return InvocationResponse{
			Stderr: []byte(`{"domain":"secret","code":"permission_denied","message":"token is not available"}`),
			Err:    errors.New("remote call failed"),
		}
	default:
		return InvocationResponse{Err: fmt.Errorf("unknown operation %s", operation)}
	}
}
