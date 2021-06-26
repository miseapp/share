package share

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// -- tests --
func TestHtml_U(t *testing.T) {
	share := NewSharedFile(
		&Source{
			Url: strp("https://httpbin.org/get"),
		},
	)

	html := share.AsHtml()
	assert.Contains(t, html, `<meta name="mise-share-url" content="https://httpbin.org/get">`)
}
