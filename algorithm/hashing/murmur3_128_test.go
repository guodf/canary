package hashing

import (
	"crypto"
	"crypto/md5"
	"fmt"
	"log"
	"testing"
)

func Test_murmur3_x64_128Hasher_Sum(t *testing.T) {
	hasher1 := murmur3_x64_128.NewHasher()
	str:="Hello World!"
	hasher1.Sum([]byte(str))
	log.Println(fmt.Sprintf("%x",hasher1.makeHash()))
	log.Println(MD5(str))
	log.Println(fmt.Sprintf("%x",md5.New().Sum([]byte(str))))
	log.Println(fmt.Sprintf("%x",crypto.MD5.New().Sum([]byte(str))))
}
