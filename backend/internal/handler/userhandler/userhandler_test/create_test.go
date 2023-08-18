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
	mockMWH domain.MiddlewareHandler,
	body *[]byte,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	userhandler.Initialize(&router.RouterGroup, mockUS, mockMWH)

	requestBody := bytes.NewReader(*body)

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, "/users/create", requestBody)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestCreate(t *testing.T) {
	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Email:        "Foo@Bar.com",
		Permission:   domain.GUEST,
		ProfileColor: "FFFFFF",
	}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("email", "Foo@Bar.com")
		ctx.Next()
	}

	t.Run("Correct", func(t *testing.T) {
		t.Parallel()
		mockMWH := &mocks.MockMiddlewareHandler{}
		mockUS := &mocks.MockUserService{}

		mockUS.
			On("Store", context.TODO(), mockUser).
			Return(nil).
			Run(func(args mock.Arguments) {
				arg, ok := args.Get(1).(*domain.User)
				assert.True(t, ok)
				arg.ID = 1
			})

		mockMWH.
			On("JWTExtractEmail").
			Return(mockAuthHF)

		byteBody, err := json.Marshal(gin.H{
			"first_name":    "Foo",
			"last_name":     "Bar",
			"profile_color": "FFFFFF",
		})
		assert.NoError(t, err)

		writer := prepareAndServeCreate(t, mockUS, mockMWH, &byteBody)

		expbody, err := json.Marshal(gin.H{
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
		assert.Equal(t, expbody, writer.Body.Bytes())
		mockUS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})

	t.Run("Fail no context", func(t *testing.T) {
		t.Parallel()

		expErr := domain.NewInternalErr()

		var emptyMockAuth gin.HandlerFunc = func(ctx *gin.Context) {}
		mockMWH := &mocks.MockMiddlewareHandler{}
		mockUS := &mocks.MockUserService{}

		mockMWH.
			On("JWTExtractEmail").
			Return(emptyMockAuth)

		byteBody, err := json.Marshal(gin.H{
			"first_name":    "Foo",
			"last_name":     "Bar",
			"profile_color": "FFFFFF",
		})
		assert.NoError(t, err)

		writer := prepareAndServeCreate(t, mockUS, mockMWH, &byteBody)

		expbody, err := json.Marshal(gin.H{
			"error": expErr.Message,
		})
		assert.NoError(t, err)

		assert.Equal(t, expErr.Status(), writer.Code)
		assert.Equal(t, expbody, writer.Body.Bytes())
		mockUS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})

	t.Run("Fail email conversion err", func(t *testing.T) {
		t.Parallel()

		expErr := domain.NewInternalErr()
		mockMWH := &mocks.MockMiddlewareHandler{}
		mockUS := &mocks.MockUserService{}

		var wrongMockAuth gin.HandlerFunc = func(ctx *gin.Context) {
			ctx.Set("email", mockUser)
			ctx.Next()
		}

		mockMWH.
			On("JWTExtractEmail").
			Return(wrongMockAuth)

		byteBody, err := json.Marshal(gin.H{
			"first_name":    "Foo",
			"last_name":     "Bar",
			"profile_color": "FFFFFF",
		})
		assert.NoError(t, err)

		writer := prepareAndServeCreate(t, mockUS, mockMWH, &byteBody)

		expbody, err := json.Marshal(gin.H{
			"error": expErr.Message,
		})
		assert.NoError(t, err)

		assert.Equal(t, expErr.Status(), writer.Code)
		assert.Equal(t, expbody, writer.Body.Bytes())
		mockUS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})

	t.Run("Fail invalid bind", func(t *testing.T) {
		t.Parallel()

		expErr := domain.NewBadRequestErr("field 'first_name' is required")
		mockMWH := &mocks.MockMiddlewareHandler{}
		mockUS := &mocks.MockUserService{}

		mockMWH.
			On("JWTExtractEmail").
			Return(mockAuthHF)

		byteBody, err := json.Marshal(gin.H{
			"first_names":   "Foo",
			"last_name":     "Bar",
			"profile_color": "Foo",
		})
		assert.NoError(t, err)

		writer := prepareAndServeCreate(t, mockUS, mockMWH, &byteBody)

		expbody, err := json.Marshal(gin.H{
			"error": expErr.Message,
		})
		assert.NoError(t, err)

		assert.Equal(t, expErr.Status(), writer.Code)
		assert.Equal(t, expbody, writer.Body.Bytes())
		mockUS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})

	t.Run("Fail createUser error", func(t *testing.T) {
		t.Parallel()
		expErr := domain.NewInternalErr()
		mockMWH := &mocks.MockMiddlewareHandler{}
		mockUS := &mocks.MockUserService{}

		mockMWH.
			On("JWTExtractEmail").
			Return(mockAuthHF)

		mockUS.On("Store", mock.Anything, mockUser).
			Return(expErr)

		byteBody, err := json.Marshal(gin.H{
			"first_name":    "Foo",
			"last_name":     "Bar",
			"profile_color": "FFFFFF",
		})
		assert.NoError(t, err)

		writer := prepareAndServeCreate(t, mockUS, mockMWH, &byteBody)

		expbody, err := json.Marshal(gin.H{
			"error": expErr.Message,
		})
		assert.NoError(t, err)

		assert.Equal(t, expErr.Status(), writer.Code)
		assert.Equal(t, expbody, writer.Body.Bytes())
		mockUS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})
}
