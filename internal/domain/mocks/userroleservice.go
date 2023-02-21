package mocks

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserRoleService struct {
	mock.Mock
}

func (m MockUserRoleService) UpdateBatch(ctx context.Context, urs *[]domain.UserRole, principal *domain.User) error {
	ret := m.Called(ctx, urs, principal)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
