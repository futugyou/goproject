package infrastructure

import (
	"context"

	"github.com/futugyousuzu/identity-server/pkg/dto"
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

func (e *EmailService) SendEmail(ctx context.Context, data dto.EmailDTO) error {
	from := mail.NewEmail(data.From, data.From)
	newEmail := mail.NewEmail(data.To, data.To)
	message := mail.NewSingleEmail(from, data.Subject, newEmail, data.Text, data.Html)
	client := sendgrid.NewSendClient(e.opts.SendgridApiKey)
	_, err := client.SendWithContext(ctx, message)
	return err
}
