package service

import (
	"context"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserRoleSetActiveBatchCorrect(t *testing.T) {
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
		On("GetByUID", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(currentUserRoles, nil)
	mockURR.
		On("UpdateBatch", mock.AnythingOfType("*context.emptyCtx"), mockUserRoles).
		Return(nil)

	URS := NewUserRoleService(mockURR)

	userroles, err := URS.SetActiveBatch(context.TODO(), mockUserRoleIDs, mockUser)
	assert.NoError(t, err)
	assert.ElementsMatch(t, *mockUserRoles, *userroles)
	mockURR.AssertExpectations(t)
}

func TestUserRoleSetActiveBatchGetByUIDErr(t *testing.T) {
	mockUser := &domain.User{
		ID: 1,
	}

	mockUserRoleIDs := []int64{1, 2}

	mockErr := domain.NewInternalErr()
	mockURR := &mocks.MockUserRoleRepository{}
	mockURR.
		On("GetByUID", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(nil, mockErr)

	mockRR := &mocks.MockRoleRepository{}

	URS := NewUserRoleService(mockURR)

	userRoles, err := URS.SetActiveBatch(context.TODO(), mockUserRoleIDs, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	assert.Nil(t, userRoles)
	mockRR.AssertExpectations(t)
	mockURR.AssertExpectations(t)

}

func TestUserRoleSetActiveBatchInvalidUserRole(t *testing.T) {
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
		On("GetByUID", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(currentUserRoles, nil)

	URS := NewUserRoleService(mockURR)

	userroles, err := URS.SetActiveBatch(context.TODO(), mockUserRoleIDs, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	assert.Nil(t, userroles)
	mockURR.AssertExpectations(t)
}

func TestUserRoleSetActiveBatchErr(t *testing.T) {
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
		On("GetByUID", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(currentUserRoles, nil)
	mockURR.
		On("UpdateBatch", mock.AnythingOfType("*context.emptyCtx"), mockUserRoles).
		Return(mockErr)

	URS := NewUserRoleService(mockURR)

	userroles, err := URS.SetActiveBatch(context.TODO(), mockUserRoleIDs, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	assert.Nil(t, userroles)
	mockURR.AssertExpectations(t)
}

func TestUserRoleSetActiveBatchNoChange(t *testing.T) {
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
		On("GetByUID", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(currentUserRoles, nil)

	URS := NewUserRoleService(mockURR)

	userroles, err := URS.SetActiveBatch(context.TODO(), mockUserRoleIDs, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	assert.Nil(t, userroles)
	mockURR.AssertExpectations(t)
}
