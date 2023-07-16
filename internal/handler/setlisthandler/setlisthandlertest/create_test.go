package setlisthandler_test

import (
	"bytes"
	"context"
	"encoding/json"
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

func prepareAndServeCreate(
	t *testing.T,
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
		http.MethodPost,
		"/setlists/create",
		requestBody,
	)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestCreateCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
	}

	expSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Foo",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().Add(time.Hour * 24).Round(0),
		UpdatedAt: time.Now().Round(0),
		Order:     datatypes.JSON([]byte(`{"order" : "1,2,3,4"}`)),
	}

	mockSetlist := &domain.Setlist{
		Name:      expSetlist.Name,
		CreatorID: expSetlist.ID,
		Deadline:  expSetlist.Deadline,
	}

	mockSetlistEntries := &[]domain.SetlistEntry{
		{
			SongID:      1,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
		{
			SongID:      2,
			Transpose:   1,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Verse 2"]}`)),
		},
	}

	expSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:          1,
			SongID:      1,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
		{
			ID:          2,
			SongID:      2,
			Transpose:   1,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Verse 2"]}`)),
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
		On("Store", mock.AnythingOfType("*context.emptyCtx"), mockSetlist, mockUser).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg, ok := args.Get(1).(*domain.Setlist)
			assert.True(t, ok)
			arg.ID = expSetlist.ID
			arg.UpdatedAt = expSetlist.UpdatedAt
			arg.Order = expSetlist.Order
		})

	mockSLES.
		On("StoreBatch", mock.AnythingOfType("*context.emptyCtx"), mockSetlistEntries, mockUser).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg, ok := args.Get(1).(*[]domain.SetlistEntry)
			assert.True(t, ok)

			for idx := range *arg {
				(*arg)[idx].ID = int64(idx + 1)
			}
		})

	byteBody, err := json.Marshal(gin.H{
		"name":            mockSetlist.Name,
		"creator_id":      mockSetlist.CreatorID,
		"deadline":        mockSetlist.Deadline,
		"created_entries": *mockSetlistEntries,
	})
	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockSL, mockSLES, mockSS, mockMWH, &byteBody)

	assert.Equal(t, http.StatusCreated, writer.Code)

	expBody, err := json.Marshal(gin.H{
		"setlist": expSetlist,
		"entries": expSetlistEntries,
	})
	assert.NoError(t, err)

	assert.Equal(t, expBody, writer.Body.Bytes())

	mockSL.AssertExpectations(t)
	mockSS.AssertExpectations(t)
	mockSLES.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestCreateBindErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
	}

	mockSetlist := &domain.Setlist{
		CreatorID: mockUser.ID,
	}

	mockSetlistEntries := &[]domain.SetlistEntry{
		{
			SongID:      1,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
		{
			SongID:      2,
			Transpose:   1,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Verse 2"]}`)),
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

	byteBody, err := json.Marshal(gin.H{
		"creator_id":      mockSetlist.CreatorID,
		"deadline":        mockSetlist.Deadline,
		"created_entries": *mockSetlistEntries,
	})

	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockSL, mockSLES, mockSS, mockMWH, &byteBody)

	assert.Equal(t, mockErr.Status(), writer.Code)
	mockSL.AssertExpectations(t)
	mockSS.AssertExpectations(t)
	mockSLES.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestCreateNoContext(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
	}

	mockSetlist := &domain.Setlist{
		Name:      "Foo",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().Add(time.Hour * 24),
	}

	mockSetlistEntries := &[]domain.SetlistEntry{
		{
			SongID:      1,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
		{
			SongID:      2,
			Transpose:   1,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Verse 2"]}`)),
		},
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

	byteBody, err := json.Marshal(gin.H{
		"name":            mockSetlist.Name,
		"creator_id":      mockSetlist.CreatorID,
		"deadline":        mockSetlist.Deadline,
		"created_entries": *mockSetlistEntries,
	})

	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockSL, mockSLES, mockSS, mockMWH, &byteBody)

	assert.Equal(t, mockErr.Status(), writer.Code)
	mockSS.AssertExpectations(t)
	mockSL.AssertExpectations(t)
	mockSLES.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestCreateSetlistStoreErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
	}

	mockSetlist := &domain.Setlist{
		Name:      "Foo",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().Add(time.Hour * 24).Round(0),
	}

	mockSetlistEntries := &[]domain.SetlistEntry{
		{
			SongID:      1,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
		{
			SongID:      2,
			Transpose:   1,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Verse 2"]}`)),
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
		On("Store", mock.AnythingOfType("*context.emptyCtx"), mockSetlist, mockUser).
		Return(mockErr)

	byteBody, err := json.Marshal(gin.H{
		"name":            mockSetlist.Name,
		"creator_id":      mockSetlist.CreatorID,
		"deadline":        mockSetlist.Deadline,
		"created_entries": *mockSetlistEntries,
	})
	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockSL, mockSLES, mockSS, mockMWH, &byteBody)

	assert.Equal(t, mockErr.Status(), writer.Code)
	mockSL.AssertExpectations(t)
	mockSS.AssertExpectations(t)
	mockSLES.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestCreateSetlistEntryStoreBatchErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
	}

	mockSetlist := &domain.Setlist{
		Name:      "Foo",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().Add(time.Hour * 24).Round(0),
	}

	mockSetlistEntries := &[]domain.SetlistEntry{
		{
			SongID:      1,
			Transpose:   0,
			Notes:       "",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Chorus 1"]}`)),
		},
		{
			SongID:      2,
			Transpose:   1,
			Notes:       "Foobar",
			Arrangement: datatypes.JSON([]byte(`{"arrangement":["Verse 1","Verse 2"]}`)),
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
		On("Store", mock.AnythingOfType("*context.emptyCtx"), mockSetlist, mockUser).
		Return(nil)

	mockSLES.
		On("StoreBatch", mock.AnythingOfType("*context.emptyCtx"), mockSetlistEntries, mockUser).
		Return(mockErr)

	byteBody, err := json.Marshal(gin.H{
		"name":            mockSetlist.Name,
		"creator_id":      mockSetlist.CreatorID,
		"deadline":        mockSetlist.Deadline,
		"created_entries": *mockSetlistEntries,
	})
	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockSL, mockSLES, mockSS, mockMWH, &byteBody)

	assert.Equal(t, mockErr.Status(), writer.Code)
	mockSS.AssertExpectations(t)
	mockSL.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
	mockSLES.AssertExpectations(t)
}
