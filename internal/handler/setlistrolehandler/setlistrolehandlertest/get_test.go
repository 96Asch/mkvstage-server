package setlistrolehandler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/handler/setlistrolehandler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func prepareAndServeGet(
	t *testing.T,
	path string,
	mockSLRS domain.SetlistRoleService,
	mockMWH domain.MiddlewareHandler,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	setlistrolehandler.Initialize(&router.RouterGroup, mockSLRS, mockMWH)

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodGet,
		fmt.Sprintf("/setlistroles%s", path),
		nil,
	)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestGetAll(t *testing.T) {
	expSetlistRoles := &[]domain.SetlistRole{
		{
			ID:        1,
			SetlistID: 1,
		},
		{
			ID:        2,
			SetlistID: 1,
		},
		{
			ID:        3,
			SetlistID: 2,
		},
	}

	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {}

	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	t.Run("Correct no filter", func(t *testing.T) {
		t.Parallel()

		var mockSetlists *[]domain.Setlist

		mockSLRS := &mocks.MockSetlistRoleService{}

		mockSLRS.
			On("Fetch", mock.AnythingOfType("*context.emptyCtx"), mockSetlists).
			Return(expSetlistRoles, nil)

		writer := prepareAndServeGet(t, "", mockSLRS, mockMWH)

		expBody, err := json.Marshal(gin.H{
			"setlistroles": expSetlistRoles,
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, writer.Code)
		assert.Equal(t, expBody, writer.Body.Bytes())

		mockSLRS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})

	t.Run("Correct filter by setlists", func(t *testing.T) {
		t.Parallel()

		mockSetlists := &[]domain.Setlist{
			{
				ID: 1,
			},
		}

		expFilteredSetlistRoles := &[]domain.SetlistRole{
			(*expSetlistRoles)[0],
			(*expSetlistRoles)[1],
		}

		mockSLRS := &mocks.MockSetlistRoleService{}

		mockSLRS.
			On("Fetch", mock.AnythingOfType("*context.emptyCtx"), mockSetlists).
			Return(expFilteredSetlistRoles, nil)

		writer := prepareAndServeGet(t, "?setlist=1", mockSLRS, mockMWH)

		expBody, err := json.Marshal(gin.H{
			"setlistroles": expFilteredSetlistRoles,
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, writer.Code)
		assert.Equal(t, expBody, writer.Body.Bytes())

		mockSLRS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})

	t.Run("Fail invalid setlist query", func(t *testing.T) {
		t.Parallel()

		expErr := domain.NewBadRequestErr("")

		mockSLRS := &mocks.MockSetlistRoleService{}

		writer := prepareAndServeGet(t, "?setlist=a", mockSLRS, mockMWH)

		expBody, err := json.Marshal(gin.H{
			"error": "strconv.Atoi: parsing \"a\": invalid syntax",
		})
		assert.NoError(t, err)

		assert.Equal(t, expErr.Status(), writer.Code)
		assert.Equal(t, expBody, writer.Body.Bytes())

		mockSLRS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})

	t.Run("Fail fetch SetlistRole err", func(t *testing.T) {
		t.Parallel()

		var mockSetlists *[]domain.Setlist

		expErr := domain.NewInternalErr()
		mockSLRS := &mocks.MockSetlistRoleService{}

		mockSLRS.
			On("Fetch", mock.AnythingOfType("*context.emptyCtx"), mockSetlists).
			Return(nil, expErr)

		writer := prepareAndServeGet(t, "", mockSLRS, mockMWH)

		expBody, err := json.Marshal(gin.H{
			"error": expErr.Error(),
		})
		assert.NoError(t, err)

		assert.Equal(t, expErr.Status(), writer.Code)
		assert.Equal(t, expBody, writer.Body.Bytes())

		mockSLRS.AssertExpectations(t)
		mockMWH.AssertExpectations(t)
	})
}
