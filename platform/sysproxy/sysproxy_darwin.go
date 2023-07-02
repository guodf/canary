package sysproxy

import (
	_ "embed"
	"fmt"
	"os/exec"
	"syscall"

	"github.com/getlantern/byteexec"
	"github.com/getlantern/elevate"
)

// Note this is a universal binary that runs on amd64 and arm64
//
//go:embed binaries/darwin/sysproxy
var sysproxy []byte

func ensureElevatedOnDarwin(be *byteexec.Exec) (err error) {
	var s syscall.Stat_t
	// we just checked its existence, not bother checking specific error again
	if err = syscall.Stat(be.Filename, &s); err != nil {
		return fmt.Errorf("error starting helper tool %s: %v", be.Filename, err)
	}
	if s.Mode&syscall.S_ISUID > 0 && s.Uid == 0 && s.Gid == 0 {
		log.Pringf("%v is already owned by root:wheel and has setuid bit on\n", be.Filename)
		return
	}
	cmd := elevate.Command(be.Filename, "setuid")
	return run(cmd)
}

func detach(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}
