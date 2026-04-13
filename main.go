package main

import (
	"os"

	"github.com/majd/ipatool/v2/internal/adapters/native"
)

func main() {
	os.Exit(native.Execute())
}
