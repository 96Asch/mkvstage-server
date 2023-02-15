package util

import (
	"testing"
	"time"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAccess(t *testing.T) {
	user := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
		Email:     "Foo@Bar.com",
	}

	now := time.Now()
	secret := "secret"

	at, err := GenerateAccessToken(user, &domain.TokenConfig{
		IAT:         now,
		ExpDuration: time.Second * time.Duration(15),
		Secret:      secret,
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, at.Access)
}

func TestGenerateRefresh(t *testing.T) {
	var ID int64 = 1

	now := time.Now()
	secret := "secret"

	at, err := GenerateRefreshToken(ID, &domain.TokenConfig{
		IAT:         now,
		ExpDuration: time.Second * time.Duration(15),
		Secret:      secret,
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, at.Refresh)
}

func TestVerifyAccessTokenCorrect(t *testing.T) {
	user := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
		Email:     "Foo@Bar.com",
	}

	now := time.Now()
	secret := "secret"

	at, err := GenerateAccessToken(user, &domain.TokenConfig{
		IAT:         now,
		ExpDuration: time.Second * time.Duration(15),
		Secret:      secret,
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, at.Access)

	claims, err := VerifyAccessToken(at.Access, secret)
	assert.NoError(t, err)

	assert.NotNil(t, claims.User)
	assert.Equal(t, jwt.NewNumericDate(now), claims.IssuedAt)
	assert.Equal(t, claims.User, user)
}

func TestVerifyAccessInvalidToken(t *testing.T) {
	now := time.Now()
	secret := []byte("secret")
	claims := refreshTokenClaims{
		UID: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(-time.Hour)),
			ID:        "Foobar",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	ss, err := token.SignedString(secret)
	assert.NoError(t, err)

	retClaims, err := VerifyAccessToken(ss, string(secret))
	assert.Error(t, err)
	assert.Nil(t, retClaims)
}

func TestVerifyRefreshTokenCorrect(t *testing.T) {
	user := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
		Email:     "Foo@Bar.com",
	}

	now := time.Now()
	secret := "secret"

	at, err := GenerateRefreshToken(user.ID, &domain.TokenConfig{
		IAT:         now,
		ExpDuration: time.Second * time.Duration(15),
		Secret:      secret,
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, at.Refresh)

	claims, err := VerifyRefreshToken(at.Refresh, secret)
	assert.NoError(t, err)

	assert.NotEmpty(t, claims.UID)
	assert.Equal(t, jwt.NewNumericDate(now), claims.IssuedAt)
	assert.Equal(t, claims.UID, user.ID)
}

func TestVerifyRefreshInvalidToken(t *testing.T) {
	now := time.Now()
	secret := []byte("secret")
	claims := refreshTokenClaims{
		UID: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(-time.Hour)),
			ID:        "Foobar",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	ss, err := token.SignedString(secret)
	assert.NoError(t, err)

	retClaims, err := VerifyRefreshToken(ss, string(secret))
	assert.Error(t, err)
	assert.Nil(t, retClaims)
}
