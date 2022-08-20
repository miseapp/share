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
func (s *SharedFile) ToBody(key string) string {
	reg, err := regexp.Compile("[\t\n]+")
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf(reg.ReplaceAllLiteralString(`
		<html>
			<head>
				<!-- data -->
				<meta name="mise-share-url" content="%s">

				<!-- preview -->
				<meta property="og:title" content="Check out this recipe!">
				<meta property="og:type" content="website">
				<meta property="og:image" content="https://images.squarespace-cdn.com/content/v1/5ffb69ddfe0aa2509285f006/1614557697674-RZGESXDJJ4CMQLTIFOZI/Stirring.png">
				<meta property="og:url" content="https://share.miseapp.co/%s">
			</head>
		</html>
	`, ""), *s.source.Url, key)
}
