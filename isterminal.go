// +build darwin

package util

import (
	"os"
	"errors"
)

/*
So far, tested on: linux windows

Examples:

    term, err := IsTerminal(os.Stdin)
    term, err := IsTerminal(os.Stdout)
    term, err := IsTerminal(os.Stderr)
*/
func IsTerminal(file *os.File) (bool, error) {
	return false, errors.New("Not implemented")
}
