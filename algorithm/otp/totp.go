package otp

import "time"

func NewTOTP(secret string, seconds int) OTP {
	tc := resetTC(seconds)
	return &tOTP{
		NewHOTP(secret, tc).(*hOTP),
	}
}

func (h *tOTP) Code(digit int) string {
	return h.hotp.Code(digit)
}

func resetTC(seconds int) uint64 {
	return uint64(time.Now().Unix() / int64(seconds))
}

func (h *tOTP) Reset(seconds int) {
	tc := resetTC(seconds)
	h.hotp.Reset(int(tc))
}

func (h *tOTP) Verify(code string) bool {
	return h.hotp.Verify(code)
}
