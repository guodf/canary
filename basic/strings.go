package basic

import (
	"strings"
	"unicode/utf16"
	"unicode/utf8"
	"unsafe"
)

func IsEmptyOrSpace(str string) (string, bool) {
	value := strings.TrimSpace(str)
	return value, len(value) == 0
}

func C_UIntPtrToUint16(ptr uintptr) []uint16 {
	i := 0
	data := []uint16{}
	for {
		b := *(*uint16)(unsafe.Pointer(ptr + uintptr(i)))
		if b == 0 {
			break
		}
		data = append(data, b)
		i += 2
	}
	return data
}

func UInt16ToUtf8Bytes(data []uint16) []byte {
	runes := utf16.Decode(data)
	utf8Bytes := make([]byte, len(runes)*3) // max length of a UTF-8 encoded rune is 3 bytes
	i := 0
	for _, r := range runes {
		i += utf8.EncodeRune(utf8Bytes[i:], r)
	}
	return utf8Bytes[0:i]
}

func C_Utf16ToUtf8(ptr uintptr) string {
	data := C_UIntPtrToUint16(ptr)
	return string(UInt16ToUtf8Bytes(data))
}
