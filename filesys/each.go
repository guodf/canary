package filesys

import (
	"io/ioutil"
	"path/filepath"

	"github.com/guodf/goutil/tuple"
)

type FileType int

const (
	File FileType = 1
	Dir  FileType = 2
)

// ScanByDepth 深度优先
func ScanByDepth(dirPath string, maxDepth int) chan *tuple.Tuple2 {
	if !IsDir(dirPath) {
		return nil
	}
	tupleChan := make(chan *tuple.Tuple2)

	go scanByDepth(dirPath, 0, maxDepth, tupleChan)

	return tupleChan
}

// 深度遍历
func scanByDepth(dirPath string, curDepth int, maxDepth int, tupleChan chan *tuple.Tuple2) {
	defer func() {
		if curDepth == 0 {
			close(tupleChan)
		}
	}()
	if maxDepth > 0 && curDepth > maxDepth {
		return
	}
	filesInfo, e := ioutil.ReadDir(dirPath)
	if e != nil {
		close(tupleChan)
		return
	}
	for _, fileInfo := range filesInfo {
		filePath := filepath.Join(dirPath, fileInfo.Name())
		if fileInfo.IsDir() {

			tupleChan <- tuple.NewTuple2(Dir, filePath)

			scanByDepth(filePath, curDepth+1, maxDepth, tupleChan)
		} else {

			tupleChan <- tuple.NewTuple2(File, filePath)

		}
	}
}

// ScanByBreadth 广度优先
func ScanByBreadth(dirPath string, maxDepth int) chan *tuple.Tuple2 {
	if !IsDir(dirPath) {
		return nil
	}
	tupleChan := make(chan *tuple.Tuple2)
	curDepth := 0
	dirs := []string{dirPath}
	go func() {
		defer close(tupleChan)
		for {
			if len(dirs) == 0 || (curDepth > maxDepth && maxDepth > 0) {
				return
			}
			curPath := dirs[0]
			filesInfo, e := ioutil.ReadDir(curPath)
			dirs = dirs[1:]
			if e != nil {
				close(tupleChan)
				return
			}
			for _, fileInfo := range filesInfo {
				filePath := filepath.Join(curPath, fileInfo.Name())
				if fileInfo.IsDir() {
					tupleChan <- tuple.NewTuple2(Dir, filePath)
					dirs = append(dirs, filePath)
				} else {
					tupleChan <- tuple.NewTuple2(File, filePath)
				}
			}
			curDepth++
		}
	}()

	return tupleChan
}
