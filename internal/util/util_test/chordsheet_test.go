package util_test

import (
	"testing"

	"github.com/96Asch/mkvstage-server/internal/util"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestChordSheetCorrect(t *testing.T) {
	t.Parallel()

	mockCS := datatypes.JSON([]byte(`{"Verse" : "Test"}`))
	err := util.ValidateChordSheet(mockCS)
	assert.NoError(t, err)
}

func TestChordSheetInvalidJSON(t *testing.T) {
	t.Parallel()

	mockCS := datatypes.JSON([]byte(`{"Verse" : "Test"`))
	err := util.ValidateChordSheet(mockCS)
	assert.Error(t, err)
}

func TestChordSheetInvalidTag(t *testing.T) {
	t.Parallel()

	mockCS := datatypes.JSON([]byte(`{"V" : "Test"}`))
	err := util.ValidateChordSheet(mockCS)
	assert.Error(t, err)
}
