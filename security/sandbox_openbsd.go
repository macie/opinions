//go:build openbsd

package security

import "golang.org/x/sys/unix"

// Sandbox restrict application access to necessary system calls needed by
// network connections and standard i/o.
//
// See: https://man.openbsd.org/pledge.2
func Sandbox() error {
	return unix.PledgePromises("stdio inet rpath")
}
