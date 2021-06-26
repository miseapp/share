package files

import (
	"crypto/md5"
	"io"
	"strings"
)

// -- types --

// a file w/ metadata
type File struct {
	Key    string
	Body   io.ReadSeeker
	Length int
	Hash   [16]byte
}

type FileContent interface {
	ToBody(key string) string
}

// -- impls --

// builds a file from an index and its content
func NewFile(i int, content FileContent) (*File, error) {
	key := Key(i)

	// try and encode the file key
	skey, err := key.Encode()
	if err != nil {
		return nil, err
	}

	// render the body
	body := content.ToBody(skey)

	// build the file
	file := File{
		Key:    skey,
		Body:   strings.NewReader(body),
		Length: len(body),
		Hash:   md5.Sum([]byte(body)),
	}

	return &file, nil
}
