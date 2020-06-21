package app

import (
	"path/filepath"
	"runtime"
)

// RootPath 获取程序启动目录
func RootPath() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(filepath.Dir(file))
}
