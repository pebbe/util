package util

import (
	"log"
	"runtime"
)

func CheckErr(err error) {
	if err != nil {
		log.SetFlags(0)
		_, filename, lineno, ok := runtime.Caller(1)
		if ok {
			log.Fatalf("%v:%v: %v\n", filename, lineno, err)
		} else {
			log.Fatalln(err)
		}
	}
}

func WarnErr(err error) {
	if err != nil {
		f := log.Flags()
		log.SetFlags(0)
		_, filename, lineno, ok := runtime.Caller(1)
		if ok {
			log.Printf("%v:%v: %v\n", filename, lineno, err)
		} else {
			log.Println(err)
		}
		log.SetFlags(f)
	}
}
