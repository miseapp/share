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

// the content of a file
type FileContent interface {
	// render the content for an id
	Render(id string) string
}

// -- impls --

// builds a file from an id and its content
func NewFile(i int, content FileContent) (*File, error) {
	k := Key(i)

	// encode the file key to an id
	key, err := k.Encode()
	if err != nil {
		return nil, err
	}

	// render the body
	body := content.Render(key)

	// build the file
	file := File{
		Key:    key,
		Body:   strings.NewReader(body),
		Length: len(body),
		Hash:   md5.Sum([]byte(body)),
	}

	return &file, nil
}
