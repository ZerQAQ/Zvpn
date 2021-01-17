package protocol

import "net"

type TCPListener struct{ l net.Listener }

func (l TCPListener) Accept() (Conn, error) { return l.l.Accept() }
func (l TCPListener) Close() error          { return l.l.Close() }

type TCPProtocol struct{}

func (TCPProtocol) Bind(addr string) (Listener, error) {
	l, err := net.Listen("tcp", addr)
	return TCPListener{l}, err
}

func (TCPProtocol) Dial(addr string) (Conn, error) {
	return net.Dial("tcp", addr)
}
