package service_test

import (
	"context"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/service"
	"github.com/96Asch/mkvstage-server/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/datatypes"
)

func TestSetlistEntryStoreBatchCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	mockSetlistEntries := &[]domain.SetlistEntry{
		{
			SongID:      1,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "C1"]}`)),
		},
		{
			SongID:      2,
			Transpose:   1,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "V2"]}`)),
		},
	}

	mockSER := &mocks.MockSetlistEntryRepository{}
	mockSLR := &mocks.MockSetlistRepository{}
	mockSR := &mocks.MockSongRepository{}

	for _, entry := range *mockSetlistEntries {
		mockSR.
			On("GetByID", mock.AnythingOfType("*context.emptyCtx"), entry.SongID).
			Return(nil, nil)

	}

	mockSER.
		On("CreateBatch", mock.AnythingOfType("*context.emptyCtx"), mockSetlistEntries).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg, ok := args.Get(1).(*[]domain.SetlistEntry)
			assert.True(t, ok)

			for idx := range *arg {
				(*arg)[idx].ID = int64(idx + 1)
			}
		})

	slr := service.NewSetlistEntryService(mockSER, mockSLR, mockSR)

	err := slr.StoreBatch(context.TODO(), mockSetlistEntries, mockUser)
	assert.NoError(t, err)

	for idx, entry := range *mockSetlistEntries {
		assert.Equal(t, int64(idx+1), entry.ID)
	}

	mockSER.AssertExpectations(t)
	mockSLR.AssertExpectations(t)
	mockSR.AssertExpectations(t)
}

func TestSetlistEntryStoreBatchClearanceErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          1,
			SongID:      1,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "C1"]}`)),
		},
		{
			ID:          2,
			SongID:      2,
			Transpose:   1,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "V2"]}`)),
		},
	}

	mockErr := domain.NewNotAuthorizedErr("")
	mockSER := &mocks.MockSetlistEntryRepository{}
	mockSLR := &mocks.MockSetlistRepository{}
	mockSR := &mocks.MockSongRepository{}

	slr := service.NewSetlistEntryService(mockSER, mockSLR, mockSR)

	err := slr.StoreBatch(context.TODO(), mockSetlistEntries, mockUser)
	assert.Error(t, err)
	assert.ErrorAs(t, err, &mockErr)
	mockSER.AssertExpectations(t)
	mockSLR.AssertExpectations(t)
	mockSR.AssertExpectations(t)
}

func TestSetlistEntryStoreBatchSongRepoErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	mockSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          1,
			SongID:      1,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "C1"]}`)),
		},
		{
			ID:          2,
			SongID:      2,
			Transpose:   1,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "V2"]}`)),
		},
	}

	mockErr := domain.NewBadRequestErr("")
	mockSER := &mocks.MockSetlistEntryRepository{}
	mockSLR := &mocks.MockSetlistRepository{}
	mockSR := &mocks.MockSongRepository{}

	mockSR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), (*mockSetlistEntries)[0].SongID).
		Return(nil, mockErr)

	slr := service.NewSetlistEntryService(mockSER, mockSLR, mockSR)

	err := slr.StoreBatch(context.TODO(), mockSetlistEntries, mockUser)
	assert.Error(t, err)
	assert.ErrorAs(t, err, &mockErr)
	mockSER.AssertExpectations(t)
	mockSLR.AssertExpectations(t)
	mockSR.AssertExpectations(t)
}

func TestSetlistEntryStoreBatchInvalidTranspose(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	mockSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          1,
			SongID:      1,
			Transpose:   util.TransposeMax + 1,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "C1"]}`)),
		},
		{
			ID:          2,
			SongID:      2,
			Transpose:   1,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "V2"]}`)),
		},
	}

	mockErr := domain.NewBadRequestErr("")
	mockSER := &mocks.MockSetlistEntryRepository{}
	mockSLR := &mocks.MockSetlistRepository{}
	mockSR := &mocks.MockSongRepository{}

	slr := service.NewSetlistEntryService(mockSER, mockSLR, mockSR)

	err := slr.StoreBatch(context.TODO(), mockSetlistEntries, mockUser)
	assert.Error(t, err)
	assert.ErrorAs(t, err, &mockErr)
	mockSER.AssertExpectations(t)
	mockSLR.AssertExpectations(t)
	mockSR.AssertExpectations(t)
}

func TestSetlistEntryStoreBatchSetlistEntryRepoErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	mockSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          1,
			SongID:      1,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "C1"]}`)),
		},
		{
			ID:          2,
			SongID:      2,
			Transpose:   1,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "V2"]}`)),
		},
	}

	mockErr := domain.NewInternalErr()
	mockSER := &mocks.MockSetlistEntryRepository{}
	mockSLR := &mocks.MockSetlistRepository{}
	mockSR := &mocks.MockSongRepository{}

	mockSR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), (*mockSetlistEntries)[0].SongID).
		Return(nil, nil)
	mockSR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), (*mockSetlistEntries)[1].SongID).
		Return(nil, nil)
	mockSER.
		On("CreateBatch", mock.AnythingOfType("*context.emptyCtx"), mockSetlistEntries).
		Return(mockErr)

	slr := service.NewSetlistEntryService(mockSER, mockSLR, mockSR)

	err := slr.StoreBatch(context.TODO(), mockSetlistEntries, mockUser)
	assert.Error(t, err)
	assert.ErrorAs(t, err, &mockErr)
	mockSER.AssertExpectations(t)
	mockSLR.AssertExpectations(t)
	mockSR.AssertExpectations(t)
}

func TestSetlistEntryFetchByIDCorrect(t *testing.T) {
	t.Parallel()

	slid := int64(1)
	mockSetlistEntry := &domain.SetlistEntry{
		ID:          1,
		SongID:      1,
		Transpose:   0,
		Notes:       "",
		Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "C1"]}`)),
	}

	mockSER := &mocks.MockSetlistEntryRepository{}
	mockSLR := &mocks.MockSetlistRepository{}
	mockSR := &mocks.MockSongRepository{}

	mockSER.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), slid).
		Return(mockSetlistEntry, nil)

	slr := service.NewSetlistEntryService(mockSER, mockSLR, mockSR)

	setlistEntry, err := slr.FetchByID(context.TODO(), slid)
	assert.NoError(t, err)
	assert.Equal(t, mockSetlistEntry, setlistEntry)

	mockSER.AssertExpectations(t)
	mockSLR.AssertExpectations(t)
	mockSR.AssertExpectations(t)
}

func TestSetlistEntryFetchByIDErr(t *testing.T) {
	t.Parallel()

	slid := int64(1)

	expErr := domain.NewInternalErr()
	mockSER := &mocks.MockSetlistEntryRepository{}
	mockSLR := &mocks.MockSetlistRepository{}
	mockSR := &mocks.MockSongRepository{}

	mockSER.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), slid).
		Return(nil, expErr)

	slr := service.NewSetlistEntryService(mockSER, mockSLR, mockSR)

	setlist, err := slr.FetchByID(context.TODO(), slid)
	assert.ErrorAs(t, err, &expErr)
	assert.Nil(t, setlist)

	mockSER.AssertExpectations(t)
	mockSLR.AssertExpectations(t)
	mockSR.AssertExpectations(t)
}

func TestSetlistEntryFetchAllCorrect(t *testing.T) {
	t.Parallel()

	mockSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          1,
			SongID:      1,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "C1"]}`)),
		},
		{
			ID:          2,
			SongID:      2,
			Transpose:   1,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "V2"]}`)),
		},
	}

	mockSER := &mocks.MockSetlistEntryRepository{}
	mockSLR := &mocks.MockSetlistRepository{}
	mockSR := &mocks.MockSongRepository{}

	mockSER.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(mockSetlistEntries, nil)

	slr := service.NewSetlistEntryService(mockSER, mockSLR, mockSR)

	setlistEntries, err := slr.FetchAll(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, mockSetlistEntries, setlistEntries)
	mockSER.AssertExpectations(t)
	mockSLR.AssertExpectations(t)
	mockSR.AssertExpectations(t)
}

func TestSetlistEntryFetchAllErr(t *testing.T) {
	t.Parallel()

	expErr := domain.NewInternalErr()
	mockSER := &mocks.MockSetlistEntryRepository{}
	mockSLR := &mocks.MockSetlistRepository{}
	mockSR := &mocks.MockSongRepository{}

	mockSER.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(nil, expErr)

	slr := service.NewSetlistEntryService(mockSER, mockSLR, mockSR)

	setlist, err := slr.FetchAll(context.TODO())
	assert.ErrorAs(t, err, &expErr)
	assert.Nil(t, setlist)
	mockSER.AssertExpectations(t)
	mockSLR.AssertExpectations(t)
	mockSR.AssertExpectations(t)
}

func TestSetlistEntryUpdateBatchCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	mockSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          1,
			SongID:      1,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "C1"]}`)),
		},
		{
			ID:          2,
			SongID:      2,
			Transpose:   1,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "V2"]}`)),
		},
	}

	mockSER := &mocks.MockSetlistEntryRepository{}
	mockSLR := &mocks.MockSetlistRepository{}
	mockSR := &mocks.MockSongRepository{}

	for _, entry := range *mockSetlistEntries {
		mockSR.
			On("GetByID", mock.AnythingOfType("*context.emptyCtx"), entry.SongID).
			Return(nil, nil)
		mockSER.
			On("GetByID", mock.AnythingOfType("*context.emptyCtx"), entry.ID).
			Return(nil, nil)
	}

	mockSER.
		On("UpdateBatch", mock.AnythingOfType("*context.emptyCtx"), mockSetlistEntries).
		Return(nil)

	slr := service.NewSetlistEntryService(mockSER, mockSLR, mockSR)

	err := slr.UpdateBatch(context.TODO(), mockSetlistEntries, mockUser)
	assert.NoError(t, err)
	mockSER.AssertExpectations(t)
	mockSLR.AssertExpectations(t)
	mockSR.AssertExpectations(t)
}

func TestSetlistEntryUpdateBatchClearanceErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          1,
			SongID:      1,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "C1"]}`)),
		},
		{
			ID:          2,
			SongID:      2,
			Transpose:   1,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "V2"]}`)),
		},
	}

	mockErr := domain.NewNotAuthorizedErr("")
	mockSER := &mocks.MockSetlistEntryRepository{}
	mockSLR := &mocks.MockSetlistRepository{}
	mockSR := &mocks.MockSongRepository{}

	slr := service.NewSetlistEntryService(mockSER, mockSLR, mockSR)

	err := slr.UpdateBatch(context.TODO(), mockSetlistEntries, mockUser)
	assert.Error(t, err)
	assert.ErrorAs(t, err, &mockErr)
	mockSER.AssertExpectations(t)
	mockSLR.AssertExpectations(t)
	mockSR.AssertExpectations(t)
}

func TestSetlistEntryUpdateBatchSongRepoErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	mockSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          1,
			SongID:      1,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "C1"]}`)),
		},
		{
			ID:          2,
			SongID:      2,
			Transpose:   1,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "V2"]}`)),
		},
	}

	mockErr := domain.NewBadRequestErr("")
	mockSER := &mocks.MockSetlistEntryRepository{}
	mockSR := &mocks.MockSongRepository{}
	mockSLR := &mocks.MockSetlistRepository{}

	mockSR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), (*mockSetlistEntries)[0].SongID).
		Return(nil, mockErr)

	slr := service.NewSetlistEntryService(mockSER, mockSLR, mockSR)

	err := slr.UpdateBatch(context.TODO(), mockSetlistEntries, mockUser)
	assert.Error(t, err)
	assert.ErrorAs(t, err, &mockErr)
	mockSR.AssertExpectations(t)
	mockSLR.AssertExpectations(t)
	mockSER.AssertExpectations(t)
}

func TestSetlistEntryUpdateBatchInvalidTranspose(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	mockSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          1,
			SongID:      1,
			Transpose:   util.TransposeMax + 1,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "C1"]}`)),
		},
		{
			ID:          2,
			SongID:      2,
			Transpose:   1,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "V2"]}`)),
		},
	}

	mockErr := domain.NewBadRequestErr("")
	mockSER := &mocks.MockSetlistEntryRepository{}
	mockSR := &mocks.MockSongRepository{}
	mockSLR := &mocks.MockSetlistRepository{}

	slr := service.NewSetlistEntryService(mockSER, mockSLR, mockSR)

	err := slr.UpdateBatch(context.TODO(), mockSetlistEntries, mockUser)
	assert.Error(t, err)
	assert.ErrorAs(t, err, &mockErr)
	mockSR.AssertExpectations(t)
	mockSLR.AssertExpectations(t)
	mockSER.AssertExpectations(t)
}

func TestSetlistEntryUpdateBatchInvalidID(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	mockSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          1,
			SongID:      1,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "C1"]}`)),
		},
		{
			ID:          2,
			SongID:      2,
			Transpose:   1,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "V2"]}`)),
		},
	}

	mockErr := domain.NewBadRequestErr("")
	mockSER := &mocks.MockSetlistEntryRepository{}
	mockSLR := &mocks.MockSetlistRepository{}

	mockSR := &mocks.MockSongRepository{}

	mockSR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), (*mockSetlistEntries)[0].SongID).
		Return(nil, nil)

	mockSER.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), (*mockSetlistEntries)[0].ID).
		Return(nil, mockErr)

	slr := service.NewSetlistEntryService(mockSER, mockSLR, mockSR)

	err := slr.UpdateBatch(context.TODO(), mockSetlistEntries, mockUser)
	assert.Error(t, err)
	assert.ErrorAs(t, err, &mockErr)
	mockSR.AssertExpectations(t)
	mockSLR.AssertExpectations(t)
	mockSER.AssertExpectations(t)
}

func TestSetlistEntryUpdateBatchSetlistEntryRepoErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	mockSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          1,
			SongID:      1,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "C1"]}`)),
		},
		{
			ID:          2,
			SongID:      2,
			Transpose:   1,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{arrangement: ["V1", "V2"]}`)),
		},
	}

	mockErr := domain.NewInternalErr()
	mockSER := &mocks.MockSetlistEntryRepository{}
	mockSR := &mocks.MockSongRepository{}
	mockSLR := &mocks.MockSetlistRepository{}

	for _, entry := range *mockSetlistEntries {
		mockSR.
			On("GetByID", mock.AnythingOfType("*context.emptyCtx"), entry.SongID).
			Return(nil, nil)
		mockSER.
			On("GetByID", mock.AnythingOfType("*context.emptyCtx"), entry.ID).
			Return(nil, nil)
	}

	mockSER.
		On("UpdateBatch", mock.AnythingOfType("*context.emptyCtx"), mockSetlistEntries).
		Return(mockErr)

	slr := service.NewSetlistEntryService(mockSER, mockSLR, mockSR)

	err := slr.UpdateBatch(context.TODO(), mockSetlistEntries, mockUser)
	assert.Error(t, err)
	assert.ErrorAs(t, err, &mockErr)
	mockSR.AssertExpectations(t)
	mockSLR.AssertExpectations(t)
	mockSER.AssertExpectations(t)
}

func TestSetlistEntryRemoveBatchCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	mockSetlistEntryIds := []int64{
		1,
		2,
	}

	mockSER := &mocks.MockSetlistEntryRepository{}
	mockSR := &mocks.MockSongRepository{}
	mockSLR := &mocks.MockSetlistRepository{}

	for _, id := range mockSetlistEntryIds {
		mockSER.
			On("GetByID", mock.AnythingOfType("*context.emptyCtx"), id).
			Return(nil, nil)
	}

	mockSER.
		On("DeleteBatch", mock.AnythingOfType("*context.emptyCtx"), mockSetlistEntryIds).
		Return(nil)

	slr := service.NewSetlistEntryService(mockSER, mockSLR, mockSR)

	err := slr.RemoveBatch(context.TODO(), mockSetlistEntryIds, mockUser)
	assert.NoError(t, err)
	mockSER.AssertExpectations(t)
	mockSLR.AssertExpectations(t)
	mockSR.AssertExpectations(t)
}

func TestSetlistEntryRemoveBatchClearanceErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockSetlistEntryIds := []int64{
		1,
		2,
	}

	mockErr := domain.NewNotAuthorizedErr("")
	mockSER := &mocks.MockSetlistEntryRepository{}
	mockSLR := &mocks.MockSetlistRepository{}
	mockSR := &mocks.MockSongRepository{}

	slr := service.NewSetlistEntryService(mockSER, mockSLR, mockSR)

	err := slr.RemoveBatch(context.TODO(), mockSetlistEntryIds, mockUser)
	assert.Error(t, err)
	assert.ErrorAs(t, err, &mockErr)
	mockSER.AssertExpectations(t)
	mockSLR.AssertExpectations(t)
	mockSR.AssertExpectations(t)
}

func TestSetlistEntryRemoveBatchInvalidID(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	mockSetlistEntryIds := []int64{
		1,
		2,
	}

	mockErr := domain.NewRecordNotFoundErr("", "")
	mockSER := &mocks.MockSetlistEntryRepository{}
	mockSLR := &mocks.MockSetlistRepository{}
	mockSR := &mocks.MockSongRepository{}

	mockSER.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockSetlistEntryIds[0]).
		Return(nil, mockErr)

	slr := service.NewSetlistEntryService(mockSER, mockSLR, mockSR)

	err := slr.RemoveBatch(context.TODO(), mockSetlistEntryIds, mockUser)
	assert.ErrorAs(t, err, &mockErr)
	mockSER.AssertExpectations(t)
	mockSLR.AssertExpectations(t)
	mockSR.AssertExpectations(t)
}
