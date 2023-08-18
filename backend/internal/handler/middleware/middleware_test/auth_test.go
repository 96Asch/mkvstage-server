package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/handler/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func prepareAndServeAuthenticate(t *testing.T, mockUS domain.UserService, mockTS domain.TokenService, token string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()
	gmh := middleware.NewGinMiddlewareHandler(mockUS, mockTS)

	router.POST("/auth", gmh.AuthenticateUser())

	ctx := gin.CreateTestContextOnly(writer, router)

	var err error
	ctx.Request, err = http.NewRequestWithContext(ctx, http.MethodPost, "/auth", nil)
	assert.NoError(t, err)

	ctx.Request.Header.Add("Authorization", token)

	router.ServeHTTP(writer, ctx.Request)

	return writer
}

func TestAuthenticateUserCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Permission:   domain.GUEST,
		ProfileColor: "FFFFFF",
	}

	mockAccess := "access-token"

	mockUS := &mocks.MockUserService{}
	mockTS := &mocks.MockTokenService{}

	mockUS.
		On("FetchByEmail", mock.AnythingOfType("*gin.Context"), mockUser.Email).
		Return(mockUser, nil)

	mockTS.
		On("ExtractEmail", mock.AnythingOfType("*gin.Context"), mockAccess).
		Return(mockUser.Email, nil)

	writer := prepareAndServeAuthenticate(t, mockUS, mockTS, mockAccess)
	assert.Equal(t, http.StatusOK, writer.Code)
	mockTS.AssertExpectations(t)
}

func TestAuthenticateUserExtractErr(t *testing.T) {
	t.Parallel()

	expErr := domain.NewBadRequestErr("")
	mockAccess := "access-token"
	mockUS := &mocks.MockUserService{}
	mockTS := &mocks.MockTokenService{}

	mockTS.
		On("ExtractEmail", mock.AnythingOfType("*gin.Context"), mockAccess).
		Return(nil, expErr)

	writer := prepareAndServeAuthenticate(t, mockUS, mockTS, mockAccess)
	assert.Equal(t, domain.Status(expErr), writer.Code)

	expBody, err := json.Marshal(gin.H{"error": expErr})
	assert.NoError(t, err)

	assert.Equal(t, expBody, writer.Body.Bytes())
	mockTS.AssertExpectations(t)
}

func TestAuthenticateUserFetchByEmailErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Permission:   domain.GUEST,
		ProfileColor: "FFFFFF",
	}

	expErr := domain.NewBadRequestErr("")
	mockAccess := "access-token"
	mockUS := &mocks.MockUserService{}
	mockTS := &mocks.MockTokenService{}

	mockUS.
		On("FetchByEmail", mock.AnythingOfType("*gin.Context"), mockUser.Email).
		Return(nil, expErr)

	mockTS.
		On("ExtractEmail", mock.AnythingOfType("*gin.Context"), mockAccess).
		Return(mockUser.Email, nil)

	writer := prepareAndServeAuthenticate(t, mockUS, mockTS, mockAccess)
	assert.Equal(t, domain.Status(expErr), writer.Code)

	expBody, err := json.Marshal(gin.H{"error": expErr})
	assert.NoError(t, err)

	assert.Equal(t, expBody, writer.Body.Bytes())
	mockTS.AssertExpectations(t)
}
