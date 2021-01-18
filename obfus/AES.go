package obfus

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"github.com/ZerQAQ/Zvpn/Zvpn"
	"github.com/ZerQAQ/Zvpn/lib"
)

const (
	maxBloNum = 64
)

type _AES struct {
	block cipher.Block
}
type _AESSession struct {
	blockMode cipher.BlockMode
}

func (s _AES) ClientHandShake(conn Zvpn.Conn) (Zvpn.Encrypter, Zvpn.Decrypter, error) {
	return s.HandShake(conn)
}
func (s _AES) ServerHandShake(conn Zvpn.Conn) (Zvpn.Encrypter, Zvpn.Decrypter, error) {
	return s.HandShake(conn)
}

func (s _AES) HandShake(conn Zvpn.Conn) (e Zvpn.Encrypter, d Zvpn.Decrypter, err error) {
	iv := lib.RandomBytes(s.block.BlockSize())
	if _, err = conn.Write(iv); err != nil {
		return
	}
	e = _AESSession{cipher.NewCBCEncrypter(s.block, iv)}

	if _, err = conn.Read(iv); err != nil {
		return
	}
	d = _AESSession{cipher.NewCBCDecrypter(s.block, iv)}
	return
}

func (s _AESSession) Write(conn Zvpn.Conn, src []byte) (l int, err error) {
	// check if len(src) is smaller than maxBloNum * blockSize
	// if not, _Write multiply time
	bloNum := (len(src)-1)/s.blockMode.BlockSize() + 1
	if bloNum > 64 {
		temp := (bloNum-1)/maxBloNum + 1
		for i := 0; i < temp; i++ {
			if i == temp-1 {
				_, err = s._Write(conn, src[maxBloNum*i:])
			} else {
				_, err = s._Write(conn, src[maxBloNum*i:maxBloNum*(i+1)])
			}
			if err != nil {
				return
			}
		}
	} else {
		return s._Write(conn, src)
	}
	return len(src), nil
}

// len(src) must smaller than maxBloNum * blockSize
func (s _AESSession) _Write(conn Zvpn.Conn, src []byte) (int, error) {
	bloNum := uint64(len(src)/s.blockMode.BlockSize() + 1)
	//first block is size
	cypherText := make([]byte, s.blockMode.BlockSize())
	for i := 0; i < 8; i++ {
		cypherText[i] = byte(bloNum & 0xff)
		bloNum >>= 8
	}
	//add plain text
	cypherText = append(cypherText, src...)
	//add padding
	paddingSize := s.blockMode.BlockSize()
	if len(src)%s.blockMode.BlockSize() != 0 {
		paddingSize = s.blockMode.BlockSize() - len(src)%s.blockMode.BlockSize()
	}
	cypherText = append(cypherText, bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)...)
	//crypt
	s.blockMode.CryptBlocks(cypherText, cypherText)
	if _, err := conn.Write(cypherText); err != nil {
		return 0, err
	}

	return len(src), nil
}

// len(dst) must be larger than maxBloNum * max(blockSize),
// which is 64 * 32 = 2048
func (s _AESSession) Read(conn Zvpn.Conn, dst []byte) (int, error) {
	//get bloNum
	temp := make([]byte, s.blockMode.BlockSize())
	if err := s.blockedRead(conn, temp); err != nil {
		return 0, err
	}
	s.blockMode.CryptBlocks(temp, temp)
	bloNum := uint64(0)
	for i := 7; i >= 0; i-- {
		bloNum <<= 8
		bloNum |= uint64(temp[i])
	}
	//receive data
	l := int(bloNum) * s.blockMode.BlockSize()
	dst = dst[:l]
	if err := s.blockedRead(conn, dst); err != nil {
		return 0, err
	}
	s.blockMode.CryptBlocks(dst, dst)

	return l - int(dst[l-1]), nil
}

func (s _AESSession) blockedRead(conn Zvpn.Conn, buf []byte) error {
	l, err := conn.Read(buf)
	if err != nil {
		return err
	}
	for l != len(buf) {
		buf = buf[l:]
		l, err = conn.Read(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

//The key argument should be the AES key, either 16, 24,
//or 32 bytes to select AES-128, AES-192, or AES-256.
func newAES(key []byte) (Zvpn.Obfuscate, error) {
	b, err := aes.NewCipher(key)
	return _AES{b}, err
}

func NewAES256(key []byte) Zvpn.Obfuscate {
	if len(key) != 32 {
		panic("key length must equal 32 bytes")
	}
	b, _ := aes.NewCipher(key)
	return _AES{b}
}
