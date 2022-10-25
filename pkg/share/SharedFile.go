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
	host   *string
}

// -- lifetime --

// inits a new shared file from a source
func NewSharedFile(source SharedSource, host *string) *SharedFile {
	return &SharedFile{
		source: source,
		host:   host,
	}
}

// -- queries --

// returns the html representation of the shared file
func (s *SharedFile) Render(key string) string {
	// strip whitespace
	reg, err := regexp.Compile("[\t\n]+")
	if err != nil {
		log.Fatal(err)
	}

	// render the html
	// TODO: pull SHARE_FILES_HOST from .env instead of hardcoding https://share.miseapp.co
	html := fmt.Sprintf(
		reg.ReplaceAllLiteralString(`
			<html>
				<head>
					<!-- data -->
					<script id="mise-share" type="%[1]s">%[2]s</script>

					<!-- browser redirect -->
					<meta http-equiv="refresh" content="0; url=miseapp://share?%[3]s/%[4]s" />

					<!-- preview -->
					<meta property="og:title" content="Check out this recipe!">
					<meta property="og:type" content="website">
					<meta property="og:image" content="https://images.squarespace-cdn.com/content/v1/5ffb69ddfe0aa2509285f006/1614557697674-RZGESXDJJ4CMQLTIFOZI/Stirring.png">
					<meta property="og:url" content="%[3]s/%[4]s">
				</head>
			</html>
		`, ""),
		s.source.Type(),
		s.source.Value(),
		*s.host,
		key,
	)

	return html
}
