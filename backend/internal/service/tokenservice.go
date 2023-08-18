package service

import (
	"context"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/96Asch/mkvstage-server/backend/internal/util"
)

type tokenService struct {
	accessSecret string
}

//revive:disable:unexported-return
func NewTokenService(accessSecret string) *tokenService {
	return &tokenService{
		accessSecret: accessSecret,
	}
}

func (ts tokenService) ExtractEmail(ctx context.Context, access string) (string, error) {
	accessToken, err := util.VerifyAccessToken(access, ts.accessSecret)

	if err != nil {
		return "", domain.NewNotAuthorizedErr(err.Error())
	}

	return accessToken.Email, nil
}
