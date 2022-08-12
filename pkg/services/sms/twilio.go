package sms

import (
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Client interface {
	Send(to, body string) error
}

type TwilioClient struct {
	client *twilio.RestClient
}

func NewTwilioClient(accountSid, apiKey, apiSecret string) *TwilioClient {
	return &TwilioClient{
		client: twilio.NewRestClientWithParams(twilio.ClientParams{
			Username:   apiKey,
			Password:   apiSecret,
			AccountSid: accountSid,
		}),
	}
}

func (c TwilioClient) Send(to, body string) error {
	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetBody(body)
	_, err := c.client.Api.CreateMessage(params)
	return err
}
