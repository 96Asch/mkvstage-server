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
	mockSS domain.SongService,
	mockMWH domain.MiddlewareHandler,
	body *[]byte,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	setlisthandler.Initialize(&router.RouterGroup, mockSL, mockSS, mockMWH)

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
		Global:    true,
		Deadline:  time.Now().Add(time.Hour * 24).Round(0),
		UpdatedAt: time.Now().Round(0),
		Order:     datatypes.JSON([]byte(`{"order" : "1,2,3,4"}`)),
	}

	mockSetlist := &domain.Setlist{
		Name:      expSetlist.Name,
		CreatorID: expSetlist.ID,
		Global:    expSetlist.Global,
		Deadline:  expSetlist.Deadline,
	}

	mockSL := &mocks.MockSetlistService{}
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

	byteBody, err := json.Marshal(mockSetlist)
	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockSL, mockSS, mockMWH, &byteBody)

	assert.Equal(t, http.StatusCreated, writer.Code)

	expBody, err := json.Marshal(gin.H{"setlist": expSetlist})
	assert.NoError(t, err)

	assert.Equal(t, expBody, writer.Body.Bytes())
	mockSL.AssertExpectations(t)
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
		Global:    true,
	}

	mockErr := domain.NewBadRequestErr("")
	mockSL := &mocks.MockSetlistService{}
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

	writer := prepareAndServeCreate(t, mockSL, mockSS, mockMWH, &byteBody)

	assert.Equal(t, mockErr.Status(), writer.Code)
	mockSL.AssertExpectations(t)
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
		Global:    true,
		Deadline:  time.Now().Add(time.Hour * 24),
	}

	mockErr := domain.NewInternalErr()
	mockSL := &mocks.MockSetlistService{}
	mockSS := &mocks.MockSongService{}

	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(mockSetlist)
	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockSL, mockSS, mockMWH, &byteBody)

	assert.Equal(t, mockErr.Status(), writer.Code)
	mockSL.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestCreateStoreErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
	}

	mockSetlist := &domain.Setlist{
		Name:      "Foo",
		CreatorID: mockUser.ID,
		Global:    true,
		Deadline:  time.Now().Add(time.Hour * 24).Round(0),
	}

	mockErr := domain.NewBadRequestErr("")
	mockSL := &mocks.MockSetlistService{}
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

	byteBody, err := json.Marshal(mockSetlist)
	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockSL, mockSS, mockMWH, &byteBody)

	assert.Equal(t, mockErr.Status(), writer.Code)
	mockSL.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}
