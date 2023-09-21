package app

import (
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
)

// RootPath 获取程序启动目录
func RootPath() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(filepath.Dir(file))
}

func AppExitListen(exit func()) {
	// 创建一个通道来接收信号
	sigs := make(chan os.Signal, 1)
	// os.Interrupt 用于windows任务管理器结束进程, 实际无效果
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// 等待信号
	<-sigs
	if exit != nil {
		exit()
	}
	// 退出进程，windows api函数，退出时不会做清理动作
	// if h, e := syscall.GetCurrentProcess(); e != nil {
	// 	syscall.TerminateProcess(h, 1)
	// }
	// 退出进程，用于所有平台，执行后defer会被正常调用
	os.Exit(0)
}
