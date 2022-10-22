package request

import (
	"log"
	"mise-share/pkg/share"
	"net/http"

	"context"

	"github.com/aws/aws-lambda-go/events"
)

// -- impls --
func Handle(
	ctx context.Context,
	req events.LambdaFunctionURLRequest,
) (
	events.LambdaFunctionURLResponse,
	error,
) {
	// decode the body
	body, err := DecodeRequestBody(req.Body, req.IsBase64Encoded)
	if err != nil {
		log.Printf("[Handle] failed to decode request - len: %d\n>>>\n%s\n<<<\n", len(req.Body), req.Body)
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

	// construct the source from variant input
	var source share.SharedSource

	if body.Source.Url != nil {
		source = &share.SourceUrl{
			Url: body.Source.Url,
		}
	} else if body.Source.Json != nil {
		source = &share.SourceJson{
			Json: body.Source.Json,
		}
	}

	// validate the source
	if source == nil {
		return EncodeFailure(
			http.StatusBadRequest,
			"the request was missing the required field: 'source.url' or 'source.json'",
		)
	}

	// init share command
	cmd, err := share.New(source)

	// if error, return failure
	if err != nil {
		return EncodeFailure(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	// run command
	url, err := cmd.Call()

	// if error, return failure
	if err != nil {
		return EncodeFailure(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	// otherwise, return success w/ the share url
	return EncodeSuccess(
		http.StatusOK,
		url,
	)
}
