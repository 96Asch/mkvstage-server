package service

import (
	"context"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserRoleUpdateBatchCorrect(t *testing.T) {
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

	mockRoles := &[]domain.Role{
		{
			ID:          1,
			Name:        "Foo",
			Description: "Foo",
		},
		{
			ID:          2,
			Name:        "Bar",
			Description: "Bar",
		},
	}

	mockURR := &mocks.MockUserRoleRepository{}
	mockURR.
		On("GetByUID", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(currentUserRoles, nil)
	mockURR.
		On("UpdateBatch", mock.AnythingOfType("*context.emptyCtx"), mockUserRoles).
		Return(nil)

	mockRR := &mocks.MockRoleRepository{}
	mockRR.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(mockRoles, nil)

	URS := NewUserRoleService(mockURR, mockRR)

	err := URS.UpdateBatch(context.TODO(), mockUserRoles, mockUser)
	assert.NoError(t, err)
	mockRR.AssertExpectations(t)
	mockURR.AssertExpectations(t)

}

func TestUserRoleUpdateBatchGetByUIDErr(t *testing.T) {
	mockUser := &domain.User{
		ID: 1,
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

	mockErr := domain.NewInternalErr()
	mockURR := &mocks.MockUserRoleRepository{}
	mockURR.
		On("GetByUID", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(nil, mockErr)

	mockRR := &mocks.MockRoleRepository{}

	URS := NewUserRoleService(mockURR, mockRR)

	err := URS.UpdateBatch(context.TODO(), mockUserRoles, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockURR.AssertExpectations(t)

}

func TestUserRoleUpdateBatchGetAllErr(t *testing.T) {
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

	mockErr := domain.NewInternalErr()
	mockURR := &mocks.MockUserRoleRepository{}
	mockURR.
		On("GetByUID", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(currentUserRoles, nil)

	mockRR := &mocks.MockRoleRepository{}
	mockRR.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(nil, mockErr)

	URS := NewUserRoleService(mockURR, mockRR)

	err := URS.UpdateBatch(context.TODO(), mockUserRoles, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockURR.AssertExpectations(t)

}

func TestUserRoleUpdateBatchNotAuth(t *testing.T) {
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
			UserID: 2,
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

	mockRoles := &[]domain.Role{
		{
			ID:          1,
			Name:        "Foo",
			Description: "Foo",
		},
		{
			ID:          2,
			Name:        "Bar",
			Description: "Bar",
		},
	}

	mockErr := domain.NewNotAuthorizedErr("")
	mockURR := &mocks.MockUserRoleRepository{}
	mockURR.
		On("GetByUID", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(currentUserRoles, nil)

	mockRR := &mocks.MockRoleRepository{}
	mockRR.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(mockRoles, nil)

	URS := NewUserRoleService(mockURR, mockRR)

	err := URS.UpdateBatch(context.TODO(), mockUserRoles, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockURR.AssertExpectations(t)

}

func TestUserRoleUpdateBatchInvalidUserRole(t *testing.T) {
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
			ID:     3,
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

	mockRoles := &[]domain.Role{
		{
			ID:          1,
			Name:        "Foo",
			Description: "Foo",
		},
		{
			ID:          2,
			Name:        "Bar",
			Description: "Bar",
		},
	}

	mockErr := domain.NewBadRequestErr("")
	mockURR := &mocks.MockUserRoleRepository{}
	mockURR.
		On("GetByUID", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(currentUserRoles, nil)

	mockRR := &mocks.MockRoleRepository{}
	mockRR.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(mockRoles, nil)

	URS := NewUserRoleService(mockURR, mockRR)

	err := URS.UpdateBatch(context.TODO(), mockUserRoles, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockURR.AssertExpectations(t)

}

func TestUserRoleUpdateBatchInvalidRole(t *testing.T) {
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
			RoleID: 3,
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
			RoleID: 3,
			Active: true,
		},
	}

	mockRoles := &[]domain.Role{
		{
			ID:          1,
			Name:        "Foo",
			Description: "Foo",
		},
		{
			ID:          2,
			Name:        "Bar",
			Description: "Bar",
		},
	}

	mockErr := domain.NewBadRequestErr("")
	mockURR := &mocks.MockUserRoleRepository{}
	mockURR.
		On("GetByUID", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(currentUserRoles, nil)

	mockRR := &mocks.MockRoleRepository{}
	mockRR.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(mockRoles, nil)

	URS := NewUserRoleService(mockURR, mockRR)

	err := URS.UpdateBatch(context.TODO(), mockUserRoles, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockURR.AssertExpectations(t)

}
