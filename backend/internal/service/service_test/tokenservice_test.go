package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractEmail(t *testing.T) {
	const (
		refreshSecret = "refresh-secret"
		accessSecret  = "access-secret"
	)

	t.Run("Correct", func(t *testing.T) {
		t.Parallel()

		assert.Fail(t, "Not implemented")
	})
}
