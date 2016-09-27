// Error types for btoverlay2
// These implement the hashicorp errwrap library specification.

package btvfs

import (
	"fmt"
)

// Error type returned for unimplemented structs (temporarily).
type ErrUnimplemented struct {}

func (this *ErrUnimplemented) Error() string {
	return "Feature unimplemented"
}

func (this *ErrUnimplemented) WrappedErrors() []error {
	return []error{}
}

// Top-level error type for btoverlay2 errors.
type ErrBittorrentOverlay2Driver struct {
	filename string
	linenumber int
	caller string

	description string
	err error
}

func (this ErrBittorrentOverlay2Driver) Error() string {
	return fmt.Sprintf("%s : inner error: %s", this.description, this.err.Error())
}

func (this *ErrBittorrentOverlay2Driver) WrappedErrors() []error {
	return []error{this.err}
}
