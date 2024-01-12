package security

import "errors"

var (
	ErrNoSandbox = errors.New("security sandbox cannot be enabled")
)
