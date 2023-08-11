package service_test

import (
	"context"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestSetlistRoleFetch(t *testing.T) {
	setlists := &[]domain.Setlist{
		{
			ID: 1,
		},
		{
			ID: 2,
		},
	}

	setlistRoles := &[]domain.SetlistRole{
		{
			ID:         1,
			SetlistID:  1,
			UserRoleID: 1,
		},
		{
			ID:         2,
			SetlistID:  1,
			UserRoleID: 2,
		},
	}

	t.Run("Correct no setlists", func(t *testing.T) {
		t.Parallel()

		mockSLRR := &mocks.MockSetlistRoleRepository{}
		mockSLR := &mocks.MockSetlistRepository{}
		mockURR := &mocks.MockUserRoleRepository{}

		mockSLRR.
			On("Get", context.TODO(), []int64{}).
			Return(setlistRoles, nil)

		setlistRoleService := service.NewSetlistRoleService(mockSLRR, mockSLR, mockURR)

		retrievedSetlistRoles, err := setlistRoleService.Fetch(context.TODO(), nil)

		assert.NoError(t, err)
		assert.Equal(t, setlistRoles, retrievedSetlistRoles)

		mockSLRR.AssertExpectations(t)
		mockSLR.AssertExpectations(t)
		mockURR.AssertExpectations(t)
	})

	t.Run("Correct with setlists", func(t *testing.T) {
		t.Parallel()

		mockSLRR := &mocks.MockSetlistRoleRepository{}
		mockSLR := &mocks.MockSetlistRepository{}
		mockURR := &mocks.MockUserRoleRepository{}

		mockSLRR.
			On("Get", context.TODO(), []int64{1, 2}).
			Return(setlistRoles, nil)

		setlistRoleService := service.NewSetlistRoleService(mockSLRR, mockSLR, mockURR)

		retrievedSetlistRoles, err := setlistRoleService.Fetch(context.TODO(), setlists)

		assert.NoError(t, err)
		assert.Equal(t, setlistRoles, retrievedSetlistRoles)

		mockSLRR.AssertExpectations(t)
		mockSLR.AssertExpectations(t)
		mockURR.AssertExpectations(t)
	})

	t.Run("Fail SetlistRole Repo Get error", func(t *testing.T) {
		t.Parallel()

		expErr := domain.NewInternalErr()
		mockSLRR := &mocks.MockSetlistRoleRepository{}
		mockSLR := &mocks.MockSetlistRepository{}
		mockURR := &mocks.MockUserRoleRepository{}

		mockSLRR.
			On("Get", context.TODO(), []int64{}).
			Return(nil, expErr)

		setlistRoleService := service.NewSetlistRoleService(mockSLRR, mockSLR, mockURR)

		retrievedSetlistRoles, err := setlistRoleService.Fetch(context.TODO(), nil)

		assert.ErrorAs(t, err, &expErr)
		assert.Nil(t, retrievedSetlistRoles)

		mockSLRR.AssertExpectations(t)
		mockSLR.AssertExpectations(t)
		mockURR.AssertExpectations(t)
	})
}

func TestSetlistRoleStore(t *testing.T) {
	admin := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	editor := &domain.User{
		ID:         2,
		Permission: domain.EDITOR,
	}

	userRoles := &[]domain.UserRole{
		{
			ID:     1,
			UserID: admin.ID,
			RoleID: 1,
		},
		{
			ID:     2,
			UserID: editor.ID,
			RoleID: 1,
		},
	}

	setlists := &[]domain.Setlist{
		{
			ID: 1,
		},
		{
			ID: 2,
		},
	}

	setlistRoles := &[]domain.SetlistRole{
		{
			ID:         1,
			SetlistID:  1,
			UserRoleID: 1,
		},
		{
			ID:         2,
			SetlistID:  1,
			UserRoleID: 2,
		},
	}

	t.Run("Correct", func(t *testing.T) {
		t.Parallel()

		mockSLRR := &mocks.MockSetlistRoleRepository{}
		mockSLR := &mocks.MockSetlistRepository{}
		mockURR := &mocks.MockUserRoleRepository{}

		mockSLRR.
			On("Create", context.TODO(), setlistRoles).
			Return(nil)

		mockSLR.
			On("GetByIDs", context.TODO(), []int64{1, 1}).
			Return(setlists, nil)

		setlistRoleService := service.NewSetlistRoleService(mockSLRR, mockSLR, mockURR)

		err := setlistRoleService.Store(context.TODO(), setlistRoles, admin)

		assert.NoError(t, err)

		mockSLRR.AssertExpectations(t)
		mockSLR.AssertExpectations(t)
		mockURR.AssertExpectations(t)
	})

	t.Run("Fail user nil", func(t *testing.T) {
		t.Parallel()

		mockSLRR := &mocks.MockSetlistRoleRepository{}
		mockSLR := &mocks.MockSetlistRepository{}
		mockURR := &mocks.MockUserRoleRepository{}

		setlistRoleService := service.NewSetlistRoleService(mockSLRR, mockSLR, mockURR)

		err := setlistRoleService.Store(context.TODO(), setlistRoles, nil)

		assert.EqualError(t, err, "No user specified")

		mockSLRR.AssertExpectations(t)
		mockSLR.AssertExpectations(t)
		mockURR.AssertExpectations(t)
	})

	t.Run("Fail setlistroles nil", func(t *testing.T) {
		t.Parallel()

		mockSLRR := &mocks.MockSetlistRoleRepository{}
		mockSLR := &mocks.MockSetlistRepository{}
		mockURR := &mocks.MockUserRoleRepository{}

		setlistRoleService := service.NewSetlistRoleService(mockSLRR, mockSLR, mockURR)

		err := setlistRoleService.Store(context.TODO(), nil, admin)

		assert.EqualError(t, err, "No setlistroles given")

		mockSLRR.AssertExpectations(t)
		mockSLR.AssertExpectations(t)
		mockURR.AssertExpectations(t)
	})

	t.Run("Fail Get UserRoles error", func(t *testing.T) {
		t.Parallel()

		expErr := domain.NewInternalErr()
		mockSLRR := &mocks.MockSetlistRoleRepository{}
		mockSLR := &mocks.MockSetlistRepository{}
		mockURR := &mocks.MockUserRoleRepository{}

		mockURR.
			On("Get", context.TODO(), []int64{1, 2}).
			Return(nil, expErr)

		setlistRoleService := service.NewSetlistRoleService(mockSLRR, mockSLR, mockURR)

		err := setlistRoleService.Store(context.TODO(), setlistRoles, editor)

		assert.EqualError(t, err, expErr.Error())

		mockSLRR.AssertExpectations(t)
		mockSLR.AssertExpectations(t)
		mockURR.AssertExpectations(t)
	})

	t.Run("Fail Not Authorized setlistrole change", func(t *testing.T) {
		t.Parallel()

		mockSLRR := &mocks.MockSetlistRoleRepository{}
		mockSLR := &mocks.MockSetlistRepository{}
		mockURR := &mocks.MockUserRoleRepository{}

		mockURR.
			On("Get", context.TODO(), []int64{1, 2}).
			Return(userRoles, nil)

		setlistRoleService := service.NewSetlistRoleService(mockSLRR, mockSLR, mockURR)

		err := setlistRoleService.Store(context.TODO(), setlistRoles, editor)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "Setlist Role of someone else")

		mockSLRR.AssertExpectations(t)
		mockSLR.AssertExpectations(t)
		mockURR.AssertExpectations(t)
	})

	t.Run("Fail Setlist GetByIDs error", func(t *testing.T) {
		t.Parallel()

		expErr := domain.NewInternalErr()
		mockSLRR := &mocks.MockSetlistRoleRepository{}
		mockSLR := &mocks.MockSetlistRepository{}
		mockURR := &mocks.MockUserRoleRepository{}

		mockSLR.
			On("GetByIDs", context.TODO(), []int64{1, 1}).
			Return(nil, expErr)

		setlistRoleService := service.NewSetlistRoleService(mockSLRR, mockSLR, mockURR)

		err := setlistRoleService.Store(context.TODO(), setlistRoles, admin)

		assert.Error(t, err)
		assert.EqualError(t, err, expErr.Error())

		mockSLRR.AssertExpectations(t)
		mockSLR.AssertExpectations(t)
		mockURR.AssertExpectations(t)
	})

	t.Run("Fail SetlistRole Create error", func(t *testing.T) {
		t.Parallel()

		expErr := domain.NewInternalErr()
		mockSLRR := &mocks.MockSetlistRoleRepository{}
		mockSLR := &mocks.MockSetlistRepository{}
		mockURR := &mocks.MockUserRoleRepository{}

		mockSLRR.
			On("Create", context.TODO(), setlistRoles).
			Return(expErr)

		mockSLR.
			On("GetByIDs", context.TODO(), []int64{1, 1}).
			Return(setlists, nil)

		setlistRoleService := service.NewSetlistRoleService(mockSLRR, mockSLR, mockURR)

		err := setlistRoleService.Store(context.TODO(), setlistRoles, admin)

		assert.Error(t, err)
		assert.EqualError(t, err, expErr.Error())

		mockSLRR.AssertExpectations(t)
		mockSLR.AssertExpectations(t)
		mockURR.AssertExpectations(t)
	})
}

func TestRemove(t *testing.T) {
	admin := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	editor := &domain.User{
		ID:         2,
		Permission: domain.EDITOR,
	}

	userRoles := &[]domain.UserRole{
		{
			ID:     1,
			UserID: admin.ID,
			RoleID: 1,
		},
		{
			ID:     2,
			UserID: editor.ID,
			RoleID: 1,
		},
	}

	// setlists := &[]domain.Setlist{
	// 	{
	// 		ID: 1,
	// 	},
	// 	{
	// 		ID: 2,
	// 	},
	// }

	setlistRoles := &[]domain.SetlistRole{
		{
			ID:         1,
			SetlistID:  1,
			UserRoleID: 1,
		},
		{
			ID:         2,
			SetlistID:  1,
			UserRoleID: 2,
		},
	}

	t.Run("Correct", func(t *testing.T) {
		t.Parallel()

		mockSLRR := &mocks.MockSetlistRoleRepository{}
		mockSLR := &mocks.MockSetlistRepository{}
		mockURR := &mocks.MockUserRoleRepository{}

		setlistRoleService := service.NewSetlistRoleService(mockSLRR, mockSLR, mockURR)

		mockSLRR.
			On("Delete", context.TODO(), []int64{1, 2}).
			Return(nil)

		mockSLRR.
			On("Get", context.TODO(), []int64{1, 2}).
			Return(nil, nil)

		err := setlistRoleService.Remove(context.TODO(), []int64{1, 2}, admin)

		assert.NoError(t, err)

		mockSLRR.AssertExpectations(t)
		mockSLR.AssertExpectations(t)
		mockURR.AssertExpectations(t)
	})

	t.Run("Fail principal is nil", func(t *testing.T) {
		t.Parallel()

		expErr := domain.NewNotAuthorizedErr("No user specified")
		mockSLRR := &mocks.MockSetlistRoleRepository{}
		mockSLR := &mocks.MockSetlistRepository{}
		mockURR := &mocks.MockUserRoleRepository{}

		setlistRoleService := service.NewSetlistRoleService(mockSLRR, mockSLR, mockURR)

		err := setlistRoleService.Remove(context.TODO(), []int64{1, 2}, nil)

		assert.Error(t, err)
		assert.EqualError(t, err, expErr.Error())
		mockSLRR.AssertExpectations(t)
		mockSLR.AssertExpectations(t)
		mockURR.AssertExpectations(t)
	})

	t.Run("Correct no setlists given", func(t *testing.T) {
		t.Parallel()

		mockSLRR := &mocks.MockSetlistRoleRepository{}
		mockSLR := &mocks.MockSetlistRepository{}
		mockURR := &mocks.MockUserRoleRepository{}

		setlistRoleService := service.NewSetlistRoleService(mockSLRR, mockSLR, mockURR)

		err := setlistRoleService.Remove(context.TODO(), []int64{}, admin)

		assert.NoError(t, err)

		mockSLRR.AssertExpectations(t)
		mockSLR.AssertExpectations(t)
		mockURR.AssertExpectations(t)
	})

	t.Run("Fail SetlistRole Get error", func(t *testing.T) {
		t.Parallel()

		expErr := domain.NewInternalErr()
		mockSLRR := &mocks.MockSetlistRoleRepository{}
		mockSLR := &mocks.MockSetlistRepository{}
		mockURR := &mocks.MockUserRoleRepository{}

		mockSLRR.
			On("Get", context.TODO(), []int64{1, 2}).
			Return(nil, expErr)

		setlistRoleService := service.NewSetlistRoleService(mockSLRR, mockSLR, mockURR)

		err := setlistRoleService.Remove(context.TODO(), []int64{1, 2}, editor)

		assert.Error(t, err)
		assert.EqualError(t, err, expErr.Error())
		mockSLRR.AssertExpectations(t)
		mockSLR.AssertExpectations(t)
		mockURR.AssertExpectations(t)
	})

	t.Run("Fail SetlistRole Get error", func(t *testing.T) {
		t.Parallel()

		expErr := domain.NewInternalErr()
		mockSLRR := &mocks.MockSetlistRoleRepository{}
		mockSLR := &mocks.MockSetlistRepository{}
		mockURR := &mocks.MockUserRoleRepository{}

		mockSLRR.
			On("Get", context.TODO(), []int64{1, 2}).
			Return(setlistRoles, nil)

		mockURR.
			On("Get", context.TODO(), []int64{1, 2}).
			Return(nil, expErr)

		setlistRoleService := service.NewSetlistRoleService(mockSLRR, mockSLR, mockURR)

		err := setlistRoleService.Remove(context.TODO(), []int64{1, 2}, editor)

		assert.Error(t, err)
		assert.EqualError(t, err, expErr.Error())
		mockSLRR.AssertExpectations(t)
		mockSLR.AssertExpectations(t)
		mockURR.AssertExpectations(t)
	})

	t.Run("Fail userrole does not match principal", func(t *testing.T) {
		t.Parallel()

		expErr := domain.NewNotAuthorizedErr("Cannot change the Setlist Role of someone else")
		mockSLRR := &mocks.MockSetlistRoleRepository{}
		mockSLR := &mocks.MockSetlistRepository{}
		mockURR := &mocks.MockUserRoleRepository{}

		mockSLRR.
			On("Get", context.TODO(), []int64{1, 2}).
			Return(setlistRoles, nil)

		mockURR.
			On("Get", context.TODO(), []int64{1, 2}).
			Return(userRoles, nil)

		setlistRoleService := service.NewSetlistRoleService(mockSLRR, mockSLR, mockURR)

		err := setlistRoleService.Remove(context.TODO(), []int64{1, 2}, editor)

		assert.Error(t, err)
		assert.EqualError(t, err, expErr.Error())
		mockSLRR.AssertExpectations(t)
		mockSLR.AssertExpectations(t)
		mockURR.AssertExpectations(t)
	})

	t.Run("Fail SetlistRole Delete error", func(t *testing.T) {
		t.Parallel()

		expErr := domain.NewInternalErr()
		mockSLRR := &mocks.MockSetlistRoleRepository{}
		mockSLR := &mocks.MockSetlistRepository{}
		mockURR := &mocks.MockUserRoleRepository{}

		mockSLRR.
			On("Get", context.TODO(), []int64{1, 2}).
			Return(setlistRoles, nil)

		mockSLRR.
			On("Delete", context.TODO(), []int64{1, 2}).
			Return(expErr)

		setlistRoleService := service.NewSetlistRoleService(mockSLRR, mockSLR, mockURR)

		err := setlistRoleService.Remove(context.TODO(), []int64{1, 2}, admin)

		assert.Error(t, err)
		assert.EqualError(t, err, expErr.Error())
		mockSLRR.AssertExpectations(t)
		mockSLR.AssertExpectations(t)
		mockURR.AssertExpectations(t)
	})
}
