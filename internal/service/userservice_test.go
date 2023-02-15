package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetchByIDUserCorrect(t *testing.T) {

	mockUser := &domain.User{
		ID:           1,
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Password:     "FooBar",
		Permission:   domain.MEMBER,
		ProfileColor: "FFFFFF",
	}

	ctx := context.TODO()

	mockUR := new(mocks.MockUserRepository)
	mockUR.On("GetByID", ctx, mockUser.ID).Return(mockUser, nil)

	US := NewUserService(mockUR)
	user, err := US.FetchByID(ctx, mockUser.ID)

	assert.NoError(t, err)
	assert.Equal(t, user, mockUser)
	mockUR.AssertExpectations(t)
}

func TestFetchAllUserCorrect(t *testing.T) {
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

	ctx := context.TODO()

	mockUR := new(mocks.MockUserRepository)
	mockUR.On("GetAll", ctx).Return(mockUsers, nil)

	US := NewUserService(mockUR)
	users, err := US.FetchAll(ctx)
	assert.NoError(t, err)
	assert.Equal(t, users, mockPublicUsers)
	mockUR.AssertExpectations(t)
}

func TestFetchAllUserInternalErr(t *testing.T) {

	ctx := context.TODO()

	expectedErr := domain.NewInternalErr()
	mockUR := new(mocks.MockUserRepository)
	mockUR.On("GetAll", ctx).Return(nil, expectedErr)

	US := NewUserService(mockUR)
	users, err := US.FetchAll(ctx)
	assert.ErrorIs(t, expectedErr, err)
	assert.Nil(t, users)
	mockUR.AssertExpectations(t)
}

func TestStoreUser(t *testing.T) {

	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Password:     "FooBar",
		Permission:   domain.MEMBER,
		ProfileColor: "FFFFFF",
	}

	ctx := context.TODO()

	mockUR := new(mocks.MockUserRepository)
	mockUR.On("Create", ctx, mockUser).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).(*domain.User)
			arg.ID = 1
		})

	US := NewUserService(mockUR)
	err := US.Store(ctx, mockUser)

	expectedUser := &domain.User{
		ID:           1,
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Password:     "FooBar",
		Permission:   domain.MEMBER,
		ProfileColor: "FFFFFF",
	}

	assert.NoError(t, err)
	assert.NotEqual(t, expectedUser.Password, mockUser.Password)

	mockUser.Password = ""
	expectedUser.Password = ""
	assert.Equal(t, expectedUser, mockUser)
	mockUR.AssertExpectations(t)
}

func TestUpdateUserCorrect(t *testing.T) {
	mockUser := &domain.User{
		ID:           1,
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Password:     "FooBar",
		Permission:   domain.MEMBER,
		ProfileColor: "FFFFFF",
	}

	ctx := context.TODO()

	mockUR := new(mocks.MockUserRepository)
	mockUR.On("Update", ctx, mockUser).Return(nil)

	US := NewUserService(mockUR)
	err := US.Update(ctx, mockUser)

	expectedUser := &domain.User{
		ID:           1,
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Password:     "FooBar",
		Permission:   domain.MEMBER,
		ProfileColor: "FFFFFF",
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, mockUser)
	mockUR.AssertExpectations(t)
}

func TestUpdateUserZeroID(t *testing.T) {
	mockUser := &domain.User{
		ID:           0,
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Password:     "FooBar",
		Permission:   domain.ADMIN,
		ProfileColor: "FFFFFF",
	}

	ctx := context.TODO()

	mockUR := new(mocks.MockUserRepository)
	mockUR.On("Update", ctx, mockUser).Return(nil)

	US := NewUserService(mockUR)
	err := US.Update(ctx, mockUser)

	expectedErr := domain.NewBadRequestErr("")
	assert.ErrorAs(t, err, &expectedErr)
	mockUR.AssertNotCalled(t, "Update")
}

func TestDeleteUserCorrectOnlyUser(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
	}

	ctx := context.TODO()

	mockUR := new(mocks.MockUserRepository)
	mockUR.On("Delete", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).Return(nil)
	mockUR.On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).Return(nil, nil)

	US := NewUserService(mockUR)
	err := US.Remove(ctx, mockUser, 0)

	assert.NoError(t, err)
	mockUR.AssertExpectations(t)
}

func TestDeleteUserCorrectOtherUser(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}
	var otherID int64 = 2

	ctx := context.TODO()

	mockUR := new(mocks.MockUserRepository)
	mockUR.On("Delete", mock.AnythingOfType("*context.emptyCtx"), otherID).Return(nil)
	mockUR.On("GetByID", mock.AnythingOfType("*context.emptyCtx"), otherID).Return(nil, nil)

	US := NewUserService(mockUR)
	err := US.Remove(ctx, mockUser, otherID)

	assert.NoError(t, err)
	mockUR.AssertExpectations(t)
}

func TestDeleteUserNotAuthorized(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
	}
	var otherID int64 = 2

	ctx := context.TODO()

	mockUR := new(mocks.MockUserRepository)

	US := NewUserService(mockUR)
	err := US.Remove(ctx, mockUser, otherID)

	expectedErr := domain.NewNotAuthorizedErr("")
	assert.ErrorAs(t, err, &expectedErr)
	mockUR.AssertNotCalled(t, "Delete")
	mockUR.AssertNotCalled(t, "GetByID")
}

func TestDeleteUserNoRecord(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}
	var otherID int64 = 2

	ctx := context.TODO()

	expectedErr := domain.NewRecordNotFoundErr("id", fmt.Sprint(otherID))
	mockUR := new(mocks.MockUserRepository)
	mockUR.On("GetByID", mock.AnythingOfType("*context.emptyCtx"), otherID).Return(nil, expectedErr)

	US := NewUserService(mockUR)
	err := US.Remove(ctx, mockUser, otherID)

	assert.ErrorAs(t, err, &expectedErr)
	mockUR.AssertCalled(t, "GetByID", ctx, otherID)
	mockUR.AssertNotCalled(t, "Delete")
}

func TestDeleteUserInternalErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}
	var otherID int64 = 2

	ctx := context.TODO()

	expectedErr := domain.NewInternalErr()
	mockUR := new(mocks.MockUserRepository)
	mockUR.On("GetByID", mock.AnythingOfType("*context.emptyCtx"), otherID).Return(nil, nil)
	mockUR.On("Delete", mock.AnythingOfType("*context.emptyCtx"), otherID).Return(expectedErr)

	US := NewUserService(mockUR)
	err := US.Remove(ctx, mockUser, otherID)

	assert.ErrorAs(t, err, &expectedErr)
	mockUR.AssertExpectations(t)
}

func TestAuthorizeCorrect(t *testing.T) {

	mockUser := &domain.User{
		Email:    "Foo@Bar.com",
		Password: "FooBar",
	}

	expectedUser := &domain.User{
		ID:           1,
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Permission:   domain.MEMBER,
		ProfileColor: "FFFFFF",
	}

	hashPass, err := util.Encrypt(mockUser.Password)
	assert.NoError(t, err)
	expectedUser.Password = hashPass

	mockUR := new(mocks.MockUserRepository)
	mockUR.
		On("GetByEmail", mock.AnythingOfType("*context.emptyCtx"), mockUser.Email).
		Return(expectedUser, nil)

	ctx := context.TODO()

	US := NewUserService(mockUR)
	user, err := US.Authorize(ctx, mockUser.Email, mockUser.Password)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockUR.AssertExpectations(t)
}

func TestAuthorizeNoUserFound(t *testing.T) {
	mockUser := &domain.User{
		Email:    "Foobar@foo.com",
		Password: "Foobar",
	}

	expectedErr := domain.NewRecordNotFoundErr("email", mockUser.Email)

	mockUR := new(mocks.MockUserRepository)
	mockUR.
		On("GetByEmail", mock.AnythingOfType("*context.emptyCtx"), mockUser.Email).
		Return(nil, expectedErr)

	ctx := context.TODO()

	US := NewUserService(mockUR)
	user, err := US.Authorize(ctx, mockUser.Email, mockUser.Password)

	assert.ErrorAs(t, err, &expectedErr)
	assert.Nil(t, user)
	mockUR.AssertExpectations(t)
}

func TestAuthorizeNotAuthorized(t *testing.T) {
	mockUser := &domain.User{
		Email:    "Foo@Bar.com",
		Password: "FooBar",
	}

	expectedUser := &domain.User{
		ID:           1,
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Permission:   domain.MEMBER,
		ProfileColor: "FFFFFF",
	}

	expectedErr := domain.NewNotAuthorizedErr("email and/or password does not exist")

	hashPass, err := util.Encrypt("FooBar2")
	assert.NoError(t, err)
	expectedUser.Password = hashPass

	mockUR := new(mocks.MockUserRepository)
	mockUR.
		On("GetByEmail", mock.AnythingOfType("*context.emptyCtx"), mockUser.Email).
		Return(expectedUser, nil)

	ctx := context.TODO()

	US := NewUserService(mockUR)
	user, err := US.Authorize(ctx, mockUser.Email, mockUser.Password)
	assert.ErrorAs(t, err, &expectedErr)
	assert.Nil(t, user)
	mockUR.AssertExpectations(t)
}
