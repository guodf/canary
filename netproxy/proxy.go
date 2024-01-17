package netproxy

import (
	"net"
)

type ProxyContext struct {
	Port    int
	Host    string
	conn    net.Conn
	tagConn net.Conn
}

func NewProxy(conn net.Conn) *ProxyContext {
	var proxyContext *ProxyContext
	// bfr := bufio.NewReader(conn)
	// firstByte, _ := bfr.ReadByte()
	// if firstByte == Version4 {
	// 	proxyContext = NewSocksV4(conn)
	// } else if firstByte == Version5 {
	// 	proxyContext = NewSocks5(conn)
	// } else {
	//  proxyContext = NewHttpProxy(conn)
	// }
	proxyContext = NewHttpProxy(conn)
	return proxyContext
}
