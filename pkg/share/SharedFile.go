package share

import (
	"fmt"
	"log"
	"regexp"
)

// -- types --

// [entity] a new shared file.
type SharedFile struct {
	source SharedSource
}

// -- lifetime --

// inits a new shared file from a source
func NewSharedFile(source SharedSource) *SharedFile {
	return &SharedFile{
		source: source,
	}
}

// -- queries --

// returns the html representation of the shared file
func (s *SharedFile) Render(id string) string {
	// strip whitespace
	reg, err := regexp.Compile("[\t\n]+")
	if err != nil {
		log.Fatal(err)
	}

	// deref the source
	// render the html
	html := fmt.Sprintf(
		reg.ReplaceAllLiteralString(`
			<html>
				<head>
					<!-- data -->
					<script id="mise-share" type="%s">%s</script>

					<!-- preview -->
					<meta property="og:title" content="Check out this recipe!">
					<meta property="og:type" content="website">
					<meta property="og:image" content="https://images.squarespace-cdn.com/content/v1/5ffb69ddfe0aa2509285f006/1614557697674-RZGESXDJJ4CMQLTIFOZI/Stirring.png">
					<meta property="og:url" content="https://share.miseapp.co/%s">
				</head>
			</html>
		`, ""),
		s.source.Type(),
		s.source.Value(),
		id,
	)

	return html
}
