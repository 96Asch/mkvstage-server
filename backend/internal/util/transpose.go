package util

const TransposeMax = 6
const TransposeMin = -5

func IsValidTranpose(transpose int16) bool {
	return transpose >= TransposeMin && transpose <= TransposeMax
}
