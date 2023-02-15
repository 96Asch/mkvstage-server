package mocks

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockTokenRepository struct {
	mock.Mock
}

func (m MockTokenRepository) GetAll(ctx context.Context, uid int64) (*[]domain.RefreshToken, error) {
	ret := m.Called(ctx, uid)

	var r0 *[]domain.RefreshToken
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*[]domain.RefreshToken)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m MockTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	ret := m.Called(ctx, token)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockTokenRepository) Delete(ctx context.Context, token *domain.RefreshToken) error {
	ret := m.Called(ctx, token)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockTokenRepository) DeleteAll(ctx context.Context, uid int64) error {
	ret := m.Called(ctx, uid)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
