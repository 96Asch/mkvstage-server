package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestChordSheetCorrect(t *testing.T) {
	mockCS := datatypes.JSON([]byte(`{"Verse" : "Test"}`))
	err := ValidateChordSheet(mockCS)
	assert.NoError(t, err)
}

func TestChordSheetInvalidJSON(t *testing.T) {
	mockCS := datatypes.JSON([]byte(`{"Verse" : "Test"`))
	err := ValidateChordSheet(mockCS)
	assert.Error(t, err)
}

func TestChordSheetInvalidTag(t *testing.T) {
	mockCS := datatypes.JSON([]byte(`{"V" : "Test"}`))
	err := ValidateChordSheet(mockCS)
	assert.Error(t, err)
}
