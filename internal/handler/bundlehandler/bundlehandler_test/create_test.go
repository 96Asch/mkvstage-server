package bundlehandler_test

import (
	"bytes"
	"context"
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

func prepareAndServeCreate(
	t *testing.T,
	mockBS domain.BundleService,
	mockMWH domain.MiddlewareHandler,
	body *[]byte,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()
	bodyReader := bytes.NewReader(*body)

	bundlehandler.Initialize(&router.RouterGroup, mockBS, mockMWH)

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, "/bundles/create", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestCreateCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "Foo@Bar.com",
		Permission: domain.MEMBER,
	}
	mockBundle := &domain.Bundle{
		Name:     "Foobar",
		ParentID: 0,
	}

	mockBS := new(mocks.MockBundleService)
	mockBS.
		On("Store", mock.AnythingOfType("*context.emptyCtx"), mockBundle, mockUser).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg, valid := args.Get(1).(*domain.Bundle)
			assert.True(t, valid)

			arg.ID = 1
		})

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH := new(mocks.MockMiddlewareHandler)
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	mockReq, err := json.Marshal(gin.H{
		"id":        int64(1),
		"name":      "Foobar",
		"parent_id": int64(0),
	})
	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockBS, mockMWH, &mockReq)

	expectedBundle, err := json.Marshal(gin.H{
		"bundle": &domain.Bundle{
			ID:       1,
			Name:     mockBundle.Name,
			ParentID: mockBundle.ParentID,
		},
	})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, writer.Code)
	assert.Equal(t, expectedBundle, writer.Body.Bytes())
	mockBS.AssertExpectations(t)
}

func TestCreateBindErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "Foo@Bar.com",
		Permission: domain.MEMBER,
	}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockBS := &mocks.MockBundleService{}
	mockMWH := &mocks.MockMiddlewareHandler{}

	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	mockReq, err := json.Marshal(gin.H{
		"name":      "",
		"parent_id": int64(0),
	})
	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockBS, mockMWH, &mockReq)

	assert.Equal(t, http.StatusBadRequest, writer.Code)
	mockBS.AssertExpectations(t)
}

func TestCreateStoreErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "Foo@Bar.com",
		Permission: domain.MEMBER,
	}
	mockBundle := &domain.Bundle{
		Name:     "Foobar",
		ParentID: 0,
	}

	mockErr := domain.NewInternalErr()
	mockBS := &mocks.MockBundleService{}

	mockBS.
		On("Store", mock.AnythingOfType("*context.emptyCtx"), mockBundle, mockUser).
		Return(mockErr)

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH := &mocks.MockMiddlewareHandler{}

	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	mockReq, err := json.Marshal(gin.H{
		"id":        int64(1),
		"name":      "Foobar",
		"parent_id": int64(0),
	})
	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockBS, mockMWH, &mockReq)

	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	mockBS.AssertExpectations(t)
}

func TestCreateNoContext(t *testing.T) {
	t.Parallel()

	mockBS := &mocks.MockBundleService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Next()
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	mockReq, err := json.Marshal(gin.H{
		"id":        int64(1),
		"name":      "Foobar",
		"parent_id": int64(0),
	})
	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockBS, mockMWH, &mockReq)

	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	mockBS.AssertExpectations(t)
}

func TestCreateNotAuth(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "Foo@Bar.com",
		Permission: domain.GUEST,
	}
	mockBundle := &domain.Bundle{
		Name:     "Foobar",
		ParentID: 0,
	}

	mockErr := domain.NewNotAuthorizedErr("")
	mockBS := &mocks.MockBundleService{}
	mockBS.
		On("Store", mock.AnythingOfType("*context.emptyCtx"), mockBundle, mockUser).
		Return(mockErr)

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	mockReq, err := json.Marshal(gin.H{
		"id":        int64(1),
		"name":      "Foobar",
		"parent_id": int64(0),
	})
	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockBS, mockMWH, &mockReq)

	assert.Equal(t, http.StatusUnauthorized, writer.Code)
	mockBS.AssertExpectations(t)
}
