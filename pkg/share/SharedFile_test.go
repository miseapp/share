package share

import (
	"mise-share/pkg/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

// -- tests --
func TestRender_U(t *testing.T) {
	share := NewSharedFile(
		&SourceUrl{
			Url: test.Str("http://unused.com"),
		},
		test.Str("http://test.host"),
	)

	html := share.Render("test")
	assert.Contains(t, html, `<meta http-equiv="refresh" content="0; url=miseapp://share?http://test.host/test"/>`)
	assert.Contains(t, html, `<meta property="og:url" content="http://test.host/test"/>`)
}

func TestRenderUrl_U(t *testing.T) {
	share := NewSharedFile(
		&SourceUrl{
			Url: test.Str("https://test.com"),
		},
		test.Str("http://unused.host"),
	)

	html := share.Render("test")
	assert.Contains(t, html, `<script id="mise-share" type="application/url">https://test.com</script>`)
}

func TestRenderJson_U(t *testing.T) {
	share := NewSharedFile(
		&SourceJson{
			Json: test.Str(`{"test":"json & w/ illegal char"}`),
		},
		test.Str("http://unused.host"),
	)

	html := share.Render("test")
	assert.Contains(t, html, `<script id="mise-share" type="application/json">{"test":"json &amp; w/ illegal char"}</script>`)
}
