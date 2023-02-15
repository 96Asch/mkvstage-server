package mehandler

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
)

func TestLogoutCorrect(t *testing.T) {
	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Password:     "FooBar",
		Email:        "Foo@Bar.com",
		Permission:   domain.GUEST,
		ProfileColor: "FFFFFF",
	}

	mockLogout := logoutReq{
		Refresh: "refresh-token",
	}

	mockTS := new(mocks.MockTokenService)
	mockTS.
		On("RemoveRefresh", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID, mockLogout.Refresh).
		Return(nil)

	mockUS := new(mocks.MockUserService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := new(mocks.MockMiddlewareHandler)
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	group := router.Group("test")
	Initialize(group, mockUS, mockTS, mockMWH)

	jsonLogout, err := json.Marshal(mockLogout)
	assert.NoError(t, err)

	reqBody := bytes.NewReader(jsonLogout)

	req, err := http.NewRequest(http.MethodDelete, "/test/me/logout", reqBody)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	mockUS.AssertExpectations(t)
}

func TestLogoutNoContext(t *testing.T) {

	mockLogout := &logoutReq{
		Refresh: "refresh-token",
	}
	mockTS := new(mocks.MockTokenService)
	mockUS := new(mocks.MockUserService)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	w := httptest.NewRecorder()

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Next()
	}
	mockMWH := new(mocks.MockMiddlewareHandler)
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	group := router.Group("test")
	Initialize(group, mockUS, mockTS, mockMWH)

	jsonLogout, err := json.Marshal(mockLogout)
	assert.NoError(t, err)

	reqBody := bytes.NewReader(jsonLogout)

	req, err := http.NewRequest(http.MethodDelete, "/test/me/logout", reqBody)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	expectedErr := domain.NewInternalErr()
	expectedRes, err := json.Marshal(gin.H{"error": expectedErr})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, expectedRes, w.Body.Bytes())
	mockUS.AssertExpectations(t)
}

func TestLogoutBindErr(t *testing.T) {
	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Password:     "FooBar",
		Email:        "Foo@Bar.com",
		Permission:   domain.GUEST,
		ProfileColor: "FFFFFF",
	}

	mockLogout := logoutReq{
		Refresh: "",
	}

	mockTS := new(mocks.MockTokenService)
	mockUS := new(mocks.MockUserService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := new(mocks.MockMiddlewareHandler)
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	group := router.Group("test")
	Initialize(group, mockUS, mockTS, mockMWH)

	jsonLogout, err := json.Marshal(mockLogout)
	assert.NoError(t, err)

	reqBody := bytes.NewReader(jsonLogout)

	req, err := http.NewRequest(http.MethodDelete, "/test/me/logout", reqBody)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUS.AssertExpectations(t)
}

func TestLogoutLogoutErr(t *testing.T) {
	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Password:     "FooBar",
		Email:        "Foo@Bar.com",
		Permission:   domain.GUEST,
		ProfileColor: "FFFFFF",
	}

	mockLogout := &logoutReq{
		Refresh: "refresh-token",
	}

	mockErr := domain.NewInternalErr()
	mockTS := new(mocks.MockTokenService)
	mockTS.
		On("RemoveRefresh", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID, mockLogout.Refresh).
		Return(mockErr)

	mockUS := new(mocks.MockUserService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := new(mocks.MockMiddlewareHandler)
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	group := router.Group("test")
	Initialize(group, mockUS, mockTS, mockMWH)

	jsonLogout, err := json.Marshal(mockLogout)
	assert.NoError(t, err)

	reqBody := bytes.NewReader(jsonLogout)

	req, err := http.NewRequest(http.MethodDelete, "/test/me/logout", reqBody)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUS.AssertExpectations(t)
}
