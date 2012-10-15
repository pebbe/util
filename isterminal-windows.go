// +build windows

package util

import (
	"os"
	"syscall"
	"unsafe"
)

/*
func IsTerminal(file *os.File) (bool, error) {
	s, e := file.Stat()
	if e != nil {
		return true, nil
	}
	m := s.Mode()
	if m&os.ModeDevice != 0 {
		return true, nil
	}
	return false, nil
}
*/

/*
Examples:

    term, err := IsTerminal(os.Stdin)
    term, err := IsTerminal(os.Stdout)
    term, err := IsTerminal(os.Stderr)
*/
func IsTerminal(file *os.File) bool {
	var st uint32
	return getConsoleMode(file.Fd(), &st) == nil
}

func getConsoleMode(hConsoleHandle syscall.Handle, lpMode *uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procGetConsoleMode.Addr(), 2, uintptr(hConsoleHandle), uintptr(unsafe.Pointer(lpMode)), 0)
	if int(r1) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
