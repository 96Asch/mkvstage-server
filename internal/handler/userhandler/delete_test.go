package userhandler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDeleteCorrect(t *testing.T) {

}

func TestDeleteNoContext(t *testing.T) {
	mockUS := new(mocks.MockUserService)

	r := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	router := gin.New()

	group := router.Group("test")
	Initialize(group, mockUS)

	req, err := http.NewRequest(http.MethodDelete, "/test/users/me/delete", nil)
	assert.NoError(t, err)

	router.ServeHTTP(r, req)

	assert.Equal(t, http.StatusInternalServerError, r.Code)
	mockUS.AssertNotCalled(t, "Remove")
}
