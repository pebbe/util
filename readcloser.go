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
	fp     *os.File
	gz     *gzip.Reader
	bz2    io.Reader
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
		rc.isGzip = true
	} else if strings.HasSuffix(filename, ".bz2") {
		rc.bz2 = bzip2.NewReader(rc.fp)
		rc.isBz2 = true
	}

	rc.isOpen = true
	return
}

func (rc *ReadCloser) Read(p []byte) (n int, err error) {
	if !rc.isOpen {
		panic("ReadCloser is closed")
	}
	if rc.isGzip {
		return rc.gz.Read(p)
	}
	if rc.isBz2 {
		return rc.bz2.Read(p)
	}
	return rc.fp.Read(p)
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
