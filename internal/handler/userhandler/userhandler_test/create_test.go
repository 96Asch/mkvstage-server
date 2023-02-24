package userhandler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/handler/userhandler"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func prepareAndServeCreate(
	t *testing.T,
	mockUS domain.UserService,
	mockTS domain.TokenService,
	body *[]byte,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	userhandler.Initialize(&router.RouterGroup, mockUS, mockTS)

	requestBody := bytes.NewReader(*body)

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, "/users/create", requestBody)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestCreateCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Password:     "FooBar",
		Email:        "Foo@Bar.com",
		Permission:   domain.GUEST,
		ProfileColor: "FFFFFF",
	}

	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}

	mockUS.
		On("Store", mock.AnythingOfType("*context.emptyCtx"), mockUser).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg, ok := args.Get(1).(*domain.User)
			assert.True(t, ok)
			arg.ID = 1
		})

	byteBody, err := json.Marshal(gin.H{
		"first_name":    "Foo",
		"last_name":     "Bar",
		"email":         "Foo@Bar.com",
		"password":      "FooBar",
		"profile_color": "FFFFFF",
	})
	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockUS, mockTS, &byteBody)

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

	assert.Equal(t, http.StatusCreated, writer.Code)
	assert.Equal(t, expectedUser, writer.Body.Bytes())
	mockUS.AssertExpectations(t)
}

func TestCreateInvalidEmail(t *testing.T) {
	t.Parallel()

	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}

	byteBody, err := json.Marshal(gin.H{
		"first_name":    "Foo",
		"last_name":     "Bar",
		"email":         "Foocom",
		"password":      "Foo",
		"profile_color": "Foo",
	})
	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockUS, mockTS, &byteBody)

	assert.Equal(t, http.StatusBadRequest, writer.Code)
	mockUS.AssertNotCalled(t, "Store")
}

func TestCreateInvalidBind(t *testing.T) {
	t.Parallel()

	mockUS := &mocks.MockUserService{}
	mockTS := &mocks.MockTokenService{}

	byteBody, err := json.Marshal(gin.H{
		"first_name":   "Foo",
		"last_name":    "Bar",
		"emails":       "Foo@bar.com",
		"password":     "Foo",
		"profileColor": "Foo",
	})
	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockUS, mockTS, &byteBody)

	assert.Equal(t, http.StatusBadRequest, writer.Code)
	mockUS.AssertNotCalled(t, "Create")
}

func TestCreateServerError(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Password:     "FooBar",
		Permission:   domain.GUEST,
		ProfileColor: "FFFFFF",
	}

	expectedErr := domain.NewInternalErr()
	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}

	mockUS.On("Store", mock.Anything, mockUser).
		Return(expectedErr)

	byteBody, err := json.Marshal(gin.H{
		"first_name":    "Foo",
		"last_name":     "Bar",
		"email":         "Foo@Bar.com",
		"password":      "FooBar",
		"permission":    "member",
		"profile_color": "FFFFFF",
	})
	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockUS, mockTS, &byteBody)

	expectedRes, err := json.Marshal(gin.H{
		"error": expectedErr,
	})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Equal(t, expectedRes, writer.Body.Bytes())
	mockUS.AssertExpectations(t)
}
