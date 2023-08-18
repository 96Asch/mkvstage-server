package mehandler_test

import (
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
	"github.com/stretchr/testify/mock"
)

func prepareAndServeMe(
	t *testing.T,
	mockUS domain.UserService,
	mockTS domain.TokenService,
	mockMWH domain.MiddlewareHandler,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	mehandler.Initialize(&router.RouterGroup, mockUS, mockTS, mockMWH)

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, "/me", nil)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestMeCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:           1,
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Permission:   domain.MEMBER,
		ProfileColor: "FFFFFF",
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}

	mockUS.
		On("FetchByID", mock.Anything, mockUser.ID).
		Return(mockUser, nil)

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	writer := prepareAndServeMe(t, mockUS, mockTS, mockMWH)
	expectBody, err := json.Marshal(gin.H{"user": mockUser})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, expectBody, writer.Body.Bytes())
	mockMWH.AssertExpectations(t)
	mockTS.AssertExpectations(t)
	mockUS.AssertExpectations(t)
}

func TestMeNoContext(t *testing.T) {
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

	writer := prepareAndServeMe(t, mockUS, mockTS, mockMWH)
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	mockMWH.AssertExpectations(t)
	mockTS.AssertExpectations(t)
	mockUS.AssertNotCalled(t, "FindByID")
}

func TestMeNotFound(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:           1,
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Permission:   domain.GUEST,
		ProfileColor: "FFFFFF",
	}

	expErr := domain.NewRecordNotFoundErr("id", "1")
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}
	mockUS.On("FetchByID", mock.Anything, mockUser.ID).Return(nil, expErr)

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	writer := prepareAndServeMe(t, mockUS, mockTS, mockMWH)

	expectBody, err := json.Marshal(gin.H{"error": expErr})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, writer.Code)
	assert.Equal(t, expectBody, writer.Body.Bytes())
	mockMWH.AssertExpectations(t)
	mockTS.AssertExpectations(t)
	mockUS.AssertExpectations(t)
}
