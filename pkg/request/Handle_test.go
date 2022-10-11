package request

import (
	"context"
	"encoding/json"
	"mise-share/pkg/test"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

// -- tests --
func TestHandleWithUrl_I(t *testing.T) {
	req, err := initRequest(RequestBody{
		Source: &RequestSource{
			Url: test.Str("https://test.com"),
		},
	})
	assert.Nil(t, err)

	res, err := Handle(context.TODO(), *req)
	assert.Nil(t, err)
	assert.Contains(t, res.Body, `http://share-files.s3.localhost.localstack.cloud:4566`)
}

func TestHandleWithJson_I(t *testing.T) {
	req, err := initRequest(RequestBody{
		Source: &RequestSource{
			Json: test.Str("https://test.com"),
		},
	})
	assert.Nil(t, err)

	res, err := Handle(context.TODO(), *req)
	assert.Nil(t, err)
	assert.Contains(t, res.Body, `http://share-files.s3.localhost.localstack.cloud:4566`)
}

func TestHandle_BadRequest_U(t *testing.T) {
	req1, err := initRequest(RequestBody{})
	assert.Nil(t, err)

	res1, err := Handle(context.TODO(), *req1)
	assert.Nil(t, err)
	assert.Equal(t, res1.StatusCode, http.StatusBadRequest)
	assert.Contains(t, res1.Body, "required field: 'source'")

	req2, err := initRequest(RequestBody{Source: &RequestSource{Url: nil, Json: nil}})
	assert.Nil(t, err)

	res2, err := Handle(context.TODO(), *req2)
	assert.Nil(t, err)
	assert.Equal(t, res2.StatusCode, http.StatusBadRequest)
	assert.Contains(t, res2.Body, "required field: 'source.url' or 'source.json'")
}

// -- helpers --

// create a lambda request from the body
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
