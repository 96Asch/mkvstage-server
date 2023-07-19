package util

import (
	"time"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"gorm.io/datatypes"
)

const DateFormat = time.DateOnly
const TimeFormat = time.DateTime

func StringToDate(timestring string) (datatypes.Date, error) {
	date, err := time.Parse(DateFormat, timestring)
	date = date.UTC().Round(0)

	if err != nil {
		return datatypes.Date{}, domain.NewInternalErr()
	}

	return datatypes.Date(date), nil
}

func StringToTime(timestring string) (time.Time, error) {
	date, err := time.Parse(TimeFormat, timestring)

	if err != nil {
		return time.Time{}, domain.NewInternalErr()
	}

	return date.Truncate(time.Minute), nil
}
