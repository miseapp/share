package files

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// -- tests --
func TestEncoding_U(t *testing.T) {
	key := Key(61)
	str, _ := key.Encode()
	assert.Equal(t, "Z", str)

	key = Key(62)
	str, _ = key.Encode()
	assert.Equal(t, "01", str)
}

func TestEncodingZero_U(t *testing.T) {
	key := Key(0)
	str, _ := key.Encode()
	assert.Equal(t, "0", str)
}
