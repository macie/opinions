// Package security contains OS specific mitigation mechanisms.

//go:build !openbsd

package security

// IsHardened reports whether security sandbox is enabled.
const IsHardened = false

// Sandbox restrict access to system resources.
func Sandbox() error {
	return nil
}
