package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/service"
	"github.com/96Asch/mkvstage-server/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	refreshSecret = "refresh-secret"
	accessSecret  = "access-secret"
)

func TestCreateAccessCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:       1,
		Email:    "Foo@Bar.com",
		Password: "FooBar",
	}

	mockTR := &mocks.MockTokenRepository{}
	mockUR := &mocks.MockUserRepository{}
	tokenConfig := &domain.TokenConfig{
		IAT:         time.Now(),
		ExpDuration: time.Minute,
		Secret:      refreshSecret,
	}

	refresh, err := util.GenerateRefreshToken(mockUser.ID, tokenConfig)
	assert.NoError(t, err)

	mockRefreshTokens := &[]domain.RefreshToken{{Refresh: refresh.Refresh}}

	mockUR.
		On("GetByID", context.TODO(), mockUser.ID).
		Return(mockUser, nil)
	mockTR.
		On("GetAll", context.TODO(), mockUser.ID).
		Return(mockRefreshTokens, nil)

	tokenService := service.NewTokenService(mockTR, mockUR, accessSecret, refreshSecret)
	ctx := context.TODO()

	accessToken, err := tokenService.CreateAccess(ctx, refresh.Refresh)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken.Access)

	claims, err := util.VerifyAccessToken(accessToken.Access, accessSecret)
	assert.NoError(t, err)
	assert.NotNil(t, claims.User)
	assert.Empty(t, claims.User.Password)

	mockUser.Password = ""
	assert.Equal(t, mockUser, claims.User)
}

func TestCreateAccessRepoErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:       1,
		Email:    "Foo@Bar.com",
		Password: "FooBar",
	}

	mockErr := domain.NewInternalErr()
	mockTR := &mocks.MockTokenRepository{}
	mockUR := &mocks.MockUserRepository{}
	tokenConfig := &domain.TokenConfig{
		IAT:         time.Now(),
		ExpDuration: time.Minute,
		Secret:      refreshSecret,
	}

	refresh, err := util.GenerateRefreshToken(mockUser.ID, tokenConfig)
	assert.NoError(t, err)

	mockUR.
		On("GetByID", context.TODO(), mockUser.ID).
		Return(mockUser, nil)
	mockTR.
		On("GetAll", context.TODO(), mockUser.ID).
		Return(nil, mockErr)

	tokenService := service.NewTokenService(mockTR, mockUR, accessSecret, refreshSecret)
	ctx := context.TODO()

	_, err = tokenService.CreateAccess(ctx, refresh.Refresh)
	assert.ErrorAs(t, err, &mockErr)
}

func TestCreateAccessRefreshNotInRepo(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:       1,
		Email:    "Foo@Bar.com",
		Password: "FooBar",
	}

	mockErr := domain.NewInternalErr()
	mockTR := &mocks.MockTokenRepository{}
	mockUR := &mocks.MockUserRepository{}
	tokenConfig := &domain.TokenConfig{
		IAT:         time.Now(),
		ExpDuration: time.Minute,
		Secret:      refreshSecret,
	}

	mockRefreshTokens := &[]domain.RefreshToken{}
	refresh, err := util.GenerateRefreshToken(mockUser.ID, tokenConfig)
	assert.NoError(t, err)

	mockUR.
		On("GetByID", context.TODO(), mockUser.ID).
		Return(mockUser, nil)
	mockTR.
		On("GetAll", context.TODO(), mockUser.ID).
		Return(mockRefreshTokens, nil)

	tokenService := service.NewTokenService(mockTR, mockUR, accessSecret, refreshSecret)
	ctx := context.TODO()

	_, err = tokenService.CreateAccess(ctx, refresh.Refresh)
	assert.ErrorAs(t, err, &mockErr)
}

func TestCreateAccessInvalidRefresh(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:       1,
		Email:    "Foo@Bar.com",
		Password: "FooBar",
	}

	mockErr := domain.NewNotAuthorizedErr("")
	mockTR := &mocks.MockTokenRepository{}
	mockUR := &mocks.MockUserRepository{}
	tokenConfig := &domain.TokenConfig{
		IAT:         time.Now(),
		ExpDuration: -time.Minute,
		Secret:      refreshSecret,
	}

	refresh, err := util.GenerateRefreshToken(mockUser.ID, tokenConfig)
	assert.NoError(t, err)

	tokenService := service.NewTokenService(mockTR, mockUR, accessSecret, refreshSecret)
	ctx := context.TODO()

	accessToken, err := tokenService.CreateAccess(ctx, refresh.Refresh)
	assert.ErrorAs(t, err, &mockErr)
	assert.Nil(t, accessToken)
}

func TestCreateAccessRefreshEmpty(t *testing.T) {
	t.Parallel()

	expectedErr := domain.NewBadRequestErr("")

	mockTR := &mocks.MockTokenRepository{}
	mockUR := &mocks.MockUserRepository{}

	tokenService := service.NewTokenService(mockTR, mockUR, accessSecret, refreshSecret)
	ctx := context.TODO()

	accessToken, err := tokenService.CreateAccess(ctx, "")
	assert.ErrorAs(t, err, &expectedErr)
	assert.Nil(t, accessToken)
}

func TestCreateRefreshCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:       1,
		Email:    "Foo@Bar.com",
		Password: "FooBar",
	}

	mockTR := &mocks.MockTokenRepository{}
	mockUR := &mocks.MockUserRepository{}

	mockTR.
		On("Create", context.TODO(), mock.AnythingOfType("*domain.RefreshToken")).
		Return(nil)

	tokenService := service.NewTokenService(mockTR, mockUR, accessSecret, refreshSecret)
	ctx := context.TODO()

	accessToken, err := tokenService.CreateRefresh(ctx, mockUser.ID, "")
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken.Refresh)

	claims, err := util.VerifyRefreshToken(accessToken.Refresh, refreshSecret)
	assert.NoError(t, err)
	assert.Equal(t, mockUser.ID, claims.UID)
}

func TestCreateRefreshDeleteErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:       1,
		Email:    "Foo@Bar.com",
		Password: "FooBar",
	}

	refreshToken := "refresh-token"
	mockErr := domain.NewInternalErr()

	mockTR := &mocks.MockTokenRepository{}
	mockUR := &mocks.MockUserRepository{}

	mockTR.
		On("Create", context.TODO(), mock.AnythingOfType("*domain.RefreshToken")).
		Return(nil)
	mockTR.
		On("Delete", context.TODO(), mockUser.ID, refreshToken).
		Return(mockErr)

	tokenService := service.NewTokenService(mockTR, mockUR, accessSecret, refreshSecret)
	ctx := context.TODO()

	_, err := tokenService.CreateRefresh(ctx, mockUser.ID, refreshToken)
	assert.ErrorAs(t, err, &mockErr)
}

func TestCreateRefreshNotExpired(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:       1,
		Email:    "Foo@Bar.com",
		Password: "FooBar",
	}

	mockTR := &mocks.MockTokenRepository{}
	mockUR := &mocks.MockUserRepository{}

	refresh, err := util.GenerateRefreshToken(mockUser.ID, &domain.TokenConfig{
		IAT:         time.Now(),
		ExpDuration: time.Minute,
		Secret:      refreshSecret,
	})
	assert.NoError(t, err)

	mockRefreshTokens := &[]domain.RefreshToken{{Refresh: refresh.Refresh}}

	mockTR.
		On("GetAll", context.TODO(), mockUser.ID).
		Return(mockRefreshTokens, nil)

	tokenService := service.NewTokenService(mockTR, mockUR, accessSecret, refreshSecret)
	ctx := context.TODO()

	accessToken, err := tokenService.CreateRefresh(ctx, mockUser.ID, refresh.Refresh)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken.Refresh)
	assert.Equal(t, refresh.Refresh, accessToken.Refresh)

	_, err = util.VerifyRefreshToken(accessToken.Refresh, refreshSecret)
	assert.NoError(t, err)
}

func TestExtractUser(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:    1,
		Email: "Foo@Bar.com",
	}

	mockTR := &mocks.MockTokenRepository{}
	mockUR := &mocks.MockUserRepository{}

	refresh, err := util.GenerateRefreshToken(mockUser.ID, &domain.TokenConfig{
		IAT:         time.Now(),
		ExpDuration: time.Minute,
		Secret:      refreshSecret,
	})
	assert.NoError(t, err)

	mockRefreshTokens := &[]domain.RefreshToken{{Refresh: refresh.Refresh}}

	mockTR.
		On("GetAll", context.TODO(), mockUser.ID).
		Return(mockRefreshTokens, nil)
	mockUR.
		On("GetByID", context.TODO(), mockUser.ID).
		Return(mockUser, nil)

	tokenService := service.NewTokenService(mockTR, mockUR, accessSecret, refreshSecret)
	ctx := context.TODO()

	accessToken, err := tokenService.CreateAccess(ctx, refresh.Refresh)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken.Access)

	user, err := tokenService.ExtractUser(ctx, accessToken.Access)
	assert.NoError(t, err)
	assert.Equal(t, mockUser, user)
}

func TestExtractUserNotAuth(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:    1,
		Email: "Foo@Bar.com",
	}

	expectedErr := domain.NewNotAuthorizedErr("")

	mockTR := &mocks.MockTokenRepository{}
	mockUR := &mocks.MockUserRepository{}

	accessToken, err := util.GenerateAccessToken(mockUser, &domain.TokenConfig{
		IAT:         time.Now(),
		ExpDuration: -time.Second,
		Secret:      accessSecret,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken.Access)

	tokenService := service.NewTokenService(mockTR, mockUR, accessSecret, refreshSecret)
	ctx := context.TODO()

	user, err := tokenService.ExtractUser(ctx, accessToken.Access)
	assert.ErrorAs(t, err, &expectedErr)
	assert.Nil(t, user)
}
