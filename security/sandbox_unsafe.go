//go:build unsafe

package security

// IsHardened reports whether security sandbox is enabled.
const IsHardened = false

// Sandbox restrict access to system resources. In unsafe builds sandbox is
// disabled.
func Sandbox() error {
	return nil
}
