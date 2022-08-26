package request

import (
	"context"
	"encoding/json"
	"mise-share/pkg/share/test"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

// -- helpers --
func initRequest(body RequestBody) (*events.LambdaFunctionURLRequest, error) {
	enc, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	res := events.LambdaFunctionURLRequest{
		Body: string(enc),
	}

	return &res, nil
}

// -- tests --
func TestHandle_I(t *testing.T) {
	req, err := initRequest(RequestBody{
		Source: &RequestSource{
			Url: test.Str("https://test.com"),
		},
	})
	assert.Nil(t, err)

	res, err := Handle(context.TODO(), *req)
	assert.Nil(t, err)
	assert.Contains(t, res.Body, `<meta name="mise-share-url" content="https://test.com">`)
	assert.Contains(t, res.Body, `<meta property="og:url" content="https://share.miseapp.co/test">`)
}

func TestHandle_BadRequest_U(t *testing.T) {
	req1, err := initRequest(RequestBody{})
	assert.Nil(t, err)

	res1, err := Handle(context.TODO(), *req1)
	assert.Nil(t, err)
	assert.Equal(t, res1.StatusCode, http.StatusBadRequest)
	assert.Contains(t, res1.Body, "required field: 'source'")

	req2, err := initRequest(RequestBody{Source: &RequestSource{Url: nil}})
	assert.Nil(t, err)

	res2, err := Handle(context.TODO(), *req2)
	assert.Nil(t, err)
	assert.Equal(t, res2.StatusCode, http.StatusBadRequest)
	assert.Contains(t, res2.Body, "required field: 'source.url'")
}
