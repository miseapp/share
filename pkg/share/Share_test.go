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

	res, err := share.Call()
	assert.Equal(t, nil, err)
	assert.Equal(t, "ok", res)
}
