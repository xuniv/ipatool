//go:build js && wasm

package main

import "github.com/majd/ipatool/v2/internal/adapters/wasm/entry"

func main() {
	entry.Start()
}
