package service

import (
	"context"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/datatypes"
)

func TestCreateSongCorrect(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	mockSong := &domain.Song{
		CreatorID:  mockUser.ID,
		Title:      "Foo",
		Subtitle:   "Bar",
		Key:        "A",
		Bpm:        120,
		ChordSheet: datatypes.JSON([]byte(`{"Verse" : "Foobar"}`)),
	}

	mockUR := new(mocks.MockUserRepository)
	mockUR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockSong.CreatorID).
		Return(mockUser, nil)
	mockSR := new(mocks.MockSongRepository)
	mockSR.
		On("Create", mock.AnythingOfType("*context.emptyCtx"), mockSong).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).(*domain.Song)
			arg.ID = 1
		})

	ss := NewSongService(mockUR, mockSR)
	ctx := context.TODO()
	err := ss.Store(ctx, mockSong, mockUser)

	assert.NoError(t, err)
	assert.NotEmpty(t, mockSong.ID)
	mockSR.AssertExpectations(t)
}

func TestCreateSongNoClearance(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockSong := &domain.Song{
		CreatorID:  mockUser.ID,
		Title:      "Foo",
		Subtitle:   "Bar",
		Key:        "A",
		Bpm:        120,
		ChordSheet: datatypes.JSON([]byte(`{"Verse" : "Foobar"}`)),
	}
	mockUR := new(mocks.MockUserRepository)
	mockSR := new(mocks.MockSongRepository)

	ss := NewSongService(mockUR, mockSR)
	ctx := context.TODO()
	err := ss.Store(ctx, mockSong, mockUser)

	mockErr := domain.NewNotAuthorizedErr("")
	assert.ErrorAs(t, err, &mockErr)
	assert.Empty(t, mockSong.ID)
	mockSR.AssertExpectations(t)
}

func TestCreateSongInvalidKey(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	mockSong := &domain.Song{
		CreatorID:  mockUser.ID,
		Title:      "Foo",
		Subtitle:   "Bar",
		Key:        "Q",
		Bpm:        120,
		ChordSheet: datatypes.JSON([]byte(`{"Verse" : "Foobar"}`)),
	}

	mockSR := new(mocks.MockSongRepository)
	mockUR := new(mocks.MockUserRepository)

	ss := NewSongService(mockUR, mockSR)
	ctx := context.TODO()
	err := ss.Store(ctx, mockSong, mockUser)

	mockErr := domain.NewBadRequestErr("")
	assert.ErrorAs(t, err, &mockErr)
	assert.Empty(t, mockSong.ID)
	mockSR.AssertExpectations(t)
}

func TestCreateSongCreatorNotExists(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	mockSong := &domain.Song{
		CreatorID:  mockUser.ID,
		Title:      "Foo",
		Subtitle:   "Bar",
		Key:        "A",
		Bpm:        120,
		ChordSheet: datatypes.JSON([]byte(`{"Verse" : "Foobar"}`)),
	}

	mockErr := domain.NewRecordNotFoundErr("", "")
	mockUR := new(mocks.MockUserRepository)
	mockUR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockSong.CreatorID).
		Return(nil, mockErr)
	mockSR := new(mocks.MockSongRepository)

	ss := NewSongService(mockUR, mockSR)
	ctx := context.TODO()
	err := ss.Store(ctx, mockSong, mockUser)

	assert.ErrorAs(t, err, &mockErr)
	assert.Empty(t, mockSong.ID)
	mockSR.AssertExpectations(t)
}

func TestCreateSongInvalidChordsheet(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	mockSong := &domain.Song{
		CreatorID:  mockUser.ID,
		Title:      "Foo",
		Subtitle:   "Bar",
		Key:        "A",
		Bpm:        120,
		ChordSheet: datatypes.JSON([]byte(`{"Verse" : "Foobar"`)),
	}

	mockUR := new(mocks.MockUserRepository)
	mockUR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockSong.CreatorID).
		Return(mockUser, nil)
	mockSR := new(mocks.MockSongRepository)

	ss := NewSongService(mockUR, mockSR)
	ctx := context.TODO()
	err := ss.Store(ctx, mockSong, mockUser)
	mockErr := domain.NewBadRequestErr("")
	assert.ErrorAs(t, err, &mockErr)
	assert.Empty(t, mockSong.ID)
	mockSR.AssertExpectations(t)
}

func TestUpdateSongCorrect(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	mockSong := &domain.Song{
		ID:         1,
		CreatorID:  mockUser.ID,
		Title:      "Foo",
		Subtitle:   "Bar",
		Key:        "A",
		Bpm:        120,
		ChordSheet: datatypes.JSON([]byte(`{"Verse" : "Foobar"}`)),
	}

	mockUR := new(mocks.MockUserRepository)
	mockUR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockSong.CreatorID).
		Return(mockUser, nil)
	mockSR := new(mocks.MockSongRepository)
	mockSR.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), mockSong).
		Return(nil)

	ss := NewSongService(mockUR, mockSR)
	ctx := context.TODO()
	err := ss.Update(ctx, mockSong, mockUser)

	assert.NoError(t, err)
	assert.NotEmpty(t, mockSong.ID)
	mockSR.AssertExpectations(t)
}

func TestUpdateSongNoClearanceNotCreator(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockSong := &domain.Song{
		ID:         1,
		CreatorID:  2,
		Title:      "Foo",
		Subtitle:   "Bar",
		Key:        "A",
		Bpm:        120,
		ChordSheet: datatypes.JSON([]byte(`{"Verse" : "Foobar"}`)),
	}
	mockUR := new(mocks.MockUserRepository)
	mockSR := new(mocks.MockSongRepository)
	mockSR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockSong.ID).
		Return(mockSong, nil)

	ss := NewSongService(mockUR, mockSR)
	ctx := context.TODO()
	err := ss.Update(ctx, mockSong, mockUser)

	mockErr := domain.NewNotAuthorizedErr("")
	assert.ErrorAs(t, err, &mockErr)
	mockSR.AssertExpectations(t)
}

func TestUpdateSongNoClearanceGetByIDErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockSong := &domain.Song{
		ID:         1,
		CreatorID:  2,
		Title:      "Foo",
		Subtitle:   "Bar",
		Key:        "A",
		Bpm:        120,
		ChordSheet: datatypes.JSON([]byte(`{"Verse" : "Foobar"}`)),
	}

	mockErr := domain.NewRecordNotFoundErr("", "")
	mockUR := new(mocks.MockUserRepository)
	mockSR := new(mocks.MockSongRepository)
	mockSR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockSong.ID).
		Return(nil, mockErr)

	ss := NewSongService(mockUR, mockSR)
	ctx := context.TODO()
	err := ss.Update(ctx, mockSong, mockUser)

	assert.ErrorAs(t, err, &mockErr)
	mockSR.AssertExpectations(t)
}

func TestSongUpdateInvalidKey(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	mockSong := &domain.Song{
		CreatorID:  mockUser.ID,
		Title:      "Foo",
		Subtitle:   "Bar",
		Key:        "Q",
		Bpm:        120,
		ChordSheet: datatypes.JSON([]byte(`{"Verse" : "Foobar"}`)),
	}

	mockSR := new(mocks.MockSongRepository)
	mockUR := new(mocks.MockUserRepository)

	ss := NewSongService(mockUR, mockSR)
	ctx := context.TODO()
	err := ss.Update(ctx, mockSong, mockUser)

	mockErr := domain.NewBadRequestErr("")
	assert.ErrorAs(t, err, &mockErr)
	assert.Empty(t, mockSong.ID)
	mockSR.AssertExpectations(t)
}

func TestUpdateSongCreatorNotExists(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	mockSong := &domain.Song{
		CreatorID:  mockUser.ID,
		Title:      "Foo",
		Subtitle:   "Bar",
		Key:        "A",
		Bpm:        120,
		ChordSheet: datatypes.JSON([]byte(`{"Verse" : "Foobar"}`)),
	}

	mockErr := domain.NewRecordNotFoundErr("", "")
	mockUR := new(mocks.MockUserRepository)
	mockUR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockSong.CreatorID).
		Return(nil, mockErr)
	mockSR := new(mocks.MockSongRepository)

	ss := NewSongService(mockUR, mockSR)
	ctx := context.TODO()
	err := ss.Update(ctx, mockSong, mockUser)

	assert.ErrorAs(t, err, &mockErr)
	assert.Empty(t, mockSong.ID)
	mockSR.AssertExpectations(t)
}

func TestUpdateSongInvalidChordsheet(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	mockSong := &domain.Song{
		ID:         1,
		CreatorID:  mockUser.ID,
		Title:      "Foo",
		Subtitle:   "Bar",
		Key:        "A",
		Bpm:        120,
		ChordSheet: datatypes.JSON([]byte(`{"Verse" : "Foobar}`)),
	}

	mockUR := new(mocks.MockUserRepository)
	mockUR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockSong.CreatorID).
		Return(mockUser, nil)
	mockSR := new(mocks.MockSongRepository)

	ss := NewSongService(mockUR, mockSR)
	ctx := context.TODO()
	err := ss.Update(ctx, mockSong, mockUser)
	mockErr := domain.NewBadRequestErr("")
	assert.ErrorAs(t, err, &mockErr)
	mockSR.AssertExpectations(t)
}

func TestRemoveSongCorrect(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockSong := &domain.Song{
		ID:         1,
		CreatorID:  1,
		Title:      "Foo",
		Subtitle:   "Bar",
		Key:        "A",
		Bpm:        120,
		ChordSheet: datatypes.JSON([]byte(`{"Verse" : "Foobar"}`)),
	}
	mockUR := new(mocks.MockUserRepository)
	mockSR := new(mocks.MockSongRepository)
	mockSR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockSong.ID).
		Return(mockSong, nil)
	mockSR.
		On("Delete", mock.AnythingOfType("*context.emptyCtx"), mockSong.ID).
		Return(nil)

	ss := NewSongService(mockUR, mockSR)
	ctx := context.TODO()
	err := ss.Remove(ctx, mockSong.ID, mockUser)

	assert.NoError(t, err)
	mockSR.AssertExpectations(t)
}

func TestRemoveSongNoClearanceNotCreator(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockSong := &domain.Song{
		ID:         1,
		CreatorID:  2,
		Title:      "Foo",
		Subtitle:   "Bar",
		Key:        "A",
		Bpm:        120,
		ChordSheet: datatypes.JSON([]byte(`{"Verse" : "Foobar"}`)),
	}
	mockUR := new(mocks.MockUserRepository)
	mockSR := new(mocks.MockSongRepository)
	mockSR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockSong.ID).
		Return(mockSong, nil)

	ss := NewSongService(mockUR, mockSR)
	ctx := context.TODO()
	err := ss.Remove(ctx, mockSong.ID, mockUser)

	mockErr := domain.NewNotAuthorizedErr("")
	assert.ErrorAs(t, err, &mockErr)
	mockSR.AssertExpectations(t)
}

func TestRemoveSongNoClearanceGetByIDErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockSongID := int64(1)

	mockErr := domain.NewRecordNotFoundErr("", "")
	mockUR := new(mocks.MockUserRepository)
	mockSR := new(mocks.MockSongRepository)
	mockSR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockSongID).
		Return(nil, mockErr)

	ss := NewSongService(mockUR, mockSR)
	ctx := context.TODO()
	err := ss.Remove(ctx, mockSongID, mockUser)

	assert.ErrorAs(t, err, &mockErr)
	mockSR.AssertExpectations(t)
}
