package main

import (
	"mise-share/pkg/share"
	"net/http"

	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// -- impls --
func handleRequest(
	ctx context.Context,
	req events.LambdaFunctionURLRequest,
) (
	events.LambdaFunctionURLResponse,
	error,
) {
	// decode the body
	body, err := DecodeRequestBody(req.Body)
	if err != nil {
		return EncodeFailure(
			http.StatusBadRequest,
			err.Error(),
		)
	}

	// validate the request
	if body.Source == nil {
		return EncodeFailure(
			http.StatusBadRequest,
			"the request was missing the required field: 'source'",
		)
	}

	if body.Source.Url == nil {
		return EncodeFailure(
			http.StatusBadRequest,
			"the request was missing the required field: 'source.url'",
		)
	}

	// init share command
	cmd, err := share.New(
		&share.Source{
			Url: body.Source.Url,
		},
	)

	if err != nil {
		return EncodeFailure(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	// run command
	url, err := cmd.Call()

	// return result structure
	if err == nil {
		return EncodeSuccess(
			http.StatusOK,
			url,
		)
	} else {
		return EncodeFailure(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
}

func main() {
	lambda.Start(handleRequest)
}
