package sysproxy

import (
	_ "embed"
	"os/exec"
	"syscall"

	"github.com/getlantern/byteexec"
)

//go:embed binaries/linux_amd64/sysproxy
var sysproxy []byte

func ensureElevatedOnDarwin(be *byteexec.Exec) (err error) {
	return nil
}

func detach(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}
