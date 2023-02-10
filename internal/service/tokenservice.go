package service

import "github.com/96Asch/mkvstage-server/internal/domain"

type tokenService struct {
	tokenRepo domain.TokenRepository
}

func NewTokenService(tr domain.TokenRepository) *tokenService {
	return &tokenService{
		tokenRepo: tr,
	}
}
