package rolehandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/handler/rolehandler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func prepareAndServeCreate(
	t *testing.T,
	mockRS domain.RoleService,
	mockMWH domain.MiddlewareHandler,
	body *[]byte,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	rolehandler.Initialize(&router.RouterGroup, mockRS, mockMWH)

	requestBody := bytes.NewReader(*body)

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, "/roles/create", requestBody)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestCreateCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
	}

	mockRole := &domain.Role{
		Name:        "Foo",
		Description: "Bar",
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockRS := &mocks.MockRoleService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)
	mockRS.
		On("Store", context.TODO(), mockRole, mockUser).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg, ok := args.Get(1).(*domain.Role)
			assert.True(t, ok)
			arg.ID = 1
		})

	byteBody, err := json.Marshal(gin.H{
		"name":        "Foo",
		"description": "Bar",
	})
	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockRS, mockMWH, &byteBody)
	assert.Equal(t, http.StatusCreated, writer.Code)

	mockRole.ID = 1

	expectedBytes, err := json.Marshal(gin.H{"role": mockRole})
	assert.NoError(t, err)
	assert.Equal(t, expectedBytes, writer.Body.Bytes())

	mockRS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestCreateNoContext(t *testing.T) {
	t.Parallel()

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockRS := &mocks.MockRoleService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(gin.H{
		"name":        "Foo",
		"description": "Bar",
	})
	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockRS, mockMWH, &byteBody)
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	mockRS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestCreateBindErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "Foo@Bar.com",
		Permission: domain.GUEST,
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockRS := &mocks.MockRoleService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(gin.H{
		"description": "Bar",
	})
	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockRS, mockMWH, &byteBody)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	mockRS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestCreateStoreErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "Foo@Bar.com",
		Permission: domain.GUEST,
	}

	mockRole := &domain.Role{
		Name:        "Foo",
		Description: "Bar",
	}

	mockErr := domain.NewBadRequestErr("")
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockRS := &mocks.MockRoleService{}

	mockRS.
		On("Store", context.TODO(), mockRole, mockUser).
		Return(mockErr)

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(gin.H{
		"name":        "Foo",
		"description": "Bar",
	})
	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockRS, mockMWH, &byteBody)
	assert.Equal(t, domain.Status(mockErr), writer.Code)
	mockRS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}
