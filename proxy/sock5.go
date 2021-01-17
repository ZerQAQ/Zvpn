package proxy

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ZerQAQ/Zvpn"
	"github.com/ZerQAQ/Zvpn/lib"
	"log"
	"net"
)

type Sock5Proxy struct{}

var ErrHandShakeFail = errors.New("sock5 handshake fail")

func (Sock5Proxy) ClientHandshake(l, r Zvpn.Conn, p Zvpn.Protocol) (Zvpn.Conn, error) {
	return r, nil
}

func (Sock5Proxy) ServerHandshake(conn Zvpn.Conn, p Zvpn.Protocol) (retConn Zvpn.Conn, retErr error) {
	retConn = nil
	retErr = ErrHandShakeFail

	recvBuf := make([]byte, 1024)
	/*
		接受sock5请求
		+----+----------+----------+
		|VER | NMETHODS | METHODS  |
		+----+----------+----------+
		| 1  |    1     | 1 to 255 |
		+----+----------+----------+
	*/
	bufLen, err := conn.Read(recvBuf)
	if err != nil {
		return nil, errors.New("err in sock5 read1: " + err.Error())
	}
	if bufLen != 3 {
		err = errors.New("err in sock5 read1, bufLen != 3")
	}
	if recvBuf[0] != 5 {
		err = errors.New("err in sock5 read1, recvBuf[0] != 5")
	}
	if err != nil {
		fmt.Println(recvBuf[:bufLen])
		return nil, err
	}

	/*
		server回应1
		+----+--------+
		|VER | METHOD |
		+----+--------+
		| 1  |   1    |
		+----+--------+
	*/
	_, err = conn.Write([]byte{5, 0})
	if err != nil {
		return nil, err
	}
	/*
		client请求代理建立连接
		+----+-----+-------+------+----------+----------+
		|VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
		+----+-----+-------+------+----------+----------+
		| 1  |  1  | X'00' |  1   | Variable |    2     |
		+----+-----+-------+------+----------+----------+
	*/
	bufLen, err = conn.Read(recvBuf)
	if lib.ErrHandler("err in sock5 read2", err) {
		return
	}
	if len(recvBuf) < 10 {
		return
	}
	if recvBuf[0] != 5 {
		return
	}

	cmd := recvBuf[1]
	aType := recvBuf[3]
	var dIp []byte
	dPort := recvBuf[bufLen-2 : bufLen]
	desAddrRaw := recvBuf[5 : bufLen-2]
	log.Println("receive request to", string(desAddrRaw))
	//解析IP
	switch aType {
	case 1:
		// IPV4
		dIp = desAddrRaw
	case 3:
		// 域名
		temp, err := net.ResolveIPAddr("ip4", string(desAddrRaw))
		if lib.ErrHandler("err when resolving ip", err) {
			return
		}
		dIp = temp.IP
	default:
		return
	}

	desAddr := net.TCPAddr{
		IP:   dIp,
		Port: int(binary.BigEndian.Uint16(dPort)),
	}

	switch cmd {
	case 1:
		//connect
		desConn, err := net.DialTCP("tcp4", nil, &desAddr)
		if lib.ErrHandler("err when dial", err) {
			return
		}
		/*
			server回应2
			+----+-----+-------+------+----------+----------+
			|VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
			+----+-----+-------+------+----------+----------+
			| 1  |  1  | X'00' |  1   | Variable |    2     |
			+----+-----+-------+------+----------+----------+
		*/
		sendBuf := []byte{5, 0, 0, 1}
		sendBuf = append(sendBuf, dIp[len(dIp)-4:]...)
		sendBuf = append(sendBuf, dPort...)
		_, err = conn.Write(sendBuf)
		if lib.ErrHandler("err in sock5 write2", err) {
			return
		}
		return desConn, nil
	default:
		return
	}
}

var Sock5 Zvpn.Proxy = &Sock5Proxy{}