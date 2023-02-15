package mehandler

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
		ID: 1,
	}

	mockTS := new(mocks.MockTokenService)
	mockUS := new(mocks.MockUserService)
	mockUS.On("Update", mock.AnythingOfType("*context.emptyCtx"), &domain.User{
		ID:           1,
		FirstName:    "Foob",
		LastName:     "Bars",
		Password:     "FooBar",
		Permission:   0,
		ProfileColor: "FFFFFF",
	}).Return(nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	w := httptest.NewRecorder()

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := new(mocks.MockMiddlewareHandler)
	mockMWH.On("AuthorizeUser").Return(mockAuthHF)

	group := router.Group("test")
	Initialize(group, mockUS, mockTS, mockMWH)

	mockByte, err := json.Marshal(gin.H{
		"first_name":    "Foob",
		"last_name":     "Bars",
		"password":      "FooBar",
		"profile_color": "FFFFFF",
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(mockByte)

	req, err := http.NewRequest(http.MethodPatch, "/test/me/update", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUS.AssertExpectations(t)
}

func TestUpdateInvalidBind(t *testing.T) {

	mockUser := &domain.User{
		ID: 1,
	}

	mockTS := new(mocks.MockTokenService)
	mockUS := new(mocks.MockUserService)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	w := httptest.NewRecorder()

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := new(mocks.MockMiddlewareHandler)
	mockMWH.On("AuthorizeUser").Return(mockAuthHF)

	group := router.Group("test")
	Initialize(group, mockUS, mockTS, mockMWH)

	mockByte, err := json.Marshal(gin.H{
		"first_names":  "Foo",
		"last_name":    "Bar",
		"password":     "FooBar",
		"permission":   "member",
		"profileColor": "FFFFFF",
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(mockByte)

	req, err := http.NewRequest(http.MethodPatch, "/test/me/update", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUS.AssertNotCalled(t, "Update")
}

func TestUpdateNoContext(t *testing.T) {
	mockUS := new(mocks.MockUserService)
	mockTS := new(mocks.MockTokenService)

	r := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	router := gin.New()

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Next()
	}
	mockMWH := new(mocks.MockMiddlewareHandler)
	mockMWH.On("AuthorizeUser").Return(mockAuthHF)

	group := router.Group("test")
	Initialize(group, mockUS, mockTS, mockMWH)

	req, err := http.NewRequest(http.MethodPatch, "/test/me/update", nil)
	assert.NoError(t, err)

	router.ServeHTTP(r, req)

	assert.Equal(t, http.StatusInternalServerError, r.Code)
	mockUS.AssertNotCalled(t, "Update")
}

func TestUpdateInternalErr(t *testing.T) {
	mockUser := &domain.User{
		ID: 1,
	}

	expectedError := domain.NewInternalErr()

	mockTS := new(mocks.MockTokenService)
	mockUS := new(mocks.MockUserService)
	mockUS.On("Update", mock.AnythingOfType("*context.emptyCtx"), &domain.User{
		ID:           1,
		FirstName:    "Foob",
		LastName:     "Bars",
		Password:     "FooBar",
		Permission:   0,
		ProfileColor: "FFFFFF",
	}).Return(expectedError)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	w := httptest.NewRecorder()

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := new(mocks.MockMiddlewareHandler)
	mockMWH.On("AuthorizeUser").Return(mockAuthHF)

	group := router.Group("test")
	Initialize(group, mockUS, mockTS, mockMWH)

	mockByte, err := json.Marshal(gin.H{
		"first_name":    "Foob",
		"last_name":     "Bars",
		"password":      "FooBar",
		"permission":    "member",
		"profile_color": "FFFFFF",
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(mockByte)

	req, err := http.NewRequest(http.MethodPatch, "/test/me/update", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	expectedBody := gin.H{"error": expectedError}
	expectedJson, err := json.Marshal(expectedBody)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, expectedJson, w.Body.Bytes())
	mockUS.AssertExpectations(t)
}
