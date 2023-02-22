package mocks

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserRoleRepository struct {
	mock.Mock
}

func (m MockUserRoleRepository) GetByID(ctx context.Context, id int64) (*domain.UserRole, error) {
	ret := m.Called(ctx, id)

	var r0 *domain.UserRole
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*domain.UserRole)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m MockUserRoleRepository) GetAll(ctx context.Context) (*[]domain.UserRole, error) {
	ret := m.Called(ctx)

	var r0 *[]domain.UserRole
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*[]domain.UserRole)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m MockUserRoleRepository) GetByUID(ctx context.Context, uid int64) (*[]domain.UserRole, error) {
	ret := m.Called(ctx, uid)

	var r0 *[]domain.UserRole
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*[]domain.UserRole)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m MockUserRoleRepository) Create(ctx context.Context, ur *domain.UserRole) error {
	ret := m.Called(ctx, ur)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockUserRoleRepository) CreateBatch(ctx context.Context, urs *[]domain.UserRole) error {
	ret := m.Called(ctx, urs)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockUserRoleRepository) Update(ctx context.Context, ur *domain.UserRole) error {
	ret := m.Called(ctx, ur)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockUserRoleRepository) UpdateBatch(ctx context.Context, urs *[]domain.UserRole) error {
	ret := m.Called(ctx, urs)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockUserRoleRepository) Delete(ctx context.Context, id int64) error {
	ret := m.Called(ctx, id)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockUserRoleRepository) DeleteBatch(ctx context.Context, ids []int64) error {
	ret := m.Called(ctx, ids)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockUserRoleRepository) DeleteByRID(ctx context.Context, rid int64) error {
	ret := m.Called(ctx, rid)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
