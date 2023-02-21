package songhandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/datatypes"
)

func TestGetByIDCorrect(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockSong := &domain.Song{
		CreatorID:  mockUser.ID,
		BundleID:   1,
		Title:      "Foo",
		Subtitle:   "Bar",
		Key:        "A",
		Bpm:        120,
		ChordSheet: datatypes.JSON([]byte(`{"Verse" : "Foobar"}`)),
	}

	mockMWH := new(mocks.MockMiddlewareHandler)
	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	mockSS := &mocks.MockSongService{}
	mockSS.
		On("FetchByID", mock.AnythingOfType("*context.emptyCtx"), mockSong.ID).
		Return(mockSong, nil)

	router := gin.New()
	Initialize(&router.RouterGroup, mockSS, mockMWH)

	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/songs/%d", mockSong.ID), nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	expectedBody, err := json.Marshal(gin.H{"song": mockSong})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedBody, w.Body.Bytes())

	mockMWH.AssertExpectations(t)
	mockSS.AssertExpectations(t)
}

func TestGetByIDInvalidParam(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockMWH := new(mocks.MockMiddlewareHandler)
	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	mockSS := &mocks.MockSongService{}

	router := gin.New()
	Initialize(&router.RouterGroup, mockSS, mockMWH)

	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/songs/a", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockMWH.AssertExpectations(t)
	mockSS.AssertExpectations(t)
}

func TestGetByIDNoRecord(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockMWH := new(mocks.MockMiddlewareHandler)
	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	mockErr := domain.NewRecordNotFoundErr("", "")

	mockSS := &mocks.MockSongService{}
	mockSS.
		On("FetchByID", mock.AnythingOfType("*context.emptyCtx"), int64(-1)).
		Return(nil, mockErr)

	router := gin.New()
	Initialize(&router.RouterGroup, mockSS, mockMWH)

	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/songs/%d", -1), nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockMWH.AssertExpectations(t)
	mockSS.AssertExpectations(t)
}

func TestGetAllCorrect(t *testing.T) {
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

	mockMWH := new(mocks.MockMiddlewareHandler)
	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	mockSS := &mocks.MockSongService{}
	mockSS.
		On("FetchAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(mockSongs, nil)

	router := gin.New()
	Initialize(&router.RouterGroup, mockSS, mockMWH)

	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/songs", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	expectedBody, err := json.Marshal(gin.H{"songs": mockSongs})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedBody, w.Body.Bytes())

	mockMWH.AssertExpectations(t)
	mockSS.AssertExpectations(t)
}

func TestGetAllFetchErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockMWH := new(mocks.MockMiddlewareHandler)
	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	mockErr := domain.NewInternalErr()

	mockSS := &mocks.MockSongService{}
	mockSS.
		On("FetchAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(nil, mockErr)

	router := gin.New()
	Initialize(&router.RouterGroup, mockSS, mockMWH)

	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/songs", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockMWH.AssertExpectations(t)
	mockSS.AssertExpectations(t)
}
