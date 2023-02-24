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

	user := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
		Email:     "Foo@Bar.com",
	}

	now := time.Now()

	accessToken, err := util.GenerateAccessToken(user, &domain.TokenConfig{
		IAT:         now,
		ExpDuration: time.Second * time.Duration(15),
		Secret:      secret,
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken.Access)
}

func TestGenerateRefresh(t *testing.T) {
	t.Parallel()

	ID := int64(1)
	now := time.Now()
	refreshToken, err := util.GenerateRefreshToken(ID, &domain.TokenConfig{
		IAT:         now,
		ExpDuration: time.Second * time.Duration(15),
		Secret:      secret,
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, refreshToken.Refresh)
}

func TestVerifyAccessTokenCorrect(t *testing.T) {
	t.Parallel()

	user := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
		Email:     "Foo@Bar.com",
	}

	now := time.Now()
	accessToken, err := util.GenerateAccessToken(user, &domain.TokenConfig{
		IAT:         now,
		ExpDuration: time.Second * time.Duration(15),
		Secret:      secret,
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken.Access)

	claims, err := util.VerifyAccessToken(accessToken.Access, secret)
	assert.NoError(t, err)

	assert.NotNil(t, claims.User)
	assert.Equal(t, jwt.NewNumericDate(now), claims.IssuedAt)
	assert.Equal(t, claims.User, user)
}

func TestVerifyAccessInvalidToken(t *testing.T) {
	t.Parallel()

	now := time.Now()
	secret := []byte("secret")
	claims := util.RefreshTokenClaims{
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

	retClaims, err := util.VerifyAccessToken(ss, string(secret))
	assert.Error(t, err)
	assert.Nil(t, retClaims)
}

func TestVerifyRefreshTokenCorrect(t *testing.T) {
	t.Parallel()

	user := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
		Email:     "Foo@Bar.com",
	}

	now := time.Now()

	refreshToken, err := util.GenerateRefreshToken(user.ID, &domain.TokenConfig{
		IAT:         now,
		ExpDuration: time.Second * time.Duration(15),
		Secret:      secret,
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, refreshToken.Refresh)

	claims, err := util.VerifyRefreshToken(refreshToken.Refresh, secret)
	assert.NoError(t, err)

	assert.NotEmpty(t, claims.UID)
	assert.Equal(t, jwt.NewNumericDate(now), claims.IssuedAt)
	assert.Equal(t, claims.UID, user.ID)
}

func TestVerifyRefreshInvalidToken(t *testing.T) {
	t.Parallel()

	now := time.Now()
	secret := []byte("secret")
	claims := util.RefreshTokenClaims{
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

	retClaims, err := util.VerifyRefreshToken(ss, string(secret))
	assert.Error(t, err)
	assert.Nil(t, retClaims)
}
