package userrolehandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/96Asch/mkvstage-server/backend/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/backend/internal/handler/userrolehandler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func prepareAndServeUpdate(
	t *testing.T,
	mockURS domain.UserRoleService,
	mockMWH domain.MiddlewareHandler,
	body *[]byte,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	userrolehandler.Initialize(&router.RouterGroup, mockURS, mockMWH)

	requestBody := bytes.NewReader(*body)

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodPatch,
		"/userroles/update",
		requestBody,
	)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestUpdateByIDCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
	}

	mockUserRoles := &[]domain.UserRole{
		{
			ID:     1,
			UserID: 1,
			RoleID: 1,
			Active: true,
		},
		{
			ID:     2,
			UserID: 1,
			RoleID: 2,
			Active: true,
		},
	}

	mockUserRoleIDs := []int64{1, 2}
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockURS := &mocks.MockUserRoleService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)
	mockURS.
		On("SetActiveBatch", context.TODO(), mockUserRoleIDs, mockUser).
		Return(mockUserRoles, nil)

	byteBody, err := json.Marshal(gin.H{
		"ids": mockUserRoleIDs,
	})
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, mockURS, mockMWH, &byteBody)

	expectedBytes, err := json.Marshal(gin.H{"user_roles": mockUserRoles})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, expectedBytes, writer.Body.Bytes())
	mockURS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDNoContext(t *testing.T) {
	t.Parallel()

	mockUserRoleIDs := []int64{1, 2}
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockURS := &mocks.MockUserRoleService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(gin.H{
		"ids": mockUserRoleIDs,
	})
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, mockURS, mockMWH, &byteBody)

	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	mockURS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDInvalidBind(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
	}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockURS := &mocks.MockUserRoleService{}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(gin.H{})
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, mockURS, mockMWH, &byteBody)

	assert.Equal(t, http.StatusBadRequest, writer.Code)

	mockURS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDSetActiveErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
	}

	mockUserRoleIDs := []int64{1, 2}
	mockErr := domain.NewInternalErr()
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockURS := &mocks.MockUserRoleService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)
	mockURS.
		On("SetActiveBatch", context.TODO(), mockUserRoleIDs, mockUser).
		Return(nil, mockErr)

	byteBody, err := json.Marshal(gin.H{
		"ids": mockUserRoleIDs,
	})
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, mockURS, mockMWH, &byteBody)

	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	mockURS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}
