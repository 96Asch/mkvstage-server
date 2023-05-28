package setlisthandler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/handler/setlisthandler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func prepareAndServeDelete(
	t *testing.T,
	paramID string,
	mockSL domain.SetlistService,
	mockSS domain.SongService,
	mockMWH domain.MiddlewareHandler,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	setlisthandler.Initialize(&router.RouterGroup, mockSL, mockSS, mockMWH)

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodDelete,
		fmt.Sprintf("/setlists/%s/delete", paramID),
		nil,
	)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestDeleteByIDCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
	}

	mockSetlistID := 1

	mockSL := &mocks.MockSetlistService{}
	mockSS := &mocks.MockSongService{}

	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	mockSL.
		On("Remove", mock.AnythingOfType("*context.emptyCtx"), int64(mockSetlistID), mockUser).
		Return(nil)

	writer := prepareAndServeDelete(t, fmt.Sprint(mockSetlistID), mockSL, mockSS, mockMWH)

	assert.Equal(t, http.StatusAccepted, writer.Code)
	mockSL.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestDeleteByIDNoContext(t *testing.T) {
	t.Parallel()

	mockSetlistID := 1

	mockErr := domain.NewInternalErr()
	mockSL := &mocks.MockSetlistService{}
	mockSS := &mocks.MockSongService{}

	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	writer := prepareAndServeDelete(t, fmt.Sprint(mockSetlistID), mockSL, mockSS, mockMWH)

	assert.Equal(t, mockErr.Status(), writer.Code)
	mockSL.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestDeleteByIDInvalidParam(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockSL := &mocks.MockSetlistService{}
	mockSS := &mocks.MockSongService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	writer := prepareAndServeDelete(t, "a", mockSL, mockSS, mockMWH)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	mockMWH.AssertExpectations(t)
	mockSL.AssertExpectations(t)
	mockSS.AssertExpectations(t)
}

func TestDeleteByIDRemoveErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
	}

	mockSetlistID := 1

	mockErr := domain.NewInternalErr()
	mockSL := &mocks.MockSetlistService{}
	mockSS := &mocks.MockSongService{}

	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	mockSL.
		On("Remove", mock.AnythingOfType("*context.emptyCtx"), int64(mockSetlistID), mockUser).
		Return(mockErr)

	writer := prepareAndServeDelete(t, fmt.Sprint(mockSetlistID), mockSL, mockSS, mockMWH)

	assert.Equal(t, mockErr.Status(), writer.Code)
	mockSL.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}
