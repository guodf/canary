package hashing

type hahser interface {
	//putByte(b byte) hahser
	Sum(bytes []byte)
	//putShort(s uint8) hahser
	//putInt(i int32) hahser
	//putLong(l int64) hahser
	//putFloat(f float32) hahser
	//putDouble(d float64) hahser
	//putBoolean(b bool) hahser
	//putChar(c uint8) hahser
	//hash() HashCode
	makeHash() []byte
}

type HashCode interface {
}
