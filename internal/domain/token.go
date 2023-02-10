package domain

import (
	"github.com/google/uuid"
)

type RefreshToken struct {
	ID      uuid.UUID `json:"-"`
	UserID  int64     `json:"-"`
	Refresh string    `json:"access"`
}

type AccessToken struct {
	Access string `json:"-"`
}

type Tokens struct {
	AccessToken
	RefreshToken
}

type TokenService interface {
}

type TokenRepository interface {
}
