package userrolehandler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/handler/userrolehandler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func prepareAndServeGet(
	t *testing.T,
	mockURS domain.UserRoleService,
	mockMWH domain.MiddlewareHandler,
	param string,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	userrolehandler.Initialize(&router.RouterGroup, mockURS, mockMWH)

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, fmt.Sprintf("/userroles%s", param), nil)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestMeCorrect(t *testing.T) {
	t.Parallel()

	mockUserRoles := &[]domain.UserRole{
		{
			ID:     1,
			UserID: 1,
			RoleID: 1,
		},
		{
			ID:     2,
			UserID: 1,
			RoleID: 2,
		},
	}

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockURS := &mocks.MockUserRoleService{}

	var mockAuth gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuth)
	mockURS.
		On("FetchByUser", context.TODO(), mockUser).
		Return(mockUserRoles, nil)

	writer := prepareAndServeGet(t, mockURS, mockMWH, "/me")

	bytebody, err := json.Marshal(gin.H{"user_roles": mockUserRoles})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, bytebody, writer.Body.Bytes())
	mockMWH.AssertExpectations(t)
	mockURS.AssertExpectations(t)
}

func TestMeNoContext(t *testing.T) {
	t.Parallel()

	mockErr := domain.NewInternalErr()
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockURS := &mocks.MockUserRoleService{}

	var mockAuth gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuth)

	writer := prepareAndServeGet(t, mockURS, mockMWH, "/me")

	assert.Equal(t, domain.Status(mockErr), writer.Code)
	mockMWH.AssertExpectations(t)
	mockURS.AssertExpectations(t)
}

func TestMeFetchErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockErr := domain.NewInternalErr()
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockURS := &mocks.MockUserRoleService{}

	var mockAuth gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuth)
	mockURS.
		On("FetchByUser", context.TODO(), mockUser).
		Return(nil, mockErr)

	writer := prepareAndServeGet(t, mockURS, mockMWH, "/me")

	assert.Equal(t, domain.Status(mockErr), writer.Code)
	mockURS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestGetAllCorrect(t *testing.T) {
	t.Parallel()

	mockUserRoles := &[]domain.UserRole{
		{
			ID:     1,
			UserID: 1,
			RoleID: 1,
		},
		{
			ID:     2,
			UserID: 1,
			RoleID: 2,
		},
		{
			ID:     3,
			UserID: 2,
			RoleID: 1,
		},
		{
			ID:     4,
			UserID: 2,
			RoleID: 2,
		},
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockURS := &mocks.MockUserRoleService{}

	mockMWH.
		On("AuthenticateUser").
		Return(nil)
	mockURS.
		On("FetchAll", context.TODO()).
		Return(mockUserRoles, nil)

	writer := prepareAndServeGet(t, mockURS, mockMWH, "")

	expBody, err := json.Marshal(gin.H{"user_roles": mockUserRoles})
	assert.NoError(t, err)

	assert.Equal(t, expBody, writer.Body.Bytes())
	assert.Equal(t, http.StatusOK, writer.Code)
	mockMWH.AssertExpectations(t)
	mockURS.AssertExpectations(t)
}

func TestGetAllFetchErr(t *testing.T) {
	t.Parallel()

	mockErr := domain.NewInternalErr()
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockURS := &mocks.MockUserRoleService{}

	mockMWH.
		On("AuthenticateUser").
		Return(nil)
	mockURS.
		On("FetchAll", context.TODO()).
		Return(nil, mockErr)

	writer := prepareAndServeGet(t, mockURS, mockMWH, "")

	assert.Equal(t, domain.Status(mockErr), writer.Code)
	mockURS.AssertExpectations(t)
}
