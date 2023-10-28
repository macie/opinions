// Package security contains OS specific mitigation mechanisms.

//go:build !(linux && amd64) && !openbsd && !unsafe

package security

import (
	"fmt"
	"runtime"
)

// IsHardened reports whether security sandbox is enabled.
const IsHardened = false

// Sandbox restrict access to system resources.
func Sandbox() error {
	return fmt.Errorf("security sandbox is unavailable on %s/%s. To use app on this platform, compile it without sandbox (with 'unsafe' flag)", runtime.GOOS, runtime.GOARCH)
}
