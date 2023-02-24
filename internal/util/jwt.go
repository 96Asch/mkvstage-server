package util

import (
	"errors"
	"fmt"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type AccessTokenClaims struct {
	User *domain.User `json:"user"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UID int64 `json:"uid"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(user *domain.User, config *domain.TokenConfig) (*domain.AccessToken, error) {
	accessUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.New("could not generate a uuid")
	}

	claims := AccessTokenClaims{
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

func GenerateRefreshToken(uid int64, config *domain.TokenConfig) (*domain.RefreshToken, error) {
	refreshUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.New("could not generate a uuid")
	}

	claims := RefreshTokenClaims{
		UID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(config.IAT),
			ExpiresAt: jwt.NewNumericDate(config.IAT.Add(config.ExpDuration)),
			ID:        refreshUUID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedRefreshToken, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return nil, errors.New("could not sign token")
	}

	return &domain.RefreshToken{
		UserID:             uid,
		ID:                 refreshUUID,
		Refresh:            signedRefreshToken,
		ExpirationDuration: config.ExpDuration,
	}, nil
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

func VerifyRefreshToken(tokenString, secret string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, errors.New(err.Error())
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok {
		return nil, fmt.Errorf("could not cast claims to accesstokenclaims")
	}

	return claims, nil
}
