package emath

import (
	"fmt"
	"strconv"
	"strings"
)

func Len(value int) int {
	count := 0
	for value != 0 {
		value = value / 10
		count++
	}
	return count
}

func SplitFloat(value float64) (int, float64) {
	left := int(value)
	rightStr := strings.Split(fmt.Sprintf("%f", value), ".")[1]
	right, _ := strconv.ParseFloat("0."+rightStr, len(rightStr))
	return left, right
}
