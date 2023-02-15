package userhandler

import (
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

	mockTS := new(mocks.MockTokenService)
	mockTS.
		On("Logout", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(nil)

	mockUS := new(mocks.MockUserService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
	})
	w := httptest.NewRecorder()

	group := router.Group("test")
	Initialize(group, mockUS, mockTS)

	req, err := http.NewRequest(http.MethodDelete, "/test/users/me/logout", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	mockUS.AssertExpectations(t)
}

func TestLogoutNoContext(t *testing.T) {

	mockTS := new(mocks.MockTokenService)

	mockUS := new(mocks.MockUserService)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	w := httptest.NewRecorder()

	group := router.Group("test")
	Initialize(group, mockUS, mockTS)

	req, err := http.NewRequest(http.MethodDelete, "/test/users/me/logout", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	expectedErr := domain.NewInternalErr()
	expectedRes, err := json.Marshal(gin.H{"error": expectedErr})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, expectedRes, w.Body.Bytes())
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

	mockErr := domain.NewInternalErr()
	mockTS := new(mocks.MockTokenService)
	mockTS.
		On("Logout", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(mockErr)

	mockUS := new(mocks.MockUserService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
	})
	w := httptest.NewRecorder()

	group := router.Group("test")
	Initialize(group, mockUS, mockTS)

	req, err := http.NewRequest(http.MethodDelete, "/test/users/me/logout", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUS.AssertExpectations(t)
}
