package protocol

import (
	"github.com/ZerQAQ/Zvpn/Zvpn"
	"net"
)

type TCPProtocol struct{}

func (TCPProtocol) Bind(addr string) (Zvpn.Listener, error) {
	l, err := net.Listen("tcp", addr)
	return Zvpn.NetListenerWrapper{L: l}, err
}

func (TCPProtocol) Dial(addr string) (Zvpn.Conn, error) {
	return net.Dial("tcp", addr)
}

var TCP Zvpn.Protocol = &TCPProtocol{}