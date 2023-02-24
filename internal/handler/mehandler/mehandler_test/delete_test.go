package mehandler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/handler/mehandler"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func prepareAndServeDelete(
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

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodDelete, "/me/delete", requestBody)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestDeleteCorrect(t *testing.T) {
	t.Parallel()

	var deleteID int64 = 1

	mockUser := &domain.User{
		ID:         deleteID,
		Permission: domain.GUEST,
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}

	mockTS.
		On("RemoveAllRefresh", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(nil)
	mockUS.
		On("Remove", mock.AnythingOfType("*context.emptyCtx"), mockUser, deleteID).
		Return(mockUser.ID, nil)

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(gin.H{
		"id": deleteID,
	})
	assert.NoError(t, err)

	writer := prepareAndServeDelete(t, mockUS, mockTS, mockMWH, &byteBody)
	assert.Equal(t, http.StatusAccepted, writer.Code)
	mockUS.AssertExpectations(t)
}

func TestDeleteNoContext(t *testing.T) {
	t.Parallel()

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) { ctx.Next() }

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	writer := prepareAndServeDelete(t, mockUS, mockTS, mockMWH, &[]byte{})
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	mockUS.AssertNotCalled(t, "Remove")
}

func TestDeleteInvalidBind(t *testing.T) {
	t.Parallel()

	var deleteID int64 = 1

	mockUser := &domain.User{
		ID:         deleteID,
		Permission: domain.GUEST,
	}

	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}
	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	mockByte, err := json.Marshal(gin.H{
		"ids": deleteID,
	})
	assert.NoError(t, err)

	writer := prepareAndServeDelete(t, mockUS, mockTS, mockMWH, &mockByte)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	mockUS.AssertNotCalled(t, "Remove")
}

func TestDeleteRemoveErr(t *testing.T) {
	t.Parallel()

	var deleteID int64 = 1

	mockUser := &domain.User{
		ID:         deleteID,
		Permission: domain.GUEST,
	}

	expectedErr := domain.NewRecordNotFoundErr("", "")
	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}
	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)
	mockUS.
		On("Remove", mock.AnythingOfType("*context.emptyCtx"), mockUser, deleteID).
		Return(nil, expectedErr)

	mockByte, err := json.Marshal(gin.H{
		"id": deleteID,
	})
	assert.NoError(t, err)

	writer := prepareAndServeDelete(t, mockUS, mockTS, mockMWH, &mockByte)

	expectedBody, err := json.Marshal(gin.H{"error": expectedErr})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, writer.Code)
	assert.Equal(t, expectedBody, writer.Body.Bytes())
	mockUS.AssertNotCalled(t, "Remove")
}

func TestDeleteRemoveRefreshErr(t *testing.T) {
	t.Parallel()

	var deleteID int64 = 1

	mockUser := &domain.User{
		ID:         deleteID,
		Permission: domain.GUEST,
	}

	mockErr := domain.NewInternalErr()
	mockUS := &mocks.MockUserService{}
	mockTS := &mocks.MockTokenService{}
	mockMWH := &mocks.MockMiddlewareHandler{}

	mockTS.
		On("RemoveAllRefresh", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(mockErr)
	mockUS.
		On("Remove", mock.AnythingOfType("*context.emptyCtx"), mockUser, deleteID).
		Return(mockUser.ID, nil)

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	mockByte, err := json.Marshal(gin.H{
		"id": deleteID,
	})
	assert.NoError(t, err)

	writer := prepareAndServeDelete(t, mockUS, mockTS, mockMWH, &mockByte)
	assert.Equal(t, mockErr.Status(), writer.Code)
	mockUS.AssertExpectations(t)
}
