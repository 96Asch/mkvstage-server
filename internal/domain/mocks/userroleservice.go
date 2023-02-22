package mocks

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserRoleService struct {
	mock.Mock
}

func (m MockUserRoleService) FetchAll(ctx context.Context) (*[]domain.UserRole, error) {
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

func (m MockUserRoleService) FetchByUser(ctx context.Context, user *domain.User) (*[]domain.UserRole, error) {
	ret := m.Called(ctx, user)

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

func (m MockUserRoleService) SetActiveBatch(ctx context.Context, urids []int64, principal *domain.User) (*[]domain.UserRole, error) {
	ret := m.Called(ctx, urids, principal)

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
