package songhandler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/handler/songhandler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func prepareAndServeDelete(
	t *testing.T,
	mockSS domain.SongService,
	mockMWH domain.MiddlewareHandler,
	param string,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	songhandler.Initialize(&router.RouterGroup, mockSS, mockMWH)

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodDelete,
		fmt.Sprintf("/songs/%s/delete", param),
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
	mockSS := &mocks.MockSongService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)
	mockSS.
		On("Remove", context.TODO(), sid, mockUser).
		Return(nil)

	writer := prepareAndServeDelete(t, mockSS, mockMWH, fmt.Sprint(sid))
	assert.Equal(t, http.StatusAccepted, writer.Code)
	mockMWH.AssertExpectations(t)
	mockSS.AssertExpectations(t)
}

func TestDeleteByIDNoContext(t *testing.T) {
	t.Parallel()

	sid := int64(1)
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockSS := &mocks.MockSongService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	writer := prepareAndServeDelete(t, mockSS, mockMWH, fmt.Sprint(sid))
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	mockSS.AssertExpectations(t)
}

func TestDeleteByIDInvalidParam(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockSS := &mocks.MockSongService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	writer := prepareAndServeDelete(t, mockSS, mockMWH, "a")
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	mockMWH.AssertExpectations(t)
	mockSS.AssertExpectations(t)
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

	sid := int64(1)
	mockErr := domain.NewInternalErr()
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockSS := &mocks.MockSongService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)
	mockSS.
		On("Remove", context.TODO(), sid, mockUser).
		Return(mockErr)

	writer := prepareAndServeDelete(t, mockSS, mockMWH, fmt.Sprint(sid))
	assert.Equal(t, domain.Status(mockErr), writer.Code)
	mockSS.AssertExpectations(t)
}
