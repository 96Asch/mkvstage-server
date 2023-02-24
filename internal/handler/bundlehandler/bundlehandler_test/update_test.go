package bundlehandler_test

import (
	"bytes"
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
	"github.com/stretchr/testify/mock"
)

func prepareAndServeUpdate(
	t *testing.T,
	param string,
	mockBS domain.BundleService,
	mockMWH domain.MiddlewareHandler,
	body *[]byte,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	bundlehandler.Initialize(&router.RouterGroup, mockBS, mockMWH)

	var request *http.Request

	if body != nil {
		bodyReader := bytes.NewReader(*body)
		req, err := http.NewRequestWithContext(
			context.TODO(),
			http.MethodPut,
			fmt.Sprintf("/bundles/%s/update", param),
			bodyReader,
		)
		assert.NoError(t, err)

		request = req
	} else {
		req, err := http.NewRequestWithContext(
			context.TODO(),
			http.MethodPut,
			fmt.Sprintf("/bundles/%s/update", param),
			nil,
		)
		assert.NoError(t, err)

		request = req
	}

	router.ServeHTTP(writer, request)

	return writer
}

func TestUpdateCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockBundle := &domain.Bundle{
		ID:       1,
		Name:     "Foo",
		ParentID: 0,
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockBS := &mocks.MockBundleService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockBS.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), mockBundle, mockUser).
		Return(nil)
	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(gin.H{
		"name":      "Foo",
		"parent_id": 0,
	})
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, "1", mockBS, mockMWH, &byteBody)

	assert.Equal(t, http.StatusOK, writer.Code)

	expectedBody, err := json.Marshal(gin.H{"bundle": mockBundle})
	assert.NoError(t, err)
	assert.Equal(t, expectedBody, writer.Body.Bytes())

	mockBS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateNoContext(t *testing.T) {
	t.Parallel()

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockBS := &mocks.MockBundleService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Next()
	}

	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	writer := prepareAndServeUpdate(t, "1", mockBS, mockMWH, nil)

	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	mockBS.AssertExpectations(t)
}

func TestUpdateParamErr(t *testing.T) {
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

	writer := prepareAndServeUpdate(t, "a", mockBS, mockMWH, nil)

	assert.Equal(t, http.StatusBadRequest, writer.Code)
	mockBS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateBindErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockErr := domain.NewBadRequestErr("")

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockBS := &mocks.MockBundleService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(gin.H{
		"names":     "Foo",
		"parent_id": 0,
	})
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, "1", mockBS, mockMWH, &byteBody)

	assert.Equal(t, domain.Status(mockErr), writer.Code)
	mockBS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockBundle := &domain.Bundle{
		ID:       1,
		Name:     "Foo",
		ParentID: 0,
	}

	mockErr := domain.NewInternalErr()
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockBS := &mocks.MockBundleService{}

	mockBS.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), mockBundle, mockUser).
		Return(mockErr)

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(gin.H{
		"name":      "Foo",
		"parent_id": 0,
	})
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, "1", mockBS, mockMWH, &byteBody)

	assert.Equal(t, domain.Status(mockErr), writer.Code)
	mockBS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}
