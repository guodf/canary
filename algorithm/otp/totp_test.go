package otp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func toUint32(bytes []byte) uint32 {
	return (uint32(bytes[0]) << 24) + (uint32(bytes[1]) << 16) +
		(uint32(bytes[2]) << 8) + uint32(bytes[3])
}

func oneTimePassword(key []byte, value []byte) uint32 {
	// sign the value using HMAC-SHA1
	hmacSha1 := hmac.New(sha1.New, key)
	hash := hmacSha1.Sum(nil)
	for _, v := range hash {
		fmt.Println(v)
	}
	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	offset := hash[len(hash)-1] & 0x0F

	// get a 32-bit (4-byte) chunk from the hash starting at offset
	hashParts := hash[offset : offset+4]

	// ignore the most significant bit as per RFC 4226
	hashParts[0] = hashParts[0] & 0x7F

	number := toUint32(hashParts)

	// size to 6 digits
	// one million is the first number with 7 digits so the remainder
	// of the division will always return < 7 digits
	pwd := number % 1000000

	return pwd
}

func Test_tOTP_Code(t *testing.T) {

	// decode the key from the first argument
	inputNoSpaces := strings.Replace("FB6HLTLQFP3M7K3L", " ", "", -1)
	inputNoSpacesUpper := strings.ToUpper(inputNoSpaces)
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(inputNoSpacesUpper)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// generate a one-time password using the time at 30-second intervals
	epochSeconds := time.Now().Unix()
	pwd := oneTimePassword(key, toBytes(epochSeconds/30))

	secondsRemaining := 30 - (epochSeconds % 30)
	fmt.Printf("%06d (%d second(s) remaining) %d\n", pwd, secondsRemaining, epochSeconds/30)

	tests := []struct {
		name    string
		secret  string
		seconds int
		digit   int
		want    string
	}{
		{
			name:    "aaa",
			secret:  "FB6HLTLQFP3M7K3L",
			seconds: 30,
			digit:   6,
		},
		{
			name:    "bbb",
			secret:  "FB6HLTLQFP3M7K3L",
			seconds: 30,
			digit:   6,
		},
	}
	for _, tt := range tests {
		secret, _ := base32.StdEncoding.DecodeString(tt.secret)
		h := NewTOTP(string(secret), tt.seconds)
		r := h.Code(tt.digit)
		t.Log(r)

	}
}
