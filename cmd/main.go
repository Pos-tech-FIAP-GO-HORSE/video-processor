package cmd

import (
	"github.com/aws/aws-lambda-go/lambda"
	"video-processor/internal/handlers"
)

func main() {
	lambda.Start(handlers.HandleRequest)
}
