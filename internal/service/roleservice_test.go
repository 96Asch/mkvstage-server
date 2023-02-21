package service

import (
	"context"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRoleFetchByIDCorrect(t *testing.T) {

	rid := int64(1)
	mockRole := &domain.Role{
		ID:          rid,
		Name:        "Foo",
		Description: "FooBar",
	}

	mockRR := &mocks.MockRoleRepository{}
	mockRR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), rid).
		Return(mockRole, nil)

	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}

	RS := NewRoleService(mockRR, mockUR, mockURR)

	ctx := context.TODO()

	role, err := RS.FetchByID(ctx, rid)

	assert.NoError(t, err)
	assert.Equal(t, mockRole, role)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRoleFetchByIDGetErr(t *testing.T) {

	rid := int64(1)
	mockErr := domain.NewRecordNotFoundErr("", "")
	mockRR := &mocks.MockRoleRepository{}
	mockRR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), rid).
		Return(nil, mockErr)

	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}

	RS := NewRoleService(mockRR, mockUR, mockURR)

	ctx := context.TODO()

	role, err := RS.FetchByID(ctx, rid)

	assert.ErrorAs(t, err, &mockErr)
	assert.Nil(t, role)

	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRoleFetchAllCorrect(t *testing.T) {

	mockRoles := &[]domain.Role{
		{
			ID:          1,
			Name:        "Foo",
			Description: "FooBar",
		},
		{
			ID:          2,
			Name:        "Bar",
			Description: "BarFoo",
		},
	}

	mockRR := &mocks.MockRoleRepository{}
	mockRR.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(mockRoles, nil)

	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}

	RS := NewRoleService(mockRR, mockUR, mockURR)

	ctx := context.TODO()

	role, err := RS.FetchAll(ctx)

	assert.NoError(t, err)
	assert.Equal(t, mockRoles, role)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRoleFetchAllGetErr(t *testing.T) {

	mockErr := domain.NewInternalErr()
	mockRR := &mocks.MockRoleRepository{}
	mockRR.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(nil, mockErr)

	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}

	RS := NewRoleService(mockRR, mockUR, mockURR)

	ctx := context.TODO()

	role, err := RS.FetchAll(ctx)

	assert.ErrorAs(t, err, &mockErr)
	assert.Nil(t, role)

	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRoleUpdateCorrect(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	rid := int64(1)
	mockRole := &domain.Role{
		ID:          rid,
		Name:        "Foo",
		Description: "FooBar",
	}

	prevMockRole := &domain.Role{
		ID:          rid,
		Name:        "Bar",
		Description: "BarFoo",
	}

	mockRR := &mocks.MockRoleRepository{}
	mockRR.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), mockRole).
		Return(nil)

	mockRR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), rid).
		Return(prevMockRole, nil)

	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}

	RS := NewRoleService(mockRR, mockUR, mockURR)

	err := RS.Update(context.TODO(), mockRole, mockUser)

	assert.NoError(t, err)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRoleUpdateIDZero(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	rid := int64(0)
	mockRole := &domain.Role{
		ID:          rid,
		Name:        "Foo",
		Description: "FooBar",
	}

	mockRR := &mocks.MockRoleRepository{}

	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}

	RS := NewRoleService(mockRR, mockUR, mockURR)

	err := RS.Update(context.TODO(), mockRole, mockUser)

	mockErr := domain.NewBadRequestErr("")
	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRoleUpdateNoPermission(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
	}

	rid := int64(1)
	mockRole := &domain.Role{
		ID:          rid,
		Name:        "Foo",
		Description: "FooBar",
	}

	mockRR := &mocks.MockRoleRepository{}

	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}

	RS := NewRoleService(mockRR, mockUR, mockURR)

	err := RS.Update(context.TODO(), mockRole, mockUser)

	mockErr := domain.NewNotAuthorizedErr("")
	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRoleUpdateNoRecord(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	rid := int64(1)
	mockRole := &domain.Role{
		ID:          rid,
		Name:        "Foo",
		Description: "FooBar",
	}

	mockErr := domain.NewRecordNotFoundErr("", "")
	mockRR := &mocks.MockRoleRepository{}

	mockRR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), rid).
		Return(nil, mockErr)

	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}

	RS := NewRoleService(mockRR, mockUR, mockURR)

	err := RS.Update(context.TODO(), mockRole, mockUser)

	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRoleUpdateErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	rid := int64(1)
	mockRole := &domain.Role{
		ID:          rid,
		Name:        "Foo",
		Description: "FooBar",
	}

	prevMockRole := &domain.Role{
		ID:          rid,
		Name:        "Bar",
		Description: "BarFoo",
	}

	mockErr := domain.NewInternalErr()
	mockRR := &mocks.MockRoleRepository{}
	mockRR.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), mockRole).
		Return(mockErr)

	mockRR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), rid).
		Return(prevMockRole, nil)

	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}

	RS := NewRoleService(mockRR, mockUR, mockURR)

	err := RS.Update(context.TODO(), mockRole, mockUser)

	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRoleStoreCorrect(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockUsers := &[]domain.User{
		{
			ID:         1,
			Permission: domain.ADMIN,
		},
		{
			ID:         2,
			Permission: domain.GUEST,
		},
	}

	mockRole := &domain.Role{
		Name:        "Foo",
		Description: "Foobar",
	}

	mockUserRoles := &[]domain.UserRole{
		{
			UserID: 1,
			RoleID: 1,
		},
		{
			UserID: 2,
			RoleID: 1,
		},
	}

	mockUR := &mocks.MockUserRepository{}
	mockUR.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(mockUsers, nil)
	mockURR := &mocks.MockUserRoleRepository{}
	mockURR.
		On("CreateBatch", mock.AnythingOfType("*context.emptyCtx"), mockUserRoles).
		Return(nil)

	mockRR := &mocks.MockRoleRepository{}
	mockRR.
		On("Create", mock.AnythingOfType("*context.emptyCtx"), mockRole).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).(*domain.Role)
			arg.ID = 1
		})

	RS := NewRoleService(mockRR, mockUR, mockURR)

	err := RS.Store(context.TODO(), mockRole, mockUser)

	assert.NoError(t, err)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRoleStoreNoPermission(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
	}

	mockRole := &domain.Role{
		Name:        "Foo",
		Description: "Foobar",
	}

	mockErr := domain.NewNotAuthorizedErr("")
	mockUR := &mocks.MockUserRepository{}

	mockURR := &mocks.MockUserRoleRepository{}

	mockRR := &mocks.MockRoleRepository{}

	RS := NewRoleService(mockRR, mockUR, mockURR)

	err := RS.Store(context.TODO(), mockRole, mockUser)

	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRoleStoreErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockRole := &domain.Role{
		Name:        "Foo",
		Description: "Foobar",
	}

	mockErr := domain.NewInternalErr()
	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}
	mockRR.
		On("Create", mock.AnythingOfType("*context.emptyCtx"), mockRole).
		Return(mockErr)

	RS := NewRoleService(mockRR, mockUR, mockURR)

	err := RS.Store(context.TODO(), mockRole, mockUser)

	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRoleStoreGetAllUserErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockRole := &domain.Role{
		Name:        "Foo",
		Description: "Foobar",
	}

	mockErr := domain.NewInternalErr()
	mockUR := &mocks.MockUserRepository{}
	mockUR.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(nil, mockErr)
	mockURR := &mocks.MockUserRoleRepository{}

	mockRR := &mocks.MockRoleRepository{}
	mockRR.
		On("Create", mock.AnythingOfType("*context.emptyCtx"), mockRole).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).(*domain.Role)
			arg.ID = 1
		})

	RS := NewRoleService(mockRR, mockUR, mockURR)

	err := RS.Store(context.TODO(), mockRole, mockUser)

	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRoleStoreCreateUserRolesErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockUsers := &[]domain.User{
		{
			ID:         1,
			Permission: domain.ADMIN,
		},
		{
			ID:         2,
			Permission: domain.GUEST,
		},
	}

	mockRole := &domain.Role{
		Name:        "Foo",
		Description: "Foobar",
	}

	mockUserRoles := &[]domain.UserRole{
		{
			UserID: 1,
			RoleID: 1,
		},
		{
			UserID: 2,
			RoleID: 1,
		},
	}

	mockErr := domain.NewInternalErr()
	mockUR := &mocks.MockUserRepository{}
	mockUR.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(mockUsers, nil)
	mockURR := &mocks.MockUserRoleRepository{}
	mockURR.
		On("CreateBatch", mock.AnythingOfType("*context.emptyCtx"), mockUserRoles).
		Return(mockErr)

	mockRR := &mocks.MockRoleRepository{}
	mockRR.
		On("Create", mock.AnythingOfType("*context.emptyCtx"), mockRole).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).(*domain.Role)
			arg.ID = 1
		})

	RS := NewRoleService(mockRR, mockUR, mockURR)

	err := RS.Store(context.TODO(), mockRole, mockUser)

	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRoleDeleteCorrect(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	rid := int64(1)

	mockUR := &mocks.MockUserRepository{}

	mockURR := &mocks.MockUserRoleRepository{}
	mockURR.
		On("DeleteByRID", mock.AnythingOfType("*context.emptyCtx"), rid).
		Return(nil)

	mockRR := &mocks.MockRoleRepository{}
	mockRR.
		On("Delete", mock.AnythingOfType("*context.emptyCtx"), rid).
		Return(nil)

	RS := NewRoleService(mockRR, mockUR, mockURR)
	err := RS.Remove(context.TODO(), rid, mockUser)

	assert.NoError(t, err)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRoleDeleteNoPermission(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
	}

	rid := int64(1)

	mockErr := domain.NewNotAuthorizedErr("")
	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	RS := NewRoleService(mockRR, mockUR, mockURR)
	err := RS.Remove(context.TODO(), rid, mockUser)

	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRoleDeleteErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	rid := int64(1)

	mockErr := domain.NewInternalErr()
	mockUR := &mocks.MockUserRepository{}

	mockURR := &mocks.MockUserRoleRepository{}

	mockRR := &mocks.MockRoleRepository{}
	mockRR.
		On("Delete", mock.AnythingOfType("*context.emptyCtx"), rid).
		Return(mockErr)

	RS := NewRoleService(mockRR, mockUR, mockURR)
	err := RS.Remove(context.TODO(), rid, mockUser)

	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRoleDeleteDeleteByRIDErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	rid := int64(1)

	mockErr := domain.NewInternalErr()
	mockUR := &mocks.MockUserRepository{}

	mockURR := &mocks.MockUserRoleRepository{}
	mockURR.
		On("DeleteByRID", mock.AnythingOfType("*context.emptyCtx"), rid).
		Return(mockErr)

	mockRR := &mocks.MockRoleRepository{}
	mockRR.
		On("Delete", mock.AnythingOfType("*context.emptyCtx"), rid).
		Return(nil)

	RS := NewRoleService(mockRR, mockUR, mockURR)
	err := RS.Remove(context.TODO(), rid, mockUser)

	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}
