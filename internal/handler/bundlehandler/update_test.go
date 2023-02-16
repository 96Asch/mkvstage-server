package bundlehandler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateCorrect(t *testing.T) {

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	uReq := &updateReq{
		Name:     "Foo",
		ParentID: 0,
	}

	mockBundle := &domain.Bundle{
		ID:       1,
		Name:     "Foo",
		ParentID: 0,
	}

	mockBS := new(mocks.MockBundleService)
	mockBS.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), mockBundle, mockUser).
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
	byteBody, err := json.Marshal(uReq)
	assert.NoError(t, err)
	req, err := http.NewRequest(http.MethodPut, "/bundles/1/update", bytes.NewReader(byteBody))
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedBody, err := json.Marshal(gin.H{"bundle": mockBundle})
	assert.NoError(t, err)
	assert.Equal(t, expectedBody, w.Body.Bytes())

	mockBS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)

}

func TestUpdateNoContext(t *testing.T) {

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
	req, err := http.NewRequest(http.MethodPut, "/bundles/1/update", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockBS.AssertExpectations(t)
}

func TestUpdateParamErr(t *testing.T) {

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
	req, err := http.NewRequest(http.MethodPut, "/bundles/a/update", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockBS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)

}

func TestUpdateBindErr(t *testing.T) {

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	uReq := &updateReq{
		Name:     "",
		ParentID: 0,
	}

	mockErr := domain.NewBadRequestErr("")

	mockBS := new(mocks.MockBundleService)

	mockMWH := new(mocks.MockMiddlewareHandler)
	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	router := gin.New()
	Initialize(&router.RouterGroup, mockBS, mockMWH)

	byteBody, err := json.Marshal(uReq)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, "/bundles/1/update", bytes.NewReader(byteBody))
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, domain.Status(mockErr), w.Code)
	mockBS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)

}

func TestUpdateErr(t *testing.T) {

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	uReq := &updateReq{
		Name:     "Foo",
		ParentID: 0,
	}

	mockBundle := &domain.Bundle{
		ID:       1,
		Name:     "Foo",
		ParentID: 0,
	}

	mockErr := domain.NewInternalErr()

	mockBS := new(mocks.MockBundleService)
	mockBS.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), mockBundle, mockUser).
		Return(mockErr)

	mockMWH := new(mocks.MockMiddlewareHandler)
	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	router := gin.New()
	Initialize(&router.RouterGroup, mockBS, mockMWH)

	byteBody, err := json.Marshal(uReq)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, "/bundles/1/update", bytes.NewReader(byteBody))
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, domain.Status(mockErr), w.Code)
	mockBS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)

}
