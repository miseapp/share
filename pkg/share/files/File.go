package files

import (
	"crypto/md5"
	"io"
	"strings"
)

// -- types --

// a created file
type File struct {
	Key string
	Url string
}

// an input file w/ metadata
type FileInput struct {
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

// creates a file from key and url
func NewFile(key string, url string) *File {
	file := File{
		Key: key,
		Url: url,
	}

	return &file
}

// creates an input file from an id and content
func NewFileInput(i int, content FileContent) (*FileInput, error) {
	k := Key(i)

	// encode the file key to an id
	key, err := k.Encode()
	if err != nil {
		return nil, err
	}

	// render the body
	body := content.Render(key)

	// build the file
	file := FileInput{
		Key:    key,
		Body:   strings.NewReader(body),
		Length: len(body),
		Hash:   md5.Sum([]byte(body)),
	}

	return &file, nil
}
