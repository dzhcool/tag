package utils

import (
	"strconv"
)

func ToInt(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return val
}
