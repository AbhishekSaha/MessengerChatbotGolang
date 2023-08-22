package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

// Main function entered by AWS Lambda
func main() {
	lambda.Start(baseHandler)
}
