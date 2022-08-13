package sms

import (
	"fmt"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Client interface {
	Send(to, body string) error
}

type TwilioClient struct {
	client     *twilio.RestClient
	fromNumber string
}

func NewTwilioClient(accountSid, apiKey, apiSecret, fromNumber string) *TwilioClient {
	return &TwilioClient{
		client: twilio.NewRestClientWithParams(twilio.ClientParams{
			Username:   apiKey,
			Password:   apiSecret,
			AccountSid: accountSid,
		}),
		fromNumber: fromNumber,
	}
}

func (c TwilioClient) Send(to, body string) error {
	params := &openapi.CreateMessageParams{}
	params.SetTo(c.parseThaiPhoneNumber(to))
	params.SetBody(body)
	params.SetFrom(c.fromNumber)
	_, err := c.client.Api.CreateMessage(params)
	return err
}

func (c TwilioClient) parseThaiPhoneNumber(phoneNumber string) string {
	return fmt.Sprintf("+66%s", phoneNumber[1:])
}
