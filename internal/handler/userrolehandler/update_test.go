package userrolehandler

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

func TestUpdateByIDCorrect(t *testing.T) {
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

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	mockURS := &mocks.MockUserRoleService{}
	mockURS.
		On("SetActiveBatch", mock.AnythingOfType("*context.emptyCtx"), mockUserRoleIDs, mockUser).
		Return(mockUserRoles, nil)

	byteBody, err := json.Marshal(gin.H{
		"ids": mockUserRoleIDs,
	})
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	reqBody := bytes.NewReader(byteBody)
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPatch, "/userroles/update", reqBody)
	assert.NoError(t, err)

	Initialize(&router.RouterGroup, mockURS, mockMWH)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedBytes, err := json.Marshal(gin.H{"user_roles": mockUserRoles})
	assert.NoError(t, err)
	assert.Equal(t, expectedBytes, w.Body.Bytes())
	mockURS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDNoContext(t *testing.T) {

	mockUserRoleIDs := []int64{1, 2}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Next()
	}
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	mockURS := &mocks.MockUserRoleService{}

	byteBody, err := json.Marshal(gin.H{
		"ids": mockUserRoleIDs,
	})
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	reqBody := bytes.NewReader(byteBody)
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPatch, "/userroles/update", reqBody)
	assert.NoError(t, err)

	Initialize(&router.RouterGroup, mockURS, mockMWH)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockURS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDInvalidBind(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
	}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	mockURS := &mocks.MockUserRoleService{}

	byteBody, err := json.Marshal(gin.H{})
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	reqBody := bytes.NewReader(byteBody)
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPatch, "/userroles/update", reqBody)
	assert.NoError(t, err)

	Initialize(&router.RouterGroup, mockURS, mockMWH)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockURS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDSetActiveErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
	}

	mockErr := domain.NewInternalErr()
	mockUserRoleIDs := []int64{1, 2}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	mockURS := &mocks.MockUserRoleService{}
	mockURS.
		On("SetActiveBatch", mock.AnythingOfType("*context.emptyCtx"), mockUserRoleIDs, mockUser).
		Return(nil, mockErr)

	byteBody, err := json.Marshal(gin.H{
		"ids": mockUserRoleIDs,
	})
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	reqBody := bytes.NewReader(byteBody)
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPatch, "/userroles/update", reqBody)
	assert.NoError(t, err)

	Initialize(&router.RouterGroup, mockURS, mockMWH)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockURS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}
