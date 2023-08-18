package rolehandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/handler/rolehandler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func prepareAndServeUpdate(
	t *testing.T,
	mockRS domain.RoleService,
	mockMWH domain.MiddlewareHandler,
	body *[]byte,
	param string,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	rolehandler.Initialize(&router.RouterGroup, mockRS, mockMWH)

	requestBody := bytes.NewReader(*body)

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodPut,
		fmt.Sprintf("/roles/%s/update", param),
		requestBody,
	)

	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestUpdateByIDCorrect(t *testing.T) {
	t.Parallel()

	rid := int64(1)
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
	}
	mockRole := &domain.Role{
		ID:          rid,
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
		On("Update", context.TODO(), mockRole, mockUser).
		Return(nil)

	byteBody, err := json.Marshal(gin.H{
		"name":        "Foo",
		"description": "Bar",
	})
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, mockRS, mockMWH, &byteBody, fmt.Sprint(rid))
	mockRole.ID = 1

	expectedBytes, err := json.Marshal(gin.H{"role": mockRole})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, expectedBytes, writer.Body.Bytes())
	mockRS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDNoContext(t *testing.T) {
	t.Parallel()

	rid := int64(1)
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

	writer := prepareAndServeUpdate(t, mockRS, mockMWH, &byteBody, fmt.Sprint(rid))

	expectedBytes, err := json.Marshal(gin.H{"error": domain.NewInternalErr()})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Equal(t, expectedBytes, writer.Body.Bytes())
	mockMWH.AssertExpectations(t)
	mockRS.AssertExpectations(t)
}

func TestUpdateByIDInvalidParam(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
	}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockRS := &mocks.MockRoleService{}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(gin.H{
		"name":        "Foo",
		"description": "Bar",
	})
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, mockRS, mockMWH, &byteBody, "a")

	assert.Equal(t, http.StatusBadRequest, writer.Code)
	mockRS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDBindErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	rid := int64(1)
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
		"name":         "Foo",
		"descriptions": "Bar",
	})
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, mockRS, mockMWH, &byteBody, fmt.Sprint(rid))
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	mockMWH.AssertExpectations(t)
	mockRS.AssertExpectations(t)
}

func TestUpdateByIDUpdateErr(t *testing.T) {
	t.Parallel()

	rid := int64(1)
	mockUser := &domain.User{
		ID:         1,
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "Foo@Bar.com",
		Permission: domain.GUEST,
	}
	mockRole := &domain.Role{
		ID:          rid,
		Name:        "Foo",
		Description: "Bar",
	}

	mockErr := domain.NewBadRequestErr("")
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockRS := &mocks.MockRoleService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockRS.
		On("Update", context.TODO(), mockRole, mockUser).
		Return(mockErr)
	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(gin.H{
		"name":        "Foo",
		"description": "Bar",
	})
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, mockRS, mockMWH, &byteBody, fmt.Sprint(rid))
	assert.Equal(t, domain.Status(mockErr), writer.Code)
	mockRS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}
