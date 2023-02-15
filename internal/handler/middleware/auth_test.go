package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthenticateUserCorrect(t *testing.T) {
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
		On("ExtractUser", mock.AnythingOfType("*context.emptyCtx"), mockAccess).
		Return(mockUser, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()
	ctx := gin.CreateTestContextOnly(w, router)

	gmh := NewGinMiddlewareHandler(mockTS)

	group := router.Group("test")
	group.POST("/auth", gmh.AuthenticateUser())

	var err error
	ctx.Request, err = http.NewRequest(http.MethodPost, "/test/auth", nil)
	assert.NoError(t, err)

	ctx.Request.Header.Add("Authorization", mockAccess)

	router.ServeHTTP(w, ctx.Request)

}
