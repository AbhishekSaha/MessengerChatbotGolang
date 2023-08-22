package main

import (
	"errors"
	"github.com/tidwall/gjson"
	"log"
)

type FacebookEvent struct {
	Field      string
	Message    string
	CustomerId string
}

// Marshals the HTTP RequestBody received from Messenger API
// into a FacebookEvent object containing
// the required fields for the Chatbot to operate
func createFacebookEvent(requestBody string) (FacebookEvent, error) {
	event := new(FacebookEvent)

	// The hard-coded JSON values were taken from the Messenger API calls,
	// retrieved from CloudWatch
	if gjson.Get(requestBody, "entry.0.changes").Exists() { //Checks if the HTTP Request is a Page Comment
		event.Field = gjson.Get(requestBody, "entry.0.changes.0.field").String()
		event.Message = gjson.Get(requestBody, "entry.0.changes.0.value.message").String()
		event.CustomerId = gjson.Get(requestBody, "entry.0.changes.0.value.from.id").String()
	} else { //Checks if the HTTP Request is a Page Message
		event.Field = "messages"
		event.Message = gjson.Get(requestBody, "entry.0.messaging.0.message.text").String()
		event.CustomerId = gjson.Get(requestBody, "entry.0.messaging.0.sender.id").String()
	}

	if len(event.Field) == 0 || len(event.Message) == 0 {
		log.Println(event)
		return FacebookEvent{}, errors.New("Received invalid Webhook Response")
	}

	return *event, nil
}

// Routes the parsed FacebookEvent to the correct workflow
func routeEvent(event FacebookEvent) error {
	if event.Field == "feed" { //Event was a comment on the Page
		//Make some call to ChatGPT asking if the comment posted on the
		//Page suggests user wants to give a review
		// callChatGPTToCheckIfMessageSuggestReview(event.Message) returns boolean
		if true { //Swap this out based on ChatGPT response
			return sendMessage("Please feel free to provide a review about our Page!",
				event.CustomerId, "UPDATE")
		}
	} else { //Event was a message sent on the Messenger chat
		// review := event.Message
		// DB Config Object- sqlCfg = sql.config(...)
		// db, err = sql.Open("mysql", cfg.FormatDSN())
		return nil
	}

	return nil
}
