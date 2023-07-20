package util

import (
	"fmt"
	"time"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"gorm.io/datatypes"
)

const DateFormat = time.DateOnly
const TimeFormat = time.RFC3339

func StringToDate(timestring string) (datatypes.Date, error) {
	date, err := time.Parse(DateFormat, timestring)

	if err != nil {
		return datatypes.Date{}, domain.NewInternalErr()
	}

	return datatypes.Date(date.UTC().Truncate(time.Minute)), nil
}

func StringToTime(timestring string) (time.Time, error) {
	if len(timestring) <= 0 {
		return time.Time{}, nil
	}

	date, err := time.Parse(TimeFormat, timestring)

	if err != nil {
		return time.Time{}, domain.NewBadRequestErr(fmt.Sprintf("%s must be in RFC3339 format", timestring))
	}

	return date.UTC().Truncate(time.Minute), nil
}
