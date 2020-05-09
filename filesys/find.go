package filesys

import "path/filepath"

// FindFilesInExt 根据扩展名查找文件
func FindFilesInExt(dirPath string, exts []string, maxDepth int) []string {
	var files []string
	if !IsDir(dirPath) {
		return files
	}

	for tuple2 := range ScanByBreadth(dirPath, maxDepth) {
		fileType := tuple2.Item1.(FileType)
		filePath := tuple2.Item2.(string)
		if fileType == File {
			for _, ext := range exts {
				if ext == filepath.Ext(filePath) {
					files = append(files, filePath)
				}
			}
		}
	}

	return files
}
