package main

import (
	"fmt"
	"github.com/pebbe/util"
	"os"
)

func main() {
	fmt.Print("stdin:  ")
	if tty, err := util.IsTerminal(os.Stdin); err == nil {
		fmt.Println(tty)
	} else {
		fmt.Println(err)
	}

	fmt.Print("stdout: ")
	if tty, err := util.IsTerminal(os.Stdout); err == nil {
		fmt.Println(tty)
	} else {
		fmt.Println(err)
	}

	fmt.Print("stderr: ")
	if tty, err := util.IsTerminal(os.Stderr); err == nil {
		fmt.Println(tty)
	} else {
		fmt.Println(err)
	}
}
