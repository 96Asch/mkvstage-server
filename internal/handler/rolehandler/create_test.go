package rolehandler

import (
	"bytes"
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

func TestCreateCorrect(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		Permission: domain.MEMBER,
	}

	mockRole := &domain.Role{
		Name:        "Foo",
		Description: "Bar",
	}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := &mocks.MockMiddlewareHandler{}
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)
	mockRS := &mocks.MockRoleService{}
	mockRS.
		On("Store", mock.AnythingOfType("*context.emptyCtx"), mockRole, mockUser).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).(*domain.Role)
			arg.ID = 1
		})

	byteBody, err := json.Marshal(gin.H{
		"name":        "Foo",
		"description": "Bar",
	})
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)

	router := gin.New()
	reqBody := bytes.NewReader(byteBody)
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/roles/create", reqBody)
	assert.NoError(t, err)

	Initialize(&router.RouterGroup, mockRS, mockMWH)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockRole.ID = 1

	expectedBytes, err := json.Marshal(gin.H{"role": mockRole})
	assert.NoError(t, err)
	assert.Equal(t, expectedBytes, w.Body.Bytes())
	mockRS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestCreateNoContext(t *testing.T) {
	mockRS := &mocks.MockRoleService{}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Next()
	}
	mockMWH := new(mocks.MockMiddlewareHandler)
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	Initialize(&router.RouterGroup, mockRS, mockMWH)

	byteBody, err := json.Marshal(gin.H{
		"title":       "Foo",
		"subtitle":    "Bar",
		"key":         "A",
		"bpm":         120,
		"bundle_id":   1,
		"creator_id":  1,
		"chord_sheet": `{"Verse" : "Foobar"}`,
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(byteBody)

	req, err := http.NewRequest(http.MethodPost, "/roles/create", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestCreateBindErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "Foo@Bar.com",
		Permission: domain.GUEST,
	}

	mockRS := &mocks.MockRoleService{}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := new(mocks.MockMiddlewareHandler)
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	Initialize(&router.RouterGroup, mockRS, mockMWH)
	byteBody, err := json.Marshal(gin.H{
		"title":       "Foo",
		"subtitle":    "Bar",
		"bpm":         120,
		"bundle_id":   1,
		"creator_id":  mockUser.ID,
		"chord_sheet": `{"Verse" : "Foobar"}`,
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(byteBody)

	req, err := http.NewRequest(http.MethodPost, "/roles/create", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockRS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestCreateStoreErr(t *testing.T) {
	mockUser := &domain.User{
		ID:         1,
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "Foo@Bar.com",
		Permission: domain.GUEST,
	}

	mockRole := &domain.Role{
		Name:        "Foo",
		Description: "Bar",
	}

	mockErr := domain.NewBadRequestErr("")
	mockRS := &mocks.MockRoleService{}
	mockRS.
		On("Store", mock.AnythingOfType("*context.emptyCtx"), mockRole, mockUser).
		Return(mockErr)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}
	mockMWH := new(mocks.MockMiddlewareHandler)
	mockMWH.On("AuthenticateUser").Return(mockAuthHF)

	Initialize(&router.RouterGroup, mockRS, mockMWH)
	byteBody, err := json.Marshal(gin.H{
		"name":        "Foo",
		"description": "Bar",
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(byteBody)

	req, err := http.NewRequest(http.MethodPost, "/roles/create", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, domain.Status(mockErr), w.Code)
	mockRS.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}
