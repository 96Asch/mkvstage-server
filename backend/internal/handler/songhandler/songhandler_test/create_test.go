package songhandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/handler/songhandler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/datatypes"
)

func prepareAndServeCreate(
	t *testing.T,
	mockSS domain.SongService,
	mockMWH domain.MiddlewareHandler,
	body *[]byte,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	songhandler.Initialize(&router.RouterGroup, mockSS, mockMWH)

	requestBody := bytes.NewReader(*body)

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodPost,
		"/songs/create",
		requestBody,
	)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestCreateCorrect(t *testing.T) {
	t.Parallel()

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
		On("Store", context.TODO(), mockSong, mockUser).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg, ok := args.Get(1).(*domain.Song)
			assert.True(t, ok)
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

	writer := prepareAndServeCreate(t, mockSS, mockMWH, &byteBody)
	mockSong.ID = 1

	expectedBytes, err := json.Marshal(gin.H{"song": mockSong})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, writer.Code)
	assert.Equal(t, expectedBytes, writer.Body.Bytes())
	mockSS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestCreateNoContext(t *testing.T) {
	t.Parallel()

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockSS := &mocks.MockSongService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

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

	writer := prepareAndServeCreate(t, mockSS, mockMWH, &byteBody)

	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	mockSS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestCreateBindErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "Foo@Bar.com",
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

	byteBody, err := json.Marshal(gin.H{
		"title":       "Foo",
		"subtitle":    "Bar",
		"bpm":         120,
		"bundle_id":   1,
		"creator_id":  mockUser.ID,
		"chord_sheet": `{"Verse" : "Foobar"}`,
	})
	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockSS, mockMWH, &byteBody)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	mockSS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestCreateStoreErr(t *testing.T) {
	t.Parallel()

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
	mockMWH := new(mocks.MockMiddlewareHandler)
	mockSS := new(mocks.MockSongService)

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)
	mockSS.
		On("Store", context.TODO(), mockSong, mockUser).
		Return(mockErr)

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

	writer := prepareAndServeCreate(t, mockSS, mockMWH, &byteBody)
	assert.Equal(t, domain.Status(mockErr), writer.Code)
	mockSS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}
