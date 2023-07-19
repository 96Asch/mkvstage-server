package util_test

import (
	"testing"
	"time"

	"github.com/96Asch/mkvstage-server/internal/util"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestStringToDate(t *testing.T) {
	t.Parallel()

	expecteds := map[string]datatypes.Date{
		"1996-02-08": datatypes.Date(time.Date(1996, 2, 8, 0, 0, 0, 0, time.UTC)),
		"2023-05-13": datatypes.Date(time.Date(2023, 5, 13, 0, 0, 0, 0, time.UTC)),
		"2000-01-01": datatypes.Date(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
	}

	for key, val := range expecteds {
		date, err := util.StringToDate(key)
		assert.NoError(t, err)
		assert.Equal(t, val, date)
	}

	fails := map[string]datatypes.Date{
		"foobar":       {},
		"2023--13":     {},
		"2000-0132-01": {},
	}

	for key, val := range fails {
		date, err := util.StringToDate(key)
		assert.Error(t, err)
		assert.Equal(t, val, date)
	}
}

func TestStringToTime(t *testing.T) {
	t.Parallel()

	expecteds := map[string]time.Time{
		"1996-02-08 13:00:00": time.Date(1996, time.February, 8, 13, 0, 0, 0, time.UTC),
		"2023-05-13 00:00:00": time.Date(2023, time.May, 13, 0, 0, 0, 0, time.UTC),
		"2000-01-01 00:00:00": time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	for key, val := range expecteds {
		date, err := util.StringToTime(key)
		assert.NoError(t, err)
		assert.Equal(t, val, date)
	}

	fails := map[string]time.Time{
		"foobar":             {},
		"2023--13":           {},
		"2000-02-01 00:00:0": {},
	}

	for key, val := range fails {
		date, err := util.StringToTime(key)
		assert.Error(t, err)
		assert.Equal(t, val, date)
	}
}
