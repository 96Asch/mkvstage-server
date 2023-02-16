package bundlehandler

import (
	"fmt"
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

func TestGetCorrect(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockBundle := &domain.Bundle{
		ID:       1,
		Name:     "Foo",
		ParentID: 0,
	}

	mockMWH := new(mocks.MockMiddlewareHandler)
	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	mockBS := new(mocks.MockBundleService)
	mockBS.
		On("FetchByID", mock.AnythingOfType("*context.emptyCtx"), mockBundle.ID).
		Return(mockBundle, nil)

	router := gin.New()
	Initialize(&router.RouterGroup, mockBS, mockMWH)

	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/bundles/%d", mockBundle.ID), nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	expectedBody, err := json.Marshal(gin.H{"bundle": mockBundle})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedBody, w.Body.Bytes())

	mockMWH.AssertExpectations(t)
	mockBS.AssertExpectations(t)
}

func TestGetInvalidParam(t *testing.T) {
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

	mockBS := new(mocks.MockBundleService)

	router := gin.New()
	Initialize(&router.RouterGroup, mockBS, mockMWH)

	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/bundles/a", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockMWH.AssertExpectations(t)
	mockBS.AssertExpectations(t)
}

func TestGetNoRecord(t *testing.T) {
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

	mockErr := domain.NewRecordNotFoundErr("", "")

	mockBS := new(mocks.MockBundleService)
	mockBS.
		On("FetchByID", mock.AnythingOfType("*context.emptyCtx"), int64(-1)).
		Return(nil, mockErr)

	router := gin.New()
	Initialize(&router.RouterGroup, mockBS, mockMWH)

	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/bundles/%d", -1), nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockMWH.AssertExpectations(t)
	mockBS.AssertExpectations(t)
}

func TestGetAllCorrect(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockBundles := &[]domain.Bundle{
		{
			ID:       1,
			Name:     "Foo",
			ParentID: 0,
		},
		{
			ID:       2,
			Name:     "Bar",
			ParentID: 1,
		},
	}

	mockMWH := new(mocks.MockMiddlewareHandler)
	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	mockBS := new(mocks.MockBundleService)
	mockBS.
		On("FetchAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(mockBundles, nil)

	router := gin.New()
	Initialize(&router.RouterGroup, mockBS, mockMWH)

	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/bundles", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	expectedBody, err := json.Marshal(gin.H{"bundles": mockBundles})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedBody, w.Body.Bytes())

	mockMWH.AssertExpectations(t)
	mockBS.AssertExpectations(t)
}

func TestGetAllFetchErr(t *testing.T) {
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

	mockErr := domain.NewInternalErr()

	mockBS := new(mocks.MockBundleService)
	mockBS.
		On("FetchAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(nil, mockErr)

	router := gin.New()
	Initialize(&router.RouterGroup, mockBS, mockMWH)

	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/bundles", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockMWH.AssertExpectations(t)
	mockBS.AssertExpectations(t)
}
