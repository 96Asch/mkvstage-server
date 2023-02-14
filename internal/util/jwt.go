package util

import (
	"errors"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func GenerateAccessToken(user *domain.User, config *domain.TokenConfig) (*domain.AccessToken, error) {

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.New("could not generate a uuid")
	}

	claims := domain.AccessTokenClaims{
		User: user,
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

	claims := domain.RefreshTokenClaims{
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
		UserID:  uid,
		ID:      id,
		Refresh: ss,
	}, nil
}

func VerifyToken[T domain.Claims](tokenString, secret string) (*T, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("incorrect signing method")
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(T)
	if !ok || !token.Valid {
		return nil, errors.New("token is not valid")
	}

	return &claims, nil
}
