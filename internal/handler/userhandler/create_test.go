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

func TestCreateCorrect(t *testing.T) {
	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Password:     "FooBar",
		Email:        "Foo@Bar.com",
		Permission:   domain.GUEST,
		ProfileColor: "FFFFFF",
	}

	mockUS := new(mocks.MockUserService)
	mockUS.On("Store", mock.AnythingOfType("*context.emptyCtx"), mockUser).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).(*domain.User)
			arg.ID = 1
		})

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	group := router.Group("test")
	Initialize(group, mockUS)

	mockByte, err := json.Marshal(gin.H{
		"first_name":    "Foo",
		"last_name":     "Bar",
		"email":         "Foo@Bar.com",
		"password":      "FooBar",
		"profile_color": "FFFFFF",
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(mockByte)

	req, err := http.NewRequest(http.MethodPost, "/test/users/create", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	expectedUser, err := json.Marshal(gin.H{
		"user": &domain.User{
			ID:           1,
			FirstName:    "Foo",
			LastName:     "Bar",
			Email:        "Foo@Bar.com",
			Permission:   domain.GUEST,
			ProfileColor: "FFFFFF",
		},
	})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, expectedUser, w.Body.Bytes())
	mockUS.AssertExpectations(t)
}

func TestCreateInvalidEmail(t *testing.T) {
	mockUS := new(mocks.MockUserService)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	w := httptest.NewRecorder()

	group := router.Group("test")
	Initialize(group, mockUS)

	mockByte, err := json.Marshal(gin.H{
		"first_name":    "Foo",
		"last_name":     "Bar",
		"email":         "Foocom",
		"password":      "Foo",
		"profile_color": "Foo",
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(mockByte)

	req, err := http.NewRequest(http.MethodPost, "/test/users/create", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUS.AssertNotCalled(t, "Store")
}

func TestCreateInvalidBind(t *testing.T) {
	mockUS := new(mocks.MockUserService)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	w := httptest.NewRecorder()

	group := router.Group("test")
	Initialize(group, mockUS)

	mockByte, err := json.Marshal(gin.H{
		"first_name":   "Foo",
		"last_name":    "Bar",
		"emails":       "Foo@bar.com",
		"password":     "Foo",
		"profileColor": "Foo",
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(mockByte)

	req, err := http.NewRequest(http.MethodPost, "/test/users/create", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUS.AssertNotCalled(t, "Create")
}

func TestCreateServerError(t *testing.T) {
	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Password:     "FooBar",
		Permission:   domain.GUEST,
		ProfileColor: "FFFFFF",
	}

	expectedErr := domain.NewInternalErr()
	mockUS := new(mocks.MockUserService)
	mockUS.On("Store", mock.Anything, mockUser).
		Return(expectedErr)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	w := httptest.NewRecorder()

	group := router.Group("test")
	Initialize(group, mockUS)

	mockByte, err := json.Marshal(gin.H{
		"first_name":    "Foo",
		"last_name":     "Bar",
		"email":         "Foo@Bar.com",
		"password":      "FooBar",
		"permission":    "member",
		"profile_color": "FFFFFF",
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(mockByte)

	req, err := http.NewRequest(http.MethodPost, "/test/users/create", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.NoError(t, err)

	expectedRes, err := json.Marshal(gin.H{
		"error": expectedErr,
	})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, expectedRes, w.Body.Bytes())
	mockUS.AssertExpectations(t)
}
