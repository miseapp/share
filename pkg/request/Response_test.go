package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// -- tests --
func TestEncodeSuccess_U(t *testing.T) {
	res, err := EncodeSuccess(200, "https://test.com")

	assert.Nil(t, err)
	assert.Equal(t, res.Body, `{"success":{"url":"https://test.com"}}`)
}

func TestEncodeFailure_U(t *testing.T) {
	res, err := EncodeFailure(500, "a test error")

	assert.Nil(t, err)
	assert.Equal(t, res.Body, `{"failure":{"message":"a test error"}}`)
}
