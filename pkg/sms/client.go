package sms

type Client interface {
	Send(to, body string) error
}
