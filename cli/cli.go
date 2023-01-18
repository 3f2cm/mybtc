/*
Package cli package provides utilities to build CLI
*/
package cli

import "io"

// Env expresses the environment of the command execution
// mainly for replacing IO for testing.
type Env struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Rand   io.Reader
}
