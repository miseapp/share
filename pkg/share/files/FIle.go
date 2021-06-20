package files

import (
	"crypto/md5"
	"io"
	"strings"
)

// -- types --

// a file w/ metadata
type File struct {
	Body   io.ReadSeeker
	Length int
	Hash   [16]byte
}

// -- impls --

// builds a file from a string
func NewFile(body string) File {
	return File{
		Body:   strings.NewReader(body),
		Length: len(body),
		Hash:   md5.Sum([]byte(body)),
	}
}
