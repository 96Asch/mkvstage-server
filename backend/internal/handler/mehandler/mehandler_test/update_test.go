package mehandler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/96Asch/mkvstage-server/backend/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/backend/internal/handler/mehandler"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
)

func prepareAndServeUpdate(
	t *testing.T,
	mockUS domain.UserService,
	mockTS domain.TokenService,
	mockMWH domain.MiddlewareHandler,
	body *[]byte,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	mehandler.Initialize(&router.RouterGroup, mockUS, mockTS, mockMWH)

	requestBody := bytes.NewReader(*body)

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodPut, "/me/update", requestBody)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestUpdateCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID: 1,
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}

	mockUS.On("Update", context.TODO(), &domain.User{
		ID:           1,
		FirstName:    "Foob",
		LastName:     "Bars",
		Permission:   0,
		ProfileColor: "FFFFFF",
	}).Return(nil)

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(gin.H{
		"first_name":    "Foob",
		"last_name":     "Bars",
		"password":      "FooBar",
		"profile_color": "FFFFFF",
	})
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, mockUS, mockTS, mockMWH, &byteBody)
	assert.Equal(t, http.StatusOK, writer.Code)
	mockMWH.AssertExpectations(t)
	mockTS.AssertExpectations(t)
	mockUS.AssertExpectations(t)
}

func TestUpdateInvalidBind(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID: 1,
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(gin.H{
		"first_names":  "Foo",
		"last_name":    "Bar",
		"password":     "FooBar",
		"permission":   "member",
		"profileColor": "FFFFFF",
	})
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, mockUS, mockTS, mockMWH, &byteBody)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	mockMWH.AssertExpectations(t)
	mockTS.AssertExpectations(t)
	mockUS.AssertNotCalled(t, "Update")
}

func TestUpdateNoContext(t *testing.T) {
	t.Parallel()

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(gin.H{
		"first_name":   "Foo",
		"last_name":    "Bar",
		"password":     "FooBar",
		"permission":   "member",
		"profileColor": "FFFFFF",
	})
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, mockUS, mockTS, mockMWH, &byteBody)
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	mockMWH.AssertExpectations(t)
	mockTS.AssertExpectations(t)
	mockUS.AssertNotCalled(t, "Update")
}

func TestUpdateInternalErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID: 1,
	}

	expectedError := domain.NewInternalErr()
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}

	mockUS.On("Update", context.TODO(), &domain.User{
		ID:           1,
		FirstName:    "Foob",
		LastName:     "Bars",
		Permission:   0,
		ProfileColor: "FFFFFF",
	}).Return(expectedError)

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(gin.H{
		"first_name":    "Foob",
		"last_name":     "Bars",
		"password":      "FooBar",
		"permission":    "member",
		"profile_color": "FFFFFF",
	})
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, mockUS, mockTS, mockMWH, &byteBody)

	expJSON, err := json.Marshal(gin.H{"error": expectedError})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Equal(t, expJSON, writer.Body.Bytes())
	mockMWH.AssertExpectations(t)
	mockTS.AssertExpectations(t)
	mockUS.AssertExpectations(t)
}
