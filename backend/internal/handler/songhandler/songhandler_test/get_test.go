package songhandler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/96Asch/mkvstage-server/backend/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/backend/internal/handler/songhandler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func prepareAndServeGet(
	t *testing.T,
	mockSS domain.SongService,
	mockMWH domain.MiddlewareHandler,
	param string,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	songhandler.Initialize(&router.RouterGroup, mockSS, mockMWH)

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodGet,
		fmt.Sprintf("/songs%s", param),
		nil,
	)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestGet(t *testing.T) {
	mockSongs := []domain.Song{
		{
			ID:    1,
			Title: "Foobar",
		},
		{
			ID:    2,
			Title: "Barfoo",
		},
	}

	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	t.Run("Correct", func(t *testing.T) {
		t.Parallel()

		mockFilterOptions := &domain.SongFilterOptions{}
		mockSS := &mocks.MockSongService{}

		mockSS.
			On("Fetch", context.TODO(), mockFilterOptions).
			Return(mockSongs, nil)

		expBody, err := json.Marshal(gin.H{
			"songs": mockSongs,
		})
		assert.NoError(t, err)

		writer := prepareAndServeGet(t, mockSS, mockMWH, "/")

		assert.Equal(t, http.StatusOK, writer.Code)
		assert.Equal(t, expBody, writer.Body.Bytes())
		mockSS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})

	t.Run("Fail Song Fetch error", func(t *testing.T) {
		t.Parallel()

		expErr := domain.NewInternalErr()
		mockFilterOptions := &domain.SongFilterOptions{}
		mockSS := &mocks.MockSongService{}

		mockSS.
			On("Fetch", context.TODO(), mockFilterOptions).
			Return(nil, expErr)

		expBody, err := json.Marshal(gin.H{
			"error": expErr.Message,
		})
		assert.NoError(t, err)

		writer := prepareAndServeGet(t, mockSS, mockMWH, "/")

		assert.Equal(t, expErr.Status(), writer.Code)
		assert.Equal(t, expBody, writer.Body.Bytes())
		mockSS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})
}
func TestGetByIDCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}
	sid := int64(1)
	mockSong := &domain.Song{
		ID:         sid,
		CreatorID:  mockUser.ID,
		BundleID:   1,
		Title:      "Foo",
		Subtitle:   "Bar",
		Key:        "A",
		Bpm:        120,
		ChordSheet: datatypes.JSON([]byte(`{"Verse" : "Foobar"}`)),
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockSS := &mocks.MockSongService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)
	mockSS.
		On("FetchByID", context.TODO(), sid).
		Return(mockSong, nil)

	writer := prepareAndServeGet(t, mockSS, mockMWH, fmt.Sprintf("/%d", sid))

	expectedBody, err := json.Marshal(gin.H{"song": mockSong})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, expectedBody, writer.Body.Bytes())
	mockMWH.AssertExpectations(t)
	mockSS.AssertExpectations(t)
}

func TestGetByIDInvalidParam(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockSS := &mocks.MockSongService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	writer := prepareAndServeGet(t, mockSS, mockMWH, "/a")

	assert.Equal(t, http.StatusBadRequest, writer.Code)

	mockMWH.AssertExpectations(t)
	mockSS.AssertExpectations(t)
}

func TestGetByIDNoRecord(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	sid := int64(-1)
	mockErr := domain.NewRecordNotFoundErr("", "")
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockSS := &mocks.MockSongService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)
	mockSS.
		On("FetchByID", context.TODO(), sid).
		Return(nil, mockErr)

	writer := prepareAndServeGet(t, mockSS, mockMWH, fmt.Sprintf("/%d", sid))

	assert.Equal(t, http.StatusNotFound, writer.Code)
	mockMWH.AssertExpectations(t)
	mockSS.AssertExpectations(t)
}
