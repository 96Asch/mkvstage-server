package service_test

import (
	"context"
	"testing"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/96Asch/mkvstage-server/backend/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/backend/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestUserRoleSetActiveBatchCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID: 1,
	}

	currentUserRoles := &[]domain.UserRole{
		{
			ID:     1,
			UserID: 1,
			RoleID: 1,
			Active: false,
		},
		{
			ID:     2,
			UserID: 1,
			RoleID: 2,
			Active: true,
		},
		{
			ID:     3,
			UserID: 1,
			RoleID: 2,
			Active: true,
		},
	}

	mockUserRoles := &[]domain.UserRole{
		{
			ID:     1,
			UserID: 1,
			RoleID: 1,
			Active: true,
		},
		{
			ID:     2,
			UserID: 1,
			RoleID: 2,
			Active: false,
		},
	}

	mockUserRoleIDs := []int64{1, 3}
	mockURR := &mocks.MockUserRoleRepository{}

	mockURR.
		On("GetByUID", context.TODO(), mockUser.ID).
		Return(currentUserRoles, nil)
	mockURR.
		On("UpdateBatch", context.TODO(), mockUserRoles).
		Return(nil)

	URS := service.NewUserRoleService(mockURR)

	userroles, err := URS.SetActiveBatch(context.TODO(), mockUserRoleIDs, mockUser)
	assert.NoError(t, err)
	assert.ElementsMatch(t, *mockUserRoles, *userroles)
	mockURR.AssertExpectations(t)
}

func TestUserRoleSetActiveBatchGetByUIDErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID: 1,
	}

	mockUserRoleIDs := []int64{1, 2}

	mockErr := domain.NewInternalErr()
	mockURR := &mocks.MockUserRoleRepository{}

	mockURR.
		On("GetByUID", context.TODO(), mockUser.ID).
		Return(nil, mockErr)

	URS := service.NewUserRoleService(mockURR)

	userRoles, err := URS.SetActiveBatch(context.TODO(), mockUserRoleIDs, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	assert.Nil(t, userRoles)
	mockURR.AssertExpectations(t)
}

func TestUserRoleSetActiveBatchInvalidUserRole(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID: 1,
	}

	currentUserRoles := &[]domain.UserRole{
		{
			ID:     1,
			UserID: 1,
			RoleID: 1,
			Active: false,
		},
		{
			ID:     2,
			UserID: 1,
			RoleID: 2,
			Active: false,
		},
	}

	mockUserRoleIDs := []int64{1, 3}
	mockErr := domain.NewBadRequestErr("")
	mockURR := &mocks.MockUserRoleRepository{}

	mockURR.
		On("GetByUID", context.TODO(), mockUser.ID).
		Return(currentUserRoles, nil)

	URS := service.NewUserRoleService(mockURR)

	userroles, err := URS.SetActiveBatch(context.TODO(), mockUserRoleIDs, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	assert.Nil(t, userroles)
	mockURR.AssertExpectations(t)
}

func TestUserRoleSetActiveBatchErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID: 1,
	}

	currentUserRoles := &[]domain.UserRole{
		{
			ID:     1,
			UserID: 1,
			RoleID: 1,
			Active: false,
		},
		{
			ID:     2,
			UserID: 1,
			RoleID: 2,
			Active: false,
		},
	}

	mockUserRoles := &[]domain.UserRole{
		{
			ID:     1,
			UserID: 1,
			RoleID: 1,
			Active: true,
		},
		{
			ID:     2,
			UserID: 1,
			RoleID: 2,
			Active: true,
		},
	}
	mockUserRoleIDs := []int64{1, 2}
	mockErr := domain.NewInternalErr()
	mockURR := &mocks.MockUserRoleRepository{}

	mockURR.
		On("GetByUID", context.TODO(), mockUser.ID).
		Return(currentUserRoles, nil)
	mockURR.
		On("UpdateBatch", context.TODO(), mockUserRoles).
		Return(mockErr)

	URS := service.NewUserRoleService(mockURR)

	userroles, err := URS.SetActiveBatch(context.TODO(), mockUserRoleIDs, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	assert.Nil(t, userroles)
	mockURR.AssertExpectations(t)
}

func TestUserRoleSetActiveBatchNoChange(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID: 1,
	}

	currentUserRoles := &[]domain.UserRole{
		{
			ID:     1,
			UserID: 1,
			RoleID: 1,
			Active: false,
		},
		{
			ID:     2,
			UserID: 1,
			RoleID: 2,
			Active: true,
		},
	}

	mockUserRoleIDs := []int64{2}
	mockErr := domain.NewBadRequestErr("")
	mockURR := &mocks.MockUserRoleRepository{}

	mockURR.
		On("GetByUID", context.TODO(), mockUser.ID).
		Return(currentUserRoles, nil)

	URS := service.NewUserRoleService(mockURR)

	userroles, err := URS.SetActiveBatch(context.TODO(), mockUserRoleIDs, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	assert.Nil(t, userroles)
	mockURR.AssertExpectations(t)
}
