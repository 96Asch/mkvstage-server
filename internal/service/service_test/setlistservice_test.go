package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/datatypes"
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

func TestSetlistUpdateCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockSetlist := &domain.Setlist{
		ID:        1,
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1),
		Name:      "Foobar",
		Order:     datatypes.JSON([]byte(`{order: [1, 2]}`)),
	}

	mockUR := &mocks.MockUserRepository{}
	mockSLR := &mocks.MockSetlistRepository{}

	mockUR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockSetlist.CreatorID).
		Return(mockUser, nil)
	mockSLR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockSetlist.ID).
		Return(mockSetlist, nil)
	mockSLR.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), mockSetlist).
		Return(mockSetlist, nil)

	sls := service.NewSetlistService(mockUR, mockSLR)

	updatedSetlist, err := sls.Update(context.TODO(), mockSetlist, mockUser)
	assert.Equal(t, mockSetlist, updatedSetlist)
	assert.NoError(t, err)
	mockUR.AssertExpectations(t)
	mockSLR.AssertExpectations(t)
}

func TestSetlistUpdateSetlistNotFound(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	mockSetlist := &domain.Setlist{
		ID:        1,
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1),
		Name:      "Foobar",
		Order:     datatypes.JSON([]byte(`{order: [1, 2]}`)),
	}

	mockErr := domain.NewRecordNotFoundErr("", "")
	mockUR := &mocks.MockUserRepository{}
	mockSLR := &mocks.MockSetlistRepository{}

	mockSLR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockSetlist.ID).
		Return(nil, mockErr)

	sls := service.NewSetlistService(mockUR, mockSLR)

	updatedSetlist, err := sls.Update(context.TODO(), mockSetlist, mockUser)
	assert.Nil(t, updatedSetlist)
	assert.ErrorAs(t, err, &mockErr)
	mockUR.AssertExpectations(t)
	mockSLR.AssertExpectations(t)
}

func TestSetlistUpdateNotAuthorized(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	mockSetlist := &domain.Setlist{
		ID:        1,
		CreatorID: 0,
		Deadline:  time.Now().AddDate(0, 0, 1),
		Name:      "Foobar",
		Order:     datatypes.JSON([]byte(`{order: [1, 2]}`)),
	}

	mockErr := domain.NewNotAuthorizedErr("")
	mockUR := &mocks.MockUserRepository{}
	mockSLR := &mocks.MockSetlistRepository{}

	mockSLR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockSetlist.ID).
		Return(mockSetlist, nil)

	sls := service.NewSetlistService(mockUR, mockSLR)

	updatedSetlist, err := sls.Update(context.TODO(), mockSetlist, mockUser)
	assert.Nil(t, updatedSetlist)
	assert.Error(t, err, &mockErr)
	mockUR.AssertExpectations(t)
	mockSLR.AssertExpectations(t)
}

func TestSetlistUpdateInvalidDeadline(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	mockSetlist := &domain.Setlist{
		ID:        1,
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, -1),
		Name:      "Foobar",
		Order:     datatypes.JSON([]byte(`{order: [1, 2]}`)),
	}

	mockErr := domain.NewBadRequestErr("")
	mockUR := &mocks.MockUserRepository{}
	mockSLR := &mocks.MockSetlistRepository{}

	mockSLR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockSetlist.ID).
		Return(mockSetlist, nil)

	sls := service.NewSetlistService(mockUR, mockSLR)

	updatedSetlist, err := sls.Update(context.TODO(), mockSetlist, mockUser)
	assert.Nil(t, updatedSetlist)
	assert.ErrorAs(t, err, &mockErr)
	mockUR.AssertExpectations(t)
	mockSLR.AssertExpectations(t)
}

func TestSetlistUpdateUserGetByIDErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockSetlist := &domain.Setlist{
		ID:        1,
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1),
		Name:      "Foobar",
		Order:     datatypes.JSON([]byte(`{order: [1, 2]}`)),
	}

	mockErr := domain.NewRecordNotFoundErr("", "")
	mockUR := &mocks.MockUserRepository{}
	mockSLR := &mocks.MockSetlistRepository{}

	mockUR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockSetlist.CreatorID).
		Return(nil, mockErr)
	mockSLR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockSetlist.ID).
		Return(mockSetlist, nil)

	sls := service.NewSetlistService(mockUR, mockSLR)

	updatedSetlist, err := sls.Update(context.TODO(), mockSetlist, mockUser)
	assert.Nil(t, updatedSetlist)
	assert.Error(t, err, &mockErr)
	mockUR.AssertExpectations(t)
	mockSLR.AssertExpectations(t)
}

func TestSetlistUpdateErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockSetlist := &domain.Setlist{
		ID:        1,
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1),
		Name:      "Foobar",
		Order:     datatypes.JSON([]byte(`{order: [1, 2]}`)),
	}

	mockErr := domain.NewInternalErr()
	mockUR := &mocks.MockUserRepository{}
	mockSLR := &mocks.MockSetlistRepository{}

	mockUR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockSetlist.CreatorID).
		Return(mockUser, nil)
	mockSLR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockSetlist.ID).
		Return(mockSetlist, nil)
	mockSLR.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), mockSetlist).
		Return(nil, mockErr)

	sls := service.NewSetlistService(mockUR, mockSLR)

	updatedSetlist, err := sls.Update(context.TODO(), mockSetlist, mockUser)
	assert.Nil(t, updatedSetlist)
	assert.ErrorAs(t, err, &mockErr)
	mockUR.AssertExpectations(t)
	mockSLR.AssertExpectations(t)
}
