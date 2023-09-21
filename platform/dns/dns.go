package dns

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/guodf/canary/filesys"
)

var assetsPath string
var dnsPath string
var confPath string

func Init(path string) error {
	if len(path) > 0 {
		assetsPath = path
	} else {
		assetsPath = getAssetsPath()
	}
	if !filesys.Exists(assetsPath) {
		os.MkdirAll(assetsPath, os.ModePerm)
	}
	dnsPath = filepath.Join(assetsPath, dnsName)
	confPath = filepath.Join(assetsPath, "dns.conf")
	os.WriteFile(dnsPath, clientData, os.ModePerm)
	os.WriteFile(confPath, confData, os.ModePerm)
	cmd := exec.Command(dnsPath, "-conf", confPath)

	return cmd.Start()
}
