package service_test

import (
	"context"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRSFetchByIDCorrect(t *testing.T) {
	t.Parallel()

	rid := int64(1)
	mockRole := &domain.Role{
		ID:          rid,
		Name:        "Foo",
		Description: "FooBar",
	}

	mockRR := &mocks.MockRoleRepository{}
	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}

	mockRR.
		On("GetByID", context.TODO(), rid).
		Return(mockRole, nil)

	RS := service.NewRoleService(mockRR, mockUR, mockURR)

	role, err := RS.FetchByID(context.TODO(), rid)
	assert.NoError(t, err)
	assert.Equal(t, mockRole, role)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRSFetchByIDGetErr(t *testing.T) {
	t.Parallel()

	rid := int64(1)
	mockErr := domain.NewRecordNotFoundErr("", "")
	mockRR := &mocks.MockRoleRepository{}
	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}

	mockRR.
		On("GetByID", context.TODO(), rid).
		Return(nil, mockErr)

	RS := service.NewRoleService(mockRR, mockUR, mockURR)

	role, err := RS.FetchByID(context.TODO(), rid)
	assert.ErrorAs(t, err, &mockErr)
	assert.Nil(t, role)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRSFetchAllCorrect(t *testing.T) {
	t.Parallel()

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
	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}

	mockRR.
		On("GetAll", context.TODO()).
		Return(mockRoles, nil)

	RS := service.NewRoleService(mockRR, mockUR, mockURR)

	role, err := RS.FetchAll(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, mockRoles, role)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRSFetchAllGetErr(t *testing.T) {
	t.Parallel()

	mockErr := domain.NewInternalErr()
	mockRR := &mocks.MockRoleRepository{}
	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}

	mockRR.
		On("GetAll", context.TODO()).
		Return(nil, mockErr)

	RS := service.NewRoleService(mockRR, mockUR, mockURR)

	role, err := RS.FetchAll(context.TODO())
	assert.ErrorAs(t, err, &mockErr)
	assert.Nil(t, role)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRSUpdateCorrect(t *testing.T) {
	t.Parallel()

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
	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}

	mockRR.
		On("Update", context.TODO(), mockRole).
		Return(nil)

	mockRR.
		On("GetByID", context.TODO(), rid).
		Return(prevMockRole, nil)

	RS := service.NewRoleService(mockRR, mockUR, mockURR)

	err := RS.Update(context.TODO(), mockRole, mockUser)

	assert.NoError(t, err)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRSUpdateIDZero(t *testing.T) {
	t.Parallel()

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

	RS := service.NewRoleService(mockRR, mockUR, mockURR)

	err := RS.Update(context.TODO(), mockRole, mockUser)
	mockErr := domain.NewBadRequestErr("")
	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRSUpdateNoPermission(t *testing.T) {
	t.Parallel()

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

	RS := service.NewRoleService(mockRR, mockUR, mockURR)

	err := RS.Update(context.TODO(), mockRole, mockUser)
	mockErr := domain.NewNotAuthorizedErr("")
	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRSUpdateNoRecord(t *testing.T) {
	t.Parallel()

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
	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}

	mockRR.
		On("GetByID", context.TODO(), rid).
		Return(nil, mockErr)

	RS := service.NewRoleService(mockRR, mockUR, mockURR)

	err := RS.Update(context.TODO(), mockRole, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRSUpdateErr(t *testing.T) {
	t.Parallel()

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
	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}

	mockRR.
		On("Update", context.TODO(), mockRole).
		Return(mockErr)

	mockRR.
		On("GetByID", context.TODO(), rid).
		Return(prevMockRole, nil)

	RS := service.NewRoleService(mockRR, mockUR, mockURR)

	err := RS.Update(context.TODO(), mockRole, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRSStoreCorrect(t *testing.T) {
	t.Parallel()

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
	mockRR := &mocks.MockRoleRepository{}
	mockURR := &mocks.MockUserRoleRepository{}

	mockUR.
		On("GetAll", context.TODO()).
		Return(mockUsers, nil)
	mockURR.
		On("CreateBatch", context.TODO(), mockUserRoles).
		Return(nil)

	mockRR.
		On("Create", context.TODO(), mockRole).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg, ok := args.Get(1).(*domain.Role)
			assert.True(t, ok)
			arg.ID = 1
		})

	RS := service.NewRoleService(mockRR, mockUR, mockURR)

	err := RS.Store(context.TODO(), mockRole, mockUser)
	assert.NoError(t, err)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRSStoreNoPermission(t *testing.T) {
	t.Parallel()

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

	RS := service.NewRoleService(mockRR, mockUR, mockURR)

	err := RS.Store(context.TODO(), mockRole, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRSStoreErr(t *testing.T) {
	t.Parallel()

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
		On("Create", context.TODO(), mockRole).
		Return(mockErr)

	RS := service.NewRoleService(mockRR, mockUR, mockURR)

	err := RS.Store(context.TODO(), mockRole, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRSStoreGetAllUserErr(t *testing.T) {
	t.Parallel()

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

	mockUR.
		On("GetAll", context.TODO()).
		Return(nil, mockErr)
	mockRR.
		On("Create", context.TODO(), mockRole).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg, ok := args.Get(1).(*domain.Role)
			assert.True(t, ok)
			arg.ID = 1
		})

	RS := service.NewRoleService(mockRR, mockUR, mockURR)

	err := RS.Store(context.TODO(), mockRole, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRSStoreCreateUserRolesErr(t *testing.T) {
	t.Parallel()

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
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	mockUR.
		On("GetAll", context.TODO()).
		Return(mockUsers, nil)
	mockURR.
		On("CreateBatch", context.TODO(), mockUserRoles).
		Return(mockErr)
	mockRR.
		On("Create", context.TODO(), mockRole).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg, ok := args.Get(1).(*domain.Role)
			assert.True(t, ok)
			arg.ID = 1
		})

	RS := service.NewRoleService(mockRR, mockUR, mockURR)

	err := RS.Store(context.TODO(), mockRole, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRSDeleteCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	rid := int64(1)

	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	mockURR.
		On("DeleteByRID", context.TODO(), rid).
		Return(nil)
	mockRR.
		On("Delete", context.TODO(), rid).
		Return(nil)

	RS := service.NewRoleService(mockRR, mockUR, mockURR)

	err := RS.Remove(context.TODO(), rid, mockUser)
	assert.NoError(t, err)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRSDeleteNoPermission(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
	}

	rid := int64(1)
	mockErr := domain.NewNotAuthorizedErr("")
	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	RS := service.NewRoleService(mockRR, mockUR, mockURR)

	err := RS.Remove(context.TODO(), rid, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRSDeleteErr(t *testing.T) {
	t.Parallel()

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
		On("Delete", context.TODO(), rid).
		Return(mockErr)

	RS := service.NewRoleService(mockRR, mockUR, mockURR)

	err := RS.Remove(context.TODO(), rid, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}

func TestRSDeleteDeleteByRIDErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	rid := int64(1)

	mockErr := domain.NewInternalErr()
	mockUR := &mocks.MockUserRepository{}
	mockURR := &mocks.MockUserRoleRepository{}
	mockRR := &mocks.MockRoleRepository{}

	mockURR.
		On("DeleteByRID", context.TODO(), rid).
		Return(mockErr)
	mockRR.
		On("Delete", context.TODO(), rid).
		Return(nil)

	RS := service.NewRoleService(mockRR, mockUR, mockURR)

	err := RS.Remove(context.TODO(), rid, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	mockRR.AssertExpectations(t)
	mockUR.AssertExpectations(t)
	mockURR.AssertExpectations(t)
}
