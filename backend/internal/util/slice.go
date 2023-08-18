package util

import (
	"strconv"
	"strings"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
)

func StringToInt64Slice(field string) ([]int64, error) {
	if field == "" {
		return []int64{}, nil
	}

	split := strings.Split(field, ",")

	if len(split) == 0 {
		return []int64{}, nil
	}

	intSlice := make([]int64, len(split))

	for idx, val := range split {
		number, err := strconv.Atoi(val)

		if err != nil {
			return []int64{}, domain.NewBadRequestErr(err.Error())
		}

		intSlice[idx] = int64(number)
	}

	return intSlice, nil
}

func StringToUintSlice(field string) ([]uint, error) {
	if field == "" {
		return []uint{}, nil
	}

	split := strings.Split(field, ",")

	if len(split) == 0 {
		return []uint{}, nil
	}

	intSlice := make([]uint, len(split))

	for idx, val := range split {
		number, err := strconv.Atoi(val)

		if err != nil {
			return []uint{}, domain.NewBadRequestErr(err.Error())
		}

		intSlice[idx] = uint(number)
	}

	return intSlice, nil
}
