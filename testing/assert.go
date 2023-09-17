// Package testing contains helpers for testable examples.
package testing

import "log"

// Must ensure no error.
func Must(err error) {
	if err != nil {
		log.Fatalf("function call returns error %v", err)
	}
}

// MustReturn ensure returning value without error.
func MustReturn[T any](val T, err error) T {
	if err != nil {
		log.Fatalf("function call returns error %v", err)
	}
	return val
}
