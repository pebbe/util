package util

import (
	"bytes"
	"io"
)

const (
	defaultReaderBufSize = 4096
	minReaderBufferSize  = 16
)

// Type implementing a robust line reader.
type Reader struct {
	rd  io.Reader
	buf []byte
	err error
	p   int // points to first unscanned byte in buffer
	n   int // points to location after last byte in buffer
}

// Create new reader with default buffer size.
func NewReader(rd io.Reader) *Reader {
	return NewReaderSize(rd, defaultReaderBufSize)
}

// Create new reader with specified buffer size.
func NewReaderSize(rd io.Reader, size int) *Reader {
	if size < minReaderBufferSize {
		size = minReaderBufferSize
	}
	r := Reader{
		buf: make([]byte, size),
		rd:  rd,
		n:   2,
	}
	r.buf[0] = '\r' // in case first line is empty
	r.buf[1] = '\n' //
	return &r
}

// Same as ReadLine(), but returns string instead of []byte.
// Result stays valid.
func (r *Reader) ReadLineString() (line string, err error) {
	b, e := r.ReadLine()
	return string(b), e
}

// Read single line from Reader, without EOL.
// EOL can be any of: \n \r \r\n \n\r.
// Last line without EOL is allowed.
// Result is only valid until next call to ReadLine() or ReadLineString()
func (r *Reader) ReadLine() (line []byte, err error) {
	var lines [][]byte
	beginning := true
	for {

		// Fill buffer if it is nearly empty, unless error was already received.
		// Part after the '||' is not necessary, but makes it (perhaps) more
		// efficient by avoiding the use of partial saves as much as possble.
		if r.n-r.p < 2 && r.err == nil ||
			r.p < r.n && r.err == nil && bytes.IndexByte(r.buf[r.p+2:r.n], '\n') < 0 && bytes.IndexByte(r.buf[r.p+2:r.n], '\r') < 0 {
			if r.p > 0 {
				// make room
				copy(r.buf, r.buf[r.p:r.n])
				r.n -= r.p
				r.p = 0
			}
			var i int
			i, r.err = r.rd.Read(r.buf[r.n:])
			r.n += i
		}

		// if at beginning of line, skip EOL of previous line
		if beginning {
			var c byte
			for i := 0; i < 2; i++ {
				// 1 or 2 times \r or \n, but not \r\r or \n\n
				if r.p < r.n && (r.buf[r.p] == '\n' || r.buf[r.p] == '\r') && r.buf[r.p] != c {
					c = r.buf[r.p]
					r.p += 1
				} else {
					break
				}
			}
			beginning = false
		}

		// if error and no more unscanned data, return saved parts if they exist, else return error
		if r.p == r.n && r.err != nil {
			if len(lines) > 0 {
				return buildline(lines, r.buf[0:0]), nil
			} else {
				return r.buf[0:0], r.err
			}
		}

		if r.p < r.n {
			// find next EOL, can start with \n or \r
			i := bytes.IndexByte(r.buf[r.p:r.n], '\n')
			if i < 0 {
				i = r.n
			}
			j := bytes.IndexByte(r.buf[r.p:r.n], '\r')
			if j < 0 {
				j = r.n
			}
			if j < i {
				i = j
			}
			if i == r.n { // no EOL found
				buf := make([]byte, r.n-r.p)
				copy(buf, r.buf[r.p:r.n])
				lines = append(lines, buf)
				r.p = r.n
			} else { // EOL found
				p := r.p
				r.p += i
				return buildline(lines, r.buf[p:r.p]), nil
			}
		}

	}
	panic("not reached")
}

func buildline(lines [][]byte, last []byte) []byte {
	if len(lines) == 0 {
		return last
	}

	i := len(last)
	for _, line := range lines {
		i += len(line)
	}
	buf := make([]byte, i)
	i = 0
	for _, line := range lines {
		copy(buf[i:], line)
		i += len(line)
	}
	copy(buf[i:], last)
	return buf
}
