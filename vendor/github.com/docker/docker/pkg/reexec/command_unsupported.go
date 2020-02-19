// +build !linux,!windows,!freebsd,!solaris,!darwin

package reexec

import (
	"os/exec"
)

// Comptcd is unsupported on operating systems apart from Linux, Windows, Solaris and Darwin.
func Comptcd(args ...string) *exec.Cmd {
	return nil
}
