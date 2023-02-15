package mocks

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockTokenService struct {
	mock.Mock
}

func (m MockTokenService) ExtractUser(ctx context.Context, token *domain.AccessToken) (*domain.User, error) {
	ret := m.Called(ctx, token)

	var r0 *domain.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*domain.User)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m MockTokenService) CreateAccess(ctx context.Context, currentRefresh string) (*domain.AccessToken, error) {
	ret := m.Called(ctx, currentRefresh)

	var r0 *domain.AccessToken
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*domain.AccessToken)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m MockTokenService) CreateRefresh(ctx context.Context, uid int64, currentRefresh string) (*domain.RefreshToken, error) {
	ret := m.Called(ctx, uid, currentRefresh)

	var r0 *domain.RefreshToken
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*domain.RefreshToken)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m MockTokenService) Logout(ctx context.Context, uid int64) error {
	ret := m.Called(ctx, uid)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
