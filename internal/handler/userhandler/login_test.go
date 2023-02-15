package userhandler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoginCorrect(t *testing.T) {
	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Password:     "FooBar",
		Email:        "Foo@Bar.com",
		Permission:   domain.GUEST,
		ProfileColor: "FFFFFF",
	}

	mockRefresh := &domain.RefreshToken{
		ID:                 uuid.New(),
		UserID:             mockUser.ID,
		Refresh:            "refresh-token",
		ExpirationDuration: time.Minute,
	}

	mockAccess := &domain.AccessToken{
		Access: "access-token",
	}

	mockTS := new(mocks.MockTokenService)
	mockTS.
		On("CreateRefresh", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID, "").
		Return(mockRefresh, nil)

	mockTS.
		On("CreateAccess", mock.AnythingOfType("*context.emptyCtx"), mockRefresh.Refresh).
		Return(mockAccess, nil)

	mockUS := new(mocks.MockUserService)
	mockUS.
		On("Authorize", mock.AnythingOfType("*context.emptyCtx"), mockUser.Email, mockUser.Password).
		Return(mockUser, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	group := router.Group("test")
	Initialize(group, mockUS, mockTS)

	mockByte, err := json.Marshal(gin.H{
		"email":    mockUser.Email,
		"password": mockUser.Password,
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(mockByte)

	req, err := http.NewRequest(http.MethodPost, "/test/users/login", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	expectedRes, err := json.Marshal(gin.H{
		"tokens": tokenPair{
			Access:  mockAccess.Access,
			Refresh: mockRefresh.Refresh,
		},
	})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedRes, w.Body.Bytes())
	mockUS.AssertExpectations(t)
}

func TestLoginBindErr(t *testing.T) {
	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Password:     "FooBar",
		Email:        "Foo@Bar.com",
		Permission:   domain.GUEST,
		ProfileColor: "FFFFFF",
	}

	mockTS := new(mocks.MockTokenService)
	mockUS := new(mocks.MockUserService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	group := router.Group("test")
	Initialize(group, mockUS, mockTS)

	mockByte, err := json.Marshal(gin.H{
		"password": mockUser.Password,
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(mockByte)

	req, err := http.NewRequest(http.MethodPost, "/test/users/login", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUS.AssertExpectations(t)
}

func TestLoginNotAuth(t *testing.T) {
	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Password:     "FooBar",
		Email:        "Foo@Bar.com",
		Permission:   domain.GUEST,
		ProfileColor: "FFFFFF",
	}

	mockErr := domain.NewNotAuthorizedErr("invalid email/password")

	mockTS := new(mocks.MockTokenService)

	mockUS := new(mocks.MockUserService)
	mockUS.
		On("Authorize", mock.AnythingOfType("*context.emptyCtx"), mockUser.Email, mockUser.Password).
		Return(nil, mockErr)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	group := router.Group("test")
	Initialize(group, mockUS, mockTS)

	mockByte, err := json.Marshal(gin.H{
		"email":    mockUser.Email,
		"password": mockUser.Password,
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(mockByte)

	req, err := http.NewRequest(http.MethodPost, "/test/users/login", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockUS.AssertExpectations(t)
}

func TestLoginRefreshErr(t *testing.T) {
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
		On("CreateRefresh", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID, "").
		Return(nil, mockErr)

	mockUS := new(mocks.MockUserService)
	mockUS.
		On("Authorize", mock.AnythingOfType("*context.emptyCtx"), mockUser.Email, mockUser.Password).
		Return(mockUser, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	group := router.Group("test")
	Initialize(group, mockUS, mockTS)

	mockByte, err := json.Marshal(gin.H{
		"email":    mockUser.Email,
		"password": mockUser.Password,
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(mockByte)

	req, err := http.NewRequest(http.MethodPost, "/test/users/login", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUS.AssertExpectations(t)
}

func TestLoginAccessErr(t *testing.T) {
	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Password:     "FooBar",
		Email:        "Foo@Bar.com",
		Permission:   domain.GUEST,
		ProfileColor: "FFFFFF",
	}

	mockRefresh := &domain.RefreshToken{
		ID:                 uuid.New(),
		UserID:             mockUser.ID,
		Refresh:            "refresh-token",
		ExpirationDuration: time.Minute,
	}

	mockErr := domain.NewInternalErr()

	mockTS := new(mocks.MockTokenService)
	mockTS.
		On("CreateRefresh", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID, "").
		Return(mockRefresh, nil)

	mockTS.
		On("CreateAccess", mock.AnythingOfType("*context.emptyCtx"), mockRefresh.Refresh).
		Return(nil, mockErr)

	mockUS := new(mocks.MockUserService)
	mockUS.
		On("Authorize", mock.AnythingOfType("*context.emptyCtx"), mockUser.Email, mockUser.Password).
		Return(mockUser, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	group := router.Group("test")
	Initialize(group, mockUS, mockTS)

	mockByte, err := json.Marshal(gin.H{
		"email":    mockUser.Email,
		"password": mockUser.Password,
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(mockByte)

	req, err := http.NewRequest(http.MethodPost, "/test/users/login", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUS.AssertExpectations(t)
}
