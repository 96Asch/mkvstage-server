package util_test

import (
	"testing"
	"time"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/util"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

const (
	secret = "secret"
)

func TestGenerateAccess(t *testing.T) {
	t.Parallel()

	now := time.Now()

	accessToken, err := util.GenerateAccessToken("foobar@barfoo.com", &domain.TokenConfig{
		IAT:         now,
		ExpDuration: time.Second * time.Duration(15),
		Secret:      secret,
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken.Access)
}

func TestVerifyAccessTokenCorrect(t *testing.T) {
	t.Parallel()

	email := "foobar@barfoo.com"

	now := time.Now()
	accessToken, err := util.GenerateAccessToken(email, &domain.TokenConfig{
		IAT:         now,
		ExpDuration: time.Second * time.Duration(15),
		Secret:      secret,
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken.Access)

	claims, err := util.VerifyAccessToken(accessToken.Access, secret)
	assert.NoError(t, err)

	assert.NotEmpty(t, claims.Email)
	assert.Equal(t, jwt.NewNumericDate(now), claims.IssuedAt)
	assert.Equal(t, claims.Email, email)
}

func TestVerifyAccessInvalidToken(t *testing.T) {
	t.Parallel()

	now := time.Now()
	secret := []byte("secret")
	claims := util.AccessTokenClaims{
		Email: "foobar",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(-time.Hour)),
			ID:        "Foobar",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	ss, err := token.SignedString(secret)
	assert.NoError(t, err)

	retClaims, err := util.VerifyAccessToken(ss, string(secret))
	assert.Error(t, err)
	assert.Nil(t, retClaims)
}
