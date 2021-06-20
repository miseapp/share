package share

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// -- tests --
func TestShare(t *testing.T) {
	share := New(
		&Source{
			Url: strp("https://httpbin.org/get"),
		},
	)

	// run command
	res, _ := share.Call()
	assert.Equal(t, "ok", res)
}
