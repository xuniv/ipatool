package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	manifestPath := flag.String("manifest", "./capability_manifest.json", "path to capability manifest")
	roundTripLatency := flag.Duration("roundtrip-latency", 50*time.Millisecond, "simulated roundtrip latency")
	timeout := flag.Duration("timeout", 2*time.Second, "timeout for host calls")
	retries := flag.Int("retries", 2, "retry attempts for host calls")
	flag.Parse()

	manifest, err := LoadCapabilityManifest(*manifestPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read manifest: %v\n", err)
		os.Exit(1)
	}

	host := NewEmbeddedHost(
		manifest,
		HostOptions{
			RoundTripLatency: *roundTripLatency,
			Timeout:          *timeout,
			Retries:          *retries,
		},
		&DemoTransport{},
	)

	fmt.Printf("loaded capabilities: network=%t storage=%t secret=%t\n",
		manifest.AllowNetwork,
		manifest.AllowStorage,
		manifest.AllowSecret,
	)

	ctx := context.Background()
	if err := runDemo(ctx, host); err != nil {
		var stdErr *StandardError
		if errors.As(err, &stdErr) {
			fmt.Fprintf(os.Stderr, "standard error: domain=%s code=%s message=%s\n", stdErr.Domain, stdErr.Code, stdErr.Message)
			os.Exit(2)
		}

		fmt.Fprintf(os.Stderr, "demo failed: %v\n", err)
		os.Exit(1)
	}
}

func runDemo(ctx context.Context, host *EmbeddedHost) error {
	hostFunctions := host.ExportedHostFunctions()
	fmt.Printf("exposed host functions: %v\n", hostFunctions)

	if err := host.NetworkCall(ctx, "https://example.com/health"); err != nil {
		return err
	}

	if err := host.StorageCall(ctx, "set:user", []byte("alice")); err != nil {
		return err
	}

	if err := host.SecretCall(ctx, "read:token"); err != nil {
		return err
	}

	result, _ := json.Marshal(map[string]string{"status": "ok"})
	fmt.Printf("wasm execution result: %s\n", result)
	return nil
}
