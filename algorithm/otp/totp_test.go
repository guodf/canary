package otp

import "testing"

func Test_tOTP_Code(t *testing.T) {
	tests := []struct {
		name    string
		secret  string
		seconds int
		digit   int
		want    string
	}{
		{
			name:    "aaa",
			secret:  "this is secret!",
			seconds: 5,
			digit:   6,
		},
		{
			name:    "bbb",
			secret:  "this is secret!",
			seconds: 5,
			digit:   6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewTOTP(tt.secret, tt.seconds)
			r := h.Code(tt.digit)
			t.Log(r)
		})
	}
}
