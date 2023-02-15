package mocks

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockTokenRepository struct {
	mock.Mock
}

func (m MockTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) {
	m.Called(ctx, token)

	// var r0 *domain.User
	// if ret.Get(0) != nil {
	// 	r0 = ret.Get(0).(*domain.User)
	// }

	// var r1 error
	// if ret.Get(1) != nil {
	// 	r1 = ret.Get(1).(error)
	// }

	// return r0, r1
}
