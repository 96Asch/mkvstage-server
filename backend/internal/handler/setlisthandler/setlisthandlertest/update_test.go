package setlisthandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/96Asch/mkvstage-server/backend/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/backend/internal/handler/setlisthandler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/datatypes"
)

func prepareAndServeUpdate(
	t *testing.T,
	param string,
	mockSL domain.SetlistService,
	mockSLES domain.SetlistEntryService,
	mockSS domain.SongService,
	mockMWH domain.MiddlewareHandler,
	body *[]byte,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	setlisthandler.Initialize(&router.RouterGroup, mockSL, mockSLES, mockSS, mockMWH)

	requestBody := bytes.NewReader(*body)

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodPut,
		fmt.Sprintf("/setlists/%s", param),
		requestBody,
	)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestUpdateByIDCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
	}

	mockPrevSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Bar",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1).Truncate(time.Minute),
		UpdatedAt: time.Now().Round(0),
	}

	expSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Foo",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1).Truncate(time.Minute),
		UpdatedAt: time.Now().Round(0),
	}

	mockSetlist := &domain.Setlist{
		Name:      expSetlist.Name,
		CreatorID: expSetlist.ID,
		Deadline:  expSetlist.Deadline,
	}

	expMockSetlist := &domain.Setlist{
		ID:        mockPrevSetlist.ID,
		Name:      expSetlist.Name,
		CreatorID: expSetlist.ID,
		Deadline:  expSetlist.Deadline,
	}

	mockCreatedSetlistEntries := &[]domain.SetlistEntry{
		{
			SongID:      1,
			Transpose:   0,
			SetlistID:   expSetlist.ID,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`["Verse 1","Verse 2"]`)),
			Rank:        1000,
		},
	}

	mockUpdatedSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          2,
			SongID:      2,
			SetlistID:   expMockSetlist.ID,
			Transpose:   2,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`["Verse 1","Chorus 1"]`)),
			Rank:        2000,
		},
	}

	expFetchedEntries := &[]domain.SetlistEntry{
		{
			ID:          1,
			SongID:      1,
			SetlistID:   expMockSetlist.ID,
			Transpose:   0,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`["Verse 1","Verse 2"]`)),
			Rank:        1000,
		},
		{
			ID:          2,
			SongID:      2,
			SetlistID:   expMockSetlist.ID,
			Transpose:   2,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`["Verse 1","Chorus 1"]`)),
			Rank:        2000,
		},
	}

	mockDeletedSetlistEntries := []int64{3}

	mockSL := &mocks.MockSetlistService{}
	mockSLES := &mocks.MockSetlistEntryService{}
	mockSS := &mocks.MockSongService{}

	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	mockSL.
		On("Update", context.TODO(), expMockSetlist, mockUser).
		Return(expSetlist, nil)

	mockSLES.
		On("StoreBatch", context.TODO(), mockCreatedSetlistEntries, mockUser).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg, ok := args.Get(1).(*[]domain.SetlistEntry)
			assert.True(t, ok)

			for idx := range *arg {
				(*arg)[idx].ID = int64(idx + 1)
			}
		})

	mockSLES.
		On("UpdateBatch", context.TODO(), mockUpdatedSetlistEntries, mockUser).
		Return(nil)

	mockSLES.
		On("RemoveBatch", context.TODO(), expMockSetlist, []int64{3}, mockUser).
		Return(nil)

	mockSLES.
		On("FetchBySetlist", context.TODO(), &[]domain.Setlist{*expMockSetlist}).
		Return(expFetchedEntries, nil)

	byteBody, err := json.Marshal(gin.H{
		"name":       mockSetlist.Name,
		"creator_id": mockSetlist.CreatorID,
		"deadline":   mockSetlist.Deadline,
		"created_entries": []gin.H{
			{
				"song_id":     1,
				"transpose":   0,
				"notes":       "Foobar",
				"arrangement": []string{"Verse 1", "Verse 2"},
				"rank":        1000,
			},
		},
		"updated_entries": []gin.H{
			{
				"id":          2,
				"song_id":     2,
				"transpose":   2,
				"setlist_id":  expSetlist.ID,
				"notes":       "",
				"arrangement": []string{"Verse 1", "Chorus 1"},
				"rank":        2000,
			},
		},
		"deleted_entries": mockDeletedSetlistEntries,
	})

	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, fmt.Sprint(mockPrevSetlist.ID), mockSL, mockSLES, mockSS, mockMWH, &byteBody)

	assert.Equal(t, http.StatusOK, writer.Code)

	type setlistResponse struct {
		*domain.Setlist
		Entries []domain.SetlistEntry `json:"entries"`
	}

	expBody, err := json.Marshal(gin.H{
		"setlist": setlistResponse{expSetlist, *expFetchedEntries},
	})
	assert.NoError(t, err)

	assert.Equal(t, expBody, writer.Body.Bytes())
	mockSL.AssertExpectations(t)
	mockSLES.AssertExpectations(t)
	mockSS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDStringedDeadlineCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
	}

	mockPrevSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Bar",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1).Truncate(time.Minute),
		UpdatedAt: time.Now().Truncate(time.Minute),
	}

	expSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Foo",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1).Truncate(time.Minute),
		UpdatedAt: time.Now().Truncate(time.Minute),
	}

	mockSetlist := &domain.Setlist{
		Name:      expSetlist.Name,
		CreatorID: expSetlist.ID,
		Deadline:  expSetlist.Deadline,
	}

	expMockSetlist := &domain.Setlist{
		ID:        mockPrevSetlist.ID,
		Name:      expSetlist.Name,
		CreatorID: expSetlist.ID,
		Deadline:  expSetlist.Deadline,
	}

	mockCreatedSetlistEntries := &[]domain.SetlistEntry{
		{
			SongID:      1,
			SetlistID:   expSetlist.ID,
			Transpose:   0,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`["Verse 1","Chorus 1"]`)),
			Rank:        1000,
		},
	}

	expFetchedEntries := &[]domain.SetlistEntry{
		{
			ID:          1,
			SongID:      1,
			SetlistID:   expMockSetlist.ID,
			Transpose:   0,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`["Verse 1","Verse 2"]`)),
			Rank:        1000,
		},
		{
			ID:          2,
			SongID:      2,
			SetlistID:   expMockSetlist.ID,
			Transpose:   2,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`["Verse 1","Chorus 1"]`)),
			Rank:        2000,
		},
	}

	mockUpdatedSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          2,
			SongID:      2,
			SetlistID:   expMockSetlist.ID,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`["Verse 1","Chorus 1"]`)),
			Rank:        1000,
		},
	}

	mockDeletedSetlistEntries := []int64{3}
	mockSL := &mocks.MockSetlistService{}
	mockSLES := &mocks.MockSetlistEntryService{}
	mockSS := &mocks.MockSongService{}

	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	mockSL.
		On("Update", context.TODO(), expMockSetlist, mockUser).
		Return(expSetlist, nil)

	mockSLES.
		On("StoreBatch", context.TODO(), mockCreatedSetlistEntries, mockUser).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg, ok := args.Get(1).(*[]domain.SetlistEntry)
			assert.True(t, ok)

			for idx := range *arg {
				(*arg)[idx].ID = int64(idx + 1)
			}
		})

	mockSLES.
		On("UpdateBatch", context.TODO(), mockUpdatedSetlistEntries, mockUser).
		Return(nil)

	mockSLES.
		On("RemoveBatch", context.TODO(), expMockSetlist, []int64{3}, mockUser).
		Return(nil)

	mockSLES.
		On("FetchBySetlist", context.TODO(), &[]domain.Setlist{*expMockSetlist}).
		Return(expFetchedEntries, nil)

	byteBody, err := json.Marshal(gin.H{
		"name":            mockSetlist.Name,
		"creator_id":      mockSetlist.CreatorID,
		"deadline":        mockSetlist.Deadline.Format(time.RFC3339),
		"created_entries": *mockCreatedSetlistEntries,
		"updated_entries": *mockUpdatedSetlistEntries,
		"deleted_entries": mockDeletedSetlistEntries,
	})

	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, fmt.Sprint(mockPrevSetlist.ID), mockSL, mockSLES, mockSS, mockMWH, &byteBody)

	assert.Equal(t, http.StatusOK, writer.Code)

	type setlistResponse struct {
		*domain.Setlist
		Entries []domain.SetlistEntry `json:"entries"`
	}

	expBody, err := json.Marshal(gin.H{
		"setlist": setlistResponse{expSetlist, *expFetchedEntries},
	})
	assert.NoError(t, err)

	assert.Equal(t, expBody, writer.Body.Bytes())
	mockSL.AssertExpectations(t)
	mockSLES.AssertExpectations(t)
	mockSS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDInvalidParam(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	expSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Foo",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1).Truncate(time.Minute),
		UpdatedAt: time.Now().Round(0),
	}

	mockSetlist := &domain.Setlist{
		Name:      expSetlist.Name,
		CreatorID: expSetlist.ID,
		Deadline:  expSetlist.Deadline,
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockSL := &mocks.MockSetlistService{}
	mockSLES := &mocks.MockSetlistEntryService{}
	mockSS := &mocks.MockSongService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(mockSetlist)
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, "A", mockSL, mockSLES, mockSS, mockMWH, &byteBody)

	assert.Equal(t, http.StatusBadRequest, writer.Code)
	mockSL.AssertExpectations(t)
	mockSLES.AssertExpectations(t)
	mockSS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDBindErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
	}

	mockPrevSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Bar",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1).Truncate(time.Minute),
		UpdatedAt: time.Now().Round(0),
	}

	mockSetlist := &domain.Setlist{
		CreatorID: mockUser.ID,
	}

	mockErr := domain.NewBadRequestErr("")
	mockSL := &mocks.MockSetlistService{}
	mockSLES := &mocks.MockSetlistEntryService{}
	mockSS := &mocks.MockSongService{}

	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(mockSetlist)
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, fmt.Sprint(mockPrevSetlist.ID), mockSL, mockSLES, mockSS, mockMWH, &byteBody)

	assert.Equal(t, mockErr.Status(), writer.Code)
	mockSL.AssertExpectations(t)
	mockSLES.AssertExpectations(t)
	mockSS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDNoContext(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
	}

	mockPrevSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Bar",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1).Truncate(time.Minute),
		UpdatedAt: time.Now().Round(0),
	}

	mockSetlist := &domain.Setlist{
		Name:      "Foo",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1).Truncate(time.Minute),
	}

	mockErr := domain.NewInternalErr()
	mockSL := &mocks.MockSetlistService{}
	mockSLES := &mocks.MockSetlistEntryService{}
	mockSS := &mocks.MockSongService{}

	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(mockSetlist)
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, fmt.Sprint(mockPrevSetlist.ID), mockSL, mockSLES, mockSS, mockMWH, &byteBody)

	assert.Equal(t, mockErr.Status(), writer.Code)
	mockSL.AssertExpectations(t)
	mockSLES.AssertExpectations(t)
	mockSS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDStoreErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
	}

	mockPrevSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Bar",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1).Truncate(time.Minute),
		UpdatedAt: time.Now().Round(0),
	}

	expSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Foo",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1).Truncate(time.Minute),
		UpdatedAt: time.Now().Round(0),
	}

	mockSetlist := &domain.Setlist{
		Name:      expSetlist.Name,
		CreatorID: expSetlist.ID,
		Deadline:  expSetlist.Deadline,
	}

	expMockSetlist := &domain.Setlist{
		ID:        mockPrevSetlist.ID,
		Name:      expSetlist.Name,
		CreatorID: expSetlist.ID,
		Deadline:  expSetlist.Deadline,
	}

	mockCreatedSetlistEntries := &[]domain.SetlistEntry{
		{
			SongID:      1,
			SetlistID:   expSetlist.ID,
			Transpose:   0,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`["Verse 1","Chorus 1"]`)),
			Rank:        1000,
		},
	}

	mockUpdatedSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          2,
			SongID:      2,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`["Verse 1","Chorus 1"]`)),
			Rank:        1000,
		},
	}

	mockDeletedSetlistEntries := []int64{3}
	mockErr := domain.NewBadRequestErr("")
	mockSL := &mocks.MockSetlistService{}
	mockSLES := &mocks.MockSetlistEntryService{}
	mockSS := &mocks.MockSongService{}
	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	mockSL.
		On("Update", context.TODO(), expMockSetlist, mockUser).
		Return(nil, mockErr)

	byteBody, err := json.Marshal(gin.H{
		"name":            mockSetlist.Name,
		"creator_id":      mockSetlist.CreatorID,
		"deadline":        mockSetlist.Deadline,
		"created_entries": *mockCreatedSetlistEntries,
		"updated_entries": *mockUpdatedSetlistEntries,
		"deleted_entries": mockDeletedSetlistEntries,
	})
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, fmt.Sprint(mockPrevSetlist.ID), mockSL, mockSLES, mockSS, mockMWH, &byteBody)

	assert.Equal(t, mockErr.Status(), writer.Code)
	mockSL.AssertExpectations(t)
	mockSLES.AssertExpectations(t)
	mockSS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDSetlistEntryStoreBatchErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
	}

	mockPrevSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Bar",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1).Truncate(time.Minute),
		UpdatedAt: time.Now().Round(0),
	}

	expSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Foo",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1).Truncate(time.Minute),
		UpdatedAt: time.Now().Round(0),
	}

	mockSetlist := &domain.Setlist{
		Name:      expSetlist.Name,
		CreatorID: expSetlist.ID,
		Deadline:  expSetlist.Deadline,
	}

	expMockSetlist := &domain.Setlist{
		ID:        mockPrevSetlist.ID,
		Name:      expSetlist.Name,
		CreatorID: expSetlist.ID,
		Deadline:  expSetlist.Deadline,
	}

	mockCreatedSetlistEntries := &[]domain.SetlistEntry{
		{
			SongID:      1,
			SetlistID:   expMockSetlist.ID,
			Transpose:   0,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`["Verse 1","Chorus 1"]`)),
			Rank:        1000,
		},
	}

	mockUpdatedSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          2,
			SongID:      2,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`["Verse 1","Chorus 1"]`)),
			Rank:        2000,
		},
	}

	mockDeletedSetlistEntries := []int64{3}
	mockErr := domain.NewBadRequestErr("")
	mockSL := &mocks.MockSetlistService{}
	mockSLES := &mocks.MockSetlistEntryService{}
	mockSS := &mocks.MockSongService{}
	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	mockSL.
		On("Update", context.TODO(), expMockSetlist, mockUser).
		Return(expSetlist, nil)

	mockSLES.
		On("StoreBatch", context.TODO(), mockCreatedSetlistEntries, mockUser).
		Return(mockErr)

	byteBody, err := json.Marshal(gin.H{
		"name":            mockSetlist.Name,
		"creator_id":      mockSetlist.CreatorID,
		"deadline":        mockSetlist.Deadline,
		"created_entries": *mockCreatedSetlistEntries,
		"updated_entries": *mockUpdatedSetlistEntries,
		"deleted_entries": mockDeletedSetlistEntries,
	})
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, fmt.Sprint(mockPrevSetlist.ID), mockSL, mockSLES, mockSS, mockMWH, &byteBody)

	assert.Equal(t, mockErr.Status(), writer.Code)
	mockSL.AssertExpectations(t)
	mockSLES.AssertExpectations(t)
	mockSS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDSetlistEntryUpdateBatchErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
	}

	mockPrevSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Bar",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1).Truncate(time.Minute),
		UpdatedAt: time.Now().Round(0),
	}

	expSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Foo",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1).Truncate(time.Minute),
		UpdatedAt: time.Now().Round(0),
	}

	mockSetlist := &domain.Setlist{
		Name:      expSetlist.Name,
		CreatorID: expSetlist.ID,
		Deadline:  expSetlist.Deadline,
	}

	expMockSetlist := &domain.Setlist{
		ID:        mockPrevSetlist.ID,
		Name:      expSetlist.Name,
		CreatorID: expSetlist.ID,
		Deadline:  expSetlist.Deadline,
	}

	mockCreatedSetlistEntries := &[]domain.SetlistEntry{
		{
			SongID:      1,
			SetlistID:   expMockSetlist.ID,
			Transpose:   0,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`["Verse 1","Chorus 1"]`)),
			Rank:        1000,
		},
	}

	mockUpdatedSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          2,
			SongID:      2,
			SetlistID:   expMockSetlist.ID,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`["Verse 1","Chorus 1"]`)),
			Rank:        1000,
		},
	}

	mockDeletedSetlistEntries := []int64{3}
	mockErr := domain.NewBadRequestErr("")
	mockSL := &mocks.MockSetlistService{}
	mockSLES := &mocks.MockSetlistEntryService{}
	mockSS := &mocks.MockSongService{}
	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	mockSL.
		On("Update", context.TODO(), expMockSetlist, mockUser).
		Return(expSetlist, nil)

	mockSLES.
		On("StoreBatch", context.TODO(), mockCreatedSetlistEntries, mockUser).
		Return(nil)

	mockSLES.
		On("UpdateBatch", context.TODO(), mockUpdatedSetlistEntries, mockUser).
		Return(mockErr)

	byteBody, err := json.Marshal(gin.H{
		"name":            mockSetlist.Name,
		"creator_id":      mockSetlist.CreatorID,
		"deadline":        mockSetlist.Deadline,
		"created_entries": *mockCreatedSetlistEntries,
		"updated_entries": *mockUpdatedSetlistEntries,
		"deleted_entries": mockDeletedSetlistEntries,
	})
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, fmt.Sprint(mockPrevSetlist.ID), mockSL, mockSLES, mockSS, mockMWH, &byteBody)

	assert.Equal(t, mockErr.Status(), writer.Code)
	mockSL.AssertExpectations(t)
	mockSLES.AssertExpectations(t)
	mockSS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDSetlistEntryRemoveBatchErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
	}

	mockPrevSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Bar",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1).Truncate(time.Minute),
		UpdatedAt: time.Now().Round(0),
	}

	expSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Foo",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1).Truncate(time.Minute),
		UpdatedAt: time.Now().Round(0),
	}

	mockSetlist := &domain.Setlist{
		Name:      expSetlist.Name,
		CreatorID: expSetlist.ID,
		Deadline:  expSetlist.Deadline,
	}

	expMockSetlist := &domain.Setlist{
		ID:        mockPrevSetlist.ID,
		Name:      expSetlist.Name,
		CreatorID: expSetlist.ID,
		Deadline:  expSetlist.Deadline,
	}

	mockCreatedSetlistEntries := &[]domain.SetlistEntry{
		{
			SongID:      1,
			SetlistID:   expMockSetlist.ID,
			Transpose:   0,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`["Verse 1","Chorus 1"]`)),
			Rank:        1000,
		},
	}

	mockUpdatedSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          2,
			SongID:      2,
			SetlistID:   expSetlist.ID,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`["Verse 1","Chorus 1"]`)),
			Rank:        2000,
		},
	}

	mockDeletedSetlistEntries := []int64{3}
	mockErr := domain.NewBadRequestErr("")
	mockSL := &mocks.MockSetlistService{}
	mockSLES := &mocks.MockSetlistEntryService{}
	mockSS := &mocks.MockSongService{}
	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	mockSL.
		On("Update", context.TODO(), expMockSetlist, mockUser).
		Return(expSetlist, nil)

	mockSLES.
		On("StoreBatch", context.TODO(), mockCreatedSetlistEntries, mockUser).
		Return(nil)

	mockSLES.
		On("UpdateBatch", context.TODO(), mockUpdatedSetlistEntries, mockUser).
		Return(nil)

	mockSLES.
		On("RemoveBatch", context.TODO(), expMockSetlist, []int64{3}, mockUser).
		Return(mockErr)

	byteBody, err := json.Marshal(gin.H{
		"name":            mockSetlist.Name,
		"creator_id":      mockSetlist.CreatorID,
		"deadline":        mockSetlist.Deadline,
		"created_entries": *mockCreatedSetlistEntries,
		"updated_entries": *mockUpdatedSetlistEntries,
		"deleted_entries": mockDeletedSetlistEntries,
	})
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, fmt.Sprint(mockPrevSetlist.ID), mockSL, mockSLES, mockSS, mockMWH, &byteBody)

	assert.Equal(t, mockErr.Status(), writer.Code)
	mockSL.AssertExpectations(t)
	mockSLES.AssertExpectations(t)
	mockSS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDSetlistEntryFetchBySetlistErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
	}

	mockPrevSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Bar",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1).Truncate(time.Minute),
		UpdatedAt: time.Now().Round(0),
	}

	expSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Foo",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1).Truncate(time.Minute),
		UpdatedAt: time.Now().Round(0),
	}

	mockSetlist := &domain.Setlist{
		Name:      expSetlist.Name,
		CreatorID: expSetlist.ID,
		Deadline:  expSetlist.Deadline,
	}

	expMockSetlist := &domain.Setlist{
		ID:        mockPrevSetlist.ID,
		Name:      expSetlist.Name,
		CreatorID: expSetlist.ID,
		Deadline:  expSetlist.Deadline,
	}

	mockCreatedSetlistEntries := &[]domain.SetlistEntry{
		{
			SongID:      1,
			SetlistID:   expMockSetlist.ID,
			Transpose:   0,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`["Verse 1","Chorus 1"]`)),
			Rank:        1000,
		},
	}

	mockUpdatedSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          2,
			SongID:      2,
			SetlistID:   expSetlist.ID,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`["Verse 1","Chorus 1"]`)),
			Rank:        2000,
		},
	}

	mockDeletedSetlistEntries := []int64{3}
	mockErr := domain.NewBadRequestErr("")
	mockSL := &mocks.MockSetlistService{}
	mockSLES := &mocks.MockSetlistEntryService{}
	mockSS := &mocks.MockSongService{}
	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	mockSL.
		On("Update", context.TODO(), expMockSetlist, mockUser).
		Return(expSetlist, nil)

	mockSLES.
		On("StoreBatch", context.TODO(), mockCreatedSetlistEntries, mockUser).
		Return(nil)

	mockSLES.
		On("UpdateBatch", context.TODO(), mockUpdatedSetlistEntries, mockUser).
		Return(nil)

	mockSLES.
		On("RemoveBatch", context.TODO(), expMockSetlist, []int64{3}, mockUser).
		Return(nil)

	mockSLES.
		On("FetchBySetlist", context.TODO(), &[]domain.Setlist{*expMockSetlist}).
		Return(nil, mockErr)

	byteBody, err := json.Marshal(gin.H{
		"name":            mockSetlist.Name,
		"creator_id":      mockSetlist.CreatorID,
		"deadline":        mockSetlist.Deadline,
		"created_entries": *mockCreatedSetlistEntries,
		"updated_entries": *mockUpdatedSetlistEntries,
		"deleted_entries": mockDeletedSetlistEntries,
	})
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, fmt.Sprint(mockPrevSetlist.ID), mockSL, mockSLES, mockSS, mockMWH, &byteBody)

	expResponse, err := json.Marshal(gin.H{
		"error": mockErr.Error(),
	})

	assert.NoError(t, err)

	assert.Equal(t, mockErr.Status(), writer.Code)
	assert.Equal(t, expResponse, writer.Body.Bytes())
	mockSL.AssertExpectations(t)
	mockSLES.AssertExpectations(t)
	mockSS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}
