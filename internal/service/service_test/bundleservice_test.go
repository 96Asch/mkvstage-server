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

func TestStoreCorrect(t *testing.T) {
	t.Parallel()

	mockBundle := &domain.Bundle{
		Name:     "Foo",
		ParentID: 0,
	}

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockBR := &mocks.MockBundleRepository{}

	mockBR.
		On("Create", mock.AnythingOfType("*context.emptyCtx"), mockBundle).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg, ok := args.Get(1).(*domain.Bundle)
			assert.True(t, ok)
			arg.ID = 1
		})

	BS := service.NewBundleService(mockBR)
	ctx := context.TODO()

	err := BS.Store(ctx, mockBundle, mockUser)
	assert.NoError(t, err)
	assert.NotEmpty(t, mockBundle.ID)
	mockBR.AssertExpectations(t)
}

func TestStoreNoClearance(t *testing.T) {
	t.Parallel()

	mockBundle := &domain.Bundle{
		Name:     "Foo",
		ParentID: 0,
	}

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockErr := domain.NewNotAuthorizedErr("")
	mockBR := &mocks.MockBundleRepository{}

	BS := service.NewBundleService(mockBR)
	ctx := context.TODO()

	err := BS.Store(ctx, mockBundle, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	mockBR.AssertExpectations(t)
}

func TestStoreNegativeParentID(t *testing.T) {
	t.Parallel()

	mockBundle := &domain.Bundle{
		Name:     "Foo",
		ParentID: -1,
	}

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockErr := domain.NewBadRequestErr("")
	mockBR := &mocks.MockBundleRepository{}

	BS := service.NewBundleService(mockBR)
	ctx := context.TODO()

	err := BS.Store(ctx, mockBundle, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	mockBR.AssertExpectations(t)
}

func TestStoreParentNotExist(t *testing.T) {
	t.Parallel()

	mockBundle := &domain.Bundle{
		Name:     "Foo",
		ParentID: 1,
	}

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockErr := domain.NewBadRequestErr("")
	mockBR := &mocks.MockBundleRepository{}

	mockBR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockBundle.ParentID).
		Return(nil, mockErr)

	BS := service.NewBundleService(mockBR)
	ctx := context.TODO()

	err := BS.Store(ctx, mockBundle, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	mockBR.AssertExpectations(t)
}

func TestFetchByID(t *testing.T) {
	t.Parallel()

	mockBundle := &domain.Bundle{
		ID:       1,
		Name:     "Foo",
		ParentID: 1,
	}

	mockBR := &mocks.MockBundleRepository{}

	mockBR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockBundle.ID).
		Return(mockBundle, nil)

	BS := service.NewBundleService(mockBR)
	ctx := context.TODO()

	bundle, err := BS.FetchByID(ctx, mockBundle.ID)
	assert.NoError(t, err)
	assert.Equal(t, mockBundle, bundle)
	mockBR.AssertExpectations(t)
}

func TestFetchAll(t *testing.T) {
	t.Parallel()

	mockBundles := &[]domain.Bundle{
		{
			ID:       1,
			Name:     "Foo",
			ParentID: 0,
		},
		{
			ID:       2,
			Name:     "Bar",
			ParentID: 0,
		},
	}

	mockBR := &mocks.MockBundleRepository{}

	mockBR.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(mockBundles, nil)

	BS := service.NewBundleService(mockBR)
	ctx := context.TODO()

	bundles, err := BS.FetchAll(ctx)
	assert.NoError(t, err)
	assert.ElementsMatch(t, *mockBundles, *bundles)
	mockBR.AssertExpectations(t)
}

func TestRemoveCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockBundle := &domain.Bundle{
		ID:       1,
		Name:     "Foo",
		ParentID: 1,
	}

	mockBundles := &[]domain.Bundle{
		*mockBundle,
		{
			ID:       2,
			Name:     "Bar",
			ParentID: 0,
		},
	}

	mockBR := &mocks.MockBundleRepository{}

	mockBR.
		On("Delete", mock.AnythingOfType("*context.emptyCtx"), mockBundle.ID).
		Return(nil)
	mockBR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockBundle.ID).
		Return(mockBundle, nil)
	mockBR.
		On("GetLeaves", mock.AnythingOfType("*context.emptyCtx")).
		Return(mockBundles, nil)

	BS := service.NewBundleService(mockBR)
	ctx := context.TODO()

	err := BS.Remove(ctx, mockBundle.ID, mockUser)
	assert.NoError(t, err)
	mockBR.AssertExpectations(t)
}

func TestRemoveNotLeaf(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockBundle := &domain.Bundle{
		ID:       1,
		Name:     "Foo",
		ParentID: 1,
	}

	mockBundles := &[]domain.Bundle{
		{
			ID:       2,
			Name:     "Bar",
			ParentID: 0,
		},
	}

	mockBR := &mocks.MockBundleRepository{}

	mockBR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockBundle.ID).
		Return(mockBundle, nil)
	mockBR.
		On("GetLeaves", mock.AnythingOfType("*context.emptyCtx")).
		Return(mockBundles, nil)

	BS := service.NewBundleService(mockBR)
	ctx := context.TODO()

	err := BS.Remove(ctx, mockBundle.ID, mockUser)
	mockErr := domain.NewBadRequestErr("")
	assert.ErrorAs(t, err, &mockErr)
	mockBR.AssertExpectations(t)
}

func TestRemoveGetLeavesErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockBundle := &domain.Bundle{
		ID:       1,
		Name:     "Foo",
		ParentID: 1,
	}

	mockErr := domain.NewInternalErr()
	mockBR := &mocks.MockBundleRepository{}

	mockBR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockBundle.ID).
		Return(mockBundle, nil)
	mockBR.
		On("GetLeaves", mock.AnythingOfType("*context.emptyCtx")).
		Return(nil, mockErr)

	BS := service.NewBundleService(mockBR)
	ctx := context.TODO()

	err := BS.Remove(ctx, mockBundle.ID, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	mockBR.AssertExpectations(t)
}

func TestRemoveNoClearance(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockBundle := &domain.Bundle{
		ID:       1,
		Name:     "Foo",
		ParentID: 1,
	}

	mockBR := &mocks.MockBundleRepository{}

	BS := service.NewBundleService(mockBR)
	ctx := context.TODO()

	err := BS.Remove(ctx, mockBundle.ID, mockUser)
	mockErr := domain.NewNotAuthorizedErr("")
	assert.ErrorAs(t, err, &mockErr)
	mockBR.AssertExpectations(t)
}

func TestDeleteNoRecord(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockBundle := &domain.Bundle{
		ID:       1,
		Name:     "Foo",
		ParentID: 1,
	}

	mockErr := domain.NewRecordNotFoundErr("", "")
	mockBR := &mocks.MockBundleRepository{}

	mockBR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockBundle.ID).
		Return(nil, mockErr)

	BS := service.NewBundleService(mockBR)
	ctx := context.TODO()

	err := BS.Remove(ctx, mockBundle.ID, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	mockBR.AssertExpectations(t)
}

func TestUpdateCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockBundle := &domain.Bundle{
		ID:       1,
		Name:     "Foo",
		ParentID: 1,
	}

	mockBR := &mocks.MockBundleRepository{}

	mockBR.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), mockBundle).
		Return(nil)
	mockBR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockBundle.ID).
		Return(mockBundle, nil)

	BS := service.NewBundleService(mockBR)
	ctx := context.TODO()

	err := BS.Update(ctx, mockBundle, mockUser)
	assert.NoError(t, err)
	mockBR.AssertExpectations(t)
}

func TestUpdateNoRecord(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockBundle := &domain.Bundle{
		ID:       1,
		Name:     "Foo",
		ParentID: 1,
	}

	mockErr := domain.NewRecordNotFoundErr("", "")
	mockBR := &mocks.MockBundleRepository{}

	mockBR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockBundle.ID).
		Return(nil, mockErr)

	BS := service.NewBundleService(mockBR)
	ctx := context.TODO()

	err := BS.Update(ctx, mockBundle, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	mockBR.AssertExpectations(t)
}

func TestUpdateNoClearance(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockBundle := &domain.Bundle{
		ID:       1,
		Name:     "Foo",
		ParentID: 1,
	}

	mockBR := &mocks.MockBundleRepository{}
	BS := service.NewBundleService(mockBR)
	ctx := context.TODO()

	err := BS.Update(ctx, mockBundle, mockUser)
	mockErr := domain.NewNotAuthorizedErr("")
	assert.ErrorAs(t, err, &mockErr)
	mockBR.AssertExpectations(t)
}
