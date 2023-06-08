package filesys

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/h2non/filetype"
)

// Ext 获取文件扩展类型
//
//	filePath: 文件必须存在，如果 filepath 不存在则使用 ext 赋值 filetype_ext
//
//	ext: 文件后缀
//	filetype_ext: 通过 filetype推断的文件后缀
//	same: ext==filetype_ext
func Ext(filePath string) (ext string, filetype_ext string, same bool) {
	ext = filepath.Ext(filePath)
	ext = strings.ToLower(ext)
	if !Exists(filePath) {
		return ext, ext, true
	}
	buf, _ := ioutil.ReadFile(filePath)

	// 检查文件类型
	kind, unknown := filetype.Match(buf)
	if unknown != nil {
		log.Println("Unknown file type:", unknown)
		return
	}
	filetype_ext = fmt.Sprintf(".%s", kind.Extension)
	return ext, filetype_ext, ext == filetype_ext
}
