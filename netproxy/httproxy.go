package netproxy

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/guodf/canary/netutil"
)

func NewHttpProxy(conn net.Conn) *ProxyContext {
	proxyContext := &ProxyContext{
		conn: conn,
	}
	bfr := bufio.NewReader(conn)
	endlineBytes, e := bfr.ReadBytes('\n')
	if len(endlineBytes) == 0 {
		return nil
	}
	if e != nil {
		log.Println("reader first line failed")
		log.Println(e)
		return nil
	}
	firstLine := []byte{}
	firstLine = append(firstLine, endlineBytes...)
	bs := bytes.TrimSuffix(firstLine, []byte{'\r', '\n'})
	firstLineArr := bytes.Split(bs, []byte{' '})
	if string(firstLineArr[0]) == "CONNECT" {
		conn.Write(firstLineArr[2])
		fmt.Fprint(conn, " 200 Connection established\r\n\r\n")
	} else {
		conn.Write(firstLine)
	}
	u := string(firstLineArr[1])
	var uri *url.URL
	if strings.Index(u, "http://") == 0 || strings.Index(u, "https://") == 0 {
		uri, _ = url.Parse(u)
	} else {
		uri, _ = url.Parse("http://" + u)
	}
	if uri.Port() == "" {
		proxyContext.Port = 80
	} else {
		p, _ := strconv.Atoi(uri.Port())
		proxyContext.Port = p
	}
	proxyContext.Host = uri.Hostname()
	tagConn, e := net.Dial("tcp", fmt.Sprintf("%s:%d", proxyContext.Host, proxyContext.Port))
	if e != nil {
		return nil
	}
	proxyContext.tagConn = tagConn
	go netutil.TcpForward(conn, tagConn)
	go netutil.TcpForward(tagConn, conn)
	return proxyContext
}
