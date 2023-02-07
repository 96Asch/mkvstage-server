package userhandler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
)

func TestUpdateInvalidBind(t *testing.T) {
	mockUS := new(mocks.MockUserService)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	w := httptest.NewRecorder()

	group := router.Group("test")
	Initialize(group, mockUS)

	mockByte, err := json.Marshal(gin.H{
		"first_names":  "Foo",
		"last_name":    "Bar",
		"password":     "Foo",
		"permission":   "Foo",
		"profileColor": "Foo",
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(mockByte)

	req, err := http.NewRequest(http.MethodPatch, "/test/users/update", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUS.AssertNotCalled(t, "Update")
}
