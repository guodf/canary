//							socks v4
//  http://ftp.icm.edu.pl/packages/socks/socks4/SOCKS4.protocol
//	握手数据
//			+----+----+----+----+----+----+----+----+----+----+....+----+
//			| VN | CD | DSTPORT |      DSTIP        | USERID       |NULL|
//			+----+----+----+----+----+----+----+----+----+----+....+----+
// 长度		   1    1      2              4            可变           1
// VN: 4
// CD:
//		1:CONNECT
//		2:BIND
// DSTPORT: 端口
// DSTIP:	ipv4
// USERID:	用户信息
// NULL:	结束位为0

//	 握手响应响应
//				+----+----+----+----+----+----+----+----+
//				| VN | CD | DSTPORT |      DSTIP        |
//				+----+----+----+----+----+----+----+----+
//
// 长度        1    1      2              4
// VN: 固定不变值0
// CD:
//
//	0x5A为允许；
//	0x5B为拒绝或失败；
//	0x5C为请求被拒绝，因为SOCKS服务器无法连接到客户端;
//	0x5D为请求被拒绝，因为USERID不匹配；
//
// DSTPORT: 端口
// DSTIP:	ipv4
package socks

import (
	"bufio"
	"errors"
	"io"
)

const (
	VN4 = 0x04
	VN5 = 0x05
)

// CMD 命令
const (
	CONNECT = 0x01
	BIND    = 0x02
)

type SocksV4 struct {
	VER      byte
	CMD      byte
	Port     [2]byte
	IP       [4]byte
	UserData []byte
	End      byte
}

var ErrorNoSocksV4 = errors.New("not socks v4")

func NewSocksV4(bfr *bufio.Reader) (*SocksV4, error) {

	socksV4 := &SocksV4{
		VER:      4,
		CMD:      0,
		Port:     [2]byte{},
		IP:       [4]byte{},
		UserData: nil,
		End:      0,
	}
	cmd, e := bfr.ReadByte()
	socksV4.CMD = cmd
	if e != nil {
		return nil, e
	}
	b := make([]byte, 2)
	length, e := bfr.Read(b)
	if e != nil || length != 2 {
		return nil, e
	}
	copy(socksV4.Port[:], b)

	b = make([]byte, 4)
	length, e = bfr.Read(b)
	if e != nil || length != 4 {
		return nil, e
	}
	copy(socksV4.IP[:], b)

	bs, e := bfr.ReadBytes(0)
	if e == io.EOF {
		socksV4.End = 0
		return socksV4, nil
	}
	if e != nil {
		return nil, e
	}
	socksV4.UserData = bs
	socksV4.End = 0
	return socksV4, nil
}

func (socksV4 *SocksV4) Accept() []byte {

	var resp []byte
	resp = append(resp, 0x00)
	resp = append(resp, 0x5A)
	for _, b := range socksV4.Port {
		resp = append(resp, b)
	}
	for _, b := range socksV4.IP {
		resp = append(resp, b)
	}
	return resp
}

func (socksV4 *SocksV4) Failed() []byte {
	var resp []byte
	resp = append(resp, 0x00)
	resp = append(resp, 0x5B)
	for _, b := range socksV4.Port {
		resp = append(resp, b)
	}
	for _, b := range socksV4.IP {
		resp = append(resp, b)
	}
	return resp
}
