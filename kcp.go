package main

import (
	"github.com/ZerQAQ/Zvpn/Zvpn"
	"github.com/xtaci/kcp-go"
	"net"
)

type KCP struct{}

type KCPListener struct{l net.Listener}

func (l KCPListener) Accept() (Zvpn.Conn, error) {
	return l.Accept()
}

func (l KCPListener) Close() error {
	return l.Close()
}

func (KCP) Listen(addr string) (Zvpn.Listener, error) {
	l, err := kcp.Listen(addr)
	return KCPListener{l}, err
}

func (KCP) Dial(addr string) (Zvpn.Conn, error) {
	return kcp.Dial(addr)
}