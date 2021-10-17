package filesys

import "path/filepath"
import "github.com/guodf/goutil/tuple"

// FindFilesInExt 根据扩展名查找文件
func FindFilesInExt(dirPath string, exts []string, maxDepth int) []string {
	var files []string
	if !IsDir(dirPath) {
		return files
	}
	tupleChan:=make(chan *tuple.Tuple2)
	go ScanByBreadth(dirPath, maxDepth,tupleChan)
	for tuple2 := range tupleChan {
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
