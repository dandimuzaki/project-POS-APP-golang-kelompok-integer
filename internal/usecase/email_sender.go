package usecase

import (
	"context"
	"project-POS-APP-golang-integer/internal/dto/request"
)

type EmailSender interface {
	Send(ctx context.Context, req request.EmailRequest) error
}
