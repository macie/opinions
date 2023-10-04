// Package security contains OS specific mitigation mechanisms.

//go:build !openbsd

package security

// Sandbox restrict access to system resources. Currently only works on OpenBSD.
func Sandbox() error {
	return nil
}
