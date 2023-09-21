package netutil

import (
	"net"
	"time"
)

func TcpForward(src, dst net.Conn) {
	dstConn := dst.(*net.TCPConn)
	dstConn.SetKeepAlive(true) // 开启保持链接
	// tagTcpConn.SetKeepAlivePeriod(time.Second * 20)          // 链接保持时间，默认15s
	// tagTcpConn.SetDeadline(time.Now().Add(time.Second * 60)) // 读写数据超时时间
	// dstConn.SetNoDelay(true) // 设置延迟发送，提高传输效率
	dstConn.SetLinger(0) // 关闭链接时是否立即关闭

	srConn := src.(*net.TCPConn)
	srConn.SetKeepAlive(true) // 开启保持链接
	// srConn.SetKeepAlivePeriod(time.Second * 20) // 链接保持时间
	// srConn.SetDeadline(time.Now().Add(time.Second * 60)) // 读写数据超时时间
	// srConn.SetNoDelay(true) // 设置延迟发送，提高传输效率
	srConn.SetLinger(0) // 关闭链接时是否立即关闭
	defer func() {
		srConn.Close()
	}()
	data := [1024]byte{}
	for {
		count, e := srConn.Read(data[:1024])
		if e != nil {
			return
		}
		if count == 0 {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		dstConn.Write(data[0:count])
	}
}
