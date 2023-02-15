package service

import (
	"context"
	"time"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/util"
)

type tokenService struct {
	tokenRepo     domain.TokenRepository
	userRepo      domain.UserRepository
	accessSecret  string
	refreshSecret string
}

func NewTokenService(tr domain.TokenRepository, ur domain.UserRepository, accessSecret, refreshSecret string) *tokenService {
	return &tokenService{
		tokenRepo:     tr,
		userRepo:      ur,
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
	}
}

func (ts tokenService) ExtractUser(ctx context.Context, token *domain.AccessToken) (*domain.User, error) {
	at, err := util.VerifyAccessToken(token.Access, ts.accessSecret)

	if err != nil {
		return nil, domain.NewNotAuthorizedErr(err.Error())
	}
	return at.User, nil
}

func (ts tokenService) CreateAccess(ctx context.Context, currentRefresh string) (*domain.AccessToken, error) {

	if currentRefresh == "" {
		return nil, domain.NewBadRequestErr("no token provided")
	}

	claims, err := util.VerifyRefreshToken(currentRefresh, ts.refreshSecret)
	if err != nil {
		return nil, domain.NewNotAuthorizedErr("token invalid")
	}

	user, err := ts.userRepo.GetByID(ctx, claims.UID)
	if err != nil {
		return nil, err
	}

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

func (ts tokenService) CreateRefresh(ctx context.Context, uid int64, currentRefresh string) (*domain.RefreshToken, error) {

	// if currentToken is provided, delete the token if the token is not valid,
	// otherwise return the given token
	if currentRefresh != "" {
		_, err := util.VerifyRefreshToken(currentRefresh, ts.refreshSecret)
		if err == nil {
			return &domain.RefreshToken{Refresh: currentRefresh}, nil
		}

		err = ts.tokenRepo.Delete(ctx, &domain.RefreshToken{Refresh: currentRefresh, UserID: uid})
		if err != nil {
			return nil, domain.NewInternalErr()
		}
	}

	config := domain.TokenConfig{
		IAT:         time.Now(),
		ExpDuration: time.Duration(72) * time.Hour,
		Secret:      ts.refreshSecret,
	}

	refreshToken, err := util.GenerateRefreshToken(uid, &config)
	if err != nil {
		return nil, domain.NewInternalErr()
	}

	err = ts.tokenRepo.Create(ctx, refreshToken)
	if err != nil {
		return nil, domain.NewInternalErr()
	}

	return refreshToken, nil
}
