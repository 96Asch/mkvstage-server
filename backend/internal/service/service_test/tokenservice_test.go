package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/96Asch/mkvstage-server/backend/internal/service"
	"github.com/96Asch/mkvstage-server/backend/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestExtractEmail(t *testing.T) {
	const (
		accessSecret = "access-secret"
	)

	expEmail := "foobar@barfoo.com"
	correctConfig := &domain.TokenConfig{
		IAT:         time.Now(),
		ExpDuration: time.Hour,
		Secret:      accessSecret,
	}

	correctAccess, err := util.GenerateAccessToken(expEmail, correctConfig)
	assert.NoError(t, err)

	wrongConfig := &domain.TokenConfig{
		IAT:         time.Now(),
		ExpDuration: -time.Hour,
		Secret:      accessSecret,
	}

	wrongAccess, err := util.GenerateAccessToken(expEmail, wrongConfig)
	assert.NoError(t, err)

	tokenService := service.NewTokenService(accessSecret)

	t.Run("Correct", func(t *testing.T) {
		t.Parallel()

		email, err := tokenService.ExtractEmail(context.TODO(), correctAccess.Access)
		assert.NoError(t, err)
		assert.Equal(t, expEmail, email)
	})

	t.Run("Fail verify err", func(t *testing.T) {
		t.Parallel()

		email, err := tokenService.ExtractEmail(context.TODO(), wrongAccess.Access)
		assert.Error(t, err)
		assert.Empty(t, email)
	})
}
