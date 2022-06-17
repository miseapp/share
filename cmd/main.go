package main

import (
	"log"
	"mise-share/pkg/share"

	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

// -- types --

// Event (value) is the api event structure
type Event struct {
	Source *share.Source `json:"source"`
}

// EventSource (value) is the the source of a shared file
type EventSource struct {
	Url *string `json:"url"`
}

// -- impls --
func handleRequest(ctx context.Context, event Event) (string, error) {
	// init share command
	share, err := share.New(
		&share.Source{
			Url: event.Source.Url,
		},
	)

	if err != nil {
		log.Println("[main.handleRequest] could not create share service", err)
	}

	// run command
	return share.Call()
}

func main() {
	lambda.Start(handleRequest)
}
