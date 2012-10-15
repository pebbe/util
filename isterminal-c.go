// +build darwin,cgo

package util

/*
#include <unistd.h>
*/
import "C"

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
	if int(C.isatty(C.int(file.Fd()))) == 0 {
		return false, nil
	}
	return true, nil
}
