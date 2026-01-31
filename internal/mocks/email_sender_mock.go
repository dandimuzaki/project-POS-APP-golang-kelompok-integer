package mocks

import (
	"context"

	"project-POS-APP-golang-integer/internal/dto/request"

	"github.com/stretchr/testify/mock"
)

type EmailSenderMock struct {
	mock.Mock
}

func (m *EmailSenderMock) Send(ctx context.Context, req request.EmailRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}
