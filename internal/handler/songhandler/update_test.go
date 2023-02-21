package songhandler

import (
	"bytes"
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

func TestUpdateByIDCorrect(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
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
	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)
	mockSS := &mocks.MockSongService{}
	mockSS.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), mockSong, mockUser).
		Return(nil)

	byteBody, err := json.Marshal(gin.H{
		"title":       "Foo",
		"subtitle":    "Bar",
		"key":         "A",
		"bpm":         120,
		"bundle_id":   1,
		"creator_id":  mockUser.ID,
		"chord_sheet": `{"Verse" : "Foobar"}`,
	})
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)

	router := gin.New()
	reqBody := bytes.NewReader(byteBody)
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/songs/%d/update", sid), reqBody)
	assert.NoError(t, err)

	Initialize(&router.RouterGroup, mockSS, mockMWH)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSong.ID = 1

	expectedBytes, err := json.Marshal(gin.H{"song": mockSong})
	assert.NoError(t, err)
	assert.Equal(t, expectedBytes, w.Body.Bytes())
	mockSS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDNoContext(t *testing.T) {
	sid := int64(1)
	mockSS := &mocks.MockSongService{}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Next()
	}
	mockMWH := new(mocks.MockMiddlewareHandler)
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	Initialize(&router.RouterGroup, mockSS, mockMWH)

	byteBody, err := json.Marshal(gin.H{
		"title":       "Foo",
		"subtitle":    "Bar",
		"key":         "A",
		"bpm":         120,
		"bundle_id":   1,
		"creator_id":  1,
		"chord_sheet": `{"Verse" : "Foobar"}`,
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(byteBody)

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/songs/%d/update", sid), bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDInvalidParam(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
	}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)
	mockSS := &mocks.MockSongService{}

	byteBody, err := json.Marshal(gin.H{
		"title":       "Foo",
		"subtitle":    "Bar",
		"key":         "A",
		"bpm":         120,
		"bundle_id":   1,
		"creator_id":  mockUser.ID,
		"chord_sheet": `{"Verse" : "Foobar"}`,
	})
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)

	router := gin.New()
	reqBody := bytes.NewReader(byteBody)
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, "/songs/a/update", reqBody)
	assert.NoError(t, err)

	Initialize(&router.RouterGroup, mockSS, mockMWH)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDBindErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "Foo@Bar.com",
		Permission: domain.GUEST,
	}
	sid := int64(1)
	mockSS := new(mocks.MockSongService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := new(mocks.MockMiddlewareHandler)
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	Initialize(&router.RouterGroup, mockSS, mockMWH)
	byteBody, err := json.Marshal(gin.H{
		"title":       "Foo",
		"subtitle":    "Bar",
		"bpm":         120,
		"bundle_id":   1,
		"creator_id":  mockUser.ID,
		"chord_sheet": `{"Verse" : "Foobar"}`,
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(byteBody)

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/songs/%d/update", sid), bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSS.AssertExpectations(t)
}

func TestUpdateByIDUpdateErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "Foo@Bar.com",
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
	mockErr := domain.NewBadRequestErr("")
	mockSS := new(mocks.MockSongService)
	mockSS.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), mockSong, mockUser).
		Return(mockErr)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := new(mocks.MockMiddlewareHandler)
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	Initialize(&router.RouterGroup, mockSS, mockMWH)
	byteBody, err := json.Marshal(gin.H{
		"title":       "Foo",
		"subtitle":    "Bar",
		"key":         "A",
		"bpm":         120,
		"bundle_id":   1,
		"creator_id":  mockUser.ID,
		"chord_sheet": `{"Verse" : "Foobar"}`,
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(byteBody)

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/songs/%d/update", sid), bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, domain.Status(mockErr), w.Code)
	mockSS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}
