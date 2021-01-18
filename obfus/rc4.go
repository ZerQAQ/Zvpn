package obfus

import (
	"crypto/rand"
	"crypto/rc4"
	"github.com/ZerQAQ/Zvpn/Zvpn"
)

const (
	RC4MaxMasterKeySize = 224
)

const (
	secondaryKeySize = 32
)

type _rc4 struct {
	MasterKey []byte
}

type _rc4Session struct {
	cipher *rc4.Cipher
}

func (s _rc4) ClientHandShake(conn Zvpn.Conn) (Zvpn.Encrypter, Zvpn.Decrypter, error) {
	return s.HandShake(conn)
}
func (s _rc4) ServerHandShake(conn Zvpn.Conn) (Zvpn.Encrypter, Zvpn.Decrypter, error) {
	return s.HandShake(conn)
}

func (s _rc4) HandShake(conn Zvpn.Conn) (Zvpn.Encrypter, Zvpn.Decrypter, error) {
	secKey := make([]byte, secondaryKeySize)
	_, err := rand.Read(secKey)
	if err != nil {
		return nil, nil, err
	}
	_, err = conn.Write(secKey)
	if err != nil {
		return nil, nil, err
	}
	e, err := rc4.NewCipher(append(s.MasterKey, secKey...))
	if err != nil {
		return nil, nil, err
	}

	_, err = conn.Read(secKey)
	if err != nil {
		return nil, nil, err
	}
	d, err := rc4.NewCipher(append(s.MasterKey, secKey...))
	if err != nil {
		return nil, nil, err
	}

	return _rc4Session{e}, _rc4Session{d}, nil
}

func (s _rc4Session) Write(conn Zvpn.Conn, src []byte) (int, error) {
	s.cipher.XORKeyStream(src, src)
	return conn.Write(src)
}

func (s _rc4Session) Read(conn Zvpn.Conn, dst []byte) (l int, err error) {
	l, err = conn.Read(dst)
	s.cipher.XORKeyStream(dst, dst[:l])
	return
}

//masKey的大小不超过RC4MaxMasterKeySize,越大越安全
//the size of masKey is not bigger than RC4MaxMasterKeySize
//bigger mean safer
func NewRC4(mstKey []byte) Zvpn.Obfuscate {
	return &_rc4{mstKey}
}
