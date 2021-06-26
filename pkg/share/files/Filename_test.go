package files

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// -- tests --
func TestEncoding_U(t *testing.T) {
	filename := Filename(61)
	str, _ := filename.String()
	assert.Equal(t, "Z.html", str)

	filename = Filename(62)
	str, _ = filename.String()
	assert.Equal(t, "01.html", str)
}

func TestEncodingZero_U(t *testing.T) {
	filename := Filename(0)
	str, _ := filename.String()
	assert.Equal(t, "0.html", str)
}
