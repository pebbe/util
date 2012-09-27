package util

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

/*
### Opening a LinesReader ###

Example 1:

    r, err := util.NewLinesReaderFromFile(filename)
    util.CheckErr(err)

Example 2:

    r = util.NewLinesReaderFromReader(os.Stdin)

Example 3:

    f, e := os.Open(filename)
    util.CheckErr(e)
    defer f.Close()

    rd, e := gzip.NewReader(f)
    util.CheckErr(e)
    defer rd.Close()

    r = util.NewLinesReaderFromReader(rd)

### Using a LinesReader ###

Example 1:

    for line := range r.ReadLines() {
        // do something with line
    }

Example 2:

    for line := range r.ReadLines() {
        // do something with line

        // if you need to stop before all lines are read:
        r.Break()
        break     // not needed at bottom of loop

        // do more things
    }
*/
type LinesReader struct {
	r         *bufio.Reader
	f         *os.File
	isOpen    bool
	needClose bool
	interrupt chan bool
}

func NewLinesReaderFromFile(filename string) (r *LinesReader, err error) {
	r = &LinesReader{interrupt: make(chan bool)}
	r.f, err = os.Open(filename)
	if err != nil {
		return
	}
	r.r = bufio.NewReader(r.f)
	r.isOpen = true
	r.needClose = true
	return
}

func NewLinesReaderFromReader(rd io.Reader) (r *LinesReader) {

	r = &LinesReader{interrupt: make(chan bool)}
	r.r = bufio.NewReader(rd)
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

func (r *LinesReader) ReadLinesBytes() <-chan []byte {
	if !r.isOpen {
		panic("LinesReader is closed")
	}
	ch := make(chan []byte)
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

			s := make([]byte, buf.Len())
			copy(s, buf.Bytes())
			if !r.isOpen && len(s) == 0 {
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
	r.isOpen = false
	if r.needClose {
		r.needClose = false
		e := r.f.Close()
		if e != nil {
			panic(e)
		}
	}
}
