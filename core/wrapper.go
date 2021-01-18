package core

import (
	"github.com/ZerQAQ/Zvpn"
)

func NewProtocol(pro Zvpn.Protocol, o Zvpn.Obfuscate) Zvpn.Protocol {
	return ProtocolWrapper{pro, o}
}

type ProtocolWrapper struct {
	pro Zvpn.Protocol
	o   Zvpn.Obfuscate
}

func (p ProtocolWrapper) Bind(addr string) (listener Zvpn.Listener, err error) {
	listener, err = p.pro.Bind(addr)
	return ListerWrapper{listener, p.o}, err
}
func (p ProtocolWrapper) Dial(addr string) (conn Zvpn.Conn, err error) {
	conn, err = p.pro.Dial(addr)
	if err != nil {
		return nil, err
	}
	e, d, err := p.o.ClientHandShake(conn)
	if err != nil {
		return nil, err
	}
	return ConnWrapper{conn, e, d}, err
}

type ListerWrapper struct {
	listener Zvpn.Listener
	o        Zvpn.Obfuscate
}

func (l ListerWrapper) Accept() (Zvpn.Conn, error) {
	conn, err := l.listener.Accept()
	if err != nil {
		return nil, err
	}
	e, d, err := l.o.ServerHandShake(conn)
	if err != nil {
		return nil, err
	}
	return ConnWrapper{conn, e, d}, nil
}
func (l ListerWrapper) Close() error { return l.listener.Close() }

type ConnWrapper struct {
	conn Zvpn.Conn
	e    Zvpn.Encrypter
	d    Zvpn.Decrypter
}

func (c ConnWrapper) Read(buf []byte) (retLen int, retErr error) {
	return c.d.Read(c.conn, buf)
}
func (c ConnWrapper) Write(buf []byte) (retLen int, retErr error) {
	return c.e.Write(c.conn, buf)
}
func (c ConnWrapper) Close() error { return c.conn.Close() }
