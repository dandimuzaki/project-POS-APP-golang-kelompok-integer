package email

import (
	"context"
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/usecase"
	"project-POS-APP-golang-integer/pkg/utils"

	"go.uber.org/zap"
)

type AsyncEmailSender struct {
	jobs   chan<- utils.EmailJob
	config utils.Configuration
	log    *zap.Logger
}

func NewAsyncEmailSender(
	jobs chan<- utils.EmailJob,
	config utils.Configuration,
	log *zap.Logger,
) usecase.EmailSender {
	return &AsyncEmailSender{
		jobs: jobs,
		config: config,
		log: log,
	}
}

func (s *AsyncEmailSender) Send(ctx context.Context, req request.EmailRequest) error {
	s.jobs <- utils.EmailJob{
		EmailContent: req,
		Config: s.config,
		Log: s.log,
	}
	return nil
}
