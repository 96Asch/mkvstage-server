package service

import (
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetchByIDUser(t *testing.T) {

	mockUser := &domain.User{
		ID:           1,
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Password:     "FooBar",
		Permission:   "member",
		ProfileColor: "FFFFFF",
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

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
			Permission:   "member",
			ProfileColor: "FFFFFF",
		},
		{
			ID:           2,
			FirstName:    "Bar",
			LastName:     "Bar",
			Email:        "Bar@Bar.com",
			Permission:   "member",
			ProfileColor: "FFFFF0",
		},
	}

	mockPublicUsers := &[]domain.User{
		{
			ID:           1,
			FirstName:    "Foo",
			LastName:     "Foo",
			Email:        "Foo@Foo.com",
			Permission:   "member",
			ProfileColor: "FFFFFF",
		},
		{
			ID:           2,
			FirstName:    "Bar",
			LastName:     "Bar",
			Email:        "Bar@Bar.com",
			Permission:   "member",
			ProfileColor: "FFFFF0",
		},
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	mockUR := new(mocks.MockUserRepository)
	mockUR.On("GetAll", ctx).Return(mockUsers, nil)

	US := NewUserService(mockUR)
	users, err := US.FetchAll(ctx)
	assert.NoError(t, err)
	assert.Equal(t, users, mockPublicUsers)
	mockUR.AssertExpectations(t)
}

func TestFetchAllUserInternalErr(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

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
		Permission:   "member",
		ProfileColor: "FFFFFF",
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

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
		Permission:   "member",
		ProfileColor: "FFFFFF",
	}

	assert.NoError(t, err)
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
		Permission:   "member",
		ProfileColor: "FFFFFF",
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

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
		Permission:   "member",
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
		Permission:   "member",
		ProfileColor: "FFFFFF",
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	mockUR := new(mocks.MockUserRepository)
	mockUR.On("Update", ctx, mockUser).Return(nil)

	US := NewUserService(mockUR)
	err := US.Update(ctx, mockUser)

	assert.Error(t, err)
	mockUR.AssertNotCalled(t, "Update")
}
