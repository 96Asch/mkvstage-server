package setlistrolehandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/96Asch/mkvstage-server/backend/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/backend/internal/handler/setlistrolehandler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func prepareAndServeCreate(
	t *testing.T,
	mockSLRS domain.SetlistRoleService,
	mockMWH domain.MiddlewareHandler,
	body *[]byte,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	setlistrolehandler.Initialize(&router.RouterGroup, mockSLRS, mockMWH)

	requestBody := bytes.NewReader(*body)
	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodPost,
		"/setlistroles",
		requestBody,
	)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestCreate(t *testing.T) {
	user := &domain.User{
		ID:         1,
		Permission: domain.ADMIN,
	}

	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", user)
		ctx.Next()
	}

	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	t.Run("Correct", func(t *testing.T) {
		t.Parallel()

		mockSLRS := &mocks.MockSetlistRoleService{}

		mockSLRS.
			On("Store",
				context.TODO(),
				&[]domain.SetlistRole{
					{
						SetlistID:  1,
						UserRoleID: 1,
					},
					{
						SetlistID:  1,
						UserRoleID: 2,
					},
				},
				user).
			Return(nil).
			Run(func(args mock.Arguments) {
				arg, ok := args.Get(1).(*[]domain.SetlistRole)
				assert.True(t, ok)

				for idx := range *arg {
					(*arg)[idx].ID = int64(idx + 1)
				}
			})

		byteBody, err := json.Marshal(gin.H{
			"userroles": []gin.H{
				{
					"setlist_id":  1,
					"userrole_id": 1,
				},
				{
					"setlist_id":  1,
					"userrole_id": 2,
				},
			},
		})
		assert.NoError(t, err)

		expResponse, err := json.Marshal(gin.H{
			"userroles": []gin.H{
				{
					"id":          1,
					"setlist_id":  1,
					"userrole_id": 1,
				},
				{
					"id":          2,
					"setlist_id":  1,
					"userrole_id": 2,
				},
			},
		})
		assert.NoError(t, err)

		writer := prepareAndServeCreate(t, mockSLRS, mockMWH, &byteBody)

		assert.Equal(t, http.StatusCreated, writer.Code)
		assert.Equal(t, expResponse, writer.Body.Bytes())
	})

	t.Run("Fail no user in context", func(t *testing.T) {
		t.Parallel()

		var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {}

		expErr := domain.NewInternalErr()
		mockMWH := &mocks.MockMiddlewareHandler{}
		mockSLRS := &mocks.MockSetlistRoleService{}

		mockMWH.On("AuthenticateUser").Return(mockAuthHF)

		byteBody, err := json.Marshal(gin.H{
			"userroles": []gin.H{
				{
					"setlist_id":  1,
					"userrole_id": 1,
				},
				{
					"setlist_id":  1,
					"userrole_id": 2,
				},
			},
		})
		assert.NoError(t, err)

		expResponse, err := json.Marshal(gin.H{
			"error": expErr.Error(),
		})
		assert.NoError(t, err)

		writer := prepareAndServeCreate(t, mockSLRS, mockMWH, &byteBody)

		assert.Equal(t, expErr.Status(), writer.Code)
		assert.Equal(t, expResponse, writer.Body.Bytes())
	})

	t.Run("Fail invalid user in context", func(t *testing.T) {
		t.Parallel()

		var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
			ctx.Set("user", 0)
			ctx.Next()
		}

		expErr := domain.NewInternalErr()
		mockMWH := &mocks.MockMiddlewareHandler{}
		mockSLRS := &mocks.MockSetlistRoleService{}

		mockMWH.On("AuthenticateUser").Return(mockAuthHF)

		byteBody, err := json.Marshal(gin.H{
			"userroles": []gin.H{
				{
					"setlist_id":  1,
					"userrole_id": 1,
				},
				{
					"setlist_id":  1,
					"userrole_id": 2,
				},
			},
		})
		assert.NoError(t, err)

		expResponse, err := json.Marshal(gin.H{
			"error": expErr.Error(),
		})
		assert.NoError(t, err)

		writer := prepareAndServeCreate(t, mockSLRS, mockMWH, &byteBody)

		assert.Equal(t, expErr.Status(), writer.Code)
		assert.Equal(t, expResponse, writer.Body.Bytes())
	})

	t.Run("Fail invalid binding", func(t *testing.T) {
		t.Parallel()

		expErr := domain.NewBadRequestErr("field 'setlist_id' is required")
		mockSLRS := &mocks.MockSetlistRoleService{}

		mockMWH.On("AuthenticateUser").Return(mockAuthHF)

		byteBody, err := json.Marshal(gin.H{
			"userroles": []gin.H{
				{
					"SetlistID":   1,
					"userrole_id": 1,
				},
				{
					"setlist_id":  1,
					"userrole_id": 2,
				},
			},
		})
		assert.NoError(t, err)

		expResponse, err := json.Marshal(gin.H{
			"error": expErr.Error(),
		})
		assert.NoError(t, err)

		writer := prepareAndServeCreate(t, mockSLRS, mockMWH, &byteBody)

		assert.Equal(t, expErr.Status(), writer.Code)
		assert.Equal(t, expResponse, writer.Body.Bytes())
	})

	t.Run("Fail SetlistRole Store err", func(t *testing.T) {
		t.Parallel()

		expErr := domain.NewInternalErr()
		mockSLRS := &mocks.MockSetlistRoleService{}

		mockMWH.On("AuthenticateUser").Return(mockAuthHF)

		mockSLRS.
			On("Store",
				context.TODO(),
				&[]domain.SetlistRole{
					{
						SetlistID:  1,
						UserRoleID: 1,
					},
					{
						SetlistID:  1,
						UserRoleID: 2,
					},
				},
				user).
			Return(expErr)

		byteBody, err := json.Marshal(gin.H{
			"userroles": []gin.H{
				{
					"setlist_id":  1,
					"userrole_id": 1,
				},
				{
					"setlist_id":  1,
					"userrole_id": 2,
				},
			},
		})
		assert.NoError(t, err)

		expResponse, err := json.Marshal(gin.H{
			"error": expErr.Error(),
		})
		assert.NoError(t, err)

		writer := prepareAndServeCreate(t, mockSLRS, mockMWH, &byteBody)

		assert.Equal(t, expErr.Status(), writer.Code)
		assert.Equal(t, expResponse, writer.Body.Bytes())
	})
}
