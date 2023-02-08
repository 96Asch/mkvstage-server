package mocks

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) FetchByID(ctx context.Context, id int64) (*domain.User, error) {
	ret := m.Called(ctx, id)

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

func (m *MockUserService) FetchAll(ctx context.Context) (*[]domain.User, error) {
	ret := m.Called(ctx)

	var r0 *[]domain.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*[]domain.User)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m *MockUserService) Store(ctx context.Context, user *domain.User) error {
	ret := m.Called(ctx, user)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m *MockUserService) Update(ctx context.Context, user *domain.User) error {
	ret := m.Called(ctx, user)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m *MockUserService) Remove(ctx context.Context, user *domain.User, id int64) error {
	ret := m.Called(ctx, user, id)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
