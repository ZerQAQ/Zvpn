package obfus

import "github.com/ZerQAQ/Zvpn/protocol"

type DoNothing struct{}

func (DoNothing) ClientHandShake(conn protocol.Conn) (Encrypter, Decrypter, error) {
	return DoNothing{}, DoNothing{}, nil
}
func (DoNothing) ServerHandShake(conn protocol.Conn) (Encrypter, Decrypter, error) {
	return DoNothing{}, DoNothing{}, nil
}
func (DoNothing) Write(conn protocol.Conn, src []byte) (int, error) { return conn.Write(src) }
func (DoNothing) Read(conn protocol.Conn, dst []byte) (int, error)  { return conn.Read(dst) }

func NewNothing(k []byte) Obfuscate {
	return DoNothing{}
}
