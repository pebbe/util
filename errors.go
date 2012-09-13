package util

import (
	"fmt"
	"os"
	"runtime"
)

func CheckErr(err error) {
	if err != nil {
		_, filename, lineno, ok := runtime.Caller(1)
		if ok {
			fmt.Fprintf(os.Stderr, "%v:%v: %v\n", filename, lineno, err)
		}
		panic(err)
	}
}
