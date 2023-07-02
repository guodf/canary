package sysproxy

import (
	_ "embed"
	"fmt"
	"log"
	"os/exec"
	"sync"

	"github.com/getlantern/byteexec"
)

var (
	mu sync.Mutex
	be *byteexec.Exec
)

func init() {
	proxy, err := byteexec.New(sysproxy, "easy_proxy")
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

func EnableProxy() bool {
	//flags := 2 | 4 | 8 // 启用所有代理：自动检测，pac，手动代理
	flags := 8
	return run(be.Command("set", fmt.Sprintf("%d", flags), "-", "-", "-"))
}

func DisableProxy() bool {
	// flags: 1:关闭所有代理 2:开启手动代理 4:开启脚本代理 8:开启自动检测代理设置
	return run(be.Command("set", "1", "-", "-", "-"))
}

func SetGlobalProxy(host, bypass string) bool {
	return run(be.Command("global", host, bypass))
}

func SetPacProxy(pac string) bool {
	return run(be.Command("pac", pac))
}
