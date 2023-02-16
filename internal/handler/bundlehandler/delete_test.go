package bundlehandler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteCorrect(t *testing.T) {

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	bid := int64(1)
	mockBS := new(mocks.MockBundleService)
	mockBS.
		On("Remove", mock.AnythingOfType("*context.emptyCtx"), bid, mockUser).
		Return(nil)

	mockMWH := new(mocks.MockMiddlewareHandler)
	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	router := gin.New()
	Initialize(&router.RouterGroup, mockBS, mockMWH)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, "/bundles/1/delete", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	mockBS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)

}

func TestDeleteNoContext(t *testing.T) {

	mockBS := new(mocks.MockBundleService)

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Next()
	}
	mockMWH := new(mocks.MockMiddlewareHandler)
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	Initialize(&router.RouterGroup, mockBS, mockMWH)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, "/bundles/1/delete", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockBS.AssertExpectations(t)
}

func TestDeleteParamErr(t *testing.T) {

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockBS := new(mocks.MockBundleService)

	mockMWH := new(mocks.MockMiddlewareHandler)
	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	router := gin.New()
	Initialize(&router.RouterGroup, mockBS, mockMWH)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, "/bundles/a/delete", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockBS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)

}

func TestDeleteRemoveErr(t *testing.T) {

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockErr := domain.NewBadRequestErr("")
	bid := int64(1)
	mockBS := new(mocks.MockBundleService)
	mockBS.
		On("Remove", mock.AnythingOfType("*context.emptyCtx"), bid, mockUser).
		Return(mockErr)

	mockMWH := new(mocks.MockMiddlewareHandler)
	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	router := gin.New()
	Initialize(&router.RouterGroup, mockBS, mockMWH)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, "/bundles/1/delete", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, domain.Status(mockErr), w.Code)
	mockBS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)

}
