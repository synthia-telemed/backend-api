package notification

import (
	"context"
	firebase "firebase.google.com/go/v4"
	messaging "firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type Config struct {
	FirebaseCredentialFilePath string `env:"FIREBASE_CRED_FILE_PATH,required"`
}

type FirebaseNotificationClient struct {
	client *messaging.Client
}

func NewFirebaseNotificationClient(ctx context.Context, cfg *Config) (*FirebaseNotificationClient, error) {
	opt := option.WithCredentialsFile(cfg.FirebaseCredentialFilePath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}
	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}
	return &FirebaseNotificationClient{client: client}, nil
}

type SendParams struct {
	Token string
	Title string
	Body  string
}

func (c FirebaseNotificationClient) Send(ctx context.Context, params SendParams, data map[string]string) error {
	_, err := c.client.Send(ctx, &messaging.Message{
		Token: params.Token,
		Notification: &messaging.Notification{
			Title: params.Title,
			Body:  params.Body,
		},
		Data: data,
	})
	return err
}
