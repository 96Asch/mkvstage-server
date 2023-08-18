package util

import (
	"errors"
	"fmt"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type AccessTokenClaims struct {
	Email string
	jwt.RegisteredClaims
}

func GenerateAccessToken(email string, config *domain.TokenConfig) (*domain.AccessToken, error) {
	accessUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.New("could not generate a uuid")
	}

	claims := AccessTokenClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(config.IAT),
			ExpiresAt: jwt.NewNumericDate(config.IAT.Add(config.ExpDuration)),
			ID:        accessUUID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return nil, errors.New("could not sign token")
	}

	return &domain.AccessToken{Access: ss}, nil
}

func VerifyAccessToken(tokenString, secret string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, errors.New(err.Error())
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok {
		return nil, fmt.Errorf("could not cast claims to accesstokenclaims")
	}

	return claims, nil
}
