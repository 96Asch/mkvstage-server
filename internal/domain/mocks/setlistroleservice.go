package mocks

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockSetlistRoleService struct {
	mock.Mock
}

func (msrs MockSetlistRoleService) Fetch(ctx context.Context, setlists *[]domain.Setlist) (*[]domain.SetlistRole, error) {
	ret := msrs.Called(ctx, setlists)

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

func (msrs MockSetlistRoleService) Store(ctx context.Context, setlistRoles *[]domain.SetlistRole) error {
	ret := msrs.Called(ctx, setlistRoles)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (msrs MockSetlistRoleService) Update(ctx context.Context, setlistRoles *[]domain.SetlistRole) error {
	ret := msrs.Called(ctx, setlistRoles)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (msrs MockSetlistRoleService) Remove(ctx context.Context, setlistRoleIDs []int64) error {
	ret := msrs.Called(ctx, setlistRoleIDs)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
