package filesys

import (
	"io/fs"
	"path/filepath"
)

// FindFilesInExt 根据扩展名查找文件
func FindFilesInExt(dirPath string, exts []string, maxDepth int) []string {
	var files []string
	if !IsDir(dirPath) {
		return files
	}
	ScanByBreadth(dirPath, maxDepth, func(fullPath string, fsInfo fs.FileInfo) bool {
		if !fsInfo.IsDir() {
			for _, ext := range exts {
				if ext == filepath.Ext(fullPath) {
					files = append(files, fullPath)
				}
			}
		}
		return true
	})

	return files
}
