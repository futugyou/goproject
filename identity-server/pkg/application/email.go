package application

import (
	"context"
)

type EmailService interface {
	SendVerifyEmail(ctx context.Context, to string) error
}
