package sysproxy

import (
	_ "embed"
	"os/exec"
	"syscall"

	"github.com/getlantern/byteexec"
)

//go:embed binaries/windows/sysproxy64.exe
var sysproxy []byte

func ensureElevatedOnDarwin(be *byteexec.Exec) (err error) {
	return nil
}

func detach(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}
