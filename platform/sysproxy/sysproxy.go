package sysproxy

import (
	_ "embed"
	"log"
	"os/exec"
	"sync"

	"github.com/getlantern/byteexec"
)

type PacProxy struct {
	Enabled bool
	PacAddr string
}
type LocalProxy struct {
	Enabled    bool
	Server     string
	BYPASSAddr string
}

type SysProxySetting struct {
	EnableAutoProxy bool
	PacProxy        PacProxy
	LocalProxy      LocalProxy
}

var (
	mu sync.Mutex
	be *byteexec.Exec
)

func init() {
	proxy, err := byteexec.New(sysproxy, "sysproxy")
	if err != nil {
		log.Printf("unable to extract helper tool: %v\n", err)
	}
	be = proxy
	ensureElevatedOnDarwin(be)
}

func run(cmd *exec.Cmd) bool {
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("unable to execute %v: %s\n%s", cmd.Path, err, string(out))
		return false
	}
	log.Println("Command %v output %v", cmd.Path, string(out))
	return true
}

func DisableProxy() bool {
	return run(be.Command("off"))
}

func SetGlobalProxy(host, bypass string) bool {
	return run(be.Command("global", host, bypass))
}

func SetPacProxy(pac string) bool {
	return run(be.Command("pac", pac))
}
