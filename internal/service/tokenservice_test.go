package service

import (
	"context"
	"testing"
	"time"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAccessCorrect(t *testing.T) {

	mockUser := &domain.User{
		ID:       1,
		Email:    "Foo@Bar.com",
		Password: "FooBar",
	}

	mockTR := new(mocks.MockTokenRepository)

	mockUR := new(mocks.MockUserRepository)
	mockUR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(mockUser, nil)
	TS := NewTokenService(mockTR, mockUR, "access-secret", "refresh-secret")

	refresh, err := util.GenerateRefreshToken(mockUser.ID, &domain.TokenConfig{
		IAT:         time.Now(),
		ExpDuration: time.Minute,
		Secret:      TS.refreshSecret,
	})
	assert.NoError(t, err)

	mockRefreshTokens := &[]domain.RefreshToken{
		{
			Refresh: refresh.Refresh,
		},
	}
	mockTR.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(mockRefreshTokens, nil)

	ctx := context.TODO()
	at, err := TS.CreateAccess(ctx, refresh.Refresh)

	assert.NoError(t, err)
	assert.NotEmpty(t, at.Access)

	claims, err := util.VerifyAccessToken(at.Access, TS.accessSecret)
	assert.NoError(t, err)
	assert.NotNil(t, claims.User)
	assert.Empty(t, claims.User.Password)
	mockUser.Password = ""
	assert.Equal(t, mockUser, claims.User)

}

func TestCreateAccessRepoErr(t *testing.T) {

	mockUser := &domain.User{
		ID:       1,
		Email:    "Foo@Bar.com",
		Password: "FooBar",
	}

	mockTR := new(mocks.MockTokenRepository)

	mockUR := new(mocks.MockUserRepository)
	mockUR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(mockUser, nil)
	TS := NewTokenService(mockTR, mockUR, "access-secret", "refresh-secret")

	refresh, err := util.GenerateRefreshToken(mockUser.ID, &domain.TokenConfig{
		IAT:         time.Now(),
		ExpDuration: time.Minute,
		Secret:      TS.refreshSecret,
	})
	assert.NoError(t, err)

	mockErr := domain.NewInternalErr()
	mockTR.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(nil, mockErr)

	ctx := context.TODO()
	_, err = TS.CreateAccess(ctx, refresh.Refresh)

	assert.ErrorAs(t, err, &mockErr)
}

func TestCreateAccessRefreshNotInRepo(t *testing.T) {

	mockUser := &domain.User{
		ID:       1,
		Email:    "Foo@Bar.com",
		Password: "FooBar",
	}

	mockTR := new(mocks.MockTokenRepository)

	mockUR := new(mocks.MockUserRepository)
	mockUR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(mockUser, nil)
	TS := NewTokenService(mockTR, mockUR, "access-secret", "refresh-secret")

	refresh, err := util.GenerateRefreshToken(mockUser.ID, &domain.TokenConfig{
		IAT:         time.Now(),
		ExpDuration: time.Minute,
		Secret:      TS.refreshSecret,
	})
	assert.NoError(t, err)

	mockErr := domain.NewInternalErr()
	mockRefreshTokens := &[]domain.RefreshToken{}
	mockTR.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(mockRefreshTokens, nil)

	ctx := context.TODO()
	_, err = TS.CreateAccess(ctx, refresh.Refresh)

	assert.ErrorAs(t, err, &mockErr)
}

func TestCreateAccessInvalidRefresh(t *testing.T) {

	mockUser := &domain.User{
		ID:       1,
		Email:    "Foo@Bar.com",
		Password: "FooBar",
	}

	mockTR := new(mocks.MockTokenRepository)
	mockUR := new(mocks.MockUserRepository)
	TS := NewTokenService(mockTR, mockUR, "access-secret", "refresh-secret")

	refresh, err := util.GenerateRefreshToken(mockUser.ID, &domain.TokenConfig{
		IAT:         time.Now(),
		ExpDuration: -time.Minute,
		Secret:      TS.refreshSecret,
	})
	assert.NoError(t, err)

	ctx := context.TODO()
	at, err := TS.CreateAccess(ctx, refresh.Refresh)

	expectedErr := domain.NewNotAuthorizedErr("")
	assert.ErrorAs(t, err, &expectedErr)
	assert.Nil(t, at)
}

func TestCreateAccessRefreshEmpty(t *testing.T) {

	mockTR := new(mocks.MockTokenRepository)
	mockUR := new(mocks.MockUserRepository)
	TS := NewTokenService(mockTR, mockUR, "access-secret", "refresh-secret")

	ctx := context.TODO()
	at, err := TS.CreateAccess(ctx, "")

	expectedErr := domain.NewBadRequestErr("")
	assert.ErrorAs(t, err, &expectedErr)
	assert.Nil(t, at)
}

func TestCreateRefreshCorrect(t *testing.T) {

	mockUser := &domain.User{
		ID:       1,
		Email:    "Foo@Bar.com",
		Password: "FooBar",
	}

	mockUR := new(mocks.MockUserRepository)
	mockTR := new(mocks.MockTokenRepository)
	mockTR.
		On("Create", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*domain.RefreshToken")).
		Return(nil)

	TS := NewTokenService(mockTR, mockUR, "access-secret", "refresh-secret")

	ctx := context.TODO()
	at, err := TS.CreateRefresh(ctx, mockUser.ID, "")

	assert.NoError(t, err)
	assert.NotEmpty(t, at.Refresh)

	claims, err := util.VerifyRefreshToken(at.Refresh, TS.refreshSecret)
	assert.NoError(t, err)
	assert.Equal(t, mockUser.ID, claims.UID)
}

func TestCreateRefreshDeleteErr(t *testing.T) {

	mockUser := &domain.User{
		ID:       1,
		Email:    "Foo@Bar.com",
		Password: "FooBar",
	}

	refreshToken := "refresh-token"
	mockErr := domain.NewInternalErr()
	mockUR := new(mocks.MockUserRepository)
	mockTR := new(mocks.MockTokenRepository)
	mockTR.
		On("Create", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*domain.RefreshToken")).
		Return(nil)
	mockTR.
		On("Delete", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID, refreshToken).
		Return(mockErr)

	TS := NewTokenService(mockTR, mockUR, "access-secret", "refresh-secret")

	ctx := context.TODO()
	_, err := TS.CreateRefresh(ctx, mockUser.ID, refreshToken)

	assert.ErrorAs(t, err, &mockErr)
}

func TestCreateRefreshNotExpired(t *testing.T) {

	mockUser := &domain.User{
		ID:       1,
		Email:    "Foo@Bar.com",
		Password: "FooBar",
	}

	mockUR := new(mocks.MockUserRepository)
	mockTR := new(mocks.MockTokenRepository)
	TS := NewTokenService(mockTR, mockUR, "access-secret", "refresh-secret")

	refresh, err := util.GenerateRefreshToken(mockUser.ID, &domain.TokenConfig{
		IAT:         time.Now(),
		ExpDuration: time.Minute,
		Secret:      TS.refreshSecret,
	})
	assert.NoError(t, err)

	mockRefreshTokens := &[]domain.RefreshToken{
		{
			Refresh: refresh.Refresh,
		},
	}

	mockTR.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(mockRefreshTokens, nil)

	ctx := context.TODO()
	at, err := TS.CreateRefresh(ctx, mockUser.ID, refresh.Refresh)

	assert.NoError(t, err)
	assert.NotEmpty(t, at.Refresh)
	assert.Equal(t, refresh.Refresh, at.Refresh)

	_, err = util.VerifyRefreshToken(at.Refresh, TS.refreshSecret)
	assert.NoError(t, err)

}

func TestExtractUser(t *testing.T) {
	mockUser := &domain.User{
		ID:    1,
		Email: "Foo@Bar.com",
	}

	mockUR := new(mocks.MockUserRepository)
	mockUR.
		On("GetByID", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(mockUser, nil)
	mockTR := new(mocks.MockTokenRepository)
	TS := NewTokenService(mockTR, mockUR, "access-secret", "refresh-secret")

	refresh, err := util.GenerateRefreshToken(mockUser.ID, &domain.TokenConfig{
		IAT:         time.Now(),
		ExpDuration: time.Minute,
		Secret:      TS.refreshSecret,
	})
	assert.NoError(t, err)

	mockRefreshTokens := &[]domain.RefreshToken{
		{
			Refresh: refresh.Refresh,
		},
	}

	mockTR.
		On("GetAll", mock.AnythingOfType("*context.emptyCtx"), mockUser.ID).
		Return(mockRefreshTokens, nil)

	ctx := context.TODO()
	at, err := TS.CreateAccess(ctx, refresh.Refresh)
	assert.NoError(t, err)
	assert.NotEmpty(t, at.Access)

	user, err := TS.ExtractUser(ctx, at.Access)
	assert.NoError(t, err)
	assert.Equal(t, mockUser, user)
}

func TestExtractUserNotAuth(t *testing.T) {
	mockUser := &domain.User{
		ID:    1,
		Email: "Foo@Bar.com",
	}

	mockUR := new(mocks.MockUserRepository)
	mockTR := new(mocks.MockTokenRepository)
	TS := NewTokenService(mockTR, mockUR, "access-secret", "refresh-secret")

	ctx := context.TODO()
	at, err := util.GenerateAccessToken(mockUser, &domain.TokenConfig{
		IAT:         time.Now(),
		ExpDuration: -time.Second,
		Secret:      TS.accessSecret,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, at.Access)

	user, err := TS.ExtractUser(ctx, at.Access)
	expectedErr := domain.NewNotAuthorizedErr("")
	assert.ErrorAs(t, err, &expectedErr)
	assert.Nil(t, user)
}
