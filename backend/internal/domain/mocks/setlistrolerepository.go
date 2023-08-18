package mocks

import (
	"context"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockSetlistRoleRepository struct {
	mock.Mock
}

func (msrs MockSetlistRoleRepository) Create(ctx context.Context, setlistRoles *[]domain.SetlistRole) error {
	ret := msrs.Called(ctx, setlistRoles)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (msrs MockSetlistRoleRepository) Get(ctx context.Context, setlistIDs []int64) (*[]domain.SetlistRole, error) {
	ret := msrs.Called(ctx, setlistIDs)

	var r0 *[]domain.SetlistRole
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*[]domain.SetlistRole)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (msrs MockSetlistRoleRepository) Update(ctx context.Context, setlistRoles *[]domain.SetlistRole) error {
	ret := msrs.Called(ctx, setlistRoles)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (msrs MockSetlistRoleRepository) Delete(ctx context.Context, setlistRoleIDs []int64) error {
	ret := msrs.Called(ctx, setlistRoleIDs)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
