package main

import (
	"context"
	"fmt"
	"mise-share/pkg/share/test"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// -- tests --
func TestHandleRequest_I(t *testing.T) {
	req := Request{
		Source: &RequestSource{
			Url: test.Str("https://test.net"),
		},
	}

	res, err := handleRequest(context.TODO(), req)
	assert.Nil(t, err)
	assert.Contains(t, res, `<meta name="mise-share-url" content="https://httpbin.org/get">`)
	assert.Contains(t, res, `<meta property="og:url" content="https://share.miseapp.co/test">`)
}

func TestHandleRequest_BadRequest_U(t *testing.T) {
	r1, err := handleRequest(context.TODO(), Request{})
	assert.Nil(t, err)
	assert.Contains(t, r1, fmt.Sprint(http.StatusBadRequest))
	assert.Contains(t, r1, "required field: 'source'")

	r2, err := handleRequest(context.TODO(), Request{Source: &RequestSource{Url: nil}})
	assert.Nil(t, err)
	assert.Contains(t, r2, fmt.Sprint(http.StatusBadRequest))
	assert.Contains(t, r2, "required field: 'source.url'")
}
