package filesys

import (
	"github.com/guodf/goutil/tuple"
	"io/ioutil"
	"log"
	"path/filepath"
)

type FileType int

const (
	File FileType = 1
	Dir  FileType = 2
)

// ScanByDepth 深度优先
func ScanByDepth(dirPath string, maxDepth int, tupleChan chan *tuple.Tuple2) {
	if !IsDir(dirPath) {
		tupleChan <- nil
		return
	}
	go scanByDepth(dirPath, 0, maxDepth, tupleChan)
}

// 深度遍历
func scanByDepth(dirPath string, curDepth int, maxDepth int, tupleChan chan *tuple.Tuple2) {
	defer func() {
		if curDepth == 0 {
			tupleChan <- nil
		}
	}()
	if maxDepth > 0 && curDepth > maxDepth {
		return
	}
	filesInfo, e := ioutil.ReadDir(dirPath)
	if e != nil {
		tupleChan <- nil

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
func ScanByBreadth(dirPath string, maxDepth int, tupleChan chan *tuple.Tuple2) {
	if !IsDir(dirPath) {
		close(tupleChan)
		return
	}

	curDepth := 0
	dirs := []string{dirPath}
	for {
		if len(dirs) == 0 || (curDepth > maxDepth && maxDepth > 0) {
			close(tupleChan)
			return
		}
		curPath := dirs[0]
		filesInfo, e := ioutil.ReadDir(curPath)
		dirs = dirs[1:]
		if e != nil {
			log.Println(e)
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
	close(tupleChan)

}
