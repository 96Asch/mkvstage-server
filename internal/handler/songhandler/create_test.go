package songhandler

import (
	"bytes"
	"encoding/json"
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

func TestCreateCorrect(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
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

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)
	mockSS := &mocks.MockSongService{}
	mockSS.
		On("Store", mock.AnythingOfType("*context.emptyCtx"), mockSong, mockUser).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).(*domain.Song)
			arg.ID = 1
		})

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
	req, err := http.NewRequest(http.MethodPost, "/songs/create", reqBody)
	assert.NoError(t, err)

	Initialize(&router.RouterGroup, mockSS, mockMWH)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSong.ID = 1

	expectedBytes, err := json.Marshal(gin.H{"song": mockSong})
	assert.NoError(t, err)
	assert.Equal(t, expectedBytes, w.Body.Bytes())
}

func TestCreateNoContext(t *testing.T) {

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

	req, err := http.NewRequest(http.MethodPost, "/songs/create", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSS.AssertExpectations(t)
}

func TestCreateBindErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "Foo@Bar.com",
		Permission: domain.GUEST,
	}

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

	req, err := http.NewRequest(http.MethodPost, "/songs/create", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSS.AssertExpectations(t)
}

func TestCreateStoreErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "Foo@Bar.com",
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

	mockErr := domain.NewBadRequestErr("")
	mockSS := new(mocks.MockSongService)
	mockSS.
		On("Store", mock.AnythingOfType("*context.emptyCtx"), mockSong, mockUser).
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

	req, err := http.NewRequest(http.MethodPost, "/songs/create", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, domain.Status(mockErr), w.Code)
	mockSS.AssertExpectations(t)
}
