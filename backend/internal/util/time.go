package util

import (
	"fmt"
	"time"

	"github.com/96Asch/mkvstage-server/internal/domain"
)

const TimeFormat = time.RFC3339

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
