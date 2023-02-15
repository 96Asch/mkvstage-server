package userhandler

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

func TestGetAllCorrect(t *testing.T) {

	mockUsers := &[]domain.User{
		{
			ID:           1,
			FirstName:    "Foo",
			LastName:     "Foo",
			Email:        "Foo@Foo.com",
			Permission:   domain.ADMIN,
			ProfileColor: "FFFFFF",
		},
		{
			ID:           2,
			FirstName:    "Bar",
			LastName:     "Bar",
			Email:        "Bar@Bar.com",
			Permission:   domain.GUEST,
			ProfileColor: "FFFFF0",
		},
	}

	mockTS := new(mocks.MockTokenService)
	mockUS := new(mocks.MockUserService)
	mockUS.On("FetchAll", mock.Anything).Return(mockUsers, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	group := router.Group("test")
	Initialize(group, mockUS, mockTS)

	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/test/users/", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	expectedRes, err := json.Marshal(gin.H{"users": mockUsers})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedRes, w.Body.Bytes())
	mockUS.AssertExpectations(t)
}

func TestGetAllInternalErr(t *testing.T) {

	expectedErr := domain.NewInternalErr()

	mockTS := new(mocks.MockTokenService)
	mockUS := new(mocks.MockUserService)
	mockUS.On("FetchAll", mock.Anything).Return(nil, expectedErr)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	group := router.Group("test")
	Initialize(group, mockUS, mockTS)

	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/test/users/", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	expectedRes, err := json.Marshal(gin.H{"error": expectedErr})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, expectedRes, w.Body.Bytes())
	mockUS.AssertExpectations(t)
}
