package hashing

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"time"
)

var good_fast_hash_seed = uint64(time.Now().UnixNano() / int64(time.Millisecond))

func MD5(str string) string {
	bytes := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", bytes)
}

func SHA1(file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	defer f.Close()
	if err != nil {
		return "", err
	}

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	bytes := h.Sum(nil)
	return hex.EncodeToString(bytes), nil
}

func SHA256(str string) string {
	bytes := sha256.New().Sum([]byte(str))
	return fmt.Sprintf("%x", bytes)
}

func SHA512(str string) string {
	bytes := sha512.New().Sum([]byte(str))
	return fmt.Sprintf("%x", bytes)
}

func Murmur3_128(str string) string {
	return ""
}

func Murmur3_32(str string) string {
	return ""
}
