package setlisthandler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/handler/setlisthandler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func prepareAndServeGet(
	t *testing.T,
	mockSL domain.SetlistService,
	mockSLES domain.SetlistEntryService,
	mockSS domain.SongService,
	mockMWH domain.MiddlewareHandler,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	setlisthandler.Initialize(&router.RouterGroup, mockSL, mockSLES, mockSS, mockMWH)

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodGet,
		"/setlists/get",
		nil,
	)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}
