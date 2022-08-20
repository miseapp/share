package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// -- types --

// the share response envelope
type Response struct {
	// the success structure, if any
	Success *ResponseSuccess `json:"success,omitempty"`

	// the failure structure, if any
	Failure *ResponseFailure `json:"failure,omitempty"`
}

// the apio response success structure
type ResponseSuccess struct {
	// the location of the shared file
	Url string `json:"url"`
}

// the api response failure structure
type ResponseFailure struct {
	// the user-facing error message
	Message string `json:"message"`
}

// -- impls --
// encodes a success result to a lambda url response
func EncodeSuccess(
	status int,
	url string,
) (
	events.LambdaFunctionURLResponse,
	error,
) {
	return encode(status, &Response{
		Success: &ResponseSuccess{Url: url},
		Failure: nil,
	})
}

// encodes a failure result to a lambda url response
func EncodeFailure(
	status int,
	message string,
) (
	events.LambdaFunctionURLResponse,
	error,
) {
	return encode(status, &Response{
		Success: nil,
		Failure: &ResponseFailure{Message: message},
	})
}

// -- i/helpers

// encode the result as a lambda url response
func encode(
	status int,
	result *Response,
) (
	events.LambdaFunctionURLResponse,
	error,
) {
	// encode body from result
	raw, err := json.Marshal(result)

	// build aws result
	if err == nil {
		return encodeWithBody(
			status,
			string(raw),
		)
	} else {
		return encodeWithBody(
			http.StatusInternalServerError,
			`{ "failure": { "message": "failed to encode result" } }`,
		)
	}

}

// encode the lambda url response
func encodeWithBody(
	status int,
	body string,
) (
	events.LambdaFunctionURLResponse,
	error,
) {
	res := events.LambdaFunctionURLResponse{
		StatusCode: status,
		Body:       body,
	}

	return res, nil
}
