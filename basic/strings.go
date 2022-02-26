package basic

import "strings"

func IsEmptyOrSpace(str string) (string, bool) {
	value := strings.TrimSpace(str)
	return value, len(value) == 0
}
