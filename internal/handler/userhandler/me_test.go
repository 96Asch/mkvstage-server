package userhandler

import (
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

func TestMeCorrect(t *testing.T) {
	mockUser := &domain.User{
		ID:           1,
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Password:     "FooBar",
		Permission:   "member",
		ProfileColor: "FFFFFF",
	}

	mockUS := new(mocks.MockUserService)
	mockUS.On("FetchByID", mock.Anything, mockUser.ID).Return(mockUser, nil)

	r := httptest.NewRecorder()

	router := gin.Default()
	router.Use(func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
	})

	group := router.Group("test")
	Initialize(group, mockUS)

	req, err := http.NewRequest(http.MethodGet, "/test/users/me", nil)
	assert.NoError(t, err)

	router.ServeHTTP(r, req)

	expectBody, err := json.Marshal(gin.H{"user": mockUser})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, expectBody, r.Body.Bytes())
	mockUS.AssertExpectations(t)
}

func TestMeNoContext(t *testing.T) {
	mockUS := new(mocks.MockUserService)

	r := httptest.NewRecorder()

	router := gin.Default()

	group := router.Group("test")
	Initialize(group, mockUS)

	req, err := http.NewRequest(http.MethodGet, "/test/users/me", nil)
	assert.NoError(t, err)

	router.ServeHTTP(r, req)

	assert.Equal(t, http.StatusInternalServerError, r.Code)
	mockUS.AssertNotCalled(t, "FindByID")
}

func TestMeNotFound(t *testing.T) {
	mockUser := &domain.User{
		ID:           1,
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Password:     "FooBar",
		Permission:   "member",
		ProfileColor: "FFFFFF",
	}

	mockUS := new(mocks.MockUserService)
	notFoundErr := domain.NewRecordNotFoundErr("id", "1")
	mockUS.On("FetchByID", mock.Anything, mockUser.ID).Return(nil, notFoundErr)

	r := httptest.NewRecorder()

	router := gin.Default()
	router.Use(func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
	})

	group := router.Group("test")
	Initialize(group, mockUS)

	req, err := http.NewRequest(http.MethodGet, "/test/users/me", nil)
	assert.NoError(t, err)

	router.ServeHTTP(r, req)

	expectBody, err := json.Marshal(gin.H{"error": notFoundErr})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, r.Code)
	assert.Equal(t, expectBody, r.Body.Bytes())
	mockUS.AssertExpectations(t)
}
