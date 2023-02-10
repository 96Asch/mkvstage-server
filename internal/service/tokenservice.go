package service

import "github.com/96Asch/mkvstage-server/internal/domain"

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
