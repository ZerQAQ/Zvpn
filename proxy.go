package Zvpn

type Proxy interface {
	// return the Conn ready to forward data
	// 本地连接: 浏览器->本地服务器
	// 远程连接: 本地服务器->远程服务器
	// 传入一个本地连接和远程连接，返回一个完成代理握手的远程连接
	// 如果本地和远程的代理协议一样，那么直接return remote即可
	ClientHandshake(local, remote Conn, p Protocol) (Conn, error)
	// take a established Conn between client and server as argument
	// return the established Conn which is the client want to established
	// 传入一个已经完成可靠网络协议握手的conn，在这个conn上进行
	// 代理协议握手，并且返回一个按客户端要求在服务端上建立的连接
	ServerHandshake(Conn, Protocol) (Conn, error)
}

//var Sock5 Proxy = &Sock5Proxy{}
