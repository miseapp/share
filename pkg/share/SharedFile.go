package share

import (
	"io"
)

// -- types --

// [entity] a new shared file.
type SharedFile struct {
	source *Source
}

type File struct {
	Body   io.ReadSeeker
	Length int
	Hash   int
}

// [value] the source of a shared file's contents.
type Source struct {
	Url *string
}

// -- impls --
// inits a new shared file from a source
func NewSharedFile(source *Source) *SharedFile {
	return &SharedFile{
		source: source,
	}
}

// returns the html representation of the shared file
func (*SharedFile) AsHtml() string {
	return "html"
}
