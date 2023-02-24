package rolehandler_test

import (
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
	"github.com/stretchr/testify/mock"
)

func prepareAndServeGet(
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
		http.MethodGet,
		fmt.Sprintf("/roles%s", param),
		nil,
	)

	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestGetAllCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockRoles := &[]domain.Role{
		{
			ID:          1,
			Name:        "Foo",
			Description: "Bar",
		},
		{
			ID:          2,
			Name:        "Bar",
			Description: "Foo",
		},
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockSS := &mocks.MockRoleService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)
	mockSS.
		On("FetchAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(mockRoles, nil)

	writer := prepareAndServeGet(t, mockSS, mockMWH, "")

	expectedBody, err := json.Marshal(gin.H{"roles": mockRoles})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, expectedBody, writer.Body.Bytes())
	mockMWH.AssertExpectations(t)
	mockSS.AssertExpectations(t)
}

func TestGetAllFetchErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockErr := domain.NewInternalErr()
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockSS := &mocks.MockRoleService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)
	mockSS.
		On("FetchAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(nil, mockErr)

	writer := prepareAndServeGet(t, mockSS, mockMWH, "")

	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	mockMWH.AssertExpectations(t)
	mockSS.AssertExpectations(t)
}
