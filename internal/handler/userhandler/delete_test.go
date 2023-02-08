package userhandler

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

func TestDeleteCorrect(t *testing.T) {
	var deleteID int64 = 1
	mockUser := &domain.User{
		ID:         deleteID,
		Permission: domain.GUEST,
	}

	mockUS := new(mocks.MockUserService)
	mockUS.On("Remove", mock.AnythingOfType("*context.emptyCtx"), mockUser, deleteID).Return(nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
	})
	w := httptest.NewRecorder()

	group := router.Group("test")
	Initialize(group, mockUS)

	mockByte, err := json.Marshal(gin.H{
		"id": deleteID,
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(mockByte)

	req, err := http.NewRequest(http.MethodDelete, "/test/users/me/delete", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	mockUS.AssertExpectations(t)
}

func TestDeleteNoContext(t *testing.T) {
	mockUS := new(mocks.MockUserService)

	r := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	router := gin.New()

	group := router.Group("test")
	Initialize(group, mockUS)

	req, err := http.NewRequest(http.MethodDelete, "/test/users/me/delete", nil)
	assert.NoError(t, err)

	router.ServeHTTP(r, req)

	assert.Equal(t, http.StatusInternalServerError, r.Code)
	mockUS.AssertNotCalled(t, "Remove")
}

func TestDeleteInvalidBind(t *testing.T) {
	var deleteID int64 = 1
	mockUser := &domain.User{
		ID:         deleteID,
		Permission: domain.GUEST,
	}

	mockUS := new(mocks.MockUserService)
	mockUS.On("Remove", mock.AnythingOfType("*context.emptyCtx"), mockUser, int64(0)).Return(nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
	})
	w := httptest.NewRecorder()

	group := router.Group("test")
	Initialize(group, mockUS)

	mockByte, err := json.Marshal(gin.H{
		"ids": deleteID,
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(mockByte)

	req, err := http.NewRequest(http.MethodDelete, "/test/users/me/delete", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUS.AssertNotCalled(t, "Remove")
}

func TestDeleteRemoveErr(t *testing.T) {
	var deleteID int64 = 1
	mockUser := &domain.User{
		ID:         deleteID,
		Permission: domain.GUEST,
	}

	expectedErr := domain.NewRecordNotFoundErr("", "")
	mockUS := new(mocks.MockUserService)
	mockUS.On("Remove", mock.AnythingOfType("*context.emptyCtx"), mockUser, deleteID).Return(expectedErr)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
	})
	w := httptest.NewRecorder()

	group := router.Group("test")
	Initialize(group, mockUS)

	mockByte, err := json.Marshal(gin.H{
		"id": deleteID,
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(mockByte)

	req, err := http.NewRequest(http.MethodDelete, "/test/users/me/delete", bodyReader)
	req.Header.Set("Content-Type", "application/xml")
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	expectedBody, err := json.Marshal(gin.H{"error": expectedErr})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, expectedBody, w.Body.Bytes())
	mockUS.AssertNotCalled(t, "Remove")
}
