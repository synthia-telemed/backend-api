package sms

import (
	"fmt"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Config struct {
	AccountSid string `env:"TWILIO_ACCOUNT_SID,required"`
	ApiKey     string `env:"TWILIO_API_KEY,required"`
	ApiSecret  string `env:"TWILIO_API_SECRET,required"`
	FromNumber string `env:"TWILIO_FROM_NUMBER,required"`
}

type TwilioClient struct {
	client     *twilio.RestClient
	fromNumber string
}

func NewTwilioClient(config *Config) *TwilioClient {
	return &TwilioClient{
		client: twilio.NewRestClientWithParams(twilio.ClientParams{
			Username:   config.ApiKey,
			Password:   config.ApiSecret,
			AccountSid: config.AccountSid,
		}),
		fromNumber: config.FromNumber,
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
