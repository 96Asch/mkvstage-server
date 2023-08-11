package mocks

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockRoleService struct {
	mock.Mock
}

func (m MockRoleService) FetchByID(ctx context.Context, id int64) (*domain.Role, error) {
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

func (m MockRoleService) FetchAll(ctx context.Context) (*[]domain.Role, error) {
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

func (m MockRoleService) Update(ctx context.Context, domain *domain.Role, principal *domain.User) error {
	ret := m.Called(ctx, domain, principal)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockRoleService) UpdateBatch(ctx context.Context, domain *domain.Role, principal *domain.User) error {
	ret := m.Called(ctx, domain, principal)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockRoleService) Store(ctx context.Context, domain *domain.Role, principal *domain.User) error {
	ret := m.Called(ctx, domain, principal)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockRoleService) Remove(ctx context.Context, id int64, principal *domain.User) error {
	ret := m.Called(ctx, id, principal)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
