package infrastructure

import (
	"context"
	"fmt"

	"github.com/futugyousuzu/identity-server/pkg/options"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService struct {
	opts *options.Options
}

func NewEmailService(opts *options.Options) *EmailService {
	return &EmailService{
		opts: opts,
	}
}

const (
	emailsubject = "Activate your account"
)

func (e *EmailService) SendVerifyEmail(ctx context.Context, to string, url string) error {
	from := mail.NewEmail(e.opts.EmailFromName, e.opts.EmailFromAddress)
	newEmail := mail.NewEmail(to, to)

	plainTextContent := fmt.Sprintf("Hello,\n\nPlease activate your account by clicking the following link: %s\n\nThank you!", url)
	htmlContent := fmt.Sprintf("<p>Hello,</p><p>Please activate your account by clicking the following link: <a href='%s'>Activate Account</a></p><p>Thank you!</p>", url)

	message := mail.NewSingleEmail(from, emailsubject, newEmail, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(e.opts.SendgridApiKey)
	_, err := client.Send(message)
	return err
}
