package service

import (
	"context"
	"testing"
	"time"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccessCorrect(t *testing.T) {

	mockUser := &domain.User{
		ID:       1,
		Email:    "Foo@Bar.com",
		Password: "FooBar",
	}

	mockTR := new(mocks.MockTokenRepository)
	TS := NewTokenService(mockTR, "access-secret", "refresh-secret")

	ctx := context.TODO()
	at, err := TS.CreateAccess(ctx, mockUser)

	assert.NoError(t, err)
	assert.NotEmpty(t, at.Access)

	claims, err := util.VerifyAccessToken(at.Access, TS.accessSecret)
	assert.NoError(t, err)
	assert.NotNil(t, claims.User)
	assert.Empty(t, claims.User.Password)
	mockUser.Password = ""
	assert.Equal(t, mockUser, claims.User)

}

func TestExtractUser(t *testing.T) {
	mockUser := &domain.User{
		ID:    1,
		Email: "Foo@Bar.com",
	}

	mockTR := new(mocks.MockTokenRepository)
	TS := NewTokenService(mockTR, "access-secret", "refresh-secret")

	ctx := context.TODO()
	at, err := TS.CreateAccess(ctx, mockUser)
	assert.NoError(t, err)
	assert.NotEmpty(t, at.Access)

	user, err := TS.ExtractUser(ctx, at)
	assert.NoError(t, err)
	assert.Equal(t, mockUser, user)
}

func TestExtractUserNotAuth(t *testing.T) {
	mockUser := &domain.User{
		ID:    1,
		Email: "Foo@Bar.com",
	}

	mockTR := new(mocks.MockTokenRepository)
	TS := NewTokenService(mockTR, "access-secret", "refresh-secret")

	ctx := context.TODO()
	at, err := util.GenerateAccessToken(mockUser, &domain.TokenConfig{
		IAT:         time.Now(),
		ExpDuration: -time.Second,
		Secret:      TS.accessSecret,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, at.Access)

	user, err := TS.ExtractUser(ctx, at)
	expectedErr := domain.NewNotAuthorizedErr("")
	assert.ErrorAs(t, err, &expectedErr)
	assert.Nil(t, user)
}
