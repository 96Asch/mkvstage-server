package tokenhandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/handler/tokenhandler"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func prepareAndServeCreate(
	t *testing.T,
	mockTS domain.TokenService,
	mockUS domain.UserService,
	body *[]byte,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	tokenhandler.Initialize(&router.RouterGroup, mockTS, mockUS)

	requestBody := bytes.NewReader(*body)

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, "/tokens/renew", requestBody)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestCreateCorrect(t *testing.T) {
	t.Parallel()

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

	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}

	mockTS.
		On("CreateAccess", context.TODO(), mockRefresh.Refresh).
		Return(mockAccess, nil)

	byteBody, err := json.Marshal(gin.H{
		"refresh": mockRefresh.Refresh,
	})
	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockTS, mockUS, &byteBody)

	expectedRes, err := json.Marshal(gin.H{
		"token": mockAccess,
	})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, expectedRes, writer.Body.Bytes())
	mockUS.AssertExpectations(t)
}

func TestCreateBindErr(t *testing.T) {
	t.Parallel()

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

	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}

	byteBody, err := json.Marshal(gin.H{
		"resfresh": mockRefresh.Refresh,
	})
	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockTS, mockUS, &byteBody)

	assert.Equal(t, http.StatusBadRequest, writer.Code)
	mockUS.AssertExpectations(t)
}

func TestCreateAccessErr(t *testing.T) {
	t.Parallel()

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
	mockTS := &mocks.MockTokenService{}
	mockUS := &mocks.MockUserService{}

	mockTS.
		On("CreateAccess", context.TODO(), mockRefresh.Refresh).
		Return(nil, mockErr)

	byteBody, err := json.Marshal(gin.H{
		"refresh": mockRefresh.Refresh,
	})
	assert.NoError(t, err)

	writer := prepareAndServeCreate(t, mockTS, mockUS, &byteBody)

	expectedRes, err := json.Marshal(gin.H{
		"error": mockErr,
	})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, writer.Code)
	assert.Equal(t, expectedRes, writer.Body.Bytes())
	mockUS.AssertExpectations(t)
}
