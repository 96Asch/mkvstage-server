package userhandler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/handler/userhandler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func prepareAndServeGet(
	t *testing.T,
	mockUS domain.UserService,
	mockTS domain.TokenService,
	param string,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	userhandler.Initialize(&router.RouterGroup, mockUS, mockTS)

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, fmt.Sprintf("/users%s", param), nil)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestGetAllCorrect(t *testing.T) {
	t.Parallel()

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

	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}

	mockUS.
		On("FetchAll", context.TODO()).
		Return(mockUsers, nil)

	writer := prepareAndServeGet(t, mockUS, mockTS, "")

	expectedRes, err := json.Marshal(gin.H{"users": mockUsers})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, expectedRes, writer.Body.Bytes())
	mockUS.AssertExpectations(t)
}

func TestGetAllInternalErr(t *testing.T) {
	t.Parallel()

	expectedErr := domain.NewInternalErr()
	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}

	mockUS.
		On("FetchAll", context.TODO()).
		Return(nil, expectedErr)

	writer := prepareAndServeGet(t, mockUS, mockTS, "")

	expectedRes, err := json.Marshal(gin.H{"error": expectedErr})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Equal(t, expectedRes, writer.Body.Bytes())
	mockUS.AssertExpectations(t)
}
