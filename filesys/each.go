package filesys

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"

	"tango/tuple"
)

type FileType int

const (
	File FileType = 1
	Dir  FileType = 2
)

type EachCallback func(fullPath string, fsInfo fs.FileInfo) bool

// ScanByDepth 深度优先
func ScanByDepth(dirPath string, maxDepth int, tupleChan chan *tuple.Tuple2) {
	if !IsDir(dirPath) {
		tupleChan <- nil
		return
	}
	go scanByDepth(dirPath, 0, maxDepth, tupleChan)
}

// ScanByBreadth 广度优先
func ScanByBreadth(dirPath string, maxDepth int, cb EachCallback) error {
	if !IsDir(dirPath) {
		return fmt.Errorf("dirPath:%s is not dir", dirPath)
	}

	curDepth := 0
	dirs := []string{dirPath}
	for {
		if len(dirs) == 0 || (maxDepth > 0 && curDepth > maxDepth) {
			return nil
		}
		curPath := dirs[0]
		filesInfo, e := ioutil.ReadDir(curPath)
		dirs = dirs[1:]
		if e != nil {
			log.Println(e)
			return e
		}

		for _, fsInfo := range filesInfo {
			fullPath := filepath.Join(curPath, fsInfo.Name())
			if fsInfo.IsDir() {
				dirs = append(dirs, fullPath)
			}
			if !cb(fullPath, fsInfo) {
				return nil
			}
		}
		curDepth++
	}
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
