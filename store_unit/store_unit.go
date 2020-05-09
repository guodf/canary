package store_unit

import "fmt"

var (
	KB = 1024
	MB = KB * 1024
	GB = MB * 1024
)

func SizeString(size int) string {
	if size > GB {
		left := size / GB
		right := size % GB
		return toString(left, right, "GB")
	} else if size > MB {
		left := size / MB
		right := size % MB
		return toString(left, right, "MB")
	}
	return toString(size/KB, size%KB, "KB")

}

func toString(left int, right int, store string) string {
	if right == 0 {
		return fmt.Sprintf("%d %s", left, store)
	}
	if right > 100 {
		right = right / 100
	}
	return fmt.Sprintf("%d.%d %s", left, right, store)
}
