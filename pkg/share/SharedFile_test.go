package share

import (
	"mise-share/pkg/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

// -- tests --
func TestRenderUrl_U(t *testing.T) {
	share := NewSharedFile(
		&SourceUrl{
			Url: test.Str("https://test.com"),
		},
	)

	html := share.Render("test")
	assert.Contains(t, html, `<meta property="og:url" content="https://share.miseapp.co/test">`)
	assert.Contains(t, html, `<script id="mise-share" type="application/url">https://test.com</script>`)
}

func TestRenderJson_U(t *testing.T) {
	share := NewSharedFile(
		&SourceJson{
			Json: test.Str(`{"test":"json"}`),
		},
	)

	html := share.Render("test")
	assert.Contains(t, html, `<meta property="og:url" content="https://share.miseapp.co/test">`)
	assert.Contains(t, html, `<script id="mise-share" type="application/json">{"test":"json"}</script>`)
}
