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

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/handler/setlisthandler"
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
		Order:     datatypes.JSON([]byte(`{"order" : "2,1,3,4"}`)),
	}

	expSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Foo",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1).Truncate(time.Minute),
		UpdatedAt: time.Now().Round(0),
		Order:     datatypes.JSON([]byte(`{"order" : "1,2,3,4"}`)),
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
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
	}

	expCreatedSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          1,
			SongID:      1,
			Transpose:   0,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
	}

	mockUpdatedSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          2,
			SongID:      2,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
	}

	mockDeletedSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          3,
			SongID:      1,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
	}

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
		On("Update", mock.AnythingOfType("*context.emptyCtx"), expMockSetlist, mockUser).
		Return(expSetlist, nil)

	mockSLES.
		On("StoreBatch", mock.AnythingOfType("*context.emptyCtx"), mockCreatedSetlistEntries, mockUser).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg, ok := args.Get(1).(*[]domain.SetlistEntry)
			assert.True(t, ok)

			for idx := range *arg {
				(*arg)[idx].ID = int64(idx + 1)
			}
		})

	mockSLES.
		On("UpdateBatch", mock.AnythingOfType("*context.emptyCtx"), mockUpdatedSetlistEntries, mockUser).
		Return(nil)

	mockSLES.
		On("RemoveBatch", mock.AnythingOfType("*context.emptyCtx"), expMockSetlist, []int64{3}, mockUser).
		Return(nil)

	byteBody, err := json.Marshal(gin.H{
		"name":            mockSetlist.Name,
		"creator_id":      mockSetlist.CreatorID,
		"deadline":        mockSetlist.Deadline,
		"created_entries": *mockCreatedSetlistEntries,
		"updated_entries": *mockUpdatedSetlistEntries,
		"deleted_entries": *mockDeletedSetlistEntries,
	})

	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, fmt.Sprint(mockPrevSetlist.ID), mockSL, mockSLES, mockSS, mockMWH, &byteBody)

	assert.Equal(t, http.StatusOK, writer.Code)

	expBody, err := json.Marshal(gin.H{
		"setlist": expSetlist,
		"entries": append(*expCreatedSetlistEntries, *mockUpdatedSetlistEntries...),
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
		Order:     datatypes.JSON([]byte(`{"order" : "2,1,3,4"}`)),
	}

	expSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Foo",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1).Truncate(time.Minute),
		UpdatedAt: time.Now().Truncate(time.Minute),
		Order:     datatypes.JSON([]byte(`{"order" : "1,2,3,4"}`)),
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
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
	}

	expCreatedSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          1,
			SongID:      1,
			Transpose:   0,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
	}

	mockUpdatedSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          2,
			SongID:      2,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
	}

	mockDeletedSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          3,
			SongID:      1,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
	}

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
		On("Update", mock.AnythingOfType("*context.emptyCtx"), expMockSetlist, mockUser).
		Return(expSetlist, nil)

	mockSLES.
		On("StoreBatch", mock.AnythingOfType("*context.emptyCtx"), mockCreatedSetlistEntries, mockUser).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg, ok := args.Get(1).(*[]domain.SetlistEntry)
			assert.True(t, ok)

			for idx := range *arg {
				(*arg)[idx].ID = int64(idx + 1)
			}
		})

	mockSLES.
		On("UpdateBatch", mock.AnythingOfType("*context.emptyCtx"), mockUpdatedSetlistEntries, mockUser).
		Return(nil)

	mockSLES.
		On("RemoveBatch", mock.AnythingOfType("*context.emptyCtx"), expMockSetlist, []int64{3}, mockUser).
		Return(nil)

	byteBody, err := json.Marshal(gin.H{
		"name":            mockSetlist.Name,
		"creator_id":      mockSetlist.CreatorID,
		"deadline":        mockSetlist.Deadline.Format(time.RFC3339),
		"created_entries": *mockCreatedSetlistEntries,
		"updated_entries": *mockUpdatedSetlistEntries,
		"deleted_entries": *mockDeletedSetlistEntries,
	})

	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, fmt.Sprint(mockPrevSetlist.ID), mockSL, mockSLES, mockSS, mockMWH, &byteBody)

	assert.Equal(t, http.StatusOK, writer.Code)

	expBody, err := json.Marshal(gin.H{
		"setlist": expSetlist,
		"entries": append(*expCreatedSetlistEntries, *mockUpdatedSetlistEntries...),
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
		Order:     datatypes.JSON([]byte(`{"order" : "1,2,3,4"}`)),
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
		Order:     datatypes.JSON([]byte(`{"order" : "2,1,3,4"}`)),
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
		Order:     datatypes.JSON([]byte(`{"order" : "2,1,3,4"}`)),
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
		Order:     datatypes.JSON([]byte(`{"order" : "2,1,3,4"}`)),
	}

	expSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Foo",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1).Truncate(time.Minute),
		UpdatedAt: time.Now().Round(0),
		Order:     datatypes.JSON([]byte(`{"order" : "1,2,3,4"}`)),
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
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
	}

	mockUpdatedSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          2,
			SongID:      2,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
	}

	mockDeletedSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          3,
			SongID:      1,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
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

	mockSL.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), expMockSetlist, mockUser).
		Return(nil, mockErr)

	byteBody, err := json.Marshal(gin.H{
		"name":            mockSetlist.Name,
		"creator_id":      mockSetlist.CreatorID,
		"deadline":        mockSetlist.Deadline,
		"created_entries": *mockCreatedSetlistEntries,
		"updated_entries": *mockUpdatedSetlistEntries,
		"deleted_entries": *mockDeletedSetlistEntries,
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
		Order:     datatypes.JSON([]byte(`{"order" : "2,1,3,4"}`)),
	}

	expSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Foo",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1).Truncate(time.Minute),
		UpdatedAt: time.Now().Round(0),
		Order:     datatypes.JSON([]byte(`{"order" : "1,2,3,4"}`)),
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
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
	}

	mockUpdatedSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          2,
			SongID:      2,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
	}

	mockDeletedSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          3,
			SongID:      1,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
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

	mockSL.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), expMockSetlist, mockUser).
		Return(nil, nil)

	mockSLES.
		On("StoreBatch", mock.AnythingOfType("*context.emptyCtx"), mockCreatedSetlistEntries, mockUser).
		Return(mockErr)

	byteBody, err := json.Marshal(gin.H{
		"name":            mockSetlist.Name,
		"creator_id":      mockSetlist.CreatorID,
		"deadline":        mockSetlist.Deadline,
		"created_entries": *mockCreatedSetlistEntries,
		"updated_entries": *mockUpdatedSetlistEntries,
		"deleted_entries": *mockDeletedSetlistEntries,
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
		Order:     datatypes.JSON([]byte(`{"order" : "2,1,3,4"}`)),
	}

	expSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Foo",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1).Truncate(time.Minute),
		UpdatedAt: time.Now().Round(0),
		Order:     datatypes.JSON([]byte(`{"order" : "1,2,3,4"}`)),
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
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
	}

	mockUpdatedSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          2,
			SongID:      2,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
	}

	mockDeletedSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          3,
			SongID:      1,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
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

	mockSL.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), expMockSetlist, mockUser).
		Return(nil, nil)

	mockSLES.
		On("StoreBatch", mock.AnythingOfType("*context.emptyCtx"), mockCreatedSetlistEntries, mockUser).
		Return(nil)

	mockSLES.
		On("UpdateBatch", mock.AnythingOfType("*context.emptyCtx"), mockUpdatedSetlistEntries, mockUser).
		Return(mockErr)

	byteBody, err := json.Marshal(gin.H{
		"name":            mockSetlist.Name,
		"creator_id":      mockSetlist.CreatorID,
		"deadline":        mockSetlist.Deadline,
		"created_entries": *mockCreatedSetlistEntries,
		"updated_entries": *mockUpdatedSetlistEntries,
		"deleted_entries": *mockDeletedSetlistEntries,
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
		Order:     datatypes.JSON([]byte(`{"order" : "2,1,3,4"}`)),
	}

	expSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Foo",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().AddDate(0, 0, 1).Truncate(time.Minute),
		UpdatedAt: time.Now().Round(0),
		Order:     datatypes.JSON([]byte(`{"order" : "1,2,3,4"}`)),
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
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
	}

	mockUpdatedSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          2,
			SongID:      2,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
	}

	mockDeletedSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          3,
			SongID:      1,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
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

	mockSL.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), expMockSetlist, mockUser).
		Return(nil, nil)

	mockSLES.
		On("StoreBatch", mock.AnythingOfType("*context.emptyCtx"), mockCreatedSetlistEntries, mockUser).
		Return(nil)

	mockSLES.
		On("UpdateBatch", mock.AnythingOfType("*context.emptyCtx"), mockUpdatedSetlistEntries, mockUser).
		Return(nil)

	mockSLES.
		On("RemoveBatch", mock.AnythingOfType("*context.emptyCtx"), expMockSetlist, []int64{3}, mockUser).
		Return(mockErr)

	byteBody, err := json.Marshal(gin.H{
		"name":            mockSetlist.Name,
		"creator_id":      mockSetlist.CreatorID,
		"deadline":        mockSetlist.Deadline,
		"created_entries": *mockCreatedSetlistEntries,
		"updated_entries": *mockUpdatedSetlistEntries,
		"deleted_entries": *mockDeletedSetlistEntries,
	})
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, fmt.Sprint(mockPrevSetlist.ID), mockSL, mockSLES, mockSS, mockMWH, &byteBody)

	assert.Equal(t, mockErr.Status(), writer.Code)
	mockSL.AssertExpectations(t)
	mockSLES.AssertExpectations(t)
	mockSS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}
