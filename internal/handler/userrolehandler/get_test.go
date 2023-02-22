package userrolehandler

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

func TestMeCorrect(t *testing.T) {

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
	var mockAuth gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH.On("AuthenticateUser").Return(mockAuth)

	mockURS := &mocks.MockUserRoleService{}
	mockURS.
		On("FetchByUser", mock.AnythingOfType("*context.emptyCtx"), mockUser).
		Return(mockUserRoles, nil)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/userroles/me", nil)
	assert.NoError(t, err)

	router := gin.New()
	Initialize(&router.RouterGroup, mockURS, mockMWH)
	gin.SetMode(gin.TestMode)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	bytebody, err := json.Marshal(gin.H{"user_roles": mockUserRoles})
	assert.NoError(t, err)
	assert.Equal(t, bytebody, w.Body.Bytes())
	mockURS.AssertExpectations(t)
}

func TestMeNoContext(t *testing.T) {
	mockMWH := &mocks.MockMiddlewareHandler{}
	var mockAuth gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Next()
	}
	mockMWH.On("AuthenticateUser").Return(mockAuth)

	mockURS := &mocks.MockUserRoleService{}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/userroles/me", nil)
	assert.NoError(t, err)

	router := gin.New()
	Initialize(&router.RouterGroup, mockURS, mockMWH)
	gin.SetMode(gin.TestMode)
	router.ServeHTTP(w, req)

	mockErr := domain.NewInternalErr()
	assert.Equal(t, domain.Status(mockErr), w.Code)
	mockURS.AssertExpectations(t)
}

func TestMeFetchErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockErr := domain.NewInternalErr()

	mockMWH := &mocks.MockMiddlewareHandler{}
	var mockAuth gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH.On("AuthenticateUser").Return(mockAuth)

	mockURS := &mocks.MockUserRoleService{}
	mockURS.
		On("FetchByUser", mock.AnythingOfType("*context.emptyCtx"), mockUser).
		Return(nil, mockErr)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/userroles/me", nil)
	assert.NoError(t, err)

	router := gin.New()
	Initialize(&router.RouterGroup, mockURS, mockMWH)
	gin.SetMode(gin.TestMode)
	router.ServeHTTP(w, req)

	assert.Equal(t, domain.Status(mockErr), w.Code)
	mockURS.AssertExpectations(t)
}

func TestGetAllCorrect(t *testing.T) {

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
	mockMWH.On("AuthenticateUser").Return(nil)
	mockURS := &mocks.MockUserRoleService{}
	mockURS.
		On("FetchAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(mockUserRoles, nil)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/userroles", nil)
	assert.NoError(t, err)

	router := gin.New()
	Initialize(&router.RouterGroup, mockURS, mockMWH)
	gin.SetMode(gin.TestMode)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	bytebody, err := json.Marshal(gin.H{"user_roles": mockUserRoles})
	assert.NoError(t, err)
	assert.Equal(t, bytebody, w.Body.Bytes())
	mockURS.AssertExpectations(t)
}

func TestGetAllFetchErr(t *testing.T) {
	mockErr := domain.NewInternalErr()

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockMWH.On("AuthenticateUser").Return(nil)

	mockURS := &mocks.MockUserRoleService{}
	mockURS.
		On("FetchAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(nil, mockErr)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/userroles", nil)
	assert.NoError(t, err)

	router := gin.New()
	Initialize(&router.RouterGroup, mockURS, mockMWH)
	gin.SetMode(gin.TestMode)
	router.ServeHTTP(w, req)

	assert.Equal(t, domain.Status(mockErr), w.Code)
	mockURS.AssertExpectations(t)
}
