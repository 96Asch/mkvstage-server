package rolehandler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteByIDCorrect(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
	}

	sid := int64(1)

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)
	mockRS := &mocks.MockRoleService{}
	mockRS.
		On("Remove", mock.AnythingOfType("*context.emptyCtx"), sid, mockUser).
		Return(nil)

	gin.SetMode(gin.TestMode)

	router := gin.New()
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/roles/%d/delete", sid), nil)
	assert.NoError(t, err)

	Initialize(&router.RouterGroup, mockRS, mockMWH)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	mockMWH.AssertExpectations(t)
	mockRS.AssertExpectations(t)
}

func TestDeleteByIDNoContext(t *testing.T) {

	mockRS := &mocks.MockRoleService{}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Next()
	}
	mockMWH := new(mocks.MockMiddlewareHandler)
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	Initialize(&router.RouterGroup, mockRS, mockMWH)

	req, err := http.NewRequest(http.MethodDelete, "/roles/1/delete", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRS.AssertExpectations(t)
}

func TestDeleteByIDInvalidParam(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockMWH := new(mocks.MockMiddlewareHandler)
	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	mockRS := &mocks.MockRoleService{}

	router := gin.New()
	Initialize(&router.RouterGroup, mockRS, mockMWH)

	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodDelete, "/roles/a/delete", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockMWH.AssertExpectations(t)
	mockRS.AssertExpectations(t)
}

func TestDeleteByIDRemoveErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "Foo@Bar.com",
		Permission: domain.GUEST,
	}

	sid := int64(1)
	mockErr := domain.NewInternalErr()
	mockRS := new(mocks.MockRoleService)
	mockRS.
		On("Remove", mock.AnythingOfType("*context.emptyCtx"), sid, mockUser).
		Return(mockErr)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := new(mocks.MockMiddlewareHandler)
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	Initialize(&router.RouterGroup, mockRS, mockMWH)

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/roles/%d/delete", sid), nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, domain.Status(mockErr), w.Code)
	mockRS.AssertExpectations(t)
}
