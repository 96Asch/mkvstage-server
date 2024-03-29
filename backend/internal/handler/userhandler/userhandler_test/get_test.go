package userhandler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/96Asch/mkvstage-server/backend/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/backend/internal/handler/userhandler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func prepareAndServeGet(
	t *testing.T,
	mockUS domain.UserService,
	mockMWH domain.MiddlewareHandler,
	param string,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	userhandler.Initialize(&router.RouterGroup, mockUS, mockMWH)

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, fmt.Sprintf("/users%s", param), nil)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestGetAll(t *testing.T) {
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

	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("email", "Foo@Bar.com")
		ctx.Next()
	}

	mockMWH.On("JWTExtractEmail").Return(mockAuthHF)

	t.Run("Correct", func(t *testing.T) {
		t.Parallel()

		mockUS := &mocks.MockUserService{}
		mockUS.
			On("FetchAll", context.TODO()).
			Return(mockUsers, nil)

		writer := prepareAndServeGet(t, mockUS, mockMWH, "")

		expectedRes, err := json.Marshal(gin.H{"users": mockUsers})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, writer.Code)
		assert.Equal(t, expectedRes, writer.Body.Bytes())
		mockUS.AssertExpectations(t)
	})

	t.Run("Fail FetchAll Error", func(t *testing.T) {
		t.Parallel()

		expectedErr := domain.NewInternalErr()

		mockUS := &mocks.MockUserService{}
		mockUS.
			On("FetchAll", context.TODO()).
			Return(nil, expectedErr)

		writer := prepareAndServeGet(t, mockUS, mockMWH, "")

		expectedRes, err := json.Marshal(gin.H{"error": expectedErr})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, writer.Code)
		assert.Equal(t, expectedRes, writer.Body.Bytes())
		mockUS.AssertExpectations(t)
	})
}
