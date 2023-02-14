package service

import (
	"context"
	"time"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/util"
)

type tokenService struct {
	tokenRepo     domain.TokenRepository
	accessSecret  string
	refreshSecret string
}

func NewTokenService(tr domain.TokenRepository, accessSecret, refreshSecret string) *tokenService {
	return &tokenService{
		tokenRepo:     tr,
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
	}
}

func (ts tokenService) ExtractUser(ctx context.Context, token *domain.AccessToken) (*domain.User, error) {
	at, err := util.VerifyToken[domain.AccessTokenClaims](token.Access, ts.accessSecret)

	if err != nil {
		return nil, domain.NewBadRequestErr(err.Error())
	}
	return at.User, nil
}

func (ts tokenService) CreateAccess(ctx context.Context, user *domain.User) (*domain.AccessToken, error) {

	config := domain.TokenConfig{
		IAT:         time.Now(),
		ExpDuration: time.Duration(15) * time.Minute,
		Secret:      ts.accessSecret,
	}

	accessToken, err := util.GenerateAccessToken(user, &config)
	if err != nil {
		return nil, domain.NewInternalErr()
	}

	return accessToken, nil
}

func (ts tokenService) CreateRefresh(ctx context.Context, user *domain.User) (*domain.RefreshToken, error) {
	config := domain.TokenConfig{
		IAT:         time.Now(),
		ExpDuration: time.Duration(72) * time.Hour,
		Secret:      ts.refreshSecret,
	}

	refreshToken, err := util.GenerateRefreshToken(user.ID, &config)
	if err != nil {
		return nil, domain.NewInternalErr()
	}

	// TODO Store Refresh in Redis

	return refreshToken, nil
}