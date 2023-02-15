package util

import (
	"errors"
	"fmt"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type accessTokenClaims struct {
	User *domain.User `json:"user"`
	jwt.RegisteredClaims
}

type refreshTokenClaims struct {
	UID int64 `json:"uid"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(user *domain.User, config *domain.TokenConfig) (*domain.AccessToken, error) {

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.New("could not generate a uuid")
	}

	claims := accessTokenClaims{
		User: &domain.User{
			ID:           user.ID,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Email:        user.Email,
			ProfileColor: user.ProfileColor,
			Permission:   user.Permission,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(config.IAT),
			ExpiresAt: jwt.NewNumericDate(config.IAT.Add(config.ExpDuration)),
			ID:        id.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return nil, errors.New("could not sign token")
	}

	return &domain.AccessToken{Access: ss}, nil
}

func GenerateRefreshToken(uid int64, config *domain.TokenConfig) (*domain.RefreshToken, error) {

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.New("could not generate a uuid")
	}

	claims := refreshTokenClaims{
		UID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(config.IAT),
			ExpiresAt: jwt.NewNumericDate(config.IAT.Add(config.ExpDuration)),
			ID:        id.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return nil, errors.New("could not sign token")
	}

	return &domain.RefreshToken{
		UserID:             uid,
		ID:                 id,
		Refresh:            ss,
		ExpirationDuration: config.ExpDuration,
	}, nil
}

func VerifyAccessToken(tokenString, secret string) (*accessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &accessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	claims, ok := token.Claims.(*accessTokenClaims)
	if !ok {
		return nil, fmt.Errorf("could not cast claims to accesstokenclaims")
	}

	return claims, nil
}
