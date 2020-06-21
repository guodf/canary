package otp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"math"
)

func NewHOTP(secret string, c uint64) OTP {
	return &hOTP{
		secret: secret,
		c:      c,
		hash:   hmac.New(sha1.New, []byte(secret)),
	}
}

func (h *hOTP) Code(digit int) string {
	cBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(cBytes, h.c)
	h.hash.Write(cBytes)
	hs := h.hash.Sum(nil)
	offset := int(hs[len(hs)-1] & 0xf)
	hs[offset] &= 0x7f
	hs[offset+1] &= 0xff
	hs[offset+2] &= 0xff
	hs[offset+3] &= 0xff
	num := binary.BigEndian.Uint32(hs[offset : offset+4])
	sNum := num % uint32(math.Pow10(digit))
	return fmt.Sprintf("%0*d", digit, sNum)
}

func (h *hOTP) Reset(c int) {
	h.hash.Reset()
	h.c = uint64(c)
}

func (h *hOTP) Verify(code string) bool {
	return h.Code(len(code)) == code
}
