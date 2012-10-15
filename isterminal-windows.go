// +build windows

package util

import (
	"os"
)

/*
So far, tested on: linux windows

Examples:

    term, err := IsTerminal(os.Stdin)
    term, err := IsTerminal(os.Stdout)
    term, err := IsTerminal(os.Stderr)
*/
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
