package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
)

// Handles the validation, parsing and routing of HTTP Request received by AWS Lambda
func baseHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Required to verify Messenger's WebHook endpoint
	// Messenger API will only send a GET HTTP request for verification,
	// in all other cases, there will be a POST
	if event.HTTPMethod == "GET" {
		return verifyHandler(ctx, event)
	}
	var response events.APIGatewayProxyResponse

	var facebookEvent, _err = createFacebookEvent(event.Body)
	if _err != nil {
		log.Fatalln(_err.Error() + "\n Http Body was: " + event.Body)
	} else {
		log.Println("Marshalled Facebook Event: " + fmt.Sprintf("%#v", facebookEvent))
	}

	_err = routeEvent(facebookEvent)
	if _err != nil {
		log.Fatalln(_err)
	} else {
		response = events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       event.Body,
		}
	}

	return response, nil
}

// Used by MessengerAPI to validate webhook integration.
// Not used in Chatbot workflow directly.
func verifyHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var queryParameters = event.QueryStringParameters
	log.Println(queryParameters)
	if queryParameters["hub.verify_token"] == "ABHISAHA_VERIFIY_META_TOKEN" {
		response := events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       queryParameters["hub.challenge"],
		}
		return response, nil
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: 403,
		Body:       "Missing/invalid token",
	}

	return response, nil
}
