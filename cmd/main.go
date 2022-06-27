package main

import (
	"mise-share/pkg/share"
	"net/http"

	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

// -- impls --
func handleRequest(ctx context.Context, req Request) (string, error) {
	// validate the request
	if req.Source == nil {
		return EncodeFailure(
			http.StatusBadRequest,
			"the request was missing the required field: 'source'",
		)
	}

	if req.Source.Url == nil {
		return EncodeFailure(
			http.StatusBadRequest,
			"the request was missing the required field: 'source.url'",
		)
	}

	// init share command
	cmd, err := share.New(
		&share.Source{
			Url: req.Source.Url,
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
