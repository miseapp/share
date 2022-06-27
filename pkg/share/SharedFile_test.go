package share

import (
	"mise-share/pkg/share/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

// -- tests --
func TestToBody_U(t *testing.T) {
	share := NewSharedFile(
		&Source{
			Url: test.Str("https://httpbin.org/get"),
		},
	)

	html := share.ToBody("test")
	assert.Contains(t, html, `<meta name="mise-share-url" content="https://httpbin.org/get">`)
	assert.Contains(t, html, `<meta property="og:url" content="https://share.miseapp.co/test">`)
}
