package service

import (
	"context"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStoreCorrect(t *testing.T) {

	mockBundle := &domain.Bundle{
		Name:     "Foo",
		ParentID: 0,
	}

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockBR := new(mocks.MockBundleRepository)
	mockBR.
		On("Create", mock.AnythingOfType("*context.emptyCtx"), mockBundle).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).(*domain.Bundle)
			arg.ID = 1
		})

	BS := NewBundleService(mockBR)
	ctx := context.TODO()
	err := BS.Store(ctx, mockBundle, mockUser)

	assert.NoError(t, err)
	assert.NotEmpty(t, mockBundle.ID)
	mockBR.AssertExpectations(t)

}

func TestStoreNoClearance(t *testing.T) {

	mockBundle := &domain.Bundle{
		Name:     "Foo",
		ParentID: 0,
	}

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockBR := new(mocks.MockBundleRepository)
	mockErr := domain.NewNotAuthorizedErr("")
	BS := NewBundleService(mockBR)
	ctx := context.TODO()
	err := BS.Store(ctx, mockBundle, mockUser)

	assert.ErrorAs(t, err, &mockErr)
	mockBR.AssertExpectations(t)
}

func TestStoreNegativeParentID(t *testing.T) {

	mockBundle := &domain.Bundle{
		Name:     "Foo",
		ParentID: -1,
	}

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockBR := new(mocks.MockBundleRepository)
	mockErr := domain.NewBadRequestErr("")
	BS := NewBundleService(mockBR)
	ctx := context.TODO()
	err := BS.Store(ctx, mockBundle, mockUser)

	assert.ErrorAs(t, err, &mockErr)
	mockBR.AssertExpectations(t)
}

func TestStoreParentNotExist(t *testing.T) {

	mockBundle := &domain.Bundle{
		Name:     "Foo",
		ParentID: 1,
	}

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockErr := domain.NewRecordNotFoundErr("", "")
	mockBR := new(mocks.MockBundleRepository)
	mockBR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockBundle.ParentID).
		Return(nil, mockErr)

	BS := NewBundleService(mockBR)
	ctx := context.TODO()
	err := BS.Store(ctx, mockBundle, mockUser)

	assert.ErrorAs(t, err, &mockErr)
	mockBR.AssertExpectations(t)
}
