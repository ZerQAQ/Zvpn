package protocol

import (
	"github.com/ZerQAQ/Zvpn"
	"net"
)

type TCPListener struct{ l net.Listener }

func (l TCPListener) Accept() (Zvpn.Conn, error) { return l.l.Accept() }
func (l TCPListener) Close() error          { return l.l.Close() }

type TCPProtocol struct{}

func (TCPProtocol) Bind(addr string) (Zvpn.Listener, error) {
	l, err := net.Listen("tcp", addr)
	return TCPListener{l}, err
}

func (TCPProtocol) Dial(addr string) (Zvpn.Conn, error) {
	return net.Dial("tcp", addr)
}

var TCP Zvpn.Protocol = &TCPProtocol{}