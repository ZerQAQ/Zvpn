package obfus

import (
	"github.com/ZerQAQ/Zvpn/Zvpn"
)

type DoNothing struct{}

func (DoNothing) ClientHandShake(conn Zvpn.Conn) (Zvpn.Encrypter, Zvpn.Decrypter, error) {
	return DoNothing{}, DoNothing{}, nil
}
func (DoNothing) ServerHandShake(conn Zvpn.Conn) (Zvpn.Encrypter, Zvpn.Decrypter, error) {
	return DoNothing{}, DoNothing{}, nil
}
func (DoNothing) Write(conn Zvpn.Conn, src []byte) (int, error) { return conn.Write(src) }
func (DoNothing) Read(conn Zvpn.Conn, dst []byte) (int, error)  { return conn.Read(dst) }

func NewNothing(k []byte) Zvpn.Obfuscate {
	return DoNothing{}
}
