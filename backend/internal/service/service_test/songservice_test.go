package service_test

import (
	"context"
	"testing"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/96Asch/mkvstage-server/backend/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/datatypes"
)

func TestSongServiceStore(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	mockSong := &domain.Song{
		CreatorID:  mockUser.ID,
		Title:      "Foo",
		Subtitle:   "Bar",
		Key:        "A",
		BundleID:   1,
		Bpm:        120,
		ChordSheet: datatypes.JSON([]byte(`{"Verse" : "Foobar"}`)),
	}

	mockBundle := &domain.Bundle{
		ID:       1,
		Name:     "foobar",
		ParentID: 0,
	}

	t.Run("Correct", func(t *testing.T) {
		t.Parallel()

		mockUR := &mocks.MockUserRepository{}
		mockSR := &mocks.MockSongRepository{}
		mockBR := &mocks.MockBundleRepository{}
		mockSong := &domain.Song{
			CreatorID:  mockUser.ID,
			Title:      "Foo",
			Subtitle:   "Bar",
			Key:        "A",
			BundleID:   1,
			Bpm:        120,
			ChordSheet: datatypes.JSON([]byte(`{"Verse" : "Foobar"}`)),
		}
		mockBR.
			On("GetByID", context.TODO(), mockSong.BundleID).
			Return(mockBundle, nil)
		mockSR.
			On("Create", context.TODO(), mockSong).
			Return(nil).
			Run(func(args mock.Arguments) {
				arg, ok := args.Get(1).(*domain.Song)
				assert.True(t, ok)
				arg.ID = 1
			})

		ss := service.NewSongService(mockUR, mockSR, mockBR)
		ctx := context.TODO()

		err := ss.Store(ctx, mockSong, mockUser)
		assert.NoError(t, err)
		assert.NotEmpty(t, mockSong.ID)
		mockSR.AssertExpectations(t)
	})

	t.Run("Fail invalid permission", func(t *testing.T) {
		t.Parallel()

		mockErr := domain.NewNotAuthorizedErr("")
		mockUR := &mocks.MockUserRepository{}
		mockSR := &mocks.MockSongRepository{}
		mockBR := &mocks.MockBundleRepository{}
		mockUser := &domain.User{ID: 1, Permission: domain.GUEST}

		ss := service.NewSongService(mockUR, mockSR, mockBR)
		ctx := context.TODO()

		mockUser.Permission = domain.GUEST
		err := ss.Store(ctx, mockSong, mockUser)
		assert.ErrorAs(t, err, &mockErr)
		mockSR.AssertExpectations(t)
	})

	t.Run("Fail invalid key", func(t *testing.T) {
		t.Parallel()

		mockErr := domain.NewBadRequestErr("")
		mockUR := &mocks.MockUserRepository{}
		mockSR := &mocks.MockSongRepository{}
		mockBR := &mocks.MockBundleRepository{}
		mockSong := &domain.Song{Key: "R"}

		ss := service.NewSongService(mockUR, mockSR, mockBR)
		ctx := context.TODO()

		err := ss.Store(ctx, mockSong, mockUser)
		assert.ErrorAs(t, err, &mockErr)
		assert.Empty(t, mockSong.ID)
		mockSR.AssertExpectations(t)
	})

	t.Run("Fail invalid chordsheet", func(t *testing.T) {
		t.Parallel()

		mockErr := domain.NewBadRequestErr("")
		mockUR := &mocks.MockUserRepository{}
		mockSR := &mocks.MockSongRepository{}
		mockBR := &mocks.MockBundleRepository{}
		mockSong := &domain.Song{Key: "A", ChordSheet: datatypes.JSON([]byte(`{"`))}

		mockUR.
			On("GetByID", context.TODO(), mockSong.CreatorID).
			Return(mockUser, nil)

		ss := service.NewSongService(mockUR, mockSR, mockBR)
		ctx := context.TODO()

		err := ss.Store(ctx, mockSong, mockUser)
		assert.ErrorAs(t, err, &mockErr)
		assert.Empty(t, mockSong.ID)
		mockSR.AssertExpectations(t)
	})

	t.Run("Fail Bundle GetByID error", func(t *testing.T) {
		t.Parallel()

		mockErr := domain.NewRecordNotFoundErr("", "")
		mockUR := &mocks.MockUserRepository{}
		mockSR := &mocks.MockSongRepository{}
		mockBR := &mocks.MockBundleRepository{}

		mockBR.
			On("GetByID", context.TODO(), mockSong.BundleID).
			Return(nil, mockErr)

		ss := service.NewSongService(mockUR, mockSR, mockBR)
		ctx := context.TODO()

		err := ss.Store(ctx, mockSong, mockUser)
		assert.ErrorAs(t, err, &mockErr)
		assert.Empty(t, mockSong.ID)
		mockSR.AssertExpectations(t)
	})
}

func TestSongServiceUpdate(t *testing.T) {
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

	mockBundle := &domain.Bundle{
		ID:       1,
		Name:     "foobar",
		ParentID: 0,
	}

	t.Run("Correct", func(t *testing.T) {
		t.Parallel()

		mockUR := &mocks.MockUserRepository{}
		mockSR := &mocks.MockSongRepository{}
		mockBR := &mocks.MockBundleRepository{}

		mockBR.
			On("GetByID", context.TODO(), mockSong.BundleID).
			Return(mockBundle, nil)

		mockUR.
			On("GetByID", context.TODO(), mockSong.CreatorID).
			Return(mockUser, nil)
		mockSR.
			On("Update", context.TODO(), mockSong).
			Return(nil)

		ss := service.NewSongService(mockUR, mockSR, mockBR)
		ctx := context.TODO()

		err := ss.Update(ctx, mockSong, mockUser)
		assert.NoError(t, err)
		assert.NotEmpty(t, mockSong.ID)
		mockSR.AssertExpectations(t)
		mockUR.AssertExpectations(t)
		mockBR.AssertExpectations(t)
	})

	t.Run("Fail invalid permission not creator", func(t *testing.T) {
		t.Parallel()

		mockErr := domain.NewNotAuthorizedErr("")
		mockUR := &mocks.MockUserRepository{}
		mockSR := &mocks.MockSongRepository{}
		mockBR := &mocks.MockBundleRepository{}
		mockUser := &domain.User{ID: 2, Permission: domain.MEMBER}

		mockSR.
			On("GetByID", context.TODO(), mockSong.ID).
			Return(mockSong, nil)

		ss := service.NewSongService(mockUR, mockSR, mockBR)
		ctx := context.TODO()

		err := ss.Update(ctx, mockSong, mockUser)
		assert.ErrorAs(t, err, &mockErr)
		mockSR.AssertExpectations(t)
		mockUR.AssertExpectations(t)
		mockBR.AssertExpectations(t)
	})

	t.Run("Fail User GetByID error", func(t *testing.T) {
		t.Parallel()

		mockErr := domain.NewRecordNotFoundErr("", "")
		mockUR := &mocks.MockUserRepository{}
		mockSR := &mocks.MockSongRepository{}
		mockBR := &mocks.MockBundleRepository{}
		mockUser := &domain.User{ID: 1, Permission: domain.MEMBER}

		mockSR.
			On("GetByID", context.TODO(), mockSong.ID).
			Return(nil, mockErr)

		ss := service.NewSongService(mockUR, mockSR, mockBR)
		ctx := context.TODO()

		err := ss.Update(ctx, mockSong, mockUser)
		assert.ErrorAs(t, err, &mockErr)
		mockSR.AssertExpectations(t)
		mockUR.AssertExpectations(t)
		mockBR.AssertExpectations(t)
	})

	t.Run("Fail invalid key", func(t *testing.T) {
		t.Parallel()

		mockErr := domain.NewBadRequestErr("")
		mockUR := &mocks.MockUserRepository{}
		mockSR := &mocks.MockSongRepository{}
		mockBR := &mocks.MockBundleRepository{}
		mockSong := &domain.Song{Key: "W"}

		ss := service.NewSongService(mockUR, mockSR, mockBR)
		ctx := context.TODO()

		err := ss.Update(ctx, mockSong, mockUser)
		assert.ErrorAs(t, err, &mockErr)
		assert.Empty(t, mockSong.ID)
		mockSR.AssertExpectations(t)
		mockUR.AssertExpectations(t)
		mockBR.AssertExpectations(t)
	})

	t.Run("Fail User GetByID error", func(t *testing.T) {
		t.Parallel()

		mockErr := domain.NewRecordNotFoundErr("", "")
		mockUR := &mocks.MockUserRepository{}
		mockSR := &mocks.MockSongRepository{}
		mockBR := &mocks.MockBundleRepository{}

		mockBR.
			On("GetByID", context.TODO(), mockSong.BundleID).
			Return(mockBundle, nil)
		mockUR.
			On("GetByID", context.TODO(), mockSong.CreatorID).
			Return(nil, mockErr)

		ss := service.NewSongService(mockUR, mockSR, mockBR)
		ctx := context.TODO()

		err := ss.Update(ctx, mockSong, mockUser)
		assert.ErrorAs(t, err, &mockErr)
		mockSR.AssertExpectations(t)
		mockUR.AssertExpectations(t)
		mockBR.AssertExpectations(t)
	})

	t.Run("Fail invalid chordsheet", func(t *testing.T) {
		t.Parallel()

		mockUR := &mocks.MockUserRepository{}
		mockSR := &mocks.MockSongRepository{}
		mockBR := &mocks.MockBundleRepository{}
		mockSong := &domain.Song{Key: "A", ChordSheet: datatypes.JSON([]byte(``))}

		ss := service.NewSongService(mockUR, mockSR, mockBR)
		ctx := context.TODO()

		err := ss.Update(ctx, mockSong, mockUser)
		mockErr := domain.NewBadRequestErr("")
		assert.ErrorAs(t, err, &mockErr)
		mockSR.AssertExpectations(t)
		mockUR.AssertExpectations(t)
		mockBR.AssertExpectations(t)
	})
}

func TestSongServiceRemove(t *testing.T) {
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

	t.Run("Correct", func(t *testing.T) {
		t.Parallel()

		mockUR := &mocks.MockUserRepository{}
		mockSR := &mocks.MockSongRepository{}
		mockBR := &mocks.MockBundleRepository{}

		mockSR.
			On("GetByID", context.TODO(), mockSong.ID).
			Return(mockSong, nil)
		mockSR.
			On("Delete", context.TODO(), mockSong.ID).
			Return(nil)

		ss := service.NewSongService(mockUR, mockSR, mockBR)
		ctx := context.TODO()

		err := ss.Remove(ctx, mockSong.ID, mockUser)
		assert.NoError(t, err)
		mockSR.AssertExpectations(t)
	})

	t.Run("Fail invalid permission", func(t *testing.T) {
		t.Parallel()

		mockUR := &mocks.MockUserRepository{}
		mockSR := &mocks.MockSongRepository{}
		mockBR := &mocks.MockBundleRepository{}
		mockUser := &domain.User{ID: 2, Permission: domain.GUEST}

		mockSR.
			On("GetByID", context.TODO(), mockSong.ID).
			Return(mockSong, nil)

		ss := service.NewSongService(mockUR, mockSR, mockBR)
		ctx := context.TODO()

		err := ss.Remove(ctx, mockSong.ID, mockUser)
		mockErr := domain.NewNotAuthorizedErr("")
		assert.ErrorAs(t, err, &mockErr)
		mockSR.AssertExpectations(t)
	})

	t.Run("Fail Song GetByID error", func(t *testing.T) {
		t.Parallel()

		mockSongID := int64(1)
		mockErr := domain.NewRecordNotFoundErr("", "")
		mockUR := &mocks.MockUserRepository{}
		mockSR := &mocks.MockSongRepository{}
		mockBR := &mocks.MockBundleRepository{}
		mockUser := &domain.User{ID: 1, Permission: domain.GUEST}

		mockSR.
			On("GetByID", context.TODO(), mockSongID).
			Return(nil, mockErr)

		ss := service.NewSongService(mockUR, mockSR, mockBR)
		ctx := context.TODO()

		err := ss.Remove(ctx, mockSongID, mockUser)
		assert.ErrorAs(t, err, &mockErr)
		mockSR.AssertExpectations(t)
	})
}
