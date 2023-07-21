package util_test

import (
	"testing"

	"github.com/96Asch/mkvstage-server/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestIsValidTranspose(t *testing.T) {
	t.Parallel()

	expecteds := map[int16]bool{
		util.TransposeMin:     true,
		0:                     true,
		util.TransposeMax:     true,
		util.TransposeMin - 1: false,
		util.TransposeMax + 1: false,
	}

	for key, val := range expecteds {
		assert.Equal(t, val, util.IsValidTranpose(key))
	}
}
