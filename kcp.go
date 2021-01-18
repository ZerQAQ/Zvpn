package main

import (
	"github.com/ZerQAQ/Zvpn/Zvpn"
	"github.com/xtaci/kcp-go"
)

var KCP Zvpn.Protocol = _KCP{}
type _KCP struct{}
func (_KCP) Bind(addr string) (Zvpn.Listener, error) {
	l, err := kcp.Listen(addr)
	return Zvpn.NetListenerWrapper{L: l}, err
}

func (_KCP) Dial(addr string) (Zvpn.Conn, error) {
	return kcp.Dial(addr)
}