package otp

import "hash"

type OTP interface {
	Code(digit int) string
	Reset(data int)
	Verify(code string) bool
}

type hOTP struct {
	secret string
	c      uint64
	hash   hash.Hash
}

type tOTP struct {
	hotp *hOTP
}