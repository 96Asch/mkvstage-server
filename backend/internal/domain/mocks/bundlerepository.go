package mocks

import (
	"context"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockBundleRepository struct {
	mock.Mock
}

func (m MockBundleRepository) GetByID(ctx context.Context, bid int64) (*domain.Bundle, error) {
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

func (m MockBundleRepository) GetAll(ctx context.Context) (*[]domain.Bundle, error) {
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

func (m MockBundleRepository) GetLeaves(ctx context.Context) (*[]domain.Bundle, error) {
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

func (m MockBundleRepository) Create(ctx context.Context, bundle *domain.Bundle) error {
	ret := m.Called(ctx, bundle)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockBundleRepository) Delete(ctx context.Context, bid int64) error {
	ret := m.Called(ctx, bid)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockBundleRepository) Update(ctx context.Context, bundle *domain.Bundle) error {
	ret := m.Called(ctx, bundle)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
