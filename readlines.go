package util

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

/*
Example usage:

    r, err := util.NewLinesReader(filename)
    util.CheckErr(err)
    for line := range r.ReadLines() {
        // do something with line

        // if you need to stop before all lines are read:
        r.Break()
        break     // not needed at end of loop

        // do something with line
    }
*/
type LinesReader struct {
	r         *bufio.Reader
	f         *os.File
	isOpen    bool
	interrupt chan bool
}

func NewLinesReader(filename string) (r *LinesReader, err error) {
	r = &LinesReader{interrupt: make(chan bool)}
	r.f, err = os.Open(filename)
	if err != nil {
		return
	}
	r.r = bufio.NewReader(r.f)
	r.isOpen = true
	return
}

func (r *LinesReader) ReadLines() <-chan string {
	if !r.isOpen {
		panic("LinesReader is closed")
	}
	ch := make(chan string)
	go func() {
	ReadLinesLoop:
		for {
			var buf bytes.Buffer
			if !r.isOpen {
				break ReadLinesLoop
			}
			for {
				line, isPrefix, err := r.r.ReadLine()
				buf.Write(line)
				if err == io.EOF {
					r.close()
					break
				}
				if err != nil {
					panic(err)
				}
				if !isPrefix {
					break
				}
			}
			s := buf.String()
			if !r.isOpen && s == "" {
				break ReadLinesLoop
			}
			select {
			case ch <- s:
			case <-r.interrupt:
				break ReadLinesLoop
			}
		}
		r.close()
		close(ch)
	}()
	return ch
}

func (r *LinesReader) Break() {
	if r.isOpen {
		r.interrupt <- true
	}
}

func (r *LinesReader) close() {
	if r.isOpen {
		r.isOpen = false
		e := r.f.Close()
		if e != nil {
			panic(e)
		}
	}
}
