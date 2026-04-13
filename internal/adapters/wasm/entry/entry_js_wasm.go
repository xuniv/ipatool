//go:build js && wasm

package entry

import "syscall/js"

// Start registers the wasm bridge entrypoint.
func Start() {
	done := make(chan struct{})
	js.Global().Set("ipatool", js.ValueOf(map[string]any{
		"status": "ready",
	}))
	<-done
}
