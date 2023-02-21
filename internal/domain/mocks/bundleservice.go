package mocks

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockBundleService struct {
	mock.Mock
}

func (m MockBundleService) FetchByID(ctx context.Context, bid int64) (*domain.Bundle, error) {
	ret := m.Called(ctx, bid)

	var r0 *domain.Bundle
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*domain.Bundle)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m MockBundleService) FetchAll(ctx context.Context) (*[]domain.Bundle, error) {
	ret := m.Called(ctx)

	var r0 *[]domain.Bundle
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*[]domain.Bundle)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m MockBundleService) Store(ctx context.Context, bundle *domain.Bundle, principal *domain.User) error {
	ret := m.Called(ctx, bundle, principal)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockBundleService) Remove(ctx context.Context, bid int64, principal *domain.User) error {
	ret := m.Called(ctx, bid, principal)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockBundleService) Update(ctx context.Context, bundle *domain.Bundle, principal *domain.User) error {
	ret := m.Called(ctx, bundle, principal)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
