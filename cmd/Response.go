package main

import (
	"encoding/json"
	"log"
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
// encodes a success result to a string response
func EncodeSuccess(status int, url string) (string, error) {
	return encode(status, &Response{
		Success: &ResponseSuccess{Url: url},
		Failure: nil,
	})
}

// encodes a failure result to a string response
func EncodeFailure(status int, message string) (string, error) {
	return encode(status, &Response{
		Success: nil,
		Failure: &ResponseFailure{Message: message},
	})
}

// -- i/helpers

// encode the result as a string response
func encode(status int, result *Response) (string, error) {
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

// encode the string response
func encodeWithBody(status int, body string) (string, error) {
	// build response wrapper
	res := &events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       body,
	}

	// marshal it; if we fail to marshal the wrapper, who knows
	raw, err := json.Marshal(res)
	if err != nil {
		log.Println("[Result.encodeWithBody] couldn't marshal `events.APIGatewayProxyRepsonse`")
	}

	return string(raw), err
}
