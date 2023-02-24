package rolehandler_test

import (
	"context"
	"fmt"
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

func prepareAndServeDelete(
	t *testing.T,
	mockRS domain.RoleService,
	mockMWH domain.MiddlewareHandler,
	param string,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	rolehandler.Initialize(&router.RouterGroup, mockRS, mockMWH)

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodDelete,
		fmt.Sprintf("/roles/%s/delete", param),
		nil,
	)

	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestDeleteByIDCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
	}

	sid := int64(1)
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
		On("Remove", mock.AnythingOfType("*context.emptyCtx"), sid, mockUser).
		Return(nil)

	writer := prepareAndServeDelete(t, mockRS, mockMWH, fmt.Sprint(sid))
	assert.Equal(t, http.StatusAccepted, writer.Code)
	mockMWH.AssertExpectations(t)
	mockRS.AssertExpectations(t)
}

func TestDeleteByIDNoContext(t *testing.T) {
	t.Parallel()

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockRS := &mocks.MockRoleService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	writer := prepareAndServeDelete(t, mockRS, mockMWH, "1")
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	mockMWH.AssertExpectations(t)
	mockRS.AssertExpectations(t)
}

func TestDeleteByIDInvalidParam(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
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

	writer := prepareAndServeDelete(t, mockRS, mockMWH, "a")
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	mockMWH.AssertExpectations(t)
	mockRS.AssertExpectations(t)
}

func TestDeleteByIDRemoveErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "Foo@Bar.com",
		Permission: domain.GUEST,
	}

	rid := int64(1)
	mockErr := domain.NewInternalErr()
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
		On("Remove", mock.AnythingOfType("*context.emptyCtx"), rid, mockUser).
		Return(mockErr)

	writer := prepareAndServeDelete(t, mockRS, mockMWH, fmt.Sprint(rid))
	assert.Equal(t, domain.Status(mockErr), writer.Code)
	mockRS.AssertExpectations(t)
}
