package obfus

import "github.com/ZerQAQ/Zvpn/protocol"

type Obfuscate interface {
	ClientHandShake(protocol.Conn) (Encrypter, Decrypter, error)
	ServerHandShake(protocol.Conn) (Encrypter, Decrypter, error)
}

type Encrypter interface {
	Write(conn protocol.Conn, src []byte) (int, error)
}

type Decrypter interface {
	Read(conn protocol.Conn, dst []byte) (int, error)
}
