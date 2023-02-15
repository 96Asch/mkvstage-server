package tokenhandler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCorrect(t *testing.T) {
	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Password:     "FooBar",
		Email:        "Foo@Bar.com",
		Permission:   domain.GUEST,
		ProfileColor: "FFFFFF",
	}

	mockRefresh := &domain.RefreshToken{
		ID:                 uuid.New(),
		UserID:             mockUser.ID,
		Refresh:            "refresh-token",
		ExpirationDuration: time.Minute,
	}

	mockAccess := &domain.AccessToken{
		Access: "access-token",
	}

	mockTS := new(mocks.MockTokenService)
	mockTS.
		On("CreateAccess", mock.AnythingOfType("*context.emptyCtx"), mockRefresh.Refresh).
		Return(mockAccess, nil)

	mockUS := new(mocks.MockUserService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	group := router.Group("test")
	Initialize(group, mockTS, mockUS)

	mockByte, err := json.Marshal(gin.H{
		"refresh": mockRefresh.Refresh,
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(mockByte)

	req, err := http.NewRequest(http.MethodPost, "/test/tokens/renewaccess", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	expectedRes, err := json.Marshal(gin.H{
		"token": mockAccess,
	})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedRes, w.Body.Bytes())
	mockUS.AssertExpectations(t)
}

func TestCreateBindErr(t *testing.T) {
	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Password:     "FooBar",
		Email:        "Foo@Bar.com",
		Permission:   domain.GUEST,
		ProfileColor: "FFFFFF",
	}

	mockRefresh := &domain.RefreshToken{
		ID:                 uuid.New(),
		UserID:             mockUser.ID,
		Refresh:            "refresh-token",
		ExpirationDuration: time.Minute,
	}

	mockTS := new(mocks.MockTokenService)
	mockUS := new(mocks.MockUserService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	group := router.Group("test")
	Initialize(group, mockTS, mockUS)

	mockByte, err := json.Marshal(gin.H{
		"resfresh": mockRefresh.Refresh,
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(mockByte)

	req, err := http.NewRequest(http.MethodPost, "/test/tokens/renewaccess", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUS.AssertExpectations(t)
}

func TestCreateAccessErr(t *testing.T) {
	mockUser := &domain.User{
		FirstName:    "Foo",
		LastName:     "Bar",
		Password:     "FooBar",
		Email:        "Foo@Bar.com",
		Permission:   domain.GUEST,
		ProfileColor: "FFFFFF",
	}

	mockRefresh := &domain.RefreshToken{
		ID:                 uuid.New(),
		UserID:             mockUser.ID,
		Refresh:            "refresh-token",
		ExpirationDuration: time.Minute,
	}

	mockErr := domain.NewNotAuthorizedErr("")
	mockTS := new(mocks.MockTokenService)
	mockTS.
		On("CreateAccess", mock.AnythingOfType("*context.emptyCtx"), mockRefresh.Refresh).
		Return(nil, mockErr)

	mockUS := new(mocks.MockUserService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()

	group := router.Group("test")
	Initialize(group, mockTS, mockUS)

	mockByte, err := json.Marshal(gin.H{
		"refresh": mockRefresh.Refresh,
	})
	assert.NoError(t, err)

	bodyReader := bytes.NewReader(mockByte)

	req, err := http.NewRequest(http.MethodPost, "/test/tokens/renewaccess", bodyReader)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	expectedRes, err := json.Marshal(gin.H{
		"error": mockErr,
	})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, expectedRes, w.Body.Bytes())
	mockUS.AssertExpectations(t)
}
