package basic

import "syscall"

func Utf16ToUtf8String(str string) string {
	data, _ := syscall.UTF16FromString(str)
	return string(UInt16ToUtf8Bytes(data))
}
func C_UTF16ToString(ptr uintptr) string {
	data := C_UIntPtrToUint16(ptr)
	return syscall.UTF16ToString(data)
}
