package mehandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/handler/mehandler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	refreshToken = "refresh-token"
)

func prepareAndServeLogout(
	t *testing.T,
	mockUS domain.UserService,
	mockTS domain.TokenService,
	mockMWH domain.MiddlewareHandler,
	body *[]byte,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	mehandler.Initialize(&router.RouterGroup, mockUS, mockTS, mockMWH)

	requestBody := bytes.NewReader(*body)

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodDelete, "/me/logout", requestBody)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestLogoutCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Password:     "FooBar",
		Email:        "Foo@Bar.com",
		Permission:   domain.GUEST,
		ProfileColor: "FFFFFF",
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}

	mockTS.
		On("RemoveRefresh", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID, refreshToken).
		Return(nil)

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(gin.H{
		"refresh": refreshToken,
	})
	assert.NoError(t, err)

	writer := prepareAndServeLogout(t, mockUS, mockTS, mockMWH, &byteBody)
	assert.Equal(t, http.StatusAccepted, writer.Code)
	mockMWH.AssertExpectations(t)
	mockTS.AssertExpectations(t)
	mockUS.AssertExpectations(t)
}

func TestLogoutNoContext(t *testing.T) {
	t.Parallel()

	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}
	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(gin.H{
		"refresh": refreshToken,
	})
	assert.NoError(t, err)

	writer := prepareAndServeLogout(t, mockUS, mockTS, mockMWH, &byteBody)

	expectedErr := domain.NewInternalErr()
	expectedRes, err := json.Marshal(gin.H{"error": expectedErr})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Equal(t, expectedRes, writer.Body.Bytes())
	mockMWH.AssertExpectations(t)
	mockTS.AssertExpectations(t)
	mockUS.AssertExpectations(t)
}

func TestLogoutBindErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Password:     "FooBar",
		Email:        "Foo@Bar.com",
		Permission:   domain.GUEST,
		ProfileColor: "FFFFFF",
	}

	mockMWH := new(mocks.MockMiddlewareHandler)
	mockTS := new(mocks.MockTokenService)
	mockUS := new(mocks.MockUserService)

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(gin.H{
		"refres": refreshToken,
	})
	assert.NoError(t, err)

	writer := prepareAndServeLogout(t, mockUS, mockTS, mockMWH, &byteBody)

	assert.Equal(t, http.StatusBadRequest, writer.Code)
	mockMWH.AssertExpectations(t)
	mockTS.AssertExpectations(t)
	mockUS.AssertExpectations(t)
}

func TestLogoutLogoutErr(t *testing.T) {
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
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}

	mockTS.
		On("RemoveRefresh", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID, refreshToken).
		Return(mockErr)

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(gin.H{
		"refresh": refreshToken,
	})
	assert.NoError(t, err)

	writer := prepareAndServeLogout(t, mockUS, mockTS, mockMWH, &byteBody)

	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	mockUS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
	mockTS.AssertExpectations(t)
	mockUS.AssertExpectations(t)
}
