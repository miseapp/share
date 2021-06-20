package share

import (
	"fmt"
	"io"
	"log"
	"regexp"
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
func (s *SharedFile) AsHtml() string {
	reg, err := regexp.Compile("[\t\n]+")
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf(reg.ReplaceAllLiteralString(`
		<html>
			<head>
				<meta name="mise-share-url" content="%s">
			</head>
		</html>
	`, ""), *s.source.Url)
}
