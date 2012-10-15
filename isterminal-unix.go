// +build linux darwin freebsd netbsd openbsd

package util

import (
	"os"
	"syscall"
	"unsafe"
)

/*
So far, tested on: linux windows

Examples:

    term, err := IsTerminal(os.Stdin)
    term, err := IsTerminal(os.Stdout)
    term, err := IsTerminal(os.Stderr)
*/
func IsTerminal(file *os.File) (bool, error) {
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, file.Fd(),
		uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	return err == 0, nil
}
