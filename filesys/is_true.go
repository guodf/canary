package filesys

import (
	"log"
	"os"
)

// Exists 判断路径是否存在，忽略权限
func Exists(path string) bool {
	_, e := os.Stat(path)
	return e == nil || os.IsExist(e)
}

// IsDir 判断是目录，忽略权限
func IsDir(path string) bool {
	fileInfo, e := os.Stat(path)
	if e != nil {
		log.Println(e)
		return false
	}
	return fileInfo.IsDir()
}
