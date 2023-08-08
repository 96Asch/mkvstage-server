package setlisthandler_test

import (
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
)

func prepareAndServeGet(
	t *testing.T,
	path string,
	mockSL domain.SetlistService,
	mockSLES domain.SetlistEntryService,
	mockSS domain.SongService,
	mockMWH domain.MiddlewareHandler,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	setlisthandler.Initialize(&router.RouterGroup, mockSL, mockSLES, mockSS, mockMWH)

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodGet,
		fmt.Sprintf("/setlists%s", path),
		nil,
	)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestGetAll(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	expSetlist := &[]domain.Setlist{
		{
			ID:        1,
			Name:      "Foobar",
			Deadline:  time.Now().UTC().AddDate(0, 0, 1).Truncate(time.Minute),
			CreatorID: mockUser.ID,
		},
		{
			ID:        2,
			Name:      "Barfoo",
			Deadline:  time.Now().UTC().AddDate(0, 0, 3).Truncate(time.Minute),
			CreatorID: 0,
		},
	}

	expSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:        1,
			SetlistID: 1,
			Rank:      1000,
		},
		{
			ID:        2,
			SetlistID: 1,
			Rank:      2000,
		},
		{
			ID:        3,
			SetlistID: 2,
			Rank:      1000,
		},
	}

	type setlistResponse struct {
		*domain.Setlist
		Entries []domain.SetlistEntry `json:"entries"`
	}

	t.Run("Correct Get All Setlists", func(t *testing.T) {
		t.Parallel()

		mockSL := &mocks.MockSetlistService{}
		mockSLES := &mocks.MockSetlistEntryService{}
		mockSS := &mocks.MockSongService{}
		mockMWH := &mocks.MockMiddlewareHandler{}

		var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {}

		mockMWH.
			On("AuthenticateUser").
			Return(mockAuthHF)

		mockSL.
			On("Fetch", mock.AnythingOfType("*context.emptyCtx"), time.Time{}, time.Time{}).
			Return(expSetlist, nil)

		mockSLES.
			On("FetchBySetlist", mock.AnythingOfType("*context.emptyCtx"), expSetlist).
			Return(expSetlistEntries, nil)

		writer := prepareAndServeGet(t, "", mockSL, mockSLES, mockSS, mockMWH)

		response := []setlistResponse{
			{
				&(*expSetlist)[0],
				[]domain.SetlistEntry{
					(*expSetlistEntries)[0],
					(*expSetlistEntries)[1],
				},
			},
			{
				&(*expSetlist)[1],
				[]domain.SetlistEntry{
					(*expSetlistEntries)[2],
				},
			},
		}

		expBody, err := json.Marshal(gin.H{
			"setlists": response,
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, writer.Code)
		assert.Equal(t, expBody, writer.Body.Bytes())
		mockSL.AssertExpectations(t)
		mockSLES.AssertExpectations(t)
		mockSS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})

	t.Run("Correct Get Setlists with timeframes", func(t *testing.T) {
		t.Parallel()

		mockSL := &mocks.MockSetlistService{}
		mockSLES := &mocks.MockSetlistEntryService{}
		mockSS := &mocks.MockSongService{}
		mockMWH := &mocks.MockMiddlewareHandler{}

		var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {}

		mockMWH.
			On("AuthenticateUser").
			Return(mockAuthHF)

		fromTime := (*expSetlist)[0].Deadline.Add(-24 * time.Hour)
		toTime := (*expSetlist)[0].Deadline.Add(-24 * time.Hour)

		mockSL.
			On("Fetch", mock.AnythingOfType("*context.emptyCtx"), fromTime, toTime).
			Return(expSetlist, nil)

		mockSLES.
			On("FetchBySetlist", mock.AnythingOfType("*context.emptyCtx"), expSetlist).
			Return(expSetlistEntries, nil)

		fmt.Println(fromTime.Format(time.RFC3339))
		writer := prepareAndServeGet(t,
			fmt.Sprintf("?from=%s&to=%s", fromTime.Format(time.RFC3339), toTime.Format(time.RFC3339)),
			mockSL, mockSLES, mockSS, mockMWH)

		expBody, err := json.Marshal(gin.H{
			"setlist": expSetlist,
			"entries": expSetlistEntries,
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, writer.Code)
		assert.Equal(t, expBody, writer.Body.Bytes())
		mockSL.AssertExpectations(t)
		mockSLES.AssertExpectations(t)
		mockSS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})

	t.Run("Fail Invalid From Field", func(t *testing.T) {
		t.Parallel()

		mockSL := &mocks.MockSetlistService{}
		mockSLES := &mocks.MockSetlistEntryService{}
		mockSS := &mocks.MockSongService{}
		mockMWH := &mocks.MockMiddlewareHandler{}

		var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {}

		mockMWH.
			On("AuthenticateUser").
			Return(mockAuthHF)

		fromTimeString := (*expSetlist)[0].Deadline.Add(-24 * time.Hour).Format(time.RFC1123)
		toTimeString := (*expSetlist)[0].Deadline.Format(time.RFC3339)
		writer := prepareAndServeGet(t, fmt.Sprintf("?from=%s&to=%s", fromTimeString, toTimeString), mockSL, mockSLES, mockSS, mockMWH)

		expBody, err := json.Marshal(gin.H{
			"error": fmt.Sprintf("%s must be in RFC3339 format", fromTimeString),
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
		assert.Equal(t, expBody, writer.Body.Bytes())
		mockSL.AssertExpectations(t)
		mockSLES.AssertExpectations(t)
		mockSS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})

	t.Run("Fail Invalid To Field", func(t *testing.T) {
		t.Parallel()

		mockSL := &mocks.MockSetlistService{}
		mockSLES := &mocks.MockSetlistEntryService{}
		mockSS := &mocks.MockSongService{}
		mockMWH := &mocks.MockMiddlewareHandler{}

		var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {}

		mockMWH.
			On("AuthenticateUser").
			Return(mockAuthHF)

		fromTimeString := (*expSetlist)[0].Deadline.Add(-24 * time.Hour).Format(time.RFC3339)
		toTimeString := (*expSetlist)[0].Deadline.Format(time.RFC1123)
		writer := prepareAndServeGet(t, fmt.Sprintf("?from=%s&to=%s", fromTimeString, toTimeString), mockSL, mockSLES, mockSS, mockMWH)

		expBody, err := json.Marshal(gin.H{
			"error": fmt.Sprintf("%s must be in RFC3339 format", toTimeString),
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
		assert.Equal(t, expBody, writer.Body.Bytes())
		mockSL.AssertExpectations(t)
		mockSLES.AssertExpectations(t)
		mockSS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})

	t.Run("Fail Fetch Setlists", func(t *testing.T) {
		t.Parallel()

		mockErr := domain.NewInternalErr()
		mockSL := &mocks.MockSetlistService{}
		mockSLES := &mocks.MockSetlistEntryService{}
		mockSS := &mocks.MockSongService{}
		mockMWH := &mocks.MockMiddlewareHandler{}

		var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {}

		mockMWH.
			On("AuthenticateUser").
			Return(mockAuthHF)

		mockSL.
			On("Fetch", mock.AnythingOfType("*context.emptyCtx"), time.Time{}, time.Time{}).
			Return(nil, mockErr)

		writer := prepareAndServeGet(t, "", mockSL, mockSLES, mockSS, mockMWH)

		expBody, err := json.Marshal(gin.H{
			"error": mockErr.Error(),
		})
		assert.NoError(t, err)

		assert.Equal(t, domain.Status(mockErr), writer.Code)
		assert.Equal(t, expBody, writer.Body.Bytes())
		mockSL.AssertExpectations(t)
		mockSLES.AssertExpectations(t)
		mockSS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})

	t.Run("Fail Fetch SetlistEntries", func(t *testing.T) {
		t.Parallel()

		mockErr := domain.NewInternalErr()
		mockSL := &mocks.MockSetlistService{}
		mockSLES := &mocks.MockSetlistEntryService{}
		mockSS := &mocks.MockSongService{}
		mockMWH := &mocks.MockMiddlewareHandler{}

		var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {}

		mockMWH.
			On("AuthenticateUser").
			Return(mockAuthHF)

		mockSL.
			On("Fetch", mock.AnythingOfType("*context.emptyCtx"), time.Time{}, time.Time{}).
			Return(expSetlist, nil)

		mockSLES.
			On("FetchBySetlist", mock.AnythingOfType("*context.emptyCtx"), expSetlist).
			Return(nil, mockErr)

		writer := prepareAndServeGet(t, "", mockSL, mockSLES, mockSS, mockMWH)

		expBody, err := json.Marshal(gin.H{
			"error": mockErr.Error(),
		})
		assert.NoError(t, err)

		assert.Equal(t, domain.Status(mockErr), writer.Code)
		assert.Equal(t, expBody, writer.Body.Bytes())
		mockSL.AssertExpectations(t)
		mockSLES.AssertExpectations(t)
		mockSS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.EDITOR,
	}

	expSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Foobar",
		Deadline:  time.Now().UTC().AddDate(0, 0, 1).Truncate(time.Minute),
		CreatorID: mockUser.ID,
	}

	expSetlistEntries := &[]domain.SetlistEntry{
		{
			ID:        1,
			SetlistID: 1,
		},
		{
			ID:        2,
			SetlistID: 1,
		},
	}

	t.Run("Correct Get Setlist By ID", func(t *testing.T) {
		t.Parallel()

		mockSL := &mocks.MockSetlistService{}
		mockSLES := &mocks.MockSetlistEntryService{}
		mockSS := &mocks.MockSongService{}
		mockMWH := &mocks.MockMiddlewareHandler{}

		var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {}

		mockMWH.
			On("AuthenticateUser").
			Return(mockAuthHF)

		mockSL.
			On("FetchByID", mock.AnythingOfType("*context.emptyCtx"), expSetlist.ID).
			Return(expSetlist, nil)

		mockSLES.
			On("FetchBySetlist", mock.AnythingOfType("*context.emptyCtx"), &[]domain.Setlist{*expSetlist}).
			Return(expSetlistEntries, nil)

		writer := prepareAndServeGet(t, fmt.Sprintf("/%d", expSetlist.ID), mockSL, mockSLES, mockSS, mockMWH)

		expBody, err := json.Marshal(gin.H{
			"setlist": expSetlist,
			"entries": expSetlistEntries,
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, writer.Code)
		assert.Equal(t, expBody, writer.Body.Bytes())
		mockSL.AssertExpectations(t)
		mockSLES.AssertExpectations(t)
		mockSS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})

	t.Run("Fail Invalid Param", func(t *testing.T) {
		t.Parallel()

		mockErr := domain.NewBadRequestErr("Could not read a")
		mockSL := &mocks.MockSetlistService{}
		mockSLES := &mocks.MockSetlistEntryService{}
		mockSS := &mocks.MockSongService{}
		mockMWH := &mocks.MockMiddlewareHandler{}

		var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {}

		mockMWH.
			On("AuthenticateUser").
			Return(mockAuthHF)

		writer := prepareAndServeGet(t, fmt.Sprintf("/%s", "a"), mockSL, mockSLES, mockSS, mockMWH)

		expBody, err := json.Marshal(gin.H{
			"error": mockErr.Error(),
		})
		assert.NoError(t, err)

		assert.Equal(t, domain.Status(mockErr), writer.Code)
		assert.Equal(t, expBody, writer.Body.Bytes())
		mockSL.AssertExpectations(t)
		mockSLES.AssertExpectations(t)
		mockSS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})

	t.Run("Fail Fetch Setlist", func(t *testing.T) {
		t.Parallel()

		mockErr := domain.NewInternalErr()
		mockSL := &mocks.MockSetlistService{}
		mockSLES := &mocks.MockSetlistEntryService{}
		mockSS := &mocks.MockSongService{}
		mockMWH := &mocks.MockMiddlewareHandler{}

		var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {}

		mockMWH.
			On("AuthenticateUser").
			Return(mockAuthHF)

		mockSL.
			On("FetchByID", mock.AnythingOfType("*context.emptyCtx"), expSetlist.ID).
			Return(nil, mockErr)

		writer := prepareAndServeGet(t, fmt.Sprintf("/%d", expSetlist.ID), mockSL, mockSLES, mockSS, mockMWH)

		expBody, err := json.Marshal(gin.H{
			"error": mockErr.Error(),
		})
		assert.NoError(t, err)

		assert.Equal(t, domain.Status(mockErr), writer.Code)
		assert.Equal(t, expBody, writer.Body.Bytes())
		mockSL.AssertExpectations(t)
		mockSLES.AssertExpectations(t)
		mockSS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})

	t.Run("Fail Fetch SetlistEntries", func(t *testing.T) {
		t.Parallel()

		mockErr := domain.NewInternalErr()
		mockSL := &mocks.MockSetlistService{}
		mockSLES := &mocks.MockSetlistEntryService{}
		mockSS := &mocks.MockSongService{}
		mockMWH := &mocks.MockMiddlewareHandler{}

		var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {}

		mockMWH.
			On("AuthenticateUser").
			Return(mockAuthHF)

		mockSL.
			On("FetchByID", mock.AnythingOfType("*context.emptyCtx"), expSetlist.ID).
			Return(expSetlist, nil)

		mockSLES.
			On("FetchBySetlist", mock.AnythingOfType("*context.emptyCtx"), &[]domain.Setlist{*expSetlist}).
			Return(nil, mockErr)

		writer := prepareAndServeGet(t, fmt.Sprintf("/%d", expSetlist.ID), mockSL, mockSLES, mockSS, mockMWH)

		expBody, err := json.Marshal(gin.H{
			"error": mockErr.Error(),
		})
		assert.NoError(t, err)

		assert.Equal(t, domain.Status(mockErr), writer.Code)
		assert.Equal(t, expBody, writer.Body.Bytes())
		mockSL.AssertExpectations(t)
		mockSLES.AssertExpectations(t)
		mockSS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})
}
