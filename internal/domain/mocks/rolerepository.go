package mocks

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockRoleRepository struct {
	mock.Mock
}

func (m MockRoleRepository) Create(ctx context.Context, role *domain.Role) error {
	ret := m.Called(ctx, role)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockRoleRepository) GetByID(ctx context.Context, id int64) (*domain.Role, error) {
	ret := m.Called(ctx, id)

	var r0 *domain.Role
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*domain.Role)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m MockRoleRepository) GetAll(ctx context.Context) (*[]domain.Role, error) {
	ret := m.Called(ctx)

	var r0 *[]domain.Role
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*[]domain.Role)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m MockRoleRepository) Update(ctx context.Context, role *domain.Role) error {
	ret := m.Called(ctx, role)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockRoleRepository) Delete(ctx context.Context, rid int64) error {
	ret := m.Called(ctx, rid)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
