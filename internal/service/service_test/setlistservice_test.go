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

func TestSetlistFetchByIDCorrect(t *testing.T) {
	t.Parallel()

	slid := int64(1)
	mockSetlist := &domain.Setlist{
		ID:        1,
		CreatorID: 1,
		Name:      "Foo",
	}

	mockUR := &mocks.MockUserRepository{}
	mockSLR := &mocks.MockSetlistRepository{}

	mockSLR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), slid).
		Return(mockSetlist, nil)

	slr := service.NewSetlistService(mockUR, mockSLR)

	setlist, err := slr.FetchByID(context.TODO(), slid)
	assert.NoError(t, err)
	assert.Equal(t, mockSetlist, setlist)
}

func TestSetlistFetchByIDErr(t *testing.T) {
	t.Parallel()

	slid := int64(1)

	expErr := domain.NewInternalErr()
	mockUR := &mocks.MockUserRepository{}
	mockSLR := &mocks.MockSetlistRepository{}

	mockSLR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), slid).
		Return(nil, expErr)

	slr := service.NewSetlistService(mockUR, mockSLR)

	setlist, err := slr.FetchByID(context.TODO(), slid)
	assert.ErrorAs(t, err, &expErr)
	assert.Nil(t, setlist)
}

func TestSetlistFetchAllCorrect(t *testing.T) {
	t.Parallel()

	mockSetlists := &[]domain.Setlist{
		{
			ID:        1,
			CreatorID: 1,
			Name:      "Foo",
		},
		{
			ID:        2,
			CreatorID: 1,
			Name:      "Bar",
		},
	}

	mockUR := &mocks.MockUserRepository{}
	mockSLR := &mocks.MockSetlistRepository{}

	mockSLR.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(mockSetlists, nil)

	slr := service.NewSetlistService(mockUR, mockSLR)

	setlists, err := slr.FetchAll(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, mockSetlists, setlists)
}

func TestSetlistFetchAllErr(t *testing.T) {
	t.Parallel()

	expErr := domain.NewInternalErr()
	mockUR := &mocks.MockUserRepository{}
	mockSLR := &mocks.MockSetlistRepository{}

	mockSLR.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(nil, expErr)

	slr := service.NewSetlistService(mockUR, mockSLR)

	setlist, err := slr.FetchAll(context.TODO())
	assert.ErrorAs(t, err, &expErr)
	assert.Nil(t, setlist)
}

func TestSetlistFetchAllGlobalCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}
	mockSetlists := &[]domain.Setlist{
		{
			ID:        1,
			CreatorID: 1,
			Name:      "Foo",
		},
		{
			ID:        2,
			CreatorID: 1,
			Name:      "Bar",
		},
	}

	mockUR := &mocks.MockUserRepository{}
	mockSLR := &mocks.MockSetlistRepository{}

	mockSLR.
		On("GetAllGlobal", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(mockSetlists, nil)

	slr := service.NewSetlistService(mockUR, mockSLR)

	setlists, err := slr.FetchAllGlobal(context.TODO(), mockUser)
	assert.NoError(t, err)
	assert.Equal(t, mockSetlists, setlists)
}

func TestSetlistFetchAllGlobalErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	expErr := domain.NewInternalErr()
	mockUR := &mocks.MockUserRepository{}
	mockSLR := &mocks.MockSetlistRepository{}

	mockSLR.
		On("GetAllGlobal", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(nil, expErr)

	slr := service.NewSetlistService(mockUR, mockSLR)

	setlist, err := slr.FetchAllGlobal(context.TODO(), mockUser)
	assert.ErrorAs(t, err, &expErr)
	assert.Nil(t, setlist)
}
