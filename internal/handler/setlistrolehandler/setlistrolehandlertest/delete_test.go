package setlistrolehandler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/handler/setlistrolehandler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func prepareAndServeDelete(
	t *testing.T,
	mockSLRS domain.SetlistRoleService,
	mockMWH domain.MiddlewareHandler,
	param string,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	setlistrolehandler.Initialize(&router.RouterGroup, mockSLRS, mockMWH)

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodDelete,
		fmt.Sprintf("/setlistroles%s", param),
		nil,
	)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestDelete(t *testing.T) {
	user := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	t.Run("Correct", func(t *testing.T) {
		t.Parallel()

		mockSLRS := &mocks.MockSetlistRoleService{}
		mockMWH := &mocks.MockMiddlewareHandler{}

		var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
			ctx.Set("user", user)
			ctx.Next()
		}

		mockMWH.
			On("AuthenticateUser").
			Return(mockAuthHF)

		mockSLRS.
			On("Remove", mock.AnythingOfType("*context.emptyCtx"), []int64{1, 2, 3, 4}, user).
			Return(nil)

		writer := prepareAndServeDelete(t, mockSLRS, mockMWH, "?ids=1,2,3,4")

		assert.Equal(t, http.StatusOK, writer.Code)
		assert.Len(t, writer.Body.Bytes(), 0)

		mockSLRS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})

	t.Run("Correct no ids given", func(t *testing.T) {
		t.Parallel()

		mockSLRS := &mocks.MockSetlistRoleService{}
		mockMWH := &mocks.MockMiddlewareHandler{}

		var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
			ctx.Set("user", user)
			ctx.Next()
		}

		mockMWH.
			On("AuthenticateUser").
			Return(mockAuthHF)

		writer := prepareAndServeDelete(t, mockSLRS, mockMWH, "?ids=")

		assert.Equal(t, http.StatusOK, writer.Code)
		assert.Len(t, writer.Body.Bytes(), 0)

		mockSLRS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})

	t.Run("Fail no user context", func(t *testing.T) {
		t.Parallel()

		expErr := domain.NewInternalErr()
		mockSLRS := &mocks.MockSetlistRoleService{}
		mockMWH := &mocks.MockMiddlewareHandler{}

		var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		}

		mockMWH.
			On("AuthenticateUser").
			Return(mockAuthHF)

		writer := prepareAndServeDelete(t, mockSLRS, mockMWH, "?ids=1,2,3,4")

		expBody, err := json.Marshal(gin.H{
			"error": expErr.Error(),
		})

		assert.NoError(t, err)
		assert.Equal(t, expErr.Status(), writer.Code)
		assert.Equal(t, expBody, writer.Body.Bytes())

		mockSLRS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})

	t.Run("Fail to cast user", func(t *testing.T) {
		t.Parallel()

		expErr := domain.NewInternalErr()
		mockSLRS := &mocks.MockSetlistRoleService{}
		mockMWH := &mocks.MockMiddlewareHandler{}

		var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
			ctx.Set("user", 0)
			ctx.Next()
		}

		mockMWH.
			On("AuthenticateUser").
			Return(mockAuthHF)

		writer := prepareAndServeDelete(t, mockSLRS, mockMWH, "?ids=1,2,3,4")

		expBody, err := json.Marshal(gin.H{
			"error": expErr.Error(),
		})

		assert.NoError(t, err)
		assert.Equal(t, expErr.Status(), writer.Code)
		assert.Equal(t, expBody, writer.Body.Bytes())

		mockSLRS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})

	t.Run("Fail parse id error", func(t *testing.T) {
		t.Parallel()

		expErr := domain.NewBadRequestErr("Could not parse a to a number")
		mockSLRS := &mocks.MockSetlistRoleService{}
		mockMWH := &mocks.MockMiddlewareHandler{}

		var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
			ctx.Set("user", user)
			ctx.Next()
		}

		mockMWH.
			On("AuthenticateUser").
			Return(mockAuthHF)

		writer := prepareAndServeDelete(t, mockSLRS, mockMWH, "?ids=1,a,3,4")

		expBody, err := json.Marshal(gin.H{
			"error": expErr.Error(),
		})

		assert.NoError(t, err)
		assert.Equal(t, expErr.Status(), writer.Code)
		assert.Equal(t, expBody, writer.Body.Bytes())

		mockSLRS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})

	t.Run("Fail SetlistRole Remove error", func(t *testing.T) {
		t.Parallel()

		expErr := domain.NewInternalErr()
		mockSLRS := &mocks.MockSetlistRoleService{}
		mockMWH := &mocks.MockMiddlewareHandler{}

		var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
			ctx.Set("user", user)
			ctx.Next()
		}

		mockMWH.
			On("AuthenticateUser").
			Return(mockAuthHF)

		mockSLRS.
			On("Remove", mock.AnythingOfType("*context.emptyCtx"), []int64{1, 2, 3, 4}, user).
			Return(expErr)

		writer := prepareAndServeDelete(t, mockSLRS, mockMWH, "?ids=1,2,3,4")

		expBody, err := json.Marshal(gin.H{
			"error": expErr.Error(),
		})

		assert.NoError(t, err)
		assert.Equal(t, expErr.Status(), writer.Code)
		assert.Equal(t, expBody, writer.Body.Bytes())

		mockSLRS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})
}
