package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/service"
	"github.com/96Asch/mkvstage-server/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetchByIDUserCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:           1,
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Password:     "FooBar",
		Permission:   domain.MEMBER,
		ProfileColor: "FFFFFF",
	}

	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	mockUR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(mockUser, nil)

	US := service.NewUserService(mockUR, mockRR, mockURR)

	user, err := US.FetchByID(context.TODO(), mockUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, user, mockUser)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestFetchAllUserCorrect(t *testing.T) {
	t.Parallel()

	mockUsers := &[]domain.User{
		{
			ID:           1,
			FirstName:    "Foo",
			LastName:     "Foo",
			Email:        "Foo@Foo.com",
			Permission:   domain.MEMBER,
			ProfileColor: "FFFFFF",
		},
		{
			ID:           2,
			FirstName:    "Bar",
			LastName:     "Bar",
			Email:        "Bar@Bar.com",
			Permission:   domain.MEMBER,
			ProfileColor: "FFFFF0",
		},
	}

	mockPublicUsers := &[]domain.User{
		{
			ID:           1,
			FirstName:    "Foo",
			LastName:     "Foo",
			Email:        "Foo@Foo.com",
			Permission:   domain.MEMBER,
			ProfileColor: "FFFFFF",
		},
		{
			ID:           2,
			FirstName:    "Bar",
			LastName:     "Bar",
			Email:        "Bar@Bar.com",
			Permission:   domain.MEMBER,
			ProfileColor: "FFFFF0",
		},
	}

	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	mockUR.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(mockUsers, nil)

	US := service.NewUserService(mockUR, mockRR, mockURR)

	users, err := US.FetchAll(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, users, mockPublicUsers)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestFetchAllUserInternalErr(t *testing.T) {
	t.Parallel()

	expectedErr := domain.NewInternalErr()
	mockUR := new(mocks.MockUserRepository)
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	mockUR.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(nil, expectedErr)

	US := service.NewUserService(mockUR, mockRR, mockURR)

	users, err := US.FetchAll(context.TODO())
	assert.ErrorIs(t, expectedErr, err)
	assert.Nil(t, users)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestStoreUserCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Password:     "FooBar",
		Permission:   domain.MEMBER,
		ProfileColor: "FFFFFF",
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

	mockUserRoles := &[]domain.UserRole{
		{
			UserID: 1,
			RoleID: 1,
		},
		{
			UserID: 1,
			RoleID: 2,
		},
	}

	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	mockUR.
		On("Create", mock.AnythingOfType("*context.emptyCtx"), mockUser).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg, ok := args.Get(1).(*domain.User)
			assert.True(t, ok)
			arg.ID = 1
		})
	mockURR.
		On("CreateBatch", mock.AnythingOfType("*context.emptyCtx"), mockUserRoles).
		Return(nil)
	mockRR.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(mockRoles, nil)

	US := service.NewUserService(mockUR, mockRR, mockURR)

	err := US.Store(context.TODO(), mockUser)
	assert.NoError(t, err)

	expectedUser := &domain.User{
		ID:           1,
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Password:     "FooBar",
		Permission:   domain.MEMBER,
		ProfileColor: "FFFFFF",
	}
	assert.NotEqual(t, expectedUser.Password, mockUser.Password)

	mockUser.Password = ""
	expectedUser.Password = ""
	assert.Equal(t, expectedUser, mockUser)
	mockUR.AssertExpectations(t)
	mockRR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestStoreUserCreateErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Password:     "FooBar",
		Permission:   domain.MEMBER,
		ProfileColor: "FFFFFF",
	}

	mockErr := domain.NewInternalErr()
	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	mockUR.
		On("Create", mock.AnythingOfType("*context.emptyCtx"), mockUser).
		Return(mockErr)

	US := service.NewUserService(mockUR, mockRR, mockURR)

	err := US.Store(context.TODO(), mockUser)
	assert.ErrorAs(t, err, &mockErr)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
	mockRR.AssertExpectations(t)
}

func TestStoreUserGetAllRoleErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Password:     "FooBar",
		Permission:   domain.MEMBER,
		ProfileColor: "FFFFFF",
	}

	mockErr := domain.NewInternalErr()
	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	mockUR.
		On("Create", mock.AnythingOfType("*context.emptyCtx"), mockUser).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg, ok := args.Get(1).(*domain.User)
			assert.True(t, ok)
			arg.ID = 1
		})
	mockRR.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(nil, mockErr)

	US := service.NewUserService(mockUR, mockRR, mockURR)

	err := US.Store(context.TODO(), mockUser)
	assert.ErrorAs(t, err, &mockErr)
	mockUR.AssertExpectations(t)
	mockRR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestStoreUserCreateBatchURErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Password:     "FooBar",
		Permission:   domain.MEMBER,
		ProfileColor: "FFFFFF",
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

	mockUserRoles := &[]domain.UserRole{
		{
			UserID: 1,
			RoleID: 1,
		},
		{
			UserID: 1,
			RoleID: 2,
		},
	}

	mockErr := domain.NewInternalErr()
	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	mockUR.
		On("Create", mock.AnythingOfType("*context.emptyCtx"), mockUser).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg, ok := args.Get(1).(*domain.User)
			assert.True(t, ok)
			arg.ID = 1
		})
	mockURR.
		On("CreateBatch", mock.AnythingOfType("*context.emptyCtx"), mockUserRoles).
		Return(mockErr)
	mockRR.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(mockRoles, nil)

	US := service.NewUserService(mockUR, mockRR, mockURR)

	err := US.Store(context.TODO(), mockUser)
	assert.ErrorAs(t, err, &mockErr)
	mockUR.AssertExpectations(t)
	mockRR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
	mockRR.AssertExpectations(t)
}

func TestUpdateUserCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:           1,
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Password:     "FooBar",
		Permission:   domain.MEMBER,
		ProfileColor: "FFFFFF",
	}

	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	mockUR.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), mockUser).
		Return(nil)

	US := service.NewUserService(mockUR, mockRR, mockURR)

	err := US.Update(context.TODO(), mockUser)
	assert.NoError(t, err)

	expectedUser := &domain.User{
		ID:           1,
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Password:     "FooBar",
		Permission:   domain.MEMBER,
		ProfileColor: "FFFFFF",
	}

	assert.Equal(t, expectedUser, mockUser)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
	mockRR.AssertExpectations(t)
}

func TestUpdateUserZeroID(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:           0,
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Password:     "FooBar",
		Permission:   domain.ADMIN,
		ProfileColor: "FFFFFF",
	}

	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	mockUR.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), mockUser).
		Return(nil)

	US := service.NewUserService(mockUR, mockRR, mockURR)

	err := US.Update(context.TODO(), mockUser)
	expectedErr := domain.NewBadRequestErr("")
	assert.ErrorAs(t, err, &expectedErr)
	mockUR.AssertNotCalled(t, "Update")
	mockURR.AssertExpectations(t)
	mockRR.AssertExpectations(t)
}

func TestDeleteUserCorrectOnlyUser(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
	}

	mockUR := new(mocks.MockUserRepository)
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	mockUR.
		On("Delete", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(nil)
	mockUR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(nil, nil)
	mockURR.
		On("DeleteByUID", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(nil)

	US := service.NewUserService(mockUR, mockRR, mockURR)

	deletedID, err := US.Remove(context.TODO(), mockUser, 0)
	assert.NoError(t, err)
	assert.Equal(t, mockUser.ID, deletedID)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
	mockRR.AssertExpectations(t)
}

func TestDeleteUserCorrectOtherUser(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	otherID := int64(2)
	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	mockUR.
		On("Delete", mock.AnythingOfType("*context.emptyCtx"), otherID).
		Return(nil)
	mockUR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), otherID).
		Return(nil, nil)
	mockURR.
		On("DeleteByUID", mock.AnythingOfType("*context.emptyCtx"), otherID).
		Return(nil)

	US := service.NewUserService(mockUR, mockRR, mockURR)

	deletedID, err := US.Remove(context.TODO(), mockUser, otherID)
	assert.NoError(t, err)
	assert.Equal(t, otherID, deletedID)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
	mockRR.AssertExpectations(t)
}

func TestDeleteUserNotAuthorized(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
	}

	otherID := int64(2)
	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	US := service.NewUserService(mockUR, mockRR, mockURR)

	_, err := US.Remove(context.TODO(), mockUser, otherID)
	expectedErr := domain.NewNotAuthorizedErr("")
	assert.ErrorAs(t, err, &expectedErr)
	mockUR.AssertNotCalled(t, "Delete")
	mockUR.AssertNotCalled(t, "GetByID")
	mockRR.AssertExpectations(t)
}

func TestDeleteUserNoRecord(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	otherID := int64(2)
	expectedErr := domain.NewRecordNotFoundErr("id", fmt.Sprint(otherID))
	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	mockUR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), otherID).
		Return(nil, expectedErr)

	US := service.NewUserService(mockUR, mockRR, mockURR)

	_, err := US.Remove(context.TODO(), mockUser, otherID)
	assert.ErrorAs(t, err, &expectedErr)
	mockUR.AssertCalled(t, "GetByID", context.TODO(), otherID)
	mockUR.AssertNotCalled(t, "Delete")
	mockRR.AssertExpectations(t)
}

func TestDeleteUserInternalErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}
	otherID := int64(2)
	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}
	expectedErr := domain.NewInternalErr()

	mockUR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), otherID).
		Return(nil, nil)
	mockUR.
		On("Delete", mock.AnythingOfType("*context.emptyCtx"), otherID).
		Return(expectedErr)

	US := service.NewUserService(mockUR, mockRR, mockURR)

	_, err := US.Remove(context.TODO(), mockUser, otherID)
	assert.ErrorAs(t, err, &expectedErr)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
	mockRR.AssertExpectations(t)
}

func TestDeleteUserDeleteUserRoleErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
	}

	mockErr := domain.NewInternalErr()
	ctx := context.TODO()

	mockUR := new(mocks.MockUserRepository)
	mockUR.On("Delete", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).Return(nil)
	mockUR.On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).Return(nil, nil)

	mockURR := &mocks.MockUserRoleRepository{}
	mockURR.On("DeleteByUID", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).Return(mockErr)

	mockRR := &mocks.MockRoleRepository{}
	US := service.NewUserService(mockUR, mockRR, mockURR)
	_, err := US.Remove(ctx, mockUser, 0)

	assert.ErrorAs(t, err, &mockErr)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestAuthorizeCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		Email:    "Foo@Bar.com",
		Password: "FooBar",
	}

	hashPass, err := util.Encrypt(mockUser.Password)
	assert.NoError(t, err)

	expectedUser := &domain.User{
		ID:           1,
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Permission:   domain.MEMBER,
		ProfileColor: "FFFFFF",
		Password:     hashPass,
	}

	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	mockUR.
		On("GetByEmail", mock.AnythingOfType("*context.emptyCtx"), mockUser.Email).
		Return(expectedUser, nil)

	US := service.NewUserService(mockUR, mockRR, mockURR)

	user, err := US.Authorize(context.TODO(), mockUser.Email, mockUser.Password)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestAuthorizeNoUserFound(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		Email:    "Foobar@foo.com",
		Password: "Foobar",
	}

	expectedErr := domain.NewRecordNotFoundErr("email", mockUser.Email)
	mockUR := new(mocks.MockUserRepository)
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	mockUR.
		On("GetByEmail", mock.AnythingOfType("*context.emptyCtx"), mockUser.Email).
		Return(nil, expectedErr)

	US := service.NewUserService(mockUR, mockRR, mockURR)

	user, err := US.Authorize(context.TODO(), mockUser.Email, mockUser.Password)
	assert.ErrorAs(t, err, &expectedErr)
	assert.Nil(t, user)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestAuthorizeNotAuthorized(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		Email:    "Foo@Bar.com",
		Password: "FooBar",
	}

	hashPass, err := util.Encrypt("FooBar2")
	assert.NoError(t, err)

	expectedUser := &domain.User{
		ID:           1,
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Permission:   domain.MEMBER,
		ProfileColor: "FFFFFF",
		Password:     hashPass,
	}

	expectedErr := domain.NewNotAuthorizedErr("email and/or password does not exist")
	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	mockUR.
		On("GetByEmail", mock.AnythingOfType("*context.emptyCtx"), mockUser.Email).
		Return(expectedUser, nil)

	US := service.NewUserService(mockUR, mockRR, mockURR)

	user, err := US.Authorize(context.TODO(), mockUser.Email, mockUser.Password)
	assert.ErrorAs(t, err, &expectedErr)
	assert.Nil(t, user)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestSetPermissionCorrect(t *testing.T) {
	t.Parallel()

	permission := domain.EDITOR
	recipient := &domain.User{
		ID:         2,
		Permission: domain.MEMBER,
	}
	principal := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}
	updatedUser := &domain.User{
		ID:         2,
		Permission: permission,
	}

	mockUR := &mocks.MockUserRepository{}
	mockRR := &mocks.MockRoleRepository{}
	mockURR := &mocks.MockUserRoleRepository{}

	mockUR.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), &domain.User{
			ID:         recipient.ID,
			Permission: permission,
		}).
		Return(nil)
	mockUR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), recipient.ID).
		Return(updatedUser, nil)

	userService := service.NewUserService(mockUR, mockRR, mockURR)

	user, err := userService.SetPermission(context.TODO(), permission, recipient, principal)
	assert.NoError(t, err)
	assert.Equal(t, updatedUser, user)
	mockUR.AssertExpectations(t)
	mockRR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestSetPermissionNotAdmin(t *testing.T) {
	t.Parallel()

	permission := domain.EDITOR
	recipient := &domain.User{
		ID:         2,
		Permission: domain.MEMBER,
	}
	principal := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	expErr := domain.NewNotAuthorizedErr("")
	mockUR := &mocks.MockUserRepository{}
	mockRR := &mocks.MockRoleRepository{}
	mockURR := &mocks.MockUserRoleRepository{}

	userService := service.NewUserService(mockUR, mockRR, mockURR)

	user, err := userService.SetPermission(context.TODO(), permission, recipient, principal)
	assert.ErrorAs(t, err, &expErr)
	assert.Nil(t, user)
	mockUR.AssertExpectations(t)
	mockRR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestSetPermissionRecIsAdmin(t *testing.T) {
	t.Parallel()

	permission := domain.EDITOR
	recipient := &domain.User{
		ID:         2,
		Permission: domain.ADMIN,
	}
	principal := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	expErr := domain.NewNotAuthorizedErr("")
	mockUR := &mocks.MockUserRepository{}
	mockRR := &mocks.MockRoleRepository{}
	mockURR := &mocks.MockUserRoleRepository{}

	userService := service.NewUserService(mockUR, mockRR, mockURR)

	user, err := userService.SetPermission(context.TODO(), permission, recipient, principal)
	assert.ErrorAs(t, err, &expErr)
	assert.Nil(t, user)
	mockUR.AssertExpectations(t)
	mockRR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestSetPermissionUpdateErr(t *testing.T) {
	t.Parallel()

	permission := domain.EDITOR
	recipient := &domain.User{
		ID:         2,
		Permission: domain.MEMBER,
	}
	principal := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	expErr := domain.NewInternalErr()
	mockUR := &mocks.MockUserRepository{}
	mockRR := &mocks.MockRoleRepository{}
	mockURR := &mocks.MockUserRoleRepository{}

	mockUR.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), &domain.User{
			ID:         recipient.ID,
			Permission: permission,
		}).
		Return(expErr)

	userService := service.NewUserService(mockUR, mockRR, mockURR)

	user, err := userService.SetPermission(context.TODO(), permission, recipient, principal)
	assert.ErrorAs(t, err, &expErr)
	assert.Nil(t, user)
	mockUR.AssertExpectations(t)
	mockRR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestSetPermissionCorrectGetErr(t *testing.T) {
	t.Parallel()

	permission := domain.EDITOR
	recipient := &domain.User{
		ID:         2,
		Permission: domain.MEMBER,
	}
	principal := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	expErr := domain.NewInternalErr()
	mockUR := &mocks.MockUserRepository{}
	mockRR := &mocks.MockRoleRepository{}
	mockURR := &mocks.MockUserRoleRepository{}

	mockUR.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), &domain.User{
			ID:         recipient.ID,
			Permission: permission,
		}).
		Return(nil)
	mockUR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), recipient.ID).
		Return(nil, expErr)

	userService := service.NewUserService(mockUR, mockRR, mockURR)

	user, err := userService.SetPermission(context.TODO(), permission, recipient, principal)
	assert.ErrorAs(t, err, &expErr)
	assert.Nil(t, user)
	mockUR.AssertExpectations(t)
	mockRR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}
