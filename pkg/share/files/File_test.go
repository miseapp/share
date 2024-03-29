package files

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

// -- mocks --
type TestContent string

func (c TestContent) Render(id string) string {
	return string(c)
}

// -- tests --
func TestNewFile_U(t *testing.T) {
	file, err := NewFileInput(1, TestContent("hello"))

	assert.Nil(t, err)
	assert.Equal(t, "1", file.Key)
	assert.Equal(t, 5, file.Length)
	assert.Equal(t, []byte{0x5d, 0x41, 0x40, 0x2a, 0xbc, 0x4b, 0x2a, 0x76, 0xb9, 0x71, 0x9d, 0x91, 0x10, 0x17, 0xc5, 0x92}, file.Hash[:])

	body, err := ioutil.ReadAll(file.Body)
	assert.Nil(t, err)
	assert.Equal(t, "hello", string(body))
}
