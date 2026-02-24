package unpack

import (
	"errors"
	"strings"
)

var errInvalidString = errors.New("invalid string")

const asciiNumShift = 48

func Unpack(s string) (string, error) {
	runes := []rune(s)
	n := len(runes)
	if n == 0 {
		return "", nil
	}
	var res strings.Builder
	var escaped bool
	for i := 0; i < n; i++ {
		r := runes[i]
		if r == '\\' {
			escaped = true
			continue
		}
		if isNum(r) && !escaped {
			return "", errInvalidString
		}
		if i+1 < n && isNum(runes[i+1]) {
			for range atoi(runes[i+1]) {
				res.WriteRune(r)
			}
			i++
		} else {
			res.WriteRune(r)
		}
		escaped = false
	}
	return res.String(), nil
}

func isNum(s rune) bool {
	return s >= '0' && s <= '9'
}

func atoi(s rune) int {
	return int(s - asciiNumShift)
}
