package Zvpn

type Obfuscate interface {
	// 传入一个已经握手完毕的基于可靠传输协议的连接
	// 你可以在上面利用非对称加密交换密秘钥或者直接
	// 使用本地配置文件中的密钥
	ClientHandShake(Conn) (Encrypter, Decrypter, error)
	ServerHandShake(Conn) (Encrypter, Decrypter, error)
}

type Encrypter interface {
	Write(conn Conn, src []byte) (int, error)
}

type Decrypter interface {
	Read(conn Conn, dst []byte) (int, error)
}
