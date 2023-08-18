package mocks

import (
	"context"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, obj *domain.User) error {
	ret := m.Called(ctx, obj)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m *MockUserRepository) CreateBatch(ctx context.Context, obj *[]domain.User) error {
	ret := m.Called(ctx, obj)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
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

func (m *MockUserRepository) GetAll(ctx context.Context) (*[]domain.User, error) {
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

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	ret := m.Called(ctx, email)

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

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	ret := m.Called(ctx, user)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64) error {
	ret := m.Called(ctx, id)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
