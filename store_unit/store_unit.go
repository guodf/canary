package store_unit

import "fmt"

var KB uint64 = 1024
var MB uint64 = KB * 1024
var GB uint64 = MB * 1024

func SizeString(size uint64) string {
	if size > GB {
		var left uint64 = size / GB
		var right uint64 = size % GB
		return toString(left, right, "GB")
	} else if size > MB {
		left := size / MB
		right := size % MB
		return toString(left, right, "MB")
	}
	return toString(size/KB, size%KB, "KB")

}

func toString(left uint64, right uint64, store string) string {
	if right == 0 {
		return fmt.Sprintf("%d %s", left, store)
	}
	if right > 100 {
		right = right / 100
	}
	return fmt.Sprintf("%d.%d %s", left, right, store)
}
