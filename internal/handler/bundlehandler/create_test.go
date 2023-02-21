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

func TestCreateCorrect(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "Foo@Bar.com",
		Permission: domain.MEMBER,
	}

	mockBundle := &domain.Bundle{
		Name:     "Foobar",
		ParentID: 0,
	}

	mockBS := new(mocks.MockBundleService)
	mockBS.
		On("Store", mock.AnythingOfType("*context.emptyCtx"), mockBundle, mockUser).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).(*domain.Bundle)
			arg.ID = 1
		})

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	group := router.Group("test")

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := new(mocks.MockMiddlewareHandler)
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	Initialize(group, mockBS, mockMWH)

	mockReq, err := json.Marshal(gin.H{
		"id":        int64(1),
		"name":      "Foobar",
		"parent_id": int64(0),
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(mockReq)

	req, err := http.NewRequest(http.MethodPost, "/test/bundles/create", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	expectedBundle, err := json.Marshal(gin.H{
		"bundle": &domain.Bundle{
			ID:       1,
			Name:     mockBundle.Name,
			ParentID: mockBundle.ParentID,
		},
	})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, expectedBundle, w.Body.Bytes())
	mockBS.AssertExpectations(t)
}

func TestCreateBindErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "Foo@Bar.com",
		Permission: domain.MEMBER,
	}

	mockBS := new(mocks.MockBundleService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	group := router.Group("test")
	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := new(mocks.MockMiddlewareHandler)
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	Initialize(group, mockBS, mockMWH)

	mockReq, err := json.Marshal(gin.H{
		"name":      "",
		"parent_id": int64(0),
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(mockReq)

	req, err := http.NewRequest(http.MethodPost, "/test/bundles/create", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockBS.AssertExpectations(t)
}

func TestCreateStoreErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "Foo@Bar.com",
		Permission: domain.MEMBER,
	}

	mockBundle := &domain.Bundle{
		Name:     "Foobar",
		ParentID: 0,
	}

	mockErr := domain.NewInternalErr()
	mockBS := new(mocks.MockBundleService)
	mockBS.
		On("Store", mock.AnythingOfType("*context.emptyCtx"), mockBundle, mockUser).
		Return(mockErr)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	group := router.Group("test")
	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := new(mocks.MockMiddlewareHandler)
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	Initialize(group, mockBS, mockMWH)

	mockReq, err := json.Marshal(gin.H{
		"id":        int64(1),
		"name":      "Foobar",
		"parent_id": int64(0),
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(mockReq)

	req, err := http.NewRequest(http.MethodPost, "/test/bundles/create", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockBS.AssertExpectations(t)
}

func TestCreateNoContext(t *testing.T) {

	mockBS := new(mocks.MockBundleService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	group := router.Group("test")

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Next()
	}
	mockMWH := new(mocks.MockMiddlewareHandler)
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	Initialize(group, mockBS, mockMWH)

	mockReq, err := json.Marshal(gin.H{
		"id":        int64(1),
		"name":      "Foobar",
		"parent_id": int64(0),
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(mockReq)

	req, err := http.NewRequest(http.MethodPost, "/test/bundles/create", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockBS.AssertExpectations(t)
}

func TestCreateNotAuth(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "Foo@Bar.com",
		Permission: domain.GUEST,
	}
	mockBundle := &domain.Bundle{
		Name:     "Foobar",
		ParentID: 0,
	}

	mockErr := domain.NewNotAuthorizedErr("")
	mockBS := new(mocks.MockBundleService)
	mockBS.
		On("Store", mock.AnythingOfType("*context.emptyCtx"), mockBundle, mockUser).
		Return(mockErr)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	group := router.Group("test")

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := new(mocks.MockMiddlewareHandler)
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	Initialize(group, mockBS, mockMWH)

	mockReq, err := json.Marshal(gin.H{
		"id":        int64(1),
		"name":      "Foobar",
		"parent_id": int64(0),
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(mockReq)

	req, err := http.NewRequest(http.MethodPost, "/test/bundles/create", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockBS.AssertExpectations(t)
}
