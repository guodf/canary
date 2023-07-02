package grandom

import (
	"math"
	"math/rand"
)

// 类型
// 范围
// 值

// 平台相关

// Int 生成int类型整数
func Int() int {
	return rand.Int()
}

// Uint8 生成uint8类型整数
func Uint8() uint8 {
	return uint8(rand.Int63n(math.MaxUint8 + 1))
}

// Uint16 生成uint16类型整数
func Uint16() uint16 {
	return uint16(rand.Int63n(math.MaxUint16 + 1))
}

// Uint32 生成uint32类型整数
func Uint32() uint32 {
	return uint32(rand.Int63n(math.MaxUint32 + 1))
}
