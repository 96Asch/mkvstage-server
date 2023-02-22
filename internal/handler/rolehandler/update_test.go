package rolehandler

import (
	"bytes"
	"encoding/json"
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

func TestUpdateByIDCorrect(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
	}

	rid := int64(1)
	mockRole := &domain.Role{
		ID:          rid,
		Name:        "Foo",
		Description: "Bar",
	}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)
	mockRS := &mocks.MockRoleService{}
	mockRS.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), mockRole, mockUser).
		Return(nil)

	byteBody, err := json.Marshal(gin.H{
		"name":        "Foo",
		"description": "Bar",
	})
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)

	router := gin.New()
	reqBody := bytes.NewReader(byteBody)
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/roles/%d/update", rid), reqBody)
	assert.NoError(t, err)

	Initialize(&router.RouterGroup, mockRS, mockMWH)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRole.ID = 1

	expectedBytes, err := json.Marshal(gin.H{"role": mockRole})
	assert.NoError(t, err)
	assert.Equal(t, expectedBytes, w.Body.Bytes())
	mockRS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDNoContext(t *testing.T) {
	sid := int64(1)
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

	byteBody, err := json.Marshal(gin.H{
		"name":        "Foo",
		"description": "Bar",
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(byteBody)

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/roles/%d/update", sid), bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDInvalidParam(t *testing.T) {
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
	mockRS := &mocks.MockRoleService{}

	byteBody, err := json.Marshal(gin.H{
		"name":        "Foo",
		"description": "Bar",
	})
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)

	router := gin.New()
	reqBody := bytes.NewReader(byteBody)
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, "/roles/a/update", reqBody)
	assert.NoError(t, err)

	Initialize(&router.RouterGroup, mockRS, mockMWH)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockRS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDBindErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "Foo@Bar.com",
		Permission: domain.GUEST,
	}
	sid := int64(1)
	mockRS := &mocks.MockRoleService{}

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
	byteBody, err := json.Marshal(gin.H{
		"name":         "Foo",
		"descriptions": "Bar",
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(byteBody)

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/roles/%d/update", sid), bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockRS.AssertExpectations(t)
}

func TestUpdateByIDUpdateErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "Foo@Bar.com",
		Permission: domain.GUEST,
	}
	rid := int64(1)
	mockRole := &domain.Role{
		ID:          rid,
		Name:        "Foo",
		Description: "Bar",
	}
	mockErr := domain.NewBadRequestErr("")
	mockRS := &mocks.MockRoleService{}
	mockRS.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), mockRole, mockUser).
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
	byteBody, err := json.Marshal(gin.H{
		"name":        "Foo",
		"description": "Bar",
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(byteBody)

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/roles/%d/update", rid), bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, domain.Status(mockErr), w.Code)
	mockRS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}
