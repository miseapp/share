package main

import (
	"mise-share/pkg/request"

	"github.com/aws/aws-lambda-go/lambda"
)

// -- impls --
func main() {
	lambda.Start(request.Handle)
}
