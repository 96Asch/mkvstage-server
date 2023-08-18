package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/service"
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
		Permission:   domain.MEMBER,
		ProfileColor: "FFFFFF",
	}

	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	mockUR.
		On("GetByID", context.TODO(), mockUser.ID).
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
		On("GetAll", context.TODO()).
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
		On("GetAll", context.TODO()).
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
		FirstName: "Foo",
		LastName:  "Bar",
		Email:     "Foo@Bar.com",

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
		On("Create", context.TODO(), mockUser).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg, ok := args.Get(1).(*domain.User)
			assert.True(t, ok)
			arg.ID = 1
		})
	mockURR.
		On("CreateBatch", context.TODO(), mockUserRoles).
		Return(nil)
	mockRR.
		On("GetAll", context.TODO()).
		Return(mockRoles, nil)

	US := service.NewUserService(mockUR, mockRR, mockURR)

	err := US.Store(context.TODO(), mockUser)
	assert.NoError(t, err)

	expectedUser := &domain.User{
		ID:           1,
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Permission:   domain.MEMBER,
		ProfileColor: "FFFFFF",
	}

	assert.Equal(t, expectedUser, mockUser)
	mockUR.AssertExpectations(t)
	mockRR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestStoreUserCreateErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		FirstName: "Foo",
		LastName:  "Bar",
		Email:     "Foo@Bar.com",

		Permission:   domain.MEMBER,
		ProfileColor: "FFFFFF",
	}

	mockErr := domain.NewInternalErr()
	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	mockUR.
		On("Create", context.TODO(), mockUser).
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
		FirstName: "Foo",
		LastName:  "Bar",
		Email:     "Foo@Bar.com",

		Permission:   domain.MEMBER,
		ProfileColor: "FFFFFF",
	}

	mockErr := domain.NewInternalErr()
	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	mockUR.
		On("Create", context.TODO(), mockUser).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg, ok := args.Get(1).(*domain.User)
			assert.True(t, ok)
			arg.ID = 1
		})
	mockRR.
		On("GetAll", context.TODO()).
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
		FirstName: "Foo",
		LastName:  "Bar",
		Email:     "Foo@Bar.com",

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
		On("Create", context.TODO(), mockUser).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg, ok := args.Get(1).(*domain.User)
			assert.True(t, ok)
			arg.ID = 1
		})
	mockURR.
		On("CreateBatch", context.TODO(), mockUserRoles).
		Return(mockErr)
	mockRR.
		On("GetAll", context.TODO()).
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
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
		Email:     "Foo@Bar.com",

		Permission:   domain.MEMBER,
		ProfileColor: "FFFFFF",
	}

	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	mockUR.
		On("Update", context.TODO(), mockUser).
		Return(nil)

	US := service.NewUserService(mockUR, mockRR, mockURR)

	err := US.Update(context.TODO(), mockUser)
	assert.NoError(t, err)

	expectedUser := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
		Email:     "Foo@Bar.com",

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
		ID:        0,
		FirstName: "Foo",
		LastName:  "Bar",
		Email:     "Foo@Bar.com",

		Permission:   domain.ADMIN,
		ProfileColor: "FFFFFF",
	}

	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	mockUR.
		On("Update", context.TODO(), mockUser).
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
		On("Delete", context.TODO(), mockUser.ID).
		Return(nil)
	mockUR.
		On("GetByID", context.TODO(), mockUser.ID).
		Return(nil, nil)
	mockURR.
		On("DeleteByUID", context.TODO(), mockUser.ID).
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
		On("Delete", context.TODO(), otherID).
		Return(nil)
	mockUR.
		On("GetByID", context.TODO(), otherID).
		Return(nil, nil)
	mockURR.
		On("DeleteByUID", context.TODO(), otherID).
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
		On("GetByID", context.TODO(), otherID).
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
		On("GetByID", context.TODO(), otherID).
		Return(nil, nil)
	mockUR.
		On("Delete", context.TODO(), otherID).
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
	mockUR.On("Delete", context.TODO(), mockUser.ID).Return(nil)
	mockUR.On("GetByID", context.TODO(), mockUser.ID).Return(nil, nil)

	mockURR := &mocks.MockUserRoleRepository{}
	mockURR.On("DeleteByUID", context.TODO(), mockUser.ID).Return(mockErr)

	mockRR := &mocks.MockRoleRepository{}
	US := service.NewUserService(mockUR, mockRR, mockURR)
	_, err := US.Remove(ctx, mockUser, 0)

	assert.ErrorAs(t, err, &mockErr)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}
