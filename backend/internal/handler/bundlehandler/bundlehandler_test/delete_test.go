package bundlehandler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/handler/bundlehandler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func prepareAndServeDelete(
	t *testing.T,
	param string,
	mockBS domain.BundleService,
	mockMWH domain.MiddlewareHandler,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	bundlehandler.Initialize(&router.RouterGroup, mockBS, mockMWH)

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodDelete,
		fmt.Sprintf("/bundles/%s/delete", param),
		nil,
	)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestDeleteCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	bid := int64(1)
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockBS := &mocks.MockBundleService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockBS.
		On("Remove", context.TODO(), bid, mockUser).
		Return(nil)
	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	writer := prepareAndServeDelete(t, fmt.Sprint(bid), mockBS, mockMWH)

	assert.Equal(t, http.StatusAccepted, writer.Code)
	mockBS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestDeleteNoContext(t *testing.T) {
	t.Parallel()

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockBS := &mocks.MockBundleService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	writer := prepareAndServeDelete(t, "1", mockBS, mockMWH)

	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	mockBS.AssertExpectations(t)
}

func TestDeleteParamErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockBS := &mocks.MockBundleService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	writer := prepareAndServeDelete(t, "a", mockBS, mockMWH)

	assert.Equal(t, http.StatusBadRequest, writer.Code)
	mockBS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestDeleteRemoveErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockErr := domain.NewBadRequestErr("")
	bid := int64(1)
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockBS := &mocks.MockBundleService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockBS.
		On("Remove", context.TODO(), bid, mockUser).
		Return(mockErr)
	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	writer := prepareAndServeDelete(t, fmt.Sprint(bid), mockBS, mockMWH)

	assert.Equal(t, domain.Status(mockErr), writer.Code)
	mockBS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}
