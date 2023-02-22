package rolehandler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllCorrect(t *testing.T) {
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
	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	mockSS := &mocks.MockRoleService{}
	mockSS.
		On("FetchAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(mockRoles, nil)

	router := gin.New()
	Initialize(&router.RouterGroup, mockSS, mockMWH)

	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/roles", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	expectedBody, err := json.Marshal(gin.H{"roles": mockRoles})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedBody, w.Body.Bytes())

	mockMWH.AssertExpectations(t)
	mockSS.AssertExpectations(t)
}

func TestGetAllFetchErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockMWH := new(mocks.MockMiddlewareHandler)
	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	mockErr := domain.NewInternalErr()

	mockSS := &mocks.MockRoleService{}
	mockSS.
		On("FetchAll", mock.AnythingOfType("*context.emptyCtx")).
		Return(nil, mockErr)

	router := gin.New()
	Initialize(&router.RouterGroup, mockSS, mockMWH)

	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/roles", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockMWH.AssertExpectations(t)
	mockSS.AssertExpectations(t)
}
