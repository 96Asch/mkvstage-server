package userhandler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/handler/userhandler"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func prepareAndServeChangePermission(
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

	userhandler.Initialize(&router.RouterGroup, mockUS, mockTS, mockMWH)

	requestBody := bytes.NewReader(*body)

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodPut, "/users/setperm", requestBody)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestChangePermissionByIDCorrect(t *testing.T) {
	t.Parallel()

	permission := domain.EDITOR
	recipientID := int64(2)
	principal := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}
	expUser := &domain.User{
		ID:         int64(recipientID),
		Permission: permission,
	}

	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}
	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuth gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", principal)
		ctx.Next()
	}

	mockUS.
		On("SetPermission", mock.AnythingOfType("*context.emptyCtx"), permission, recipientID, principal).
		Return(expUser, nil)
	mockMWH.
		On("AuthenticateUser").
		Return(mockAuth)

	byteBody, err := json.Marshal(gin.H{
		"recipient_id": recipientID,
		"permission":   permission,
	})
	assert.NoError(t, err)

	writer := prepareAndServeChangePermission(t, mockUS, mockTS, mockMWH, &byteBody)

	expBody, err := json.Marshal(gin.H{"user": expUser})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, expBody, writer.Body.Bytes())
	mockUS.AssertExpectations(t)
	mockTS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestChangePermissionByIDNoContext(t *testing.T) {
	t.Parallel()

	permission := domain.EDITOR
	recipientID := int64(2)

	expErr := domain.NewInternalErr()
	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}
	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuth gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuth)

	byteBody, err := json.Marshal(gin.H{
		"recipient_id": recipientID,
		"permission":   permission,
	})
	assert.NoError(t, err)

	writer := prepareAndServeChangePermission(t, mockUS, mockTS, mockMWH, &byteBody)

	expBody, err := json.Marshal(gin.H{"error": expErr})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Equal(t, expBody, writer.Body.Bytes())
	mockUS.AssertExpectations(t)
	mockTS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestChangePermissionByCastErr(t *testing.T) {
	t.Parallel()

	permission := domain.EDITOR
	recipientID := int64(2)

	expErr := domain.NewInternalErr()
	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}
	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuth gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", "principal")
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuth)

	byteBody, err := json.Marshal(gin.H{
		"recipient_id": recipientID,
		"permission":   permission,
	})
	assert.NoError(t, err)

	writer := prepareAndServeChangePermission(t, mockUS, mockTS, mockMWH, &byteBody)

	expBody, err := json.Marshal(gin.H{"error": expErr})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Equal(t, expBody, writer.Body.Bytes())
	mockUS.AssertExpectations(t)
	mockTS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestChangePermissionByIDBindErr(t *testing.T) {
	t.Parallel()

	permission := domain.EDITOR
	principal := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}
	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuth gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", principal)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuth)

	byteBody, err := json.Marshal(gin.H{
		"permission": permission,
	})
	assert.NoError(t, err)

	writer := prepareAndServeChangePermission(t, mockUS, mockTS, mockMWH, &byteBody)

	assert.Equal(t, http.StatusBadRequest, writer.Code)
	mockUS.AssertExpectations(t)
	mockTS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestChangePermissionBySetPermErr(t *testing.T) {
	t.Parallel()

	permission := domain.EDITOR
	recipientID := int64(2)
	principal := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	expErr := domain.NewNotAuthorizedErr("")
	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}
	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuth gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", principal)
		ctx.Next()
	}

	mockUS.
		On("SetPermission", mock.AnythingOfType("*context.emptyCtx"), permission, recipientID, principal).
		Return(nil, expErr)
	mockMWH.
		On("AuthenticateUser").
		Return(mockAuth)

	byteBody, err := json.Marshal(gin.H{
		"recipient_id": recipientID,
		"permission":   permission,
	})
	assert.NoError(t, err)

	writer := prepareAndServeChangePermission(t, mockUS, mockTS, mockMWH, &byteBody)

	expBody, err := json.Marshal(gin.H{"error": expErr})
	assert.NoError(t, err)

	assert.Equal(t, domain.Status(expErr), writer.Code)
	assert.Equal(t, expBody, writer.Body.Bytes())
	mockUS.AssertExpectations(t)
	mockTS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}
