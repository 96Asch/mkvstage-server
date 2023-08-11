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
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
)

func prepareAndServeGet(
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
		http.MethodGet,
		fmt.Sprintf("/bundles%s", param),
		nil,
	)

	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestGetByIDCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockBundle := &domain.Bundle{
		ID:       1,
		Name:     "Foo",
		ParentID: 0,
	}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockBS := &mocks.MockBundleService{}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)
	mockBS.
		On("FetchByID", context.TODO(), mockBundle.ID).
		Return(mockBundle, nil)

	writer := prepareAndServeGet(t, "/1", mockBS, mockMWH)

	expectedBody, err := json.Marshal(gin.H{"bundle": mockBundle})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, expectedBody, writer.Body.Bytes())

	mockMWH.AssertExpectations(t)
	mockBS.AssertExpectations(t)
}

func TestGetInvalidParam(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockBS := &mocks.MockBundleService{}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	writer := prepareAndServeGet(t, "/a", mockBS, mockMWH)

	assert.Equal(t, http.StatusBadRequest, writer.Code)

	mockMWH.AssertExpectations(t)
	mockBS.AssertExpectations(t)
}

func TestGetNoRecord(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockErr := domain.NewRecordNotFoundErr("", "")
	mockBS := &mocks.MockBundleService{}
	mockMWH := &mocks.MockMiddlewareHandler{}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)
	mockBS.
		On("FetchByID", context.TODO(), int64(-1)).
		Return(nil, mockErr)

	writer := prepareAndServeGet(t, "/-1", mockBS, mockMWH)

	assert.Equal(t, http.StatusNotFound, writer.Code)

	mockMWH.AssertExpectations(t)
	mockBS.AssertExpectations(t)
}

func TestGetAllCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}
	mockBundles := &[]domain.Bundle{
		{
			ID:       1,
			Name:     "Foo",
			ParentID: 0,
		},
		{
			ID:       2,
			Name:     "Bar",
			ParentID: 1,
		},
	}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockBS := &mocks.MockBundleService{}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)
	mockBS.
		On("FetchAll", context.TODO()).
		Return(mockBundles, nil)

	writer := prepareAndServeGet(t, "", mockBS, mockMWH)

	expectedBody, err := json.Marshal(gin.H{"bundles": mockBundles})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, expectedBody, writer.Body.Bytes())

	mockMWH.AssertExpectations(t)
	mockBS.AssertExpectations(t)
}

func TestGetAllFetchErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockErr := domain.NewInternalErr()
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockBS := &mocks.MockBundleService{}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)
	mockBS.
		On("FetchAll", context.TODO()).
		Return(nil, mockErr)

	writer := prepareAndServeGet(t, "", mockBS, mockMWH)

	assert.Equal(t, http.StatusInternalServerError, writer.Code)

	mockMWH.AssertExpectations(t)
	mockBS.AssertExpectations(t)
}
