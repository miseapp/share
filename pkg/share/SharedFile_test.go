package share

import (
	"mise-share/pkg/share/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

// -- tests --
func TestRender_U(t *testing.T) {
	share := NewSharedFile(
		&Source{
			Url: test.Str("https://test.com"),
		},
	)

	html := share.Render("test")
	assert.Contains(t, html, `<meta name="mise-share-url" content="https://test.com">`)
	assert.Contains(t, html, `<meta property="og:url" content="https://share.miseapp.co/test">`)
}
