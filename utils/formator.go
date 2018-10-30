package utils

import (
	"strconv"
)

func AddZeroForInt(in int, length int) string {
	s := strconv.Itoa(in)
	for ; len(s) < length; {
		s = "0" + s
	}
	return s
}

func AddZeroForStr(s string, length int) string {
	for ; len(s) < length; {
		s = "0" + s
	}
	return s
}
