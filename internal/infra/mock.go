package infra

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockTxManager struct {
	mock.Mock
}

func (m *MockTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx)
	if fn != nil {
		return fn(ctx)
	}
	return args.Error(0)
}
