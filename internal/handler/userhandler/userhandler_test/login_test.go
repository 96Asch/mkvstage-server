package userhandler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/handler/userhandler"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func prepareAndServeLogin(
	t *testing.T,
	mockUS domain.UserService,
	mockTS domain.TokenService,
	body *[]byte,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockMWH.
		On("AuthenticateUser").
		Return(nil)

	userhandler.Initialize(&router.RouterGroup, mockUS, mockTS, mockMWH)

	requestBody := bytes.NewReader(*body)

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, "/users/login", requestBody)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestLoginCorrect(t *testing.T) {
	t.Parallel()

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

	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}

	mockTS.
		On("CreateRefresh", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID, "").
		Return(mockRefresh, nil)
	mockTS.
		On("CreateAccess", mock.AnythingOfType("*context.emptyCtx"), mockRefresh.Refresh).
		Return(mockAccess, nil)
	mockUS.
		On("Authorize", mock.AnythingOfType("*context.emptyCtx"), mockUser.Email, mockUser.Password).
		Return(mockUser, nil)

	byteBody, err := json.Marshal(gin.H{
		"email":    mockUser.Email,
		"password": mockUser.Password,
	})
	assert.NoError(t, err)

	writer := prepareAndServeLogin(t, mockUS, mockTS, &byteBody)

	expectedRes, err := json.Marshal(gin.H{
		"tokens": gin.H{
			"access_token":  mockAccess.Access,
			"refresh_token": mockRefresh.Refresh,
		},
	})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, expectedRes, writer.Body.Bytes())
	mockTS.AssertExpectations(t)
	mockUS.AssertExpectations(t)
}

func TestLoginBindErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Password:     "FooBar",
		Email:        "Foo@Bar.com",
		Permission:   domain.GUEST,
		ProfileColor: "FFFFFF",
	}

	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}

	byteBody, err := json.Marshal(gin.H{
		"password": mockUser.Password,
	})
	assert.NoError(t, err)

	writer := prepareAndServeLogin(t, mockUS, mockTS, &byteBody)

	assert.Equal(t, http.StatusBadRequest, writer.Code)
	mockTS.AssertExpectations(t)
	mockUS.AssertExpectations(t)
}

func TestLoginNotAuth(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Password:     "FooBar",
		Email:        "Foo@Bar.com",
		Permission:   domain.GUEST,
		ProfileColor: "FFFFFF",
	}

	mockErr := domain.NewNotAuthorizedErr("invalid email/password")
	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}

	mockUS.
		On("Authorize", mock.AnythingOfType("*context.emptyCtx"), mockUser.Email, mockUser.Password).
		Return(nil, mockErr)

	byteBody, err := json.Marshal(gin.H{
		"email":    mockUser.Email,
		"password": mockUser.Password,
	})
	assert.NoError(t, err)

	writer := prepareAndServeLogin(t, mockUS, mockTS, &byteBody)

	assert.Equal(t, http.StatusUnauthorized, writer.Code)
	mockTS.AssertExpectations(t)
	mockUS.AssertExpectations(t)
}

func TestLoginRefreshErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Password:     "FooBar",
		Email:        "Foo@Bar.com",
		Permission:   domain.GUEST,
		ProfileColor: "FFFFFF",
	}

	mockErr := domain.NewInternalErr()
	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}

	mockTS.
		On("CreateRefresh", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID, "").
		Return(nil, mockErr)
	mockUS.
		On("Authorize", mock.AnythingOfType("*context.emptyCtx"), mockUser.Email, mockUser.Password).
		Return(mockUser, nil)

	byteBody, err := json.Marshal(gin.H{
		"email":    mockUser.Email,
		"password": mockUser.Password,
	})
	assert.NoError(t, err)

	writer := prepareAndServeLogin(t, mockUS, mockTS, &byteBody)

	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	mockUS.AssertExpectations(t)
}

func TestLoginAccessErr(t *testing.T) {
	t.Parallel()

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
	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}

	mockTS.
		On("CreateRefresh", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID, "").
		Return(mockRefresh, nil)
	mockTS.
		On("CreateAccess", mock.AnythingOfType("*context.emptyCtx"), mockRefresh.Refresh).
		Return(nil, mockErr)
	mockUS.
		On("Authorize", mock.AnythingOfType("*context.emptyCtx"), mockUser.Email, mockUser.Password).
		Return(mockUser, nil)

	byteBody, err := json.Marshal(gin.H{
		"email":    mockUser.Email,
		"password": mockUser.Password,
	})
	assert.NoError(t, err)

	writer := prepareAndServeLogin(t, mockUS, mockTS, &byteBody)

	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	mockUS.AssertExpectations(t)
}
