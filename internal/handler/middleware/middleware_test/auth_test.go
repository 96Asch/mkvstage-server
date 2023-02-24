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

func prepareAndServeAuthenticate(t *testing.T, mockTS domain.TokenService, token string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()
	gmh := middleware.NewGinMiddlewareHandler(mockTS)

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
		Password:     "FooBar",
		Email:        "Foo@Bar.com",
		Permission:   domain.GUEST,
		ProfileColor: "FFFFFF",
	}

	mockAccess := "access-token"

	mockTS := new(mocks.MockTokenService)
	mockTS.
		On("ExtractUser", mock.AnythingOfType("*gin.Context"), mockAccess).
		Return(mockUser, nil)

	writer := prepareAndServeAuthenticate(t, mockTS, mockAccess)
	assert.Equal(t, http.StatusOK, writer.Code)
	mockTS.AssertExpectations(t)
}

func TestAuthenticateUserExtractErr(t *testing.T) {
	t.Parallel()

	expErr := domain.NewBadRequestErr("")
	mockAccess := "access-token"
	mockTS := &mocks.MockTokenService{}

	mockTS.
		On("ExtractUser", mock.AnythingOfType("*gin.Context"), mockAccess).
		Return(nil, expErr)

	writer := prepareAndServeAuthenticate(t, mockTS, mockAccess)
	assert.Equal(t, domain.Status(expErr), writer.Code)

	expBody, err := json.Marshal(gin.H{"error": expErr})
	assert.NoError(t, err)

	assert.Equal(t, expBody, writer.Body.Bytes())
	mockTS.AssertExpectations(t)
}
