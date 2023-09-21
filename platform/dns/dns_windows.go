package dns

import (
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed assets/windows/dns.exe
var clientData []byte

//go:embed assets/dns.conf
var confData []byte

var dnsName = "dns.exe"

func getAssetsPath() string {
	return filepath.Join(os.Getenv("APPDATA"), "canary")
}
