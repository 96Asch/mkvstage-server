package songhandler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteByIDCorrect(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
	}

	sid := int64(1)

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)
	mockSS := &mocks.MockSongService{}
	mockSS.
		On("Remove", mock.AnythingOfType("*context.emptyCtx"), sid, mockUser).
		Return(nil)

	gin.SetMode(gin.TestMode)

	router := gin.New()
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/songs/%d/delete", sid), nil)
	assert.NoError(t, err)

	Initialize(&router.RouterGroup, mockSS, mockMWH)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	mockMWH.AssertExpectations(t)
	mockSS.AssertExpectations(t)
}

func TestDeleteByIDNoContext(t *testing.T) {

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

	req, err := http.NewRequest(http.MethodDelete, "/songs/1/delete", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSS.AssertExpectations(t)
}

func TestDeleteByIDInvalidParam(t *testing.T) {
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

	req, err := http.NewRequest(http.MethodDelete, "/songs/a/delete", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockMWH.AssertExpectations(t)
	mockSS.AssertExpectations(t)
}

func TestDeleteByIDRemoveErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "Foo@Bar.com",
		Permission: domain.GUEST,
	}

	sid := int64(1)
	mockErr := domain.NewInternalErr()
	mockSS := new(mocks.MockSongService)
	mockSS.
		On("Remove", mock.AnythingOfType("*context.emptyCtx"), sid, mockUser).
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

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/songs/%d/delete", sid), nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, domain.Status(mockErr), w.Code)
	mockSS.AssertExpectations(t)
}
