package application

import (
	"context"

	"github.com/futugyousuzu/identity-server/pkg/dto"
)

type EmailService interface {
	SendEmail(ctx context.Context, data dto.EmailDTO) error
}
