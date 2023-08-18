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

func TestGetAllCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockSongs := &[]domain.Song{
		{
			CreatorID:  mockUser.ID,
			BundleID:   0,
			Title:      "Foo",
			Subtitle:   "Bar",
			Key:        "A",
			Bpm:        120,
			ChordSheet: datatypes.JSON([]byte(`{"Verse" : "Foobar"}`)),
		},
		{
			CreatorID:  mockUser.ID,
			BundleID:   1,
			Title:      "Bar",
			Subtitle:   "Foo",
			Key:        "B",
			Bpm:        122,
			ChordSheet: datatypes.JSON([]byte(`{"Verse" : "BarFoo"}`)),
		},
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
		On("FetchAll", context.TODO()).
		Return(mockSongs, nil)

	writer := prepareAndServeGet(t, mockSS, mockMWH, "")

	expectedBody, err := json.Marshal(gin.H{"songs": mockSongs})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, expectedBody, writer.Body.Bytes())
	mockMWH.AssertExpectations(t)
	mockSS.AssertExpectations(t)
}

func TestGetAllFetchErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockErr := domain.NewInternalErr()
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
		On("FetchAll", context.TODO()).
		Return(nil, mockErr)

	writer := prepareAndServeGet(t, mockSS, mockMWH, "")

	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	mockMWH.AssertExpectations(t)
	mockSS.AssertExpectations(t)
}
