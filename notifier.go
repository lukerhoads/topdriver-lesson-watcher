package main

import (
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Notifier struct {
	phoneNumber string
	client      *twilio.RestClient
}

func NewNotifier(accountSid, authToken, phoneNumber string) *Notifier {
	client := twilio.NewRestClientWithParams(twilio.RestClientParams{
		Username: accountSid,
		Password: authToken,
	})

	return &Notifier{
		phoneNumber: phoneNumber,
		client:      client,
	}
}

func (n *Notifier) SendText(to, message string) error {
	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(n.phoneNumber)
	params.SetBody(message)

	_, err := n.client.ApiV2010.CreateMessage(params)
	if err != nil {
		return err
	}

	return nil
}
