package util

import (
	"compress/bzip2"
	"compress/gzip"
	"io"
	"os"
	"strings"
)

// Implements io.ReadCloser
type ReadCloser struct {
	r      io.Reader
	fp     *os.File
	gz     *gzip.Reader
	isOpen bool
	isGzip bool
	isBz2  bool
}

// Opens for reading a plain file, gzip'ed file (extension .gz), or bzip2'ed file (extension .bz2)
func Open(filename string) (rc *ReadCloser, err error) {
	rc = &ReadCloser{}

	rc.fp, err = os.Open(filename)
	if err != nil {
		return
	}

	if strings.HasSuffix(filename, ".gz") {
		rc.gz, err = gzip.NewReader(rc.fp)
		if err != nil {
			rc.fp.Close()
			return
		}
		rc.r = rc.gz
		rc.isGzip = true
	} else if strings.HasSuffix(filename, ".bz2") {
		rc.r = bzip2.NewReader(rc.fp)
		rc.isBz2 = true
	} else {
		rc.r = rc.fp
	}

	rc.isOpen = true
	return
}

func (rc *ReadCloser) Read(p []byte) (n int, err error) {
	if !rc.isOpen {
		panic("ReadCloser is closed")
	}
	return rc.r.Read(p)
}

func (rc *ReadCloser) Close() (err error) {
	if rc.isOpen {
		rc.isOpen = false
		if rc.isGzip {
			err = rc.gz.Close()
		}
		err2 := rc.fp.Close()
		if err == nil {
			err = err2
		}
	}
	return
}
