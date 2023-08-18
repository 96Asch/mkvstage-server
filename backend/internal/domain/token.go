package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID                 uuid.UUID     `json:"-"`
	UserID             int64         `json:"-"`
	Refresh            string        `json:"refresh"`
	ExpirationDuration time.Duration `json:"-"`
}

type AccessToken struct {
	Access string `json:"access"`
}

type Tokens struct {
	AccessToken
	RefreshToken
}

type TokenConfig struct {
	IAT         time.Time
	ExpDuration time.Duration
	Secret      string
}

type TokenService interface {
	ExtractEmail(ctx context.Context, access string) (string, error)
}
